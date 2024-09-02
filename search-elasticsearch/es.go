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
	"embed"
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-answer-plugins/util"
	"strings"
	"sync"

	"github.com/apache/incubator-answer-plugins/search-elasticsearch/i18n"
	"github.com/apache/incubator-answer/plugin"
	"github.com/olivere/elastic/v7"
	"github.com/segmentfault/pacman/log"
)

//go:embed  info.yaml
var Info embed.FS

type SearchEngine struct {
	Config   *SearchEngineConfig
	Operator *Operator
	syncer   plugin.SearchSyncer
	syncing  bool
	lock     sync.Mutex
}

type SearchEngineConfig struct {
	Endpoints string `json:"endpoints"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func init() {
	plugin.Register(&SearchEngine{
		Config: &SearchEngineConfig{},
		lock:   sync.Mutex{},
	})
}

func (s *SearchEngine) Info() plugin.Info {
	info := &util.Info{}
	info.GetInfo(Info)

	return plugin.Info{
		Name:        plugin.MakeTranslator(i18n.InfoName),
		SlugName:    info.SlugName,
		Description: plugin.MakeTranslator(i18n.InfoDescription),
		Author:      info.Author,
		Version:     info.Version,
		Link:        info.Link,
	}
}

func (s *SearchEngine) Description() plugin.SearchDesc {
	return plugin.SearchDesc{}
}

func (s *SearchEngine) SearchContents(
	ctx context.Context, cond *plugin.SearchBasicCond) (
	res []plugin.SearchResult, total int64, err error) {
	if s.Operator == nil {
		return nil, 0, fmt.Errorf("es client not init")
	}
	resp, err := s.Operator.QueryDoc(ctx, s.getIndexName(),
		s.buildQuery(cond), s.buildSort(cond), s.buildCols(), cond.Page, cond.PageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("es query error: %w", err)
	}
	if resp == nil {
		return nil, 0, nil
	}
	return s.warpResult(resp)
}

func (s *SearchEngine) SearchQuestions(
	ctx context.Context, cond *plugin.SearchBasicCond) (
	res []plugin.SearchResult, total int64, err error) {
	if s.Operator == nil {
		return nil, 0, fmt.Errorf("es client not init")
	}
	query := s.buildQuery(cond)
	query.Must(elastic.NewTermQuery("type", "question"))
	resp, err := s.Operator.QueryDoc(ctx, s.getIndexName(),
		query, s.buildSort(cond), s.buildCols(), cond.Page, cond.PageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("es query error: %w", err)
	}
	if resp == nil {
		return nil, 0, nil
	}
	return s.warpResult(resp)
}

func (s *SearchEngine) SearchAnswers(
	ctx context.Context, cond *plugin.SearchBasicCond) (
	res []plugin.SearchResult, total int64, err error) {
	if s.Operator == nil {
		return nil, 0, fmt.Errorf("es client not init")
	}
	query := s.buildQuery(cond)
	query.Must(elastic.NewTermQuery("type", "answer"))
	resp, err := s.Operator.QueryDoc(ctx, s.getIndexName(),
		query, s.buildSort(cond), s.buildCols(), cond.Page, cond.PageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("es query error: %w", err)
	}
	if resp == nil {
		return nil, 0, nil
	}
	return s.warpResult(resp)
}

func (s *SearchEngine) UpdateContent(ctx context.Context, content *plugin.SearchContent) error {
	if s.Operator == nil {
		return fmt.Errorf("es client not init")
	}
	return s.Operator.SaveDoc(ctx, s.getIndexName(), content.ObjectID, CreateDocFromSearchContent(content.ObjectID, content))
}

func (s *SearchEngine) DeleteContent(ctx context.Context, contentID string) error {
	if s.Operator == nil {
		return fmt.Errorf("es client not init")
	}
	return s.Operator.DeleteDoc(ctx, s.getIndexName(), contentID)
}

func (s *SearchEngine) RegisterSyncer(ctx context.Context, syncer plugin.SearchSyncer) {
	s.syncer = syncer
	s.sync()
}

func (s *SearchEngine) warpResult(resp *elastic.SearchResult) ([]plugin.SearchResult, int64, error) {
	res := make([]plugin.SearchResult, 0)
	for _, hit := range resp.Hits.Hits {
		docByte, err := hit.Source.MarshalJSON()
		if err != nil {
			log.Errorf("es unmarshal error: %v", err)
			continue
		}

		var content AnswerPostDoc
		err = json.Unmarshal(docByte, &content)
		if err != nil {
			log.Errorf("es unmarshal error: %v", err)
			continue
		}

		res = append(res, plugin.SearchResult{
			ID:   hit.Id,
			Type: content.Type,
		})
	}
	log.Debugf("search result: %d", len(res))
	return res, resp.TotalHits(), nil
}

func (s *SearchEngine) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Name:        "endpoints",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigEndpointsTitle),
			Description: plugin.MakeTranslator(i18n.ConfigEndpointsDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: s.Config.Endpoints,
		},
		{
			Name:        "username",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigUsernameTitle),
			Description: plugin.MakeTranslator(i18n.ConfigUsernameDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: s.Config.Username,
		},
		{
			Name:        "password",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigPasswordTitle),
			Description: plugin.MakeTranslator(i18n.ConfigPasswordDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: s.Config.Password,
		},
	}
}

func (s *SearchEngine) ConfigReceiver(config []byte) error {
	conf := &SearchEngineConfig{}
	_ = json.Unmarshal(config, conf)
	s.Config = conf

	log.Debugf("try to init es client: %s", conf.Endpoints)

	operator, err := NewOperator(strings.Split(conf.Endpoints, ","), conf.Username, conf.Password)
	if err != nil {
		return fmt.Errorf("init es client error: %w", err)
	}
	s.Operator = operator
	err = s.Operator.CreateIndex(context.Background(), s.getIndexName(), indexJson)
	if err != nil {
		return fmt.Errorf("create index error: %w", err)
	}
	return nil
}

func (s *SearchEngine) getIndexName() string {
	return "answer_post"
}

func (s *SearchEngine) buildSort(cond *plugin.SearchBasicCond) (sort *elastic.FieldSort) {
	switch cond.Order {
	case plugin.SearchNewestOrder:
		return elastic.NewFieldSort("created").Desc()
	case plugin.SearchActiveOrder:
		return elastic.NewFieldSort("active").Desc()
	case plugin.SearchScoreOrder:
		return elastic.NewFieldSort("score").Desc()
	default:
		return nil
	}
}

func (s *SearchEngine) buildCols() (cols *elastic.FetchSourceContext) {
	return elastic.NewFetchSourceContext(true).Include("id", "type")
}

func (s *SearchEngine) buildQuery(cond *plugin.SearchBasicCond) (
	query *elastic.BoolQuery) {

	log.Debugf("build query: %+v", cond)

	q := elastic.NewBoolQuery()
	for _, tagGroup := range cond.TagIDs {
		if len(tagGroup) > 0 {
			q.Must(elastic.NewTermsQuery("tags", convertToInterfaceSlice(tagGroup)...))
		}
	}
	if len(cond.UserID) > 0 {
		q.Must(elastic.NewTermQuery("user_id", cond.UserID))
	}
	if len(cond.QuestionID) > 0 {
		q.Must(elastic.NewTermQuery("question_id", cond.QuestionID))
	}
	if cond.VoteAmount > 0 {
		q.Must(elastic.NewRangeQuery("score").Gte(cond.VoteAmount))
	}
	if cond.ViewAmount > 0 {
		q.Must(elastic.NewRangeQuery("views").Gte(cond.ViewAmount))
	}
	if cond.AnswerAmount > 0 {
		q.Must(elastic.NewRangeQuery("answers").Gte(cond.AnswerAmount))
	}
	if cond.AnswerAccepted == plugin.AcceptedCondTrue {
		q.Must(elastic.NewTermQuery("has_accepted", true))
	} else if cond.AnswerAccepted == plugin.AcceptedCondFalse {
		q.MustNot(elastic.NewTermQuery("has_accepted", true))
	}
	if cond.QuestionAccepted == plugin.AcceptedCondTrue {
		q.MustNot(elastic.NewTermQuery("has_accepted", true))
	} else if cond.QuestionAccepted == plugin.AcceptedCondFalse {
		q.Must(elastic.NewTermQuery("has_accepted", false))
	}
	if len(cond.Words) > 0 {
		q.Must(elastic.NewMultiMatchQuery(strings.Join(cond.Words, " "), "title", "content"))
	}
	q.Must(elastic.NewTermQuery("status", plugin.SearchContentStatusAvailable))
	return q
}

func convertToInterfaceSlice(slice []string) []interface{} {
	s := make([]interface{}, len(slice))
	for i, v := range slice {
		s[i] = v
	}
	return s
}
