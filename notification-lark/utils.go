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

package lark

import (
	"crypto/rand"
	"math/big"
)

type GenerateRandomStringArgs struct {
	Length     uint64
	StringPool string
}

// GenerateRandomString use crypto to generate a random string
func GenerateRandomString(args *GenerateRandomStringArgs) string {
	// check args
	if args.Length <= 0 || args.StringPool == "" {
		return ""
	}

	// generate random string
	b := make([]byte, args.Length)
	for i := uint64(0); i < args.Length; i++ {
		idx := RandomInt(0, int64(len(args.StringPool)))
		b[i] = args.StringPool[idx]
	}

	return string(b)
}

func RandomInt(min, max int64) int64 {
	result, _ := rand.Int(rand.Reader, big.NewInt(max-min))
	return result.Int64() + min
}

func PtrBool(b bool) *bool {
	return &b
}
