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
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

// initSettings update algolia search settings
func (s *SearchAlgolia) initSettings() (err error) {
	var (
		settings = search.Settings{}
	)
	err = settings.UnmarshalJSON(AlgoliaSearchServerConfig)
	if err != nil {
		return
	}

	// point virtual index to sort
	settings.Replicas = opt.Replicas(
		"virtual("+s.getIndexName(NewestIndex)+")",
		"virtual("+s.getIndexName(ActiveIndex)+")",
		"virtual("+s.getIndexName(ScoreIndex)+")",
	)

	_, err = s.getIndex("").SetSettings(settings, opt.ForwardToReplicas(true))
	if err != nil {
		return
	}
	err = s.initVirtualReplicaSetting()
	return
}

// initVirtualReplicaSetting init virtual index replica setting
func (s *SearchAlgolia) initVirtualReplicaSetting() (err error) {

	_, err = s.getIndex(NewestIndex).SetSettings(search.Settings{
		CustomRanking: opt.CustomRanking(
			"desc(created)",
			"desc(content)",
			"desc(title)"),
	})
	if err != nil {
		return
	}

	_, err = s.getIndex(ActiveIndex).SetSettings(search.Settings{
		CustomRanking: opt.CustomRanking(
			"desc(active)",
			"desc(content)",
			"desc(title)"),
	})
	if err != nil {
		return
	}

	_, err = s.getIndex(ScoreIndex).SetSettings(search.Settings{
		CustomRanking: opt.CustomRanking(
			"desc(score)",
			"desc(content)",
			"desc(title)"),
	})
	if err != nil {
		return
	}
	return
}
