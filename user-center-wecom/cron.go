package wecom

import (
	"time"

	"github.com/segmentfault/pacman/log"
)

func (uc *UserCenter) CronSyncData() {
	go func() {
		ticker := time.NewTicker(time.Hour)
		for {
			select {
			case <-ticker.C:
				log.Infof("user center try to sync data")
				uc.syncCompany()
			}
		}
	}()
}
