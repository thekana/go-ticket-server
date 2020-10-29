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
	"reflect"
	"strings"
	"ticket-reservation/db/model"
	"time"

	customError "ticket-reservation/custom_error"
	"ticket-reservation/db"
	log "ticket-reservation/log"
	"ticket-reservation/utils"
)

type MyStruct struct {
	QueueChan chan *ReservationQueueElem
	Signal    chan struct{}
	Timer     *time.Ticker
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
		QueueChan: make(chan *ReservationQueueElem, 100),
		Signal:    make(chan struct{}),
		Timer:     time.NewTicker(time.Millisecond * 20),
	}
	return app, err
}

func (app *App) SpinTaskWorker() {
	for {
		select {
		case <-app.My.Signal:
		case <-app.My.Timer.C:
			app.WorkerPerformTask()
		}
	}
}

// Update Batch Database
func (app *App) WorkerPerformTask() {
	for task := range app.My.QueueChan {
		var ticket *model.ReservationDetail
		ticket, err := app.DB.MakeReservation(task.UserID, task.EventID, task.Amount)
		task.c <- ReservationQueueResult{
			ticket: ticket,
			err:    err,
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
