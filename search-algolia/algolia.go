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

package algolia

import (
	"context"
	"embed"
	"github.com/apache/incubator-answer-plugins/util"
	"strconv"
	"strings"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/apache/incubator-answer-plugins/search-algolia/i18n"
	"github.com/apache/incubator-answer/plugin"
)

//go:embed  info.yaml
var Info embed.FS

type SearchAlgolia struct {
	Config *AlgoliaSearchConfig
	client *search.Client
	syncer plugin.SearchSyncer
}

func init() {
	uc := &SearchAlgolia{Config: &AlgoliaSearchConfig{}}
	plugin.Register(uc)
}

func (s *SearchAlgolia) Info() plugin.Info {
	info := &util.Info{}
	info.GetInfo(Info)

	return plugin.Info{
		Name:        plugin.MakeTranslator(i18n.InfoName),
		SlugName:    info.SlugName,
		Description: plugin.MakeTranslator(i18n.InfoDescription),
		Version:     info.Version,
		Author:      info.Author,
		Link:        info.Link,
	}
}

func (s *SearchAlgolia) Description() plugin.SearchDesc {
	desc := plugin.SearchDesc{}
	if s.Config.ShowLogo {
		desc.Icon = icon
	}
	return desc
}

func (s *SearchAlgolia) RegisterSyncer(ctx context.Context, syncer plugin.SearchSyncer) {
	s.syncer = syncer
	s.sync()
}

func (s *SearchAlgolia) SearchContents(ctx context.Context, cond *plugin.SearchBasicCond) (res []plugin.SearchResult, total int64, err error) {
	var (
		filters      = "status<10"
		tagFilters   []string
		userIDFilter string
		votesFilter  string
	)
	if len(cond.TagIDs) > 0 {
		for _, tagGroup := range cond.TagIDs {
			var tagsIn []string
			if len(tagGroup) > 0 {
				for _, tagID := range tagGroup {
					tagsIn = append(tagsIn, "tags:"+tagID)
				}
			}
			tagFilters = append(tagFilters, "("+strings.Join(tagsIn, " OR ")+")")
		}
		if len(tagFilters) > 0 {
			filters += " AND " + strings.Join(tagFilters, " AND ")
		}
	}
	if len(cond.UserID) > 0 {
		userIDFilter = "userID:" + cond.UserID
		filters += " AND " + userIDFilter
	}
	if cond.VoteAmount == 0 {
		votesFilter = "votes=" + strconv.Itoa(cond.VoteAmount)
		filters += " AND " + votesFilter
	} else if cond.VoteAmount > 0 {
		votesFilter = "votes>=" + strconv.Itoa(cond.VoteAmount)
		filters += " AND " + votesFilter
	}

	var (
		query = strings.TrimSpace(strings.Join(cond.Words, " "))
		opts  = []interface{}{
			opt.AttributesToRetrieve("objectID", "type"),
			opt.Filters(filters),
			opt.Page(cond.Page - 1),
			opt.HitsPerPage(cond.PageSize),
		}
		qres search.QueryRes
	)

	qres, err = s.getIndex(string(cond.Order)).Search(query, opts...)
	for _, hit := range qres.Hits {
		res = append(res, plugin.SearchResult{
			ID:   hit["objectID"].(string),
			Type: hit["type"].(string),
		})
	}
	total = int64(qres.NbHits)
	return res, total, err
}

func (s *SearchAlgolia) SearchQuestions(ctx context.Context, cond *plugin.SearchBasicCond) (res []plugin.SearchResult, total int64, err error) {
	var (
		filters       = "status<10 AND type:question"
		tagFilters    []string
		userIDFilter  string
		viewsFilter   string
		answersFilter string
	)
	if len(cond.TagIDs) > 0 {
		for _, tagGroup := range cond.TagIDs {
			var tagsIn []string
			if len(tagGroup) > 0 {
				for _, tagID := range tagGroup {
					tagsIn = append(tagsIn, "tags:"+tagID)
				}
			}
			tagFilters = append(tagFilters, "("+strings.Join(tagsIn, " OR ")+")")
		}
		if len(tagFilters) > 0 {
			filters += " AND " + strings.Join(tagFilters, " AND ")
		}
	}
	if cond.QuestionAccepted == plugin.AcceptedCondFalse {
		userIDFilter = "hasAccepted:false"
		filters += " AND " + userIDFilter
	}

	if cond.ViewAmount > -1 {
		viewsFilter = "views>=" + strconv.Itoa(cond.ViewAmount)
		filters += " AND " + viewsFilter
	}

	// check answers
	if cond.AnswerAmount == 0 {
		answersFilter = "answers=0"
		filters += " AND " + answersFilter
	} else if cond.AnswerAmount > 0 {
		answersFilter = "answers>=" + strconv.Itoa(cond.AnswerAmount)
		filters += " AND " + answersFilter
	}

	var (
		query = strings.TrimSpace(strings.Join(cond.Words, " "))
		opts  = []interface{}{
			opt.AttributesToRetrieve("objectID", "type"),
			opt.Filters(filters),
			opt.Page(cond.Page - 1),
			opt.HitsPerPage(cond.PageSize),
		}
		qres search.QueryRes
	)

	qres, err = s.getIndex(string(cond.Order)).Search(query, opts...)
	for _, hit := range qres.Hits {
		res = append(res, plugin.SearchResult{
			ID:   hit["objectID"].(string),
			Type: hit["type"].(string),
		})
	}

	total = int64(qres.NbHits)
	return res, total, err
}

func (s *SearchAlgolia) SearchAnswers(ctx context.Context, cond *plugin.SearchBasicCond) (res []plugin.SearchResult, total int64, err error) {
	var (
		filters          = "status<10 AND type:answer"
		tagFilters       []string
		userIDFilter     string
		questionIDFilter string
	)
	if len(cond.TagIDs) > 0 {
		for _, tagGroup := range cond.TagIDs {
			var tagsIn []string
			if len(tagGroup) > 0 {
				for _, tagID := range tagGroup {
					tagsIn = append(tagsIn, "tags:"+tagID)
				}
			}
			tagFilters = append(tagFilters, "("+strings.Join(tagsIn, " OR ")+")")
		}
		if len(tagFilters) > 0 {
			filters += " AND " + strings.Join(tagFilters, " AND ")
		}
	}
	if cond.AnswerAccepted == plugin.AcceptedCondTrue {
		userIDFilter = "hasAccepted=true"
		filters += " AND " + userIDFilter
	}

	if len(cond.QuestionID) > 0 {
		questionIDFilter = "questionID=" + cond.QuestionID
		filters += questionIDFilter
	}

	var (
		query = strings.TrimSpace(strings.Join(cond.Words, " "))
		opts  = []interface{}{
			opt.AttributesToRetrieve("objectID", "type"),
			opt.Filters(filters),
			opt.Page(cond.Page - 1),
			opt.HitsPerPage(cond.PageSize),
		}
		qres search.QueryRes
	)

	qres, err = s.getIndex(string(cond.Order)).Search(query, opts...)
	for _, hit := range qres.Hits {
		res = append(res, plugin.SearchResult{
			ID:   hit["objectID"].(string),
			Type: hit["type"].(string),
		})
	}
	total = int64(qres.NbHits)
	return res, total, err
}

// UpdateContent updates the content to algolia server
func (s *SearchAlgolia) UpdateContent(ctx context.Context, content *plugin.SearchContent) (err error) {
	_, err = s.getIndex("").SaveObject(content)
	return
}

// DeleteContent deletes the content
func (s *SearchAlgolia) DeleteContent(ctx context.Context, contentID string) (err error) {
	_, err = s.getIndex("").DeleteObject(contentID)
	return
}

// connect connect to algolia server
func (s *SearchAlgolia) connect() (err error) {
	s.client = search.NewClient(s.Config.APPID, s.Config.APIKey)
	return
}

// init or create index
func (s *SearchAlgolia) getIndex(order string) (index *search.Index) {
	idx := s.getIndexName(order)
	return s.client.InitIndex(idx)
}

func (s *SearchAlgolia) getIndexName(order string) string {
	// main index
	var idx = s.Config.Index
	switch order {
	case NewestIndex:
		// the index of sort results by newest
		idx = idx + "_" + NewestIndex
	case ActiveIndex:
		// the index of sort results by active
		idx = idx + "_" + ActiveIndex
	case ScoreIndex:
		// the index of sort results by score
		idx = idx + "_" + ScoreIndex
	}
	return idx
}
