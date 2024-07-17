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

package meilisearch

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-answer-plugins/util"
	"strings"
	"sync"

	"github.com/apache/incubator-answer-plugins/search-meilisearch/i18n"
	"github.com/apache/incubator-answer/plugin"
	"github.com/meilisearch/meilisearch-go"
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
)

//go:embed  info.yaml
var Info embed.FS

const (
	primaryKey       = "objectID"
	defaultIndexName = "answer_post"
)

var (
	configuredErr = fmt.Errorf("meilisearch is not configured correctly")
)

type Search struct {
	Config  *SearchConfig
	Client  *meilisearch.Client
	syncer  plugin.SearchSyncer
	syncing bool
	lock    sync.Mutex
}

type SearchConfig struct {
	Host      string `json:"host"`
	ApiKey    string `json:"api_key"`
	IndexName string `json:"index_name"`
	Async     bool   `json:"async"`
}

func init() {
	plugin.Register(&Search{
		Config: &SearchConfig{},
		lock:   sync.Mutex{},
	})
}

func (s *Search) Info() plugin.Info {
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

func (s *Search) Description() plugin.SearchDesc {
	return plugin.SearchDesc{Icon: "PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjMwIiB2aWV3Qm94PSIwIDAgNDk1IDc0IiBmaWxsPSJub25lIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8cGF0aCBkPSJNMTgxLjg0IDQyLjUzNDdDMTgxLjg0IDM3LjYxMzYgMTg0LjE5OSAzNC43MTQ5IDE4OC43MTYgMzQuNzE0OUMxOTIuOTYzIDM0LjcxNDkgMTk0LjM3OCAzNy43NDg0IDE5NC4zNzggNDEuNjU4NFY2Mi42MjM3SDIwMy45NTFWNDAuNTc5OEMyMDMuOTUxIDMyLjM1NTQgMTk5LjYzNyAyNi40OTA2IDE5MS4xNDMgMjYuNDkwNkMxODYuMDg3IDI2LjQ5MDYgMTgyLjUxNCAyOC4wNDEgMTc5LjQxMyAzMS40NzkxQzE3Ny4zOSAyOC4zNzgxIDE3My45NTIgMjYuNDkwNiAxNjkuMTY2IDI2LjQ5MDZDMTY0LjExIDI2LjQ5MDYgMTYwLjYwNSAyOC41ODA0IDE1OC45ODcgMzEuNjEzOVYyNy4yOTk1SDE1MC4xNTZWNjIuNjIzN0gxNTkuNzI4VjQyLjMzMjVDMTU5LjcyOCAzNy42MTM2IDE2Mi4xNTUgMzQuNzE0OSAxNjYuNjA0IDM0LjcxNDlDMTcwLjg1MSAzNC43MTQ5IDE3Mi4yNjcgMzcuNzQ4NCAxNzIuMjY3IDQxLjY1ODRWNjIuNjIzN0gxODEuODRWNDIuNTM0N1oiIGZpbGw9IiMyMTAwNEIiLz4KPHBhdGggZD0iTTI0My4yNDIgNDcuNzI1NUMyNDMuMjQyIDQ3LjcyNTUgMjQzLjM3NyA0Ni40NDQ3IDI0My4zNzcgNDQuODk0MkMyNDMuMzc3IDM0LjQ0NTIgMjM2LjI5OSAyNi40OTA2IDIyNS44NSAyNi40OTA2QzIxNS40MDEgMjYuNDkwNiAyMDguMTIgMzQuNDQ1MiAyMDguMTIgNDQuODk0MkMyMDguMTIgNTUuNzQ3NiAyMTUuNDY4IDYzLjQzMjYgMjI1LjkxNyA2My40MzI2QzIzNC4wNzQgNjMuNDMyNiAyNDAuNTQ2IDU4LjUxMTUgMjQyLjYzNiA1MS4zNjU4SDIzMi45OTZDMjMxLjg1IDUzLjkyNzQgMjI5LjA4NiA1NS4yMDgzIDIyNi4xODcgNTUuMjA4M0MyMjEuNDAxIDU1LjIwODMgMjE4LjMgNTIuNTc5MiAyMTcuNjI2IDQ3LjcyNTVIMjQzLjI0MlpNMjI1Ljc4MyAzNC4xNzU2QzIzMC4yMzIgMzQuMTc1NiAyMzMuMTMxIDM2Ljg3MjEgMjMzLjgwNSA0MC44NDk0SDIxNy43NkMyMTguNTY5IDM2LjgwNDcgMjIxLjQwMSAzNC4xNzU2IDIyNS43ODMgMzQuMTc1NloiIGZpbGw9IiMyMTAwNEIiLz4KPHBhdGggZD0iTTI0NC43ODkgMzUuNTIzOEgyNDkuMDM2VjYyLjYyMzdIMjU4LjYwOFYyNy4yOTk1SDI0NC43ODlWMzUuNTIzOFpNMjUzLjgyMiAyMi43MTU1QzI1Ny4xOTMgMjIuNzE1NSAyNTkuNjE5IDIwLjM1NiAyNTkuNjE5IDE2Ljk4NTRDMjU5LjYxOSAxMy42MTQ4IDI1Ny4xOTMgMTEuMTg3OSAyNTMuODIyIDExLjE4NzlDMjUwLjQ1MSAxMS4xODc5IDI0OC4wMjQgMTMuNjE0OCAyNDguMDI0IDE2Ljk4NTRDMjQ4LjAyNCAyMC4zNTYgMjUwLjQ1MSAyMi43MTU1IDI1My44MjIgMjIuNzE1NVoiIGZpbGw9IiMyMTAwNEIiLz4KPHBhdGggZD0iTTI3OC40MyA1NC4zOTkzQzI3OC4xNiA1NC4zOTkzIDI3Ny43NTYgNTQuNDY2NyAyNzcuMTQ5IDU0LjQ2NjdDMjc0Ljk5MiA1NC40NjY3IDI3NC43MjIgNTMuNDU1NiAyNzQuNzIyIDUxLjk3MjVWMTIuMDY0M0gyNjUuMTVWNTIuNjQ2NkMyNjUuMTUgNTkuNjU3NSAyNjcuODQ2IDYyLjc1ODUgMjc1LjQ2NCA2Mi43NTg1QzI3Ni43NDUgNjIuNzU4NSAyNzcuOTU4IDYyLjYyMzcgMjc4LjQzIDYyLjU1NjJWNTQuMzk5M1oiIGZpbGw9IiMyMTAwNEIiLz4KPHBhdGggZD0iTTI3OS41MTkgMzUuNTIzOEgyODMuNzY2VjYyLjYyMzdIMjkzLjMzOVYyNy4yOTk1SDI3OS41MTlWMzUuNTIzOFpNMjg4LjU1MyAyMi43MTU1QzI5MS45MjMgMjIuNzE1NSAyOTQuMzUgMjAuMzU2IDI5NC4zNSAxNi45ODU0QzI5NC4zNSAxMy42MTQ4IDI5MS45MjMgMTEuMTg3OSAyODguNTUzIDExLjE4NzlDMjg1LjE4MiAxMS4xODc5IDI4Mi43NTUgMTMuNjE0OCAyODIuNzU1IDE2Ljk4NTRDMjgyLjc1NSAyMC4zNTYgMjg1LjE4MiAyMi43MTU1IDI4OC41NTMgMjIuNzE1NVoiIGZpbGw9IiMyMTAwNEIiLz4KPHBhdGggZD0iTTMxMi41NTcgNjIuOTkzOUMzMjEuODYgNjIuOTkzOSAzMjYuMjQyIDU4LjA3MjggMzI2LjI0MiA1Mi44ODJDMzI2LjI0MiAzOC40NTU3IDMwNS4wMDcgNDYuNDc3OCAzMDUuMDA3IDM2Ljk3MjZDMzA1LjAwNyAzMy44NzE3IDMwNy42MzYgMzEuMjQyNiAzMTIuOTYyIDMxLjI0MjZDMzE4LjQyMiAzMS4yNDI2IDMyMC45ODQgMzQuMjA4NyAzMjEuMzg4IDM3LjkxNjRIMzI2LjE3NUMzMjUuNzcgMzMuMjY1IDMyMi42MDIgMjcuMDYzIDMxMy4wOTcgMjcuMDYzQzMwNC45NCAyNy4wNjMgMzAwLjM1NiAzMS45MTY3IDMwMC4zNTYgMzcuMTc0OUMzMDAuMzU2IDUxLjI2NDEgMzIxLjU5MSA0My4xNzQ2IDMyMS41OTEgNTMuMDE2OEMzMjEuNTkxIDU2LjQ1NDggMzE4LjM1NSA1OC44MTQzIDMxMi41NTcgNTguODE0M0MzMDYuNjI1IDU4LjgxNDMgMzAzLjY1OSA1NS44NDgxIDMwMy4zMjIgNTEuNDY2M0gyOTguNDY4QzI5OC44NzIgNTcuNDY2IDMwMi42NDggNjIuOTkzOSAzMTIuNTU3IDYyLjk5MzlaIiBmaWxsPSIjMjEwMDRCIi8+CjxwYXRoIGQ9Ik0zNjQuMjU2IDQ2LjQxMDRDMzY0LjI1NiA0Ni40MTA0IDM2NC4zMjQgNDUuMzMxOCAzNjQuMzI0IDQ0LjU5MDNDMzY0LjMyNCAzNC44ODI5IDM1OC4wNTQgMjcuMDYzIDM0Ny44MDggMjcuMDYzQzMzNy40OTQgMjcuMDYzIDMzMC45NTUgMzUuNDg5NiAzMzAuOTU1IDQ0Ljk5NDdDMzMwLjk1NSA1NC42MzQ3IDMzNy4wMjIgNjIuOTkzOSAzNDcuODc1IDYyLjk5MzlDMzU2LjAzMiA2Mi45OTM5IDM2MS42OTUgNTguMDA1MyAzNjMuNzE3IDUxLjQ2NjNIMzU4LjcyOEMzNTcuMjQ1IDU1LjY0NTkgMzUzLjIwMSA1OC42Nzk1IDM0Ny45NDIgNTguNjc5NUMzNDAuNzI5IDU4LjY3OTUgMzM2LjIxMyA1My4zNTM5IDMzNS43NDEgNDYuNDEwNEgzNjQuMjU2Wk0zNDcuODA4IDMxLjM3NzRDMzU0LjU0OSAzMS4zNzc0IDM1OC45MzEgMzUuODk0IDM1OS41MzcgNDIuNTAwNUgzMzUuODc2QzMzNi42ODUgMzYuMTYzNyAzNDEuMTM0IDMxLjM3NzQgMzQ3LjgwOCAzMS4zNzc0WiIgZmlsbD0iIzIxMDA0QiIvPgo8cGF0aCBkPSJNMzk0LjAzNyA0NS44NzExVjQ5LjEwNjlDMzk0LjAzNyA1NC45NzE4IDM4OS43OSA1OS4wMTY1IDM4MS42MzMgNTkuMDE2NUMzNzYuNTc4IDU5LjAxNjUgMzczLjgxNCA1Ni45MjY3IDM3My44MTQgNTIuNDEwMUMzNzMuODE0IDUwLjExODEgMzc0Ljg5MiA0OC4zNjU0IDM3Ni41NzggNDcuNDIxNkMzNzguMzMgNDYuNDc3OCAzODAuNjkgNDUuODcxMSAzOTQuMDM3IDQ1Ljg3MTFaTTM4MS4wOTQgNjIuOTkzOUMzODcuMDI2IDYyLjk5MzkgMzkxLjgxMyA2MS4xMDYzIDM5NC4yNCA1Ny4xOTY0VjYyLjE4NDlIMzk4LjgyNFYzOS43MzY2QzM5OC44MjQgMzIuMTE4OSAzOTQuNDQyIDI3LjA2MyAzODQuNTMyIDI3LjA2M0MzNzUuMDI3IDI3LjA2MyAzNzAuODQ3IDMxLjg0OTMgMzY5Ljk3MSAzNy45ODM4SDM3NC42MjNDMzc1LjU2NiAzMy4xMzAxIDM3OS4yNzQgMzEuMTc1MiAzODQuMzMgMzEuMTc1MkMzOTAuODAyIDMxLjE3NTIgMzk0LjAzNyAzMy44NzE3IDM5NC4wMzcgMzkuNjY5MVY0MS44OTM4QzM4My4xODQgNDEuODkzOCAzNzguNjY3IDQyLjA5NiAzNzUuMjk3IDQzLjQ0NDJDMzcxLjM4NyA0NC45OTQ3IDM2OS4wOTUgNDguNDMyOCAzNjkuMDk1IDUyLjU0NDlDMzY5LjA5NSA1OC41NDQ2IDM3Mi45MzcgNjIuOTkzOSAzODEuMDk0IDYyLjk5MzlaIiBmaWxsPSIjMjEwMDRCIi8+CjxwYXRoIGQ9Ik00MjQuOTkxIDI3LjYwMjNDNDI0Ljk5MSAyNy42MDIzIDQyNC4xODIgMjcuNTM0OSA0MjMuODQ1IDI3LjUzNDlDNDE3LjUwOCAyNy41MzQ5IDQxNC4xMzggMzAuODM4MSA0MTIuODU3IDMzLjE5NzVWMjcuODcySDQwOC4yNzNWNjIuMTg0OUg0MTMuMDU5VjQyLjcwMjdDNDEzLjA1OSAzNS41NTcgNDE3LjQ0MSAzMi4wNTE1IDQyMy4zMDYgMzIuMDUxNUM0MjQuMTgyIDMyLjA1MTUgNDI0Ljk5MSAzMi4xMTg5IDQyNC45OTEgMzIuMTE4OVYyNy42MDIzWiIgZmlsbD0iIzIxMDA0QiIvPgo8cGF0aCBkPSJNNDI1LjgwOSA0NS4wNjIxQzQyNS44MDkgNTQuNDMyNSA0MzIuMjggNjIuOTkzOSA0NDIuNzI5IDYyLjk5MzlDNDUyLjAzMiA2Mi45OTM5IDQ1Ny40MjUgNTYuNzkxOSA0NTguNzczIDQ5Ljk4MzJINDUzLjkyQzQ1Mi41MDQgNTUuMzA4OCA0NDguNTk0IDU4LjY3OTUgNDQyLjcyOSA1OC42Nzk1QzQzNS41MTYgNTguNjc5NSA0MzAuNjYyIDUyLjk0OTQgNDMwLjY2MiA0NS4wNjIxQzQzMC42NjIgMzcuMTA3NSA0MzUuNTE2IDMxLjM3NzQgNDQyLjcyOSAzMS4zNzc0QzQ0OC41OTQgMzEuMzc3NCA0NTIuNTA0IDM0Ljc0OCA0NTMuOTIgNDAuMDczNkg0NTguNzczQzQ1Ny40MjUgMzMuMjY1IDQ1Mi4wMzIgMjcuMDYzIDQ0Mi43MjkgMjcuMDYzQzQzMi4yOCAyNy4wNjMgNDI1LjgwOSAzNS42MjQ0IDQyNS44MDkgNDUuMDYyMVoiIGZpbGw9IiMyMTAwNEIiLz4KPHBhdGggZD0iTTQ3MC4wNDEgMTEuNjI1NUg0NjUuMjU1VjYyLjE4NDlINDcwLjA0MVY0MS44OTM4QzQ3MC4wNDEgMzQuODgyOSA0NzQuNTU4IDMxLjI0MjYgNDgwLjM1NSAzMS4yNDI2QzQ4Ni40OSAzMS4yNDI2IDQ4OS4zODkgMzUuMDE3NyA0ODkuMzg5IDQxLjIxOTZWNjIuMTg0OUg0OTQuMTc1VjQwLjI3NTlDNDk0LjE3NSAzMi42NTgyIDQ4OS42NTggMjcuMDYzIDQ4MS4xNjQgMjcuMDYzQzQ3NC43NiAyNy4wNjMgNDcxLjI1NSAzMC41Njg1IDQ3MC4wNDEgMzIuNjU4MlYxMS42MjU1WiIgZmlsbD0iIzIxMDA0QiIvPgo8cGF0aCBkPSJNMC44MjQ5NTEgNzMuOTkzTDI0LjA2ODggMTQuNTIyNEMyNy4zNDQzIDYuMTQxNzkgMzUuNDIyMyAwLjYyNTk3NyA0NC40MjAyIDAuNjI1OTc3SDU4LjQzMzZMMzUuMTg5OCA2MC4wOTY2QzMxLjkxNDMgNjguNDc3MiAyMy44MzYzIDczLjk5MyAxNC44MzgzIDczLjk5M0gwLjgyNDk1MVoiIGZpbGw9InVybCgjcGFpbnQwX2xpbmVhcl8wXzE1KSIvPgo8cGF0aCBkPSJNMzQuOTI0NiA3My45OTMyTDU4LjE2ODQgMTQuNTIyNkM2MS40NDM5IDYuMTQxOTcgNjkuNTIxOSAwLjYyNjE1MiA3OC41MTk5IDAuNjI2MTUySDkyLjUzMzJMNjkuMjg5NCA2MC4wOTY4QzY2LjAxMzkgNjguNDc3NCA1Ny45MzU5IDczLjk5MzIgNDguOTM3OSA3My45OTMySDM0LjkyNDZaIiBmaWxsPSJ1cmwoI3BhaW50MV9saW5lYXJfMF8xNSkiLz4KPHBhdGggZD0iTTY5LjAyNjIgNzMuOTkzMkw5Mi4yNyAxNC41MjI2Qzk1LjU0NTUgNi4xNDE5NyAxMDMuNjIzIDAuNjI2MTUyIDExMi42MjEgMC42MjYxNTJIMTI2LjYzNUwxMDMuMzkxIDYwLjA5NjhDMTAwLjExNSA2OC40Nzc0IDkyLjAzNzUgNzMuOTkzMiA4My4wMzk1IDczLjk5MzJINjkuMDI2MloiIGZpbGw9InVybCgjcGFpbnQyX2xpbmVhcl8wXzE1KSIvPgo8ZGVmcz4KPGxpbmVhckdyYWRpZW50IGlkPSJwYWludDBfbGluZWFyXzBfMTUiIHgxPSIxMjYuNjM1IiB5MT0iLTQuOTc3OTkiIHgyPSIwLjgyNDk1MiIgeTI9IjY2LjA5NzgiIGdyYWRpZW50VW5pdHM9InVzZXJTcGFjZU9uVXNlIj4KPHN0b3Agc3RvcC1jb2xvcj0iI0ZGNUNBQSIvPgo8c3RvcCBvZmZzZXQ9IjEiIHN0b3AtY29sb3I9IiNGRjRFNjIiLz4KPC9saW5lYXJHcmFkaWVudD4KPGxpbmVhckdyYWRpZW50IGlkPSJwYWludDFfbGluZWFyXzBfMTUiIHgxPSIxMjYuNjM1IiB5MT0iLTQuOTc3OTkiIHgyPSIwLjgyNDk1MiIgeTI9IjY2LjA5NzgiIGdyYWRpZW50VW5pdHM9InVzZXJTcGFjZU9uVXNlIj4KPHN0b3Agc3RvcC1jb2xvcj0iI0ZGNUNBQSIvPgo8c3RvcCBvZmZzZXQ9IjEiIHN0b3AtY29sb3I9IiNGRjRFNjIiLz4KPC9saW5lYXJHcmFkaWVudD4KPGxpbmVhckdyYWRpZW50IGlkPSJwYWludDJfbGluZWFyXzBfMTUiIHgxPSIxMjYuNjM1IiB5MT0iLTQuOTc3OTkiIHgyPSIwLjgyNDk1MiIgeTI9IjY2LjA5NzgiIGdyYWRpZW50VW5pdHM9InVzZXJTcGFjZU9uVXNlIj4KPHN0b3Agc3RvcC1jb2xvcj0iI0ZGNUNBQSIvPgo8c3RvcCBvZmZzZXQ9IjEiIHN0b3AtY29sb3I9IiNGRjRFNjIiLz4KPC9saW5lYXJHcmFkaWVudD4KPC9kZWZzPgo8L3N2Zz4="}
}

func (s *Search) SearchContents(_ context.Context, cond *plugin.SearchBasicCond) (
	res []plugin.SearchResult, total int64, err error) {
	if s.Client == nil {
		return nil, 0, configuredErr
	}
	query, searchRequest := s.buildQuery(cond)
	searchRequest.Filter = s.buildFilter(cond)

	index := s.Client.Index(s.Config.IndexName)
	searchResult, err := index.Search(query, searchRequest)
	if err != nil {
		log.Errorf("meilisearch error: %s", err.Error())
		return nil, 0, err
	}
	return s.warpResult(searchResult)
}

func (s *Search) SearchQuestions(_ context.Context, cond *plugin.SearchBasicCond) (
	res []plugin.SearchResult, total int64, err error) {
	if s.Client == nil {
		return nil, 0, configuredErr
	}
	query, searchRequest := s.buildQuery(cond)

	filter := s.buildFilter(cond)
	filter = append(filter, "type = question")
	searchRequest.Filter = filter

	index := s.Client.Index(s.Config.IndexName)
	searchResult, err := index.Search(query, searchRequest)
	if err != nil {
		log.Errorf("search error: %s", err.Error())
		return nil, 0, err
	}
	return s.warpResult(searchResult)
}

func (s *Search) SearchAnswers(_ context.Context, cond *plugin.SearchBasicCond) (
	res []plugin.SearchResult, total int64, err error) {
	if s.Client == nil {
		return nil, 0, configuredErr
	}

	query, searchRequest := s.buildQuery(cond)
	filter := s.buildFilter(cond)
	filter = append(filter, "type = answer")
	searchRequest.Filter = filter

	index := s.Client.Index(s.Config.IndexName)
	searchResult, err := index.Search(query, searchRequest)
	if err != nil {
		log.Errorf("search error: %s", err.Error())
		return nil, 0, err
	}
	return s.warpResult(searchResult)
}

func (s *Search) UpdateContent(_ context.Context, content *plugin.SearchContent) error {
	if s.Client == nil {
		return configuredErr
	}

	index := s.Client.Index(s.Config.IndexName)
	if s.Config.Async {
		_, err := index.AddDocuments([]*plugin.SearchContent{content}, primaryKey)
		return err
	} else {
		resp, err := index.AddDocuments([]*plugin.SearchContent{content}, primaryKey)
		if err != nil {
			return err
		}
		return waitForTask(s.Client, resp)
	}
}

func (s *Search) DeleteContent(_ context.Context, contentID string) error {
	if s.Client == nil {
		return configuredErr
	}

	index := s.Client.Index(s.Config.IndexName)
	if s.Config.Async {
		_, err := index.DeleteDocument(contentID)
		return err
	} else {
		resp, err := index.DeleteDocument(contentID)
		if err != nil {
			return err
		}
		return waitForTask(s.Client, resp)
	}
}

func (s *Search) RegisterSyncer(ctx context.Context, syncer plugin.SearchSyncer) {
	s.syncer = syncer
	go s.sync(ctx)
}

func (s *Search) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Name:        "host",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigHostTitle),
			Description: plugin.MakeTranslator(i18n.ConfigHostDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: s.Config.Host,
		},
		{
			Name:        "api_key",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigApiKeyTitle),
			Description: plugin.MakeTranslator(i18n.ConfigApiKeyDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: s.Config.ApiKey,
		},
		{
			Name:        "index_name",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigIndexTitle),
			Description: plugin.MakeTranslator(i18n.ConfigIndexDescription),
			Required:    false,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: s.Config.IndexName,
		},
		{
			Name:        "async",
			Type:        plugin.ConfigTypeSwitch,
			Title:       plugin.MakeTranslator(i18n.ConfigAsyncTitle),
			Description: plugin.MakeTranslator(i18n.ConfigAsyncDescription),
			Required:    false,
			Value:       s.Config.Async,
		},
	}
}

func (s *Search) ConfigReceiver(config []byte) error {
	conf := &SearchConfig{}
	_ = json.Unmarshal(config, conf)

	// if index name is empty, use default index name
	if conf.IndexName == "" {
		conf.IndexName = defaultIndexName
	}
	s.Config = conf

	log.Debugf("try to init meilisearch client: %s", conf.Host)

	s.Client = meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   conf.Host,
		APIKey: conf.ApiKey,
	})

	s.tryToCreateIndex()

	index := s.Client.Index(conf.IndexName)
	_, err := index.UpdateSearchableAttributes(&[]string{"title", "content"})
	if err != nil {
		log.Errorf("update searchable attributes error: %s", err.Error())
		return err
	}
	_, err = index.UpdateFilterableAttributes(&[]string{"title", "content", "tags", "status", "answers", "type", "questionID", "userID", "views", "created", "active", "score", "hasAccepted"})
	if err != nil {
		log.Errorf("update filterable attributes error: %s", err.Error())
		return err
	}
	_, err = index.UpdateSortableAttributes(&[]string{"active", "created", "active", "score"})
	if err != nil {
		log.Errorf("update sortable attributes error: %s", err.Error())
		return err
	}
	_, err = index.UpdateDisplayedAttributes(&[]string{"title", "content", "objectID", "type"})
	if err != nil {
		log.Errorf("update displayed attributes error: %s", err.Error())
		return err
	}
	return nil
}

func (s *Search) warpResult(resp *meilisearch.SearchResponse) ([]plugin.SearchResult, int64, error) {
	res := make([]plugin.SearchResult, 0)
	for _, hit := range resp.Hits {

		var content plugin.SearchContent
		bytes, err := json.Marshal(hit)
		if err != nil {
			log.Errorf("marshal hit error: %s", err.Error())
			return nil, 0, err
		}
		err = json.Unmarshal(bytes, &content)
		if err != nil {
			log.Errorf("unmarshal hit error: %s", err.Error())
			return nil, 0, err
		}

		res = append(res, plugin.SearchResult{
			ID:   content.ObjectID,
			Type: content.Type,
		})
	}
	log.Debugf("search result: %d", len(res))
	return res, resp.TotalHits, nil
}

func (s *Search) buildQuery(cond *plugin.SearchBasicCond) (string, *meilisearch.SearchRequest) {
	searchRequest := meilisearch.SearchRequest{}

	// page
	if cond.Page > 0 {
		searchRequest.Page = int64(cond.Page)
	}
	if cond.PageSize > 0 {
		searchRequest.HitsPerPage = int64(cond.PageSize)
	}

	// order
	switch cond.Order {
	case plugin.SearchNewestOrder:
		searchRequest.Sort = []string{"created:desc"}
	case plugin.SearchActiveOrder:
		searchRequest.Sort = []string{"created:desc"}
	case plugin.SearchScoreOrder:
		searchRequest.Sort = []string{"score:desc"}
	}

	var query string
	if cond.Words != nil && len(cond.Words) > 0 {
		query = strings.Join(cond.Words, " ")
	}
	return query, &searchRequest
}

func (s *Search) buildFilter(cond *plugin.SearchBasicCond) []string {
	var filter []string
	if cond.TagIDs != nil && len(cond.TagIDs) > 0 {
		for _, tagGroup := range cond.TagIDs {
			if len(tagGroup) > 0 {
				filter = append(filter, fmt.Sprintf("tags IN [%s]", strings.Join(tagGroup, ",")))
			}
		}
	}
	if cond.UserID != "" {
		filter = append(filter, fmt.Sprintf("userID = %s", cond.UserID))
	}
	// QuestionAccepted
	if cond.QuestionAccepted == plugin.AcceptedCondTrue {
		filter = append(filter, "hasAccepted = true")
	} else if cond.QuestionAccepted == plugin.AcceptedCondFalse {
		filter = append(filter, "hasAccepted = false")
	}

	// AnswerAccepted
	if cond.AnswerAccepted == plugin.AcceptedCondTrue {
		filter = append(filter, "hasAccepted = true")
	} else if cond.AnswerAccepted == plugin.AcceptedCondFalse {
		filter = append(filter, "hasAccepted = false")
	}

	// QuestionID
	if cond.QuestionID != "" {
		filter = append(filter, fmt.Sprintf("questionID = %s", cond.QuestionID))
	}

	// VoteAmount
	if cond.VoteAmount > 0 {
		filter = append(filter, fmt.Sprintf("voteAmount >= %d", cond.VoteAmount))
	}

	// ViewAmount
	if cond.ViewAmount > 0 {
		filter = append(filter, fmt.Sprintf("viewAmount >= %d", cond.ViewAmount))
	}

	// AnswerAmount
	if cond.AnswerAmount > 0 {
		filter = append(filter, fmt.Sprintf("answerAmount >= %d", cond.AnswerAmount))
	}
	return filter
}

func waitForTask(client *meilisearch.Client, resp *meilisearch.TaskInfo) error {
	task, err := client.WaitForTask(resp.TaskUID)
	if err != nil {
		return err
	}
	if task.Status != meilisearch.TaskStatusSucceeded {
		err = errors.InternalServer(task.Error.Code).WithMsg("invoke meili failed:" + task.Error.Message).WithStack()
		return err
	}
	return nil
}
