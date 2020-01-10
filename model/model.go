package model

import (
	"os"
	"time"

	"bitbucket.org/auzty/gobackup/archive"
	"bitbucket.org/auzty/gobackup/compressor"
	"bitbucket.org/auzty/gobackup/config"
	"bitbucket.org/auzty/gobackup/database"
	"bitbucket.org/auzty/gobackup/encryptor"
	"bitbucket.org/auzty/gobackup/logger"
	"bitbucket.org/auzty/gobackup/notification"
	"bitbucket.org/auzty/gobackup/storage"
)

// Model class
type Model struct {
	Config config.ModelConfig
}
type Report struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  string
}

// Perform model
func (ctx Model) Perform() {
	var laporan Report
	//laporan.StartTime, _ = time.Parse("2 Jan 2006 15:04", time.Now().String())
	laporan.StartTime = time.Now()
	logger.Info("##====== " + ctx.Config.Name + " ========")
	logger.Info("WorkDir:", ctx.Config.DumpPath+"\n")

	err := database.Run(ctx.Config)
	if err != nil {
		logger.Error(err)
		return
	}

	if ctx.Config.Archive != nil {
		err = archive.Run(ctx.Config)
		if err != nil {
			logger.Error(err)
			return
		}
	}

	archivePath, err := compressor.Run(ctx.Config)
	if err != nil {
		logger.Error(err)
		return
	}

	archivePath, err = encryptor.Run(archivePath, ctx.Config)
	if err != nil {
		logger.Error(err)
		return
	}

	err = storage.Run(ctx.Config, archivePath)
	if err != nil {
		logger.Error(err)
		return
	}

	defer ctx.cleanup(archivePath, laporan)

}

// Cleanup model temp files
func (ctx Model) cleanup(archivePath string, laporan Report) {
	logger.Info("Cleanup temp dir:" + config.TempPath + "...\n")
	err := os.RemoveAll(config.TempPath)
	if err != nil {
		logger.Error("Cleanup temp dir "+config.TempPath+" error:", err)
	}
	logger.Info("======= End " + ctx.Config.Name + " =======\n\n")

	format := "2 Jan 2006 15:04:05 MST"
	//	laporan.EndTime, err = time.Parse("2 Jan 2006 15:04", time.Since(laporan.StartTime))
	tduration := time.Since(laporan.StartTime)
	laporan.Duration = tduration.String()
	laporan.EndTime = laporan.StartTime.Add(tduration)
	//	laporan.Duration = laporan.EndTime.Sub(laporan.StartTime).String()
	logger.Info("start : ", laporan.StartTime.Format(format))
	logger.Info("end : ", laporan.EndTime.Format(format))
	logger.Info("durasi : ", laporan.Duration)
	_, err = notification.Run(ctx.Config, archivePath)
	if err != nil {
		logger.Error(err)
		//		return
	}

}
