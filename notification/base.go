package notification

import (
	"fmt"

	"github.com/huacnlee/gobackup/config"
	"github.com/huacnlee/gobackup/logger"
	"github.com/spf13/viper"
)

// Base compressor
type Base struct {
	name  string
	model config.ModelConfig
	viper *viper.Viper
}

// Context compressor
type Context interface {
	perform()
}

func newBase(model config.ModelConfig) (base Base) {
	fmt.Println("heheheh")
	base = Base{
		name:  model.Name,
		model: model,
		viper: model.Notifications.Viper,
	}
	return
}

// Run compressor
func Run(model config.ModelConfig) (archivePath string, err error) {
	//base := newBase(model)
	var ctx Context
	ctx.perform()

	/*
		var ctx Context
		switch model.Name.Type {
		case "slack":
			ctx = &Slack{Base: base}
		default:
			ctx = &Slack{}
		}
	*/

	logger.Info("------------ Notification -------------")

	logger.Info("------------ -------------\n")

	return
}
