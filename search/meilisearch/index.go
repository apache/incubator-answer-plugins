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
