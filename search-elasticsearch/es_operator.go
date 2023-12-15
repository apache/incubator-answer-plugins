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
	"context"
	"github.com/olivere/elastic/v7"
	"github.com/segmentfault/pacman/log"
	"net/http"
)

// Operator elasticsearch client
// only support basic functions, not support bulk operations
type Operator struct {
	C *elastic.Client
}

func NewOperator(urls []string, username, password string) (c *Operator, err error) {
	doer := LoggingHttpClient{
		c: http.Client{},
	}
	esClient, err := elastic.NewSimpleClient(
		elastic.SetHttpClient(doer),
		elastic.SetURL(urls...),
		elastic.SetBasicAuth(username, password),
		elastic.SetSniff(false),
		elastic.SetErrorLog(NewErrLogger()))
	if err != nil {
		return nil, err
	}
	c = &Operator{C: esClient}
	return c, nil
}

func (op *Operator) CreateIndex(ctx context.Context, indexName string, mapping string) (err error) {
	log.Debugf("try to create index: %s", indexName)
	exist, err := op.C.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}
	if exist {
		log.Debugf("index %s already exists", indexName)
		return nil
	}
	_, err = op.C.CreateIndex(indexName).BodyString(mapping).Do(ctx)
	if err != nil {
		log.Errorf("create index %s failed: %s", indexName, err.Error())
		return err
	}
	return nil
}

func (op *Operator) QueryDoc(ctx context.Context, indexName string,
	query elastic.Query, sort *elastic.FieldSort, cols *elastic.FetchSourceContext,
	page, size int) (
	result *elastic.SearchResult, err error) {
	log.Debugf("try to query doc from index: %s, %d, %d", indexName, page, size)
	from := (page - 1) * size
	service := op.C.Search().Index(indexName).Query(query).From(from).Size(size)
	if cols != nil {
		service = service.FetchSourceContext(cols)
	}
	if sort != nil {
		service = service.SortBy(sort)
	}
	result, err = service.Do(ctx)
	if err != nil {
		log.Errorf("query doc from index %s failed: %s", indexName, err.Error())
		return nil, err
	}
	return result, nil
}

func (op *Operator) SaveDoc(ctx context.Context, indexName string, id string, doc interface{}) (err error) {
	log.Debugf("try to save doc to index: %s, %s", indexName, id)
	_, err = op.C.Update().Index(indexName).Id(id).Refresh("false").DocAsUpsert(true).Doc(doc).Upsert(doc).Do(ctx)
	if err != nil {
		log.Errorf("save doc to index %s failed: %s", indexName, err.Error())
		return err
	}
	return nil
}

func (op *Operator) DeleteDoc(ctx context.Context, indexName string, id string) (err error) {
	log.Debugf("try to delete doc from index: %s, %s", indexName, id)
	_, err = op.C.Delete().Index(indexName).Id(id).Refresh("false").Do(ctx)
	if err != nil {
		log.Errorf("delete doc from index %s failed: %s", indexName, err.Error())
		return err
	}
	return nil
}
