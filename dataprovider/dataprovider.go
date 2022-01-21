package dataprovider

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Dataprovider struct {
	Faileds map[string]*Failed
	// map并发安全读写锁
	lock *sync.RWMutex
}

type Failed struct {
	StartTime int64
	First     bool
	Message   string
}

func NewDataprovider() *Dataprovider {
	return &Dataprovider{
		Faileds: make(map[string]*Failed),
		lock:    &sync.RWMutex{},
	}
}

func (d *Dataprovider) Add(id, msg string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.Faileds[id] == nil {
		d.Faileds[id] = &Failed{
			StartTime: time.Now().Unix(),
			First:     true,
			Message:   msg,
		}
		return
	}
}
func (d *Dataprovider) Del(id string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.Faileds[id] != nil {
		delete(d.Faileds, id)
	}
}

func (d *Dataprovider) Trigger(failedCh chan Failed, interval time.Duration) {
	go func() {
		for {
			for id, failed := range d.Faileds {
				// 第一次触发告警
				if !failed.First {
					continue
				}
				logrus.Info("first trigger")
				failedCh <- *failed
				d.lock.RLock()
				d.Faileds[id].First = false
				d.lock.RUnlock()
			}
			time.Sleep(150 * time.Millisecond)
		}
	}()
	for range time.Tick(interval) {
		// 后续定期告警
		for _, failed := range d.Faileds {
			if failed.First {
				continue
			}
			logrus.Info("periodic trigger")
			count := (time.Now().Unix() - failed.StartTime) / 60
			failedCh <- Failed{
				Id:        failed.Id,
				StartTime: failed.StartTime,
				First:     failed.First,
				Message:   failed.Message + fmt.Sprintf("持续时间%d/分钟", count),
			}
		}
	}
}
