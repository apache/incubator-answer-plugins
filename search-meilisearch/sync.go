package meilisearch

import (
	"context"
	"github.com/answerdev/answer/plugin"
	"github.com/segmentfault/pacman/log"
)

const (
	MaxGetPageSize = 1000
	MaxPutPerSize  = 100
)

// sync data that already exist in Answer to meilisearch
func (s *Search) sync(ctx context.Context) {
	log.Infof("start to sync question data to meilisearch")
	if s.syncing {
		log.Warnf("syncing is running, skip")
		return
	}

	syncFns := []func(ctx context.Context, page, pageSize int) (answerList []*plugin.SearchContent, err error){
		s.syncer.GetQuestionsPage,
		s.syncer.GetAnswersPage,
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.syncing {
		log.Warnf("syncing is running, skip")
		return
	}

	s.syncing = true
	for _, fn := range syncFns {
		s.syncQuestionAndAnswerData(ctx, fn)
	}
	s.syncing = false
}

func (s *Search) syncQuestionAndAnswerData(ctx context.Context,
	syncFunc func(ctx context.Context, page, pageSize int) (answerList []*plugin.SearchContent, err error)) {
	log.Infof("start to sync data to meilisearch")
	for page, pageSize := 1, MaxGetPageSize; ; page++ {
		log.Infof("start to sync page %d", page)
		dataList, err := syncFunc(ctx, page, pageSize)
		if err != nil {
			log.Errorf("get data failed %s", err)
			return
		}
		if len(dataList) == 0 {
			log.Infof("get page %d success, no other data, stop sync", page)
			break
		}
		log.Infof("get page %d success, record count %d", page, len(dataList))
		for i := 0; i < len(dataList); i += 100 {
			end := i + MaxPutPerSize
			if i+MaxPutPerSize > len(dataList) {
				end = len(dataList)
			}
			resp, err := s.Client.Index(s.Config.IndexName).AddDocuments(
				dataList[i:end], primaryKey)
			if err != nil {
				log.Errorf("add documents failed %s", err)
				return
			}
			if err := waitForTask(s.Client, resp); err != nil {
				log.Errorf("wait for task failed %s", err)
			}
			log.Infof("sync page %d, progress %d/%d", page, end, len(dataList))
		}
	}
}
