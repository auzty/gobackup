package notification

import (
	"time"

	"bitbucket.org/auzty/gobackup/config"
	"bitbucket.org/auzty/gobackup/logger"
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
	StartTime time.Time
	EndTime   time.Time
	Duration  string
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
	var ctx Context
	switch model.Notifications.Type {
	case "slack":
		ctx = &Slack{Base: base}
	default:
		logger.Info("error default")
	}

	logger.Info("------------ Notification -------------")
	ctx.perform(backupPath)

	logger.Info("------------ -------------\n")

	return
}
