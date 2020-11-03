package app

import (
	"crypto/rsa"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"sync"
	customError "ticket-reservation/custom_error"
	"ticket-reservation/db"
	"ticket-reservation/db/model"
	log "ticket-reservation/log"
	"ticket-reservation/utils"
	"time"
)

type RWMap struct {
	sync.RWMutex
	m map[int]int
}

func (r *RWMap) Get(key int) (int, bool) {
	r.RLock()
	defer r.RUnlock()
	item, found := r.m[key]
	return item, found
}

func (r *RWMap) Set(key int, item int) {
	r.Lock()
	defer r.Unlock()
	r.m[key] = item
}

type MyStruct struct {
	QueueChan     chan *ReservationQueueElem
	Signal        chan struct{}
	Timer         *time.Ticker
	EventQuotaMap *RWMap
	Batch         chan *ReservationQueueElem
}

type App struct {
	Logger                log.Logger
	Config                *Config
	TokenSignerPrivateKey *rsa.PrivateKey
	TokenSignerPublicKey  *rsa.PublicKey
	DB                    db.DB
	My                    *MyStruct
}

const (
	BatchSize int           = 50
	TickTime  time.Duration = time.Millisecond * 100
	Full      string        = `
 _______  __   __  ___      ___     
|       ||  | |  ||   |    |   |    
|    ___||  | |  ||   |    |   |    
|   |___ |  |_|  ||   |    |   |    
|    ___||       ||   |___ |   |___ 
|   |    |       ||       ||       |
|___|    |_______||_______||_______|
`
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

	app.My = &MyStruct{
		QueueChan: make(chan *ReservationQueueElem, BatchSize),
		Batch:     make(chan *ReservationQueueElem, BatchSize),
		Signal:    make(chan struct{}),
		Timer:     time.NewTicker(TickTime),
		EventQuotaMap: &RWMap{
			m: make(map[int]int),
		},
	}
	// Query events and put in EventQuotaMap
	//fmt.Println("Loading event quotas to memory")
	events, err := app.DB.ViewAllEvents(false, 0)
	if checkPostgresErrorCode(err, pgerrcode.UndefinedTable) {
		//fmt.Println("Retrying")
		events, err = app.DB.ViewAllEvents(false, 0)
	}
	if events != nil {
		for _, e := range events {
			app.My.EventQuotaMap.Set(e.EventID, e.RemainingQuota)
		}
	}
	//fmt.Println("Loaded event quotas to memory")
	return app, err
}

func (app *App) SpinWorker() {
	go app.AddTasks()
	for {
		select {
		case <-app.My.Timer.C:
			// fmt.Println(app.My.EventQuotaMap)
			// Waiting for a signal from ticker
			app.WorkerPerformBatchTask()
		case <-app.My.Signal:
			// fmt.Print(Full)
			// Waiting for a signal from AddTasks()
			app.WorkerPerformBatchTask()
		}
	}
}

// To optimize performance we must update DB in batches
func (app *App) AddTasks() {
	for task := range app.My.QueueChan {
		quota, found := app.My.EventQuotaMap.Get(task.EventID)
		if !found {
			task.c <- ReservationQueueResult{
				ticket: nil,
				err: &customError.UserError{
					Code:           70,
					Message:        "Event not found",
					HTTPStatusCode: http.StatusNotFound,
				},
			}
			// Return early and skip this one
			continue
		}
		newQuota := quota - task.Amount
		if newQuota < 0 {
			task.c <- ReservationQueueResult{
				ticket: nil,
				err: &customError.UserError{
					Code:           10,
					Message:        "Not Enough Quota",
					HTTPStatusCode: http.StatusBadRequest,
				},
			}
			// Return early and skip this one
			continue
		}
		app.My.EventQuotaMap.Set(task.EventID, newQuota)
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
