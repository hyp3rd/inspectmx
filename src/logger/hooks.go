package logger

import "github.com/sirupsen/logrus"

type DefaultFieldHook struct {
	GetValue func() string
}

func (h *DefaultFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *DefaultFieldHook) Fire(e *logrus.Entry) error {
	e.Data["application"] = h.GetValue()
	return nil
}
