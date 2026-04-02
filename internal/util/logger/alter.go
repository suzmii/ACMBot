package logger

import (
	serverchansdk "github.com/easychen/serverchan-sdk-golang"
	"github.com/sirupsen/logrus"
)

type AlertHook struct {
	SendKey string
}

func (h *AlertHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func (h *AlertHook) Fire(entry *logrus.Entry) error {
	go func(msg string, sendKey string) {
		_, err := serverchansdk.ScSend(sendKey, "ACMBot", msg, &serverchansdk.ScSendOptions{
			Tags: "服务器报警|ACMBot",
		})
		if err != nil {
			logrus.Info("failed to send alert to serverchan: %v", err)
		}
		logrus.Info(">>> [ALERT] 发送提示: ", msg)
	}(entry.Message, h.SendKey)
	return nil
}
