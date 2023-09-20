package es

import (
	"github.com/olivere/elastic/v7"
	"github.com/segmentfault/pacman/log"
	"net/http"
	"net/http/httputil"
)

type LoggingHttpClient struct {
	c http.Client
}

func (l LoggingHttpClient) Do(r *http.Request) (*http.Response, error) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Errorf("dump request failed: %s", err.Error())
		return nil, err
	}
	log.Debugf("es search request: %s", string(requestDump))
	return l.c.Do(r)
}

type ErrLogger struct {
}

func (l ErrLogger) Printf(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

func NewErrLogger() elastic.Logger {
	return &ErrLogger{}
}
