/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

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
