package notification

import (
	"time"

	"github.com/auzty/gobackup/config"
	"github.com/auzty/gobackup/logger"
	"github.com/spf13/viper"
)

// Base notification
type Base struct {
	name  string
	model config.ModelConfig
	viper *viper.Viper
	lapor Report
}

type Report struct {
	StartTime     time.Time
	EndTime       time.Time
	Duration      string
	MessageString string
	BackupStatus  string
}

// Context notification
type Context interface {
	perform(backupPath string) (archivePath string, err error)
}

func newBase(model config.ModelConfig, inputlaporan Report) (base Base, laporan Report) {
	base = Base{
		name:  model.Name,
		model: model,
		viper: model.Notifications.Viper,
		lapor: inputlaporan,
	}
	return
}

// Run notification
func Run(model config.ModelConfig, backupPath string, laporan Report) (archivePath string, err error) {
	base, laporan := newBase(model, laporan)

	//	logger.Info(model.Notifications, "######")
	logger.Info("------------ Notification -------------")
	var ctx Context
	switch model.Notifications.Type {
	case "slack":
		ctx = &Slack{Base: base}
		ctx.perform(backupPath)
	default:
		logger.Info("No Notification Set")
	}
	logger.Info("------------ -------------\n")

	return
}
