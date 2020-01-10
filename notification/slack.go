package notification

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"bitbucket.org/auzty/gobackup/logger"
)

// notification
//
// type: slack
type Slack struct {
	Base
	webHook string
}

func (ctx *Slack) perform(backupPath string) (archivePath string, err error) {
	// get total size of backup
	fopen, err := os.Stat(backupPath)
	if err != nil {
		logger.Error(err)
	}
	//	logger.Info(ctx.lapor, "@@@@@@@@@@@@@@@@@@@@")
	logger.Info("Backup Size (compressed) : ", getFileSizeReadable(fopen.Size()))

	format := "2 Jan 2006 15:04:05 MST"

	ctx.webHook = ctx.model.Notifications.Viper.GetString("webhook")
	if len(ctx.webHook) == 0 {
		logger.Warn("Webhook URL not found, please review your gobackup.yml")
	} else {
		// sending data to mattermost

		messagestruct := `
		{
		"attachments": [
			{
				"color": "#00ff00",
				"pretext": "##### Backup Report for test.example.com",
				"fields": [
					{
						"title": "Started",
						"value": "` + ctx.lapor.StartTime.Format(format) + `",
						"short": true
					},
					{
						"title": "Finished",
						"value": "` + ctx.lapor.EndTime.Format(format) + `",
						"short": true
					},
					{
						"title": "Duration",
						"value": "` + ctx.lapor.Duration + `",
						"short": false
					},
					{
						"title": "Backup Size",
						"value": "` + getFileSizeReadable(fopen.Size()) + ` (compressed)",
						"short": false
					}
				]
			}
		]
		}
		`

		var jsonStr = []byte(messagestruct)
		req, err := http.NewRequest("POST", ctx.webHook, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Error(err)
		}
		defer resp.Body.Close()
		if strings.Index(resp.Status, "200") < 0 {
			body, _ := ioutil.ReadAll(resp.Body)
			logger.Error(string(body))
		}
		logger.Info("response body : ", resp.Status)
	}

	return
}

func getFileSizeReadable(input int64) string {
	var kb = int64(1024)
	var mb = int64(kb * 1024)
	var gb = int64(mb * 1024)

	if input < mb {
		result := strconv.FormatInt(input/kb, 10)
		return result + " KB"
	}
	if (input >= mb) && (input <= gb) {
		result := strconv.FormatInt(input/mb, 10)
		return result + " MB"
	} else {
		result := strconv.FormatInt(input/gb, 10)
		return result + " GB"

	}
}
