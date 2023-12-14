package algolia

import (
	"context"
	"github.com/apache/incubator-answer/plugin"
	"github.com/segmentfault/pacman/log"
)

func (s *SearchAlgolia) sync() {
	var page, pageSize = 1, 100
	go func() {
		log.Info("algolia: start sync questions...")
		page = 1
		for {
			log.Infof("algolia: sync question page %d, page size %d", page, pageSize)
			questionList, err := s.syncer.GetQuestionsPage(context.TODO(), page, pageSize)
			if err != nil {
				log.Error("algolia: sync questions error", err)
				break
			}
			if len(questionList) == 0 {
				break
			}
			err = s.batchUpdateContent(context.TODO(), questionList)
			if err != nil {
				log.Error("algolia: sync questions error", err)
			}
			page += 1
		}

		log.Info("algolia: start sync answers...")
		page = 1
		for {
			log.Infof("algolia: sync answer page %d, page size %d", page, pageSize)
			answerList, err := s.syncer.GetAnswersPage(context.TODO(), page, pageSize)
			if err != nil {
				log.Error("algolia: sync answers error", err)
				break
			}

			if len(answerList) == 0 {
				break
			}

			err = s.batchUpdateContent(context.TODO(), answerList)
			if err != nil {
				log.Error("algolia: sync answers error", err)
			}

			page += 1
		}
		log.Info("algolia: sync done")
	}()
}

func (s *SearchAlgolia) batchUpdateContent(ctx context.Context, contents []*plugin.SearchContent) (err error) {
	res, err := s.getIndex("").SaveObjects(contents)
	if err != nil {
		return
	}
	err = res.Wait()
	return
}
