package notification

import (
	"bitbucket.org/auzty/gobackup/config"
	"bitbucket.org/auzty/gobackup/logger"
	"github.com/spf13/viper"
)

// Base notification
type Base struct {
	name  string
	model config.ModelConfig
	viper *viper.Viper
}

// Context notification
type Context interface {
	perform(backupPath string) (archivePath string, err error)
}

func newBase(model config.ModelConfig) (base Base) {
	base = Base{
		name:  model.Name,
		model: model,
		viper: model.Notifications.Viper,
	}
	return
}

// Run notification
func Run(model config.ModelConfig, backupPath string) (archivePath string, err error) {
	base := newBase(model)

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
