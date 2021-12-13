package engine

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"otteralter/dataprovider"
	"otteralter/utils"
	"strconv"
	"time"
)

var httpClient = &http.Client{
	Transport: &http.Transport{
		Proxy:             http.ProxyFromEnvironment,
		DisableKeepAlives: true,
	},
}

var notification = utils.DingTalkNotification{
	MessageType: "markdown",
	Markdown:    &utils.DingTalkNotificationMarkdown{},
	At: &utils.DingTalkNotificationAt{
		IsAtAll: true,
	},
}

const explode = "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"

var header = fmt.Sprintf("## otter告警--数据库同步错误\n ![警报 图标](%s) \n--- \n**=====侦测到故障=====** ", explode)

func SendNotification(targetURL *utils.URL, dingTalkSecret string, failed dataprovider.Failed) {
	if dingTalkSecret != "" {
		timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
		stringToSign := []byte(timestamp + "\n" + dingTalkSecret)
		mac := hmac.New(sha256.New, []byte(dingTalkSecret))
		mac.Write(stringToSign) // nolint: errcheck
		signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		qs := targetURL.Query()
		qs.Set("timestamp", timestamp)
		qs.Set("sign", signature)
		targetURL.RawQuery = qs.Encode()
	}
	notification.Markdown.Title = "otter告警--数据库同步错误"
	notification.Markdown.Text = createMsgText(failed)

	body, err := json.Marshal(&notification)
	if err != nil {
		logrus.Error(err, "error encoding DingTalk request")
		return
	}
	httpReq, err := http.NewRequest("POST", targetURL.String(), bytes.NewReader(body))
	if err != nil {
		logrus.Error(err, "error building DingTalk request")
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		logrus.Error(err, "error sending notification to DingTalk")
		return
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != 200 {
		logrus.Errorf("unacceptable response code %d", resp.StatusCode)
		return
	}
	var robotResp utils.DingTalkNotificationResponse
	enc := json.NewDecoder(resp.Body)
	if err = enc.Decode(&robotResp); err != nil {
		logrus.Error(err, "error decoding response from DingTalk")
		return
	}
	if robotResp.ErrorCode != 0 {
		logrus.Error("Failed to send notification to DingTalk  respCode ", robotResp.ErrorCode, " respMsg ", robotResp.ErrorMessage)
		return
	}
	logrus.Info("message sent successfully")
}

func createMsgText(failed dataprovider.Failed) string {
	body := fmt.Sprint(
		"\n- 产生时间: ", time.Unix(failed.StartTime, 0).Format("2006-01-02 15:04:05"),
		"\n\n", failed.Message,
	)
	return header + body
}
