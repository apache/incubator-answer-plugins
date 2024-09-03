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
	"github.com/apache/incubator-answer/plugin"
	"github.com/segmentfault/pacman/log"
)

func (s *SearchEngine) sync() {
	var page, pageSize = 1, 100
	if s.syncing {
		log.Warnf("es: syncing is running, skip")
		return
	}

	go func() {
		s.lock.Lock()
		defer s.lock.Unlock()
		if s.syncing {
			log.Warnf("es: syncing is running, skip")
			return
		}

		s.syncing = true
		log.Info("es: start sync questions...")
		page = 1
		for {
			log.Infof("es: sync question page %d, page size %d", page, pageSize)
			questionList, err := s.syncer.GetQuestionsPage(context.TODO(), page, pageSize)
			if err != nil {
				log.Error("es: sync questions error", err)
				break
			}
			if len(questionList) == 0 {
				break
			}
			err = s.batchUpdateContent(context.TODO(), questionList)
			if err != nil {
				log.Error("es: sync questions error", err)
			}
			page += 1
		}

		log.Info("es: start sync answers...")
		page = 1
		for {
			log.Infof("es: sync answer page %d, page size %d", page, pageSize)
			answerList, err := s.syncer.GetAnswersPage(context.TODO(), page, pageSize)
			if err != nil {
				log.Error("es: sync answers error", err)
				break
			}

			if len(answerList) == 0 {
				break
			}

			err = s.batchUpdateContent(context.TODO(), answerList)
			if err != nil {
				log.Error("es: sync answers error", err)
			}

			page += 1
		}
		s.syncing = false
		log.Info("es: sync done")
	}()
}

func (s *SearchEngine) batchUpdateContent(ctx context.Context, contents []*plugin.SearchContent) (err error) {
	for _, content := range contents {
		err = s.Operator.SaveDoc(ctx, s.getIndexName(), content.ObjectID, CreateDocFromSearchContent(content.ObjectID, content))
		if err != nil {
			return
		}
	}

	return nil
}
