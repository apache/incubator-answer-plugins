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
	"github.com/meilisearch/meilisearch-go"
	"github.com/segmentfault/pacman/log"
)

// try to create index if not exist
func (s *Search) tryToCreateIndex() {
	index, err := s.Client.GetIndex(s.Config.IndexName)
	if index != nil {
		log.Infof("index %s already exist, skip create", s.Config.IndexName)
		return
	}
	if err != nil && index == nil {
		log.Infof("get index failed %s, maybe not exist, try to create", err)
	}

	log.Infof("try to create index %s", s.Config.IndexName)
	resp, err := s.Client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        s.Config.IndexName,
		PrimaryKey: primaryKey,
	})
	if err != nil {
		log.Errorf("create index error: %s", err.Error())
		return
	}
	if err = waitForTask(s.Client, resp); err != nil {
		log.Errorf("create index error: %s", err.Error())
	} else {
		log.Infof("create index %s success", s.Config.IndexName)
	}
	return
}
