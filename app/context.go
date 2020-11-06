package app

import (
	"crypto/rsa"
	"ticket-reservation/db"
	log "ticket-reservation/log"
	"ticket-reservation/redis_cache"
)

type Context struct {
	Logger                log.Logger
	Config                *Config
	RemoteAddress         string
	TokenSignerPrivateKey *rsa.PrivateKey
	TokenSignerPublicKey  *rsa.PublicKey
	DB                    db.DB
	My                    *MyStruct
	RedisCache            redis_cache.Cache
}

func (app *App) NewContext() *Context {
	return &Context{
		Logger:                app.Logger,
		Config:                app.Config,
		TokenSignerPrivateKey: app.TokenSignerPrivateKey,
		TokenSignerPublicKey:  app.TokenSignerPublicKey,
		DB:                    app.DB,
		My:                    app.My,
		RedisCache:            app.RedisCache,
	}
}

func (ctx *Context) WithLogger(logger log.Logger) *Context {
	ret := *ctx
	ret.Logger = logger
	return &ret
}

func (ctx *Context) WithRemoteAddress(address string) *Context {
	ret := *ctx
	ret.RemoteAddress = address
	return &ret
}

func (ctx *Context) getLogger() log.Logger {
	return ctx.Logger.WithFields(log.Fields{
		"module": "app",
	})
}
