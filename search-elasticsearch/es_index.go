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

import "github.com/apache/incubator-answer/plugin"

var indexJson = `
{
    "settings": {
        "number_of_shards": 3,
        "number_of_replicas": 1
    },
    "mappings": {
        "properties": {
            "id": {
                "type": "keyword",
                "doc_values": false,
                "norms": false,
                "similarity": "boolean"
            },
			"object_id": {
				"type": "keyword"
			},
            "title": {
                "type": "text"
            },
            "type": {
                "type": "text"
            },
            "content": {
                "type": "text"
            },
            "user_id": {
                "type": "keyword"
            },
            "question_id": {
                "type": "keyword"
            },
            "answers": {
                "type": "long"
            },
            "status": {
                "type": "long"
            },
            "views": {
                "type": "long"
            },
            "created": {
                "type": "long"
            },
            "active": {
                "type": "long"
            },
            "score": {
                "type": "long"
            },
            "has_accepted": {
                "type": "boolean"
            },
            "tags": {
                "type": "text",
                "fields": {
                    "keyword": {
                        "type": "keyword"
                    }
                }
            }
        }
    }
}
`

type AnswerPostDoc struct {
	Id          string   `json:"id"`
	ObjectID    string   `json:"object_id"`
	Title       string   `json:"title"`
	Type        string   `json:"type"`
	Content     string   `json:"content"`
	UserID      string   `json:"user_id"`
	QuestionID  string   `json:"question_id"`
	Answers     int64    `json:"answers"`
	Status      int64    `json:"status"`
	Views       int64    `json:"views"`
	Created     int64    `json:"created"`
	Active      int64    `json:"active"`
	Score       int64    `json:"score"`
	HasAccepted bool     `json:"has_accepted"`
	Tags        []string `json:"tags"`
}

func CreateDocFromSearchContent(id string, content *plugin.SearchContent) (doc *AnswerPostDoc) {
	doc = &AnswerPostDoc{}
	doc.Id = id
	doc.ObjectID = content.ObjectID
	doc.Title = content.Title
	doc.Type = content.Type
	doc.Content = content.Content
	doc.UserID = content.UserID
	doc.QuestionID = content.QuestionID
	doc.Answers = content.Answers
	doc.Status = int64(content.Status)
	doc.Views = content.Views
	doc.Created = content.Created
	doc.Active = content.Active
	doc.Score = content.Score
	doc.HasAccepted = content.HasAccepted
	doc.Tags = content.Tags
	return
}

type SearchContent struct {
	ObjectID    string   `json:"objectID"`
	Title       string   `json:"title"`
	Type        string   `json:"type"`
	Content     string   `json:"content"`
	Answers     int64    `json:"answers"`
	Status      int64    `json:"status"`
	Tags        []string `json:"tags"`
	QuestionID  string   `json:"questionID"`
	UserID      string   `json:"userID"`
	Views       int64    `json:"views"`
	Created     int64    `json:"created"`
	Active      int64    `json:"active"`
	Score       int64    `json:"score"`
	HasAccepted bool     `json:"hasAccepted"`
}
