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
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-answer-plugins/cdn-s3/i18n"
	"github.com/apache/incubator-answer-plugins/util"
	"github.com/apache/incubator-answer/ui"
	"github.com/segmentfault/pacman/log"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-answer/plugin"
)

var (
	staticPath = os.Getenv("ANSWER_STATIC_PATH")
	enable     = false
)

//go:embed  info.yaml
var Info embed.FS

const (
	// 10MB
	defaultMaxFileSize int64 = 10 * 1024 * 1024
)

type CDN struct {
	Config *CDNConfig
	Client *Client
}

type CDNConfig struct {
	Endpoint        string `json:"endpoint"`
	BucketName      string `json:"bucket_name"`
	ObjectKeyPrefix string `json:"object_key_prefix"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	AccessToken     string `json:"access_token"`
	VisitUrlPrefix  string `json:"visit_url_prefix"`
	MaxFileSize     string `json:"max_file_size"`
	Region          string `json:"region"`
	DisableSSL      bool   `json:"disable_ssl"`
}

func init() {
	plugin.Register(&CDN{
		Config: &CDNConfig{},
	})
}

func (c *CDN) Info() plugin.Info {
	info := util.Info{}
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

// GetStaticPrefix get static prefix
func (c *CDN) GetStaticPrefix() string {
	if !enable {
		return ""
	}
	return c.Config.VisitUrlPrefix + c.Config.ObjectKeyPrefix
}

// scanFiles scan all the static files in the build directory
func (c *CDN) scanFiles() {
	if staticPath == "" {
		err := c.scanEmbedFiles("build")
		if err != nil {
			enable = false
			log.Error("failed: scan embed files: ", err)
			return
		}
		log.Info("complete: scan embed files")
		enable = true
		return
	}

	err := c.scanStaticPathFiles(staticPath)
	if err != nil {
		enable = false
		log.Error("fialed: scan static path files: ", err)
		return
	}
	enable = true
	log.Info("complete: scan static path files")
}

// scanStaticPathFiles scan static path files
func (c *CDN) scanStaticPathFiles(fileName string) (err error) {
	if len(fileName) == 0 {
		return
	}

	// scan static path files
	entry, err := os.ReadDir(fileName)
	if err != nil {
		log.Error("read static dir failed: %v", err)
		return
	}
	for _, info := range entry {
		if info.IsDir() {
			err = c.scanStaticPathFiles(filepath.Join(fileName, info.Name()))
			if err != nil {
				return
			}
		}

		var file *os.File
		filePath := filepath.Join(fileName, info.Name())
		fi, _ := info.Info()
		size := fi.Size()
		file, err = os.Open(filePath)
		if err != nil {
			log.Error("open file failed: %v", err)
			return
		}

		suffix := staticPath[:1]
		if suffix != "/" {
			suffix = ""
		}
		filePath = strings.TrimPrefix(filePath, staticPath+suffix)

		// rebuild custom io.Reader
		ns := strings.Split(info.Name(), ".")
		if info.Name() == "asset-manifest.json" {
			err = c.Upload(filePath, c.rebuildReader(file, map[string]string{
				"\"/static": "",
			}), size)
			if err != nil {
				return
			}
			continue
		}
		if ns[0] == "main" {
			ext := strings.ToLower(filepath.Ext(filePath))
			if ext == ".js" || ext == ".map" {
				err = c.Upload(filePath, c.rebuildReader(file, map[string]string{
					"\"static": "",
					"=\"/\",":  "=\"\",",
				}), size)
				if err != nil {
					return
				}
				continue
			}

			if ext == ".css" {
				err = c.Upload(filePath, c.rebuildReader(file, map[string]string{
					"url(/static": "url(../../static",
				}), size)
				if err != nil {
					return
				}
				continue
			}
		}

		err = c.Upload(filePath, file, size)
		if err != nil {
			return
		}
	}
	return
}

func (c *CDN) scanEmbedFiles(fileName string) (err error) {
	entry, err := ui.Build.ReadDir(fileName)
	if err != nil {
		log.Error("read static dir failed: %v", err)
		return
	}
	for _, info := range entry {
		if info.IsDir() {
			err = c.scanEmbedFiles(filepath.Join(fileName, info.Name()))
			if err != nil {
				return
			}
			continue
		}

		var file fs.File
		filePath := filepath.Join(fileName, info.Name())
		fi, _ := info.Info()
		size := fi.Size()
		file, err = ui.Build.Open(filePath)
		defer file.Close()
		if err != nil {
			log.Error("open file failed: %v", err)
			return
		}

		filePath = strings.TrimPrefix(filePath, "build/")

		// rebuild custom io.Reader
		ns := strings.Split(info.Name(), ".")
		if info.Name() == "asset-manifest.json" {
			err = c.Upload(filePath, c.rebuildReader(file, map[string]string{
				"\"/static": "",
			}), size)
			if err != nil {
				return
			}
			continue
		}

		if ns[0] == "main" {
			ext := strings.ToLower(filepath.Ext(filePath))
			if ext == ".js" || ext == ".map" {
				err = c.Upload(filePath, c.rebuildReader(file, map[string]string{
					"\"static": "",
					"=\"/\",":  "=\"\",",
				}), size)
				if err != nil {
					return
				}
				continue
			}

			if ext == ".css" {
				err = c.Upload(filePath, c.rebuildReader(file, map[string]string{
					"url(/static": "url(../../static",
				}), size)
				if err != nil {
					return
				}
				continue
			}
		}

		err = c.Upload(filePath, c.rebuildReader(file, nil), size)
		if err != nil {
			return
		}
	}
	return
}

func (c *CDN) rebuildReader(file io.Reader, replaceMap map[string]string) io.ReadSeeker {
	var (
		bufr = make([]byte, 0)
		res  string
	)

	for {
		buf := make([]byte, 1024)
		n, err := file.Read(buf)
		if err != nil {
			break
		}
		bufr = append(bufr, buf[:n]...)
	}

	res = string(bufr)

	if replaceMap != nil {
		for oldStr, newStr := range replaceMap {
			if oldStr != "" {
				if newStr == "" {
					prefix := c.Config.VisitUrlPrefix + c.Config.ObjectKeyPrefix
					if prefix[len(prefix)-1:] == "/" {
						prefix = strings.TrimSuffix(prefix, "/")
					}
					newStr = "\"" + prefix + "/static"
				}
				res = strings.ReplaceAll(res, oldStr, newStr)
			}
		}
	}

	return strings.NewReader(res)
}
func (c *CDN) Upload(filePath string, file io.ReadSeeker, size int64) (err error) {

	if !c.CheckFileType(filePath) {
		log.Error(plugin.MakeTranslator(i18n.ErrUnsupportedFileType), filePath)
		return
	}

	if size > c.maxFileSizeLimit() {
		log.Error(plugin.MakeTranslator(i18n.ErrOverFileSizeLimit))
		return
	}

	objectKey := c.createObjectKey(filePath)

	err = c.Client.PutObject(objectKey, strings.ToLower(filepath.Ext(filePath)), file)
	if err != nil {
		log.Error(plugin.MakeTranslator(i18n.ErrUploadFileFailed), err)
		return
	}
	return c.checkCDNAvailable(objectKey)
}

func (c *CDN) checkCDNAvailable(objectKey string) error {
	url := c.Config.VisitUrlPrefix + objectKey
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Error("check error:", url)
		return fmt.Errorf("failed to get object, %s", response.Status)
	}
	return nil
}

func (c *CDN) createObjectKey(filePath string) string {
	return c.Config.ObjectKeyPrefix + filePath
}

func (c *CDN) randomObjectKey() string {
	bytes := make([]byte, 4)
	_, _ = rand.Read(bytes)
	return fmt.Sprintf("%d", time.Now().UnixNano()) + hex.EncodeToString(bytes)
}

func (c *CDN) CheckFileType(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	if _, ok := plugin.DefaultCDNFileType[ext]; ok {
		return true
	}
	return false
}

func (c *CDN) maxFileSizeLimit() int64 {
	if len(c.Config.MaxFileSize) == 0 {
		return defaultMaxFileSize
	}
	limit, _ := strconv.Atoi(c.Config.MaxFileSize)
	if limit <= 0 {
		return defaultMaxFileSize
	}
	return int64(limit) * 1024 * 1024
}

func (c *CDN) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Name:        "endpoint",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigEndpointTitle),
			Description: plugin.MakeTranslator(i18n.ConfigEndpointDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: c.Config.Endpoint,
		},
		{
			Name:        "bucket_name",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigBucketNameTitle),
			Description: plugin.MakeTranslator(i18n.ConfigBucketNameDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: c.Config.BucketName,
		},
		{
			Name:        "object_key_prefix",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigObjectKeyPrefixTitle),
			Description: plugin.MakeTranslator(i18n.ConfigObjectKeyPrefixDescription),
			Required:    false,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: c.Config.ObjectKeyPrefix,
		},
		{
			Name:        "access_key_id",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigAccessKeyIdTitle),
			Description: plugin.MakeTranslator(i18n.ConfigAccessKeyIdDescription),
			Required:    false,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: c.Config.AccessKeyID,
		},
		{
			Name:        "access_key_secret",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigAccessKeySecretTitle),
			Description: plugin.MakeTranslator(i18n.ConfigAccessKeySecretDescription),
			Required:    false,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: c.Config.AccessKeySecret,
		},
		{
			Name:        "access_token",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigAccessTokenTitle),
			Description: plugin.MakeTranslator(i18n.ConfigAccessTokenDescription),
			Required:    false,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: c.Config.AccessToken,
		},
		{
			Name:        "visit_url_prefix",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigVisitUrlPrefixTitle),
			Description: plugin.MakeTranslator(i18n.ConfigVisitUrlPrefixDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: c.Config.VisitUrlPrefix,
		},
		{
			Name:        "max_file_size",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigMaxFileSizeTitle),
			Description: plugin.MakeTranslator(i18n.ConfigMaxFileSizeDescription),
			Required:    false,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeNumber,
			},
			Value: c.Config.MaxFileSize,
		},
		{
			Name:        "region",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigRegionTitle),
			Description: plugin.MakeTranslator(i18n.ConfigRegionDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: c.Config.Region,
		},
		{
			Name:  "disable_ssl",
			Type:  plugin.ConfigTypeSwitch,
			Title: plugin.MakeTranslator(i18n.ConfigDisableSSLTitle),
			Value: c.Config.DisableSSL,
			UIOptions: plugin.ConfigFieldUIOptions{
				Label: plugin.MakeTranslator(i18n.ConfigDisableSSLDescription),
			},
		},
	}
}

func (c *CDN) ConfigReceiver(config []byte) error {
	cfg := &CDNConfig{}
	_ = json.Unmarshal(config, cfg)
	c.Config = cfg
	c.Client = NewS3Client(
		c.Config.AccessKeyID,
		c.Config.AccessKeySecret,
		c.Config.AccessToken,
		c.Config.Endpoint,
		c.Config.Region,
		c.Config.BucketName,
		c.Config.DisableSSL,
	)
	go c.scanFiles()
	return nil
}
