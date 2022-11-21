package common

import (
	"log"
	"net/http"
	"time"

	"syscall"

	"github.com/SimonRichardson/colossus/teleprinter/logs"
	"github.com/fvbock/endless"
)

type Callback func()

type ServerTimeout struct {
	Read, Write time.Duration
}

func ListenAndServe(addr string,
	timeout ServerTimeout,
	logger logs.Logger,
	handler http.Handler,
	callback Callback,
) error {
	server := endless.NewServer(addr, handler)
	server.ReadTimeout = timeout.Read
	server.WriteTimeout = timeout.Write
	server.ErrorLog = log.New(logger, "", 0)

	server.RegisterSignalHook(endless.PRE_SIGNAL, syscall.SIGINT, callback)
	server.RegisterSignalHook(endless.PRE_SIGNAL, syscall.SIGQUIT, callback)
	server.RegisterSignalHook(endless.PRE_SIGNAL, syscall.SIGTERM, callback)

	return server.ListenAndServe()
}
