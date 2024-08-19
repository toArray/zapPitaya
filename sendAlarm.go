package zapPitaya

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/topfreegames/pitaya/v2/logger"
	"io/ioutil"
	"net/http"
	"time"
)

// MarkdownContent 企业微信消息 Markdown 数据
type MarkdownContent struct {
	Content string `json:"content"`
}

// MarkdownType 企业微信消息Markdown格式数据
type MarkdownType struct {
	MsgType  string           `json:"msgtype"`
	Markdown *MarkdownContent `json:"markdown"`
}

// TrySendAlarmOfMarkDown 执行告警,失败重试(markdown格式)
// @Param 	webHook string	企业微信机器人地址
// @Param  	retryCount int	重试次数
// @Param  	content *MarkdownContent	告警内容
func TrySendAlarmOfMarkDown(webHook string, retryCount int, content *MarkdownContent) {

	// 序列化数据
	data := new(MarkdownType)
	data.MsgType = "markdown"
	data.Markdown = &MarkdownContent{Content: content.Content}
	dataByte, err := json.Marshal(data)
	if err != nil {
		logger.Log.Errorf("send alarm to wechat fialed for json.Marshal err: %v", err)
		return
	}

	// 发送
	for i := 0; i < retryCount; i++ {
		if err = sendAlarm(webHook, dataByte); err != nil {
			if i < retryCount-1 {
				time.Sleep(2 * time.Second)
			}
		} else {
			break
		}
	}
}

// sendAlarm 执行post请求发送消息到企业微信机器人
// @Param webHook string 机器人地址
// @Param data []byte 消息实体
func sendAlarm(webHook string, data []byte) (err error) {

	// 创建请求失败
	client := &http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest("POST", webHook, bytes.NewBuffer(data))
	if err != nil {
		logger.Log.Errorf("send alarm to wechat fialed for create http request err: %v", err)
		return
	}

	// 发送请求失败
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Errorf("send alarm to wechat fialed for send http request err: %v", err)
		return
	}
	defer resp.Body.Close()

	// 读取body失败
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Errorf("send alarm to wechat fialed for read body err: %v", err)
		return
	}

	// 响应code不正确
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("send alarm to wechat fialed for response code != 200 code: %d", resp.StatusCode)
		logger.Log.Error(err)
		return
	}

	logger.Log.Infof("send alarm to wechat success. code: %d, body: %s", resp.StatusCode, string(body))
	return
}
