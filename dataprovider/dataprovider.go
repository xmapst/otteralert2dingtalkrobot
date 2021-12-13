package dataprovider

import (
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
	for {
		d.lock.RLock()
		for id, failed := range d.Faileds {
			if failed.First {
				failedCh <- *failed
				logrus.Infof("id %s generation time %s", id, time.Unix(failed.StartTime, 0).Format("2006-01-02 15:04:05"))
				d.lock.RUnlock()
				d.lock.Lock()
				d.Faileds[id].First = false
				d.lock.Unlock()
				d.lock.RLock()
				continue
			}
			if (time.Now().Unix() - failed.StartTime) > int64(interval/time.Second) {
				failedCh <- *failed
				logrus.Infof("id %s generation time %s", id, time.Unix(failed.StartTime, 0).Format("2006-01-02 15:04:05"))
				d.lock.RUnlock()
				d.lock.Lock()
				delete(d.Faileds, id)
				d.lock.Unlock()
				d.lock.RLock()
			}
		}
		d.lock.RUnlock()
		time.Sleep(150 * time.Millisecond)
	}
}
