package model

import (
	"os"
	"time"

	"github.com/auzty/gobackup/archive"
	"github.com/auzty/gobackup/compressor"
	"github.com/auzty/gobackup/config"
	"github.com/auzty/gobackup/database"
	"github.com/auzty/gobackup/encryptor"
	"github.com/auzty/gobackup/logger"
	"github.com/auzty/gobackup/notification"
	"github.com/auzty/gobackup/storage"
)

// Model class
type Model struct {
	Config config.ModelConfig
}

/*
type Report struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  string
}
*/

// Perform model
func (ctx Model) Perform() {
	var laporan notification.Report
	archiveLocation := ""
	//laporan.StartTime, _ = time.Parse("2 Jan 2006 15:04", time.Now().String())
	laporan.StartTime = time.Now()
	logger.Info("====== " + ctx.Config.Name + " ========")
	logger.Info("WorkDir:", ctx.Config.DumpPath+"\n")

	err := database.Run(ctx.Config)
	if err != nil {
		logger.Error(err)
		//return
	}

	if ctx.Config.Archive != nil {
		err = archive.Run(ctx.Config)
		if err != nil {
			logger.Error(err)
			laporan.BackupStatus = "error"
			laporan.MessageString = "Archiving Error \n" + laporan.MessageString + err.Error()
			//return
		} else {
			laporan.BackupStatus = "ok"
		}
		archivePath, err := compressor.Run(ctx.Config)
		if err != nil {
			logger.Error(err)
			laporan.BackupStatus = "error"
			laporan.MessageString = "Compressing Error \n" + laporan.MessageString + err.Error()
			//return
		} else {
			laporan.BackupStatus = "ok"
		}

		archivePath, err = encryptor.Run(archivePath, ctx.Config)
		if err != nil {
			logger.Error(err)
			laporan.BackupStatus = "error"
			laporan.MessageString = "Encrypting Error \n" + laporan.MessageString + err.Error()
			//return
		} else {
			laporan.BackupStatus = "ok"
		}

		err = storage.Run(ctx.Config, archivePath)
		if err != nil {
			logger.Error(err)
			laporan.BackupStatus = "error"
			laporan.MessageString = "Storing using " + ctx.Config.StoreWith.Type + " Error \n```\n" + laporan.MessageString + err.Error() + "\n```"
			//return
		} else {
			laporan.BackupStatus = "ok"
		}
		logger.Info("####################\n", laporan.MessageString)
		archiveLocation = archivePath

	}

	defer ctx.cleanup(archiveLocation, laporan)

}

// Cleanup model temp files
func (ctx Model) cleanup(archivePath string, laporan notification.Report) {
	logger.Info("Cleanup temp dir:" + config.TempPath + "...\n")
	err := os.RemoveAll(config.TempPath)
	if err != nil {
		logger.Error("Cleanup temp dir "+config.TempPath+" error:", err)
	}
	logger.Info("======= End " + ctx.Config.Name + " =======\n\n")

	//	format := "2 Jan 2006 15:04:05 MST"
	tduration := time.Since(laporan.StartTime)
	laporan.Duration = tduration.String()
	laporan.EndTime = laporan.StartTime.Add(tduration)
	//	logger.Info("start : ", laporan.StartTime.Format(format))
	//	logger.Info("end : ", laporan.EndTime.Format(format))
	//	logger.Info("durasi : ", laporan.Duration)

	_, err = notification.Run(ctx.Config, archivePath, laporan)
	if err != nil {
		logger.Error(err)
		//		return
	}

}
