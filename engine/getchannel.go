package engine

import (
	"bytes"
	"fmt"
	"github.com/go-zookeeper/zk"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"otteralter/dataprovider"
	"strconv"
)

var chPrefix = "/otter/channel"
var json = jsoniter.ConfigCompatibleWithStandardLibrary

type ChannelPipelineState struct {
	Active     bool   `json:"active"`
	PipelineID int    `json:"pipelineId"`
	Status     string `json:"status"`
}

var ChState = map[string]int{
	"START": 0,
	"STOP":  1,
	"PAUSE": 2,
	"NONE":  3,
}

func GetChannel(conn *zk.Conn, db *dataprovider.Dataprovider) {
	list, _, err := conn.Children(chPrefix)
	if err != nil {
		logrus.Error("zk连接异常,获取不到节点信息")
		return
	}
	for _, cID := range list {
		chPath := chPrefix + "/" + cID
		cStatus, _, err := conn.Get(chPath)
		if err != nil {
			logrus.Error(err)
			continue
		}
		status := string(bytes.ReplaceAll(cStatus, []byte(`"`), []byte("")))
		if ChState[status] != 0 {
			db.Add(cID, fmt.Sprintf("通道%s异常, 当前状态: %s", cID, status))
		} else {
			db.Del(cID)
		}
		pipelineList, _, err := conn.Children(chPath)
		if err != nil {
			logrus.Error(err)
			continue
		}
		for _, pID := range pipelineList {
			var cps = ChannelPipelineState{}
			res, _, err := conn.Get(chPath + "/" + pID + "/" + "mainstem")
			if err != nil {
				if err != zk.ErrNoNode {
					logrus.Error(err)
				}
				continue
			}
			err = json.Unmarshal(res, &cps)
			if err != nil {
				logrus.Error(err)
				continue
			}
			if !cps.Active {
				db.Add(strconv.Itoa(cps.PipelineID), fmt.Sprintf("通道%s下管道%d异常, 当前状态: %s", cID, cps.PipelineID, cps.Status))
			} else {
				db.Del(strconv.Itoa(cps.PipelineID))
			}
		}
	}
}
