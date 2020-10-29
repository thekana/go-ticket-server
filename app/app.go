package app

import (
	"crypto/rsa"
	"fmt"
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
	"sync"
	"ticket-reservation/db/model"
	"time"

	customError "ticket-reservation/custom_error"
	"ticket-reservation/db"
	log "ticket-reservation/log"
	"ticket-reservation/utils"
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
}

type App struct {
	Logger                log.Logger
	Config                *Config
	TokenSignerPrivateKey *rsa.PrivateKey
	TokenSignerPublicKey  *rsa.PublicKey
	DB                    db.DB
	My                    *MyStruct
}

var (
	uni       *ut.UniversalTranslator
	trans     ut.Translator
	validate  *validator.Validate
	BATCHSIZE int           = 50
	TICKTIME  time.Duration = time.Second * 10
	m         sync.Mutex
	ascii     string = `
 _______  __   __  ___      ___     
|       ||  | |  ||   |    |   |    
|    ___||  | |  ||   |    |   |    
|   |___ |  |_|  ||   |    |   |    
|    ___||       ||   |___ |   |___ 
|   |    |       ||       ||       |
|___|    |_______||_______||_______|
`
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

type AppNewOptions struct {
}

func New(logger log.Logger, options *AppNewOptions) (app *App, err error) {
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
		QueueChan: make(chan *ReservationQueueElem, BATCHSIZE),
		Signal:    make(chan struct{}),
		Timer:     time.NewTicker(TICKTIME),
		EventQuotaMap: &RWMap{
			m: make(map[int]int),
		},
	}
	// Query events and put in EventQuotaMap
	events, _ := app.DB.ViewAllEvents(false, 0)
	if events != nil {
		for _, e := range events {
			app.My.EventQuotaMap.Set(e.EventID, e.RemainingQuota)
		}
	}
	fmt.Println(app.My.EventQuotaMap)
	return app, err
}

func (app *App) SpinWorker() {
	batch := make(chan *ReservationQueueElem, BATCHSIZE)
	go app.WorkerPerformTask(batch)
	for {
		select {
		case <-app.My.Timer.C:
			fmt.Println(app.My.EventQuotaMap)
			go app.WorkerPerformBatchTask(batch)
		}
	}
}

// To optimize performance we must update DB in batches
func (app *App) WorkerPerformTask(batch chan *ReservationQueueElem) {
	// Query the current quotas of all events into memory
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
		}
		app.My.EventQuotaMap.Set(task.EventID, newQuota)
		batch <- task
		if len(batch) >= BATCHSIZE {
			fmt.Println(ascii)
			app.WorkerPerformBatchTask(batch)
		}
		//select {
		//case batch <- task:
		//	fmt.Println("Keep filling")
		//default:
		//	fmt.Println(ascii)
		//	go app.WorkerPerformBatchTask(batch)
		//}
	}
}

func (app *App) WorkerPerformBatchTask(batch chan *ReservationQueueElem) {
	//m.Lock()
	//defer m.Unlock()
	var jobs []*model.ReservationRequest
	var returnChan []chan ReservationQueueResult
	deductQuotaMap := make(map[int]int)

	for i := 0; i < BATCHSIZE; i++ {
		select {
		// Perform 100 jobs at most
		case item := <-batch:
			jobs = append(jobs, &model.ReservationRequest{
				EventID: item.EventID,
				UserID:  item.UserID,
				Amount:  item.Amount,
			})
			deductQuotaMap[item.EventID] += item.Amount
			returnChan = append(returnChan, item.c)
		case <-time.After(time.Millisecond * 2):
			// Just in case batch not full
			break
		}
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

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func checkPostgresErrorCode(err error, code string) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == code {
			return true
		}
	}
	return false
}
