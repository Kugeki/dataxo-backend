package restapi

import (
	"log"
	"net/http"
	"time"
)

type Opt func(s *http.Server) error

func WithAddr(addr string) Opt {
	return func(s *http.Server) error {
		s.Addr = addr
		return nil
	}
}

func WithReadTimeout(timeout time.Duration) Opt {
	return func(s *http.Server) error {
		s.ReadTimeout = timeout
		return nil
	}
}

func WithWriteTimeout(timeout time.Duration) Opt {
	return func(s *http.Server) error {
		s.WriteTimeout = timeout
		return nil
	}
}

func WithErrorLog(errorLog *log.Logger) Opt {
	return func(s *http.Server) error {
		s.ErrorLog = errorLog
		return nil
	}
}
