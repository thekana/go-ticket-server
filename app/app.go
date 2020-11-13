package app

import (
	"crypto/rsa"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/jackc/pgconn"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	customError "ticket-reservation/custom_error"
	"ticket-reservation/db"
	"ticket-reservation/db/model"
	"ticket-reservation/log"
	"ticket-reservation/redis_cache"
	"ticket-reservation/utils"
	"time"
)

type MyStruct struct {
	QueueChan           chan *ReservationQueueElem
	Signal              chan struct{}
	ClearBatchTicker    *time.Ticker
	UpdateDBEventTicker *time.Ticker
	Batch               chan *ReservationQueueElem
}

type App struct {
	Logger                log.Logger
	Config                *Config
	TokenSignerPrivateKey *rsa.PrivateKey
	TokenSignerPublicKey  *rsa.PublicKey
	DB                    db.DB
	RedisCache            redis_cache.Cache
	My                    *MyStruct
}

const (
	BatchSize int = 50
	TickTime      = time.Millisecond * 100
)

var (
	uni      *ut.UniversalTranslator
	trans    ut.Translator
	validate *validator.Validate
)

func init() {
	en := en.New()
	uni = ut.New(en, en)

	// this is usually known or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	translator, found := uni.GetTranslator("en")
	if !found {
		panic("translator not found")
	}

	validate = validator.New()

	if err := entranslations.RegisterDefaultTranslations(validate, translator); err != nil {
		panic(err)
	}

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	trans = translator
}

func New(logger log.Logger) (app *App, err error) {
	app = &App{
		Logger: logger,
	}

	app.Config, err = InitConfig()
	if err != nil {
		return nil, err
	}

	app.TokenSignerPrivateKey, err = readRSAPrivateKey(app.Config.TokenSignerPrivateKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load token signer private key")
	}
	app.TokenSignerPublicKey, err = readRSAPublicKey(app.Config.TokenSignerPublicKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load token signer public key")
	}

	dbConfig, err := db.InitConfig()
	if err != nil {
		return nil, err
	}

	app.DB, err = db.New(dbConfig, logger)
	if err != nil {
		return nil, err
	}

	redisCacheConfig, err := redis_cache.InitConfig()
	if err != nil {
		return nil, err
	}

	app.RedisCache, err = redis_cache.New(redisCacheConfig, logger)
	if err != nil {
		return nil, err
	}

	app.My = &MyStruct{
		QueueChan:           make(chan *ReservationQueueElem, BatchSize),
		Batch:               make(chan *ReservationQueueElem, BatchSize),
		Signal:              make(chan struct{}),
		ClearBatchTicker:    time.NewTicker(TickTime),
		UpdateDBEventTicker: time.NewTicker(time.Second * 30),
	}

	return app, err
}

func (app *App) SpinWorker() {
	go app.AddTasks()
	for {
		select {
		case <-app.My.ClearBatchTicker.C:
			// Waiting for a signal from ticker
			go app.WorkerPerformBatchTask()
		case <-app.My.Signal:
			// Waiting for a signal from AddTasks()
			go app.WorkerPerformBatchTask()
		case <-app.My.UpdateDBEventTicker.C:
			go app.QueryWorker()
		}
	}
}

func (app *App) QueryWorker() {
	mutex := app.RedisCache.GetLockInstance().NewMutex("refresh-quotas-lock")
	if err := mutex.Lock(); err != nil {
		app.Logger.Debugf("Too bad! Refresh next time")
		return
	}
	_ = app.DB.RefreshEventQuotasFromEntryInReservationsTable()
	_, _ = mutex.Unlock()
}

// To optimize performance we must update DB in batches
func (app *App) AddTasks() {
	for task := range app.My.QueueChan {
		// check if in cache
		found, err := app.RedisCache.GetEventQuota(task.EventID)
		if err != nil {
			task.c <- ReservationQueueResult{
				ticket: nil,
				err:    err,
			}
			// Return early and skip this one
			continue
		}
		if found == -1 {
			// Not in cache we fetch from DB
			thisEvent, err := app.DB.ViewEventDetail(task.EventID)
			if err != nil {
				task.c <- ReservationQueueResult{
					ticket: nil,
					err: &customError.UserError{
						Code:           customError.EventNotFound,
						Message:        "Event not found",
						HTTPStatusCode: http.StatusNotFound,
					},
				}
				continue
			}
			// Put in cache
			err = app.RedisCache.SetNXEventQuota(thisEvent.EventID, thisEvent.RemainingQuota)
			if err != nil {
				task.c <- ReservationQueueResult{
					ticket: nil,
					err:    err,
				}
				// Return early and skip this one
				continue
			}
		}
		err = app.RedisCache.DecEventQuota(task.EventID, task.Amount)
		if err != nil {
			task.c <- ReservationQueueResult{
				ticket: nil,
				err:    err,
			}
			continue
		}
		app.My.Batch <- task
		// When len channel is over a certain amount
		// Send a signal to perform task
		if len(app.My.Batch) >= 50 {
			app.My.Signal <- struct{}{}
		}
	}
}

func (app *App) WorkerPerformBatchTask() {
	var jobs []*model.ReservationRequest
	var returnChan []chan ReservationQueueResult
	deductQuotaMap := make(map[int]int)
	size := len(app.My.Batch)
	for i := 0; i < size; i++ {
		item := <-app.My.Batch
		jobs = append(jobs, &model.ReservationRequest{
			EventID: item.EventID,
			UserID:  item.UserID,
			Amount:  item.Amount,
		})
		deductQuotaMap[item.EventID] += item.Amount
		returnChan = append(returnChan, item.c)
	}
	results, err := app.DB.MakeReservationBatch(jobs, deductQuotaMap)
	if err != nil {
		for _, c := range returnChan {
			c <- ReservationQueueResult{
				ticket: nil,
				err:    err,
			}
		}
		for _, j := range jobs {
			_ = app.RedisCache.IncEventQuota(j.EventID, j.Amount)
		}
		return
	}
	for i := 0; i < len(returnChan); i++ {
		returnChan[i] <- ReservationQueueResult{
			ticket: results[i],
			err:    nil,
		}
	}
}

func (app *App) Close() error {
	if err := app.DB.Close(); err != nil {
		return err
	}
	if err := app.RedisCache.Close(); err != nil {
		return err
	}
	return nil
}

func readRSAPrivateKey(filepath string) (*rsa.PrivateKey, error) {
	privateKeyBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	privateKey := utils.BytesToPrivateKey(privateKeyBytes)
	return privateKey, nil
}

func readRSAPublicKey(filepath string) (*rsa.PublicKey, error) {
	publicKeyBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	publicKey := utils.BytesToPublicKey(publicKeyBytes)
	return publicKey, nil
}

func validateInput(input interface{}) *customError.ValidationError {
	err := validate.Struct(input)
	if err != nil {
		messages := make([]string, 0)
		for _, e := range err.(validator.ValidationErrors) {
			messages = append(messages, e.Translate(trans))
		}
		errMessage := strings.Join(messages, ", ")
		return &customError.ValidationError{
			Code:    customError.InputValidationError,
			Message: errMessage,
		}
	}
	return nil
}

type RequesterMetadata struct {
	Token string `json:"token" validate:"required,min=1"`
}

func checkPostgresErrorCode(err error, code string) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == code {
			return true
		}
	}
	return false
}
