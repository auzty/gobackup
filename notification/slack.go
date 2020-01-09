package notification

import (
	"github.com/huacnlee/gobackup/helper"
	"github.com/huacnlee/gobackup/logger"
)

// Tgz .tar.gz compressor
//
// type: tgz
type Slack struct {
	Base
}

func (ctx *Slack) perform() (archivePath string, err error) {

	logger.Info("slack cuy")

	return
}

func (ctx *Slack) options() (opts []string) {
	if helper.IsGnuTar {
		opts = append(opts, "--ignore-failed-read")
	}
	opts = append(opts, "-zcf")

	return
}
