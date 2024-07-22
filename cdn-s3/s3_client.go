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

package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"strings"
)

type Client struct {
	s3Config *aws.Config
	bucket   string
}

func NewS3Client(id, secret, token, endpoint, region, bucket string, disableSSL bool) *Client {
	s3Client := &Client{
		s3Config: &aws.Config{
			Credentials:      credentials.NewStaticCredentials(id, secret, token),
			Endpoint:         aws.String(endpoint),
			Region:           aws.String(region),
			DisableSSL:       aws.Bool(disableSSL),
			S3ForcePathStyle: aws.Bool(true),
		},
		bucket: bucket,
	}
	return s3Client
}

func (s *Client) PutObject(key, ext string, file io.ReadSeeker) (err error) {
	newSession, err := session.NewSession(s.s3Config)
	if err != nil {
		return fmt.Errorf("failed to create session, %s", err.Error())
	}

	extType := strings.TrimPrefix(ext, ".")
	contentType := ""
	switch extType {
	case "jpg", "jpeg", "png":
		contentType = "image/" + extType
	case "svg":
		contentType = "image/svg+xml"
	case "js":
		contentType = "application/javascript"
	case "css":
		contentType = "text/css"
	case "map":
		contentType = "application/json"
	case "woff":
		contentType = "application/font-woff"
	case "woff2":
		contentType = "application/font-woff2"
	}

	_, err = s3.New(newSession).PutObject(&s3.PutObjectInput{
		Body:        file,
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("failed to put object, %s", err.Error())
	}
	return err
}
