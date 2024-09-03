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

package lark_test

import (
	"testing"

	answer "github.com/apache/incubator-answer-plugins/notification-lark"
)

func TestGenRandomStr(t *testing.T) {
	rs := answer.GenerateRandomString(&answer.GenerateRandomStringArgs{
		Length:     10,
		StringPool: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	})

	t.Logf("Random string: %s", rs)
}

func TestPtrBool(t *testing.T) {
	b := true
	bp := answer.PtrBool(b)
	if *bp != b {
		t.Errorf("PtrBool failed")
	}
}
