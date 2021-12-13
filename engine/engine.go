package engine

import (
	"github.com/go-zookeeper/zk"
	"github.com/sirupsen/logrus"
	"otteralter/dataprovider"
	"otteralter/utils"
	"time"
)

var (
	Addr           []string
	DingTalkUrl    string
	DingTalkSecret string
	Interval       string
)

func Run() {
	_interval, err := time.ParseDuration(Interval)
	if err != nil {
		logrus.Fatalln(err)
	}
	dingURL, err := utils.ParseURL(DingTalkUrl)
	if err != nil {
		logrus.Fatalln(err)
	}
	conn, _, err := zk.Connect(Addr, time.Second*1, zk.WithLogInfo(false))
	if err != nil {
		logrus.Fatalln(err)
	}
	defer conn.Close()
	failedCh := make(chan dataprovider.Failed)
	db := dataprovider.NewDataprovider()

	// 发送告警
	go func() {
		for failed := range failedCh {
			go SendNotification(dingURL, DingTalkSecret, failed)
		}
	}()
	// 数据采样
	go func() {
		for range time.Tick(150 * time.Millisecond) {
			GetChannel(conn, db)
		}
	}()
	// 数据触发
	db.Trigger(failedCh, _interval)
}
