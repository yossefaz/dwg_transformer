package utils

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/yossefazoulay/go_utils/queue"
	globalUtils "github.com/yossefazoulay/go_utils/utils"
	"os"
	"os/exec"
	"strings"
	"transformer/config"
)


func HandleError(err error, msg string, exit bool) {
	if err != nil {
		config.Logger.Log.Error(fmt.Sprintf("%s: %s", msg, err))
	}
	if exit {
		os.Exit(1)
	}
}

func MessageReceiver(m amqp.Delivery, rmq queue.Rabbitmq)  {
	resultConfig := getResultConfig()
	log := config.Logger.Log
	pFIle := &globalUtils.PickFile{}
	err := json.Unmarshal(m.Body, pFIle)
	HandleError(err, "Error decoding message", false)
	if pFIle.From !=  resultConfig.From{
		if err := m.Ack(false); err != nil {
			log.Error("Error acknowledging message : %s", err)
		} else {
			res:= execute(pFIle, config.LocalConfig.OutputFormat)
			mess, err1 := rmq.SendMessage(res, resultConfig.Success)
			HandleError(err1, "message sending error", false)
			config.Logger.Log.Info(mess)

		}
	}
}

func getResultConfig() globalUtils.Result {
	return config.LocalConfig.Queue.Rabbitmq.Result
}

func execute(pfile *globalUtils.PickFile, output string) []byte{
	resultConfig := getResultConfig()
	outpath, convert := getOutputPath(pfile.Path, output)
	if convert {
		cmd := exec.Command("dwgread", pfile.Path, "-O", output, "-o", outpath)
		err := cmd.Run()
		if err != nil {
			return setResult(pfile, pfile.Path, resultConfig.From, resultConfig.Fail, true)
		}
		return setResult(pfile, outpath, resultConfig.From, resultConfig.Success, false)
	}
	return setResult(pfile, outpath, resultConfig.From, resultConfig.Success,  true)
}

func setResult(pfile *globalUtils.PickFile, path string, from string, to string, error bool)[]byte {
	execRes := 1
	if error {
		execRes = 0
	}
	keys := make([]string, 0, len(pfile.Result))
	for k := range pfile.Result {
		keys = append(keys, k)
	}
	mess, err := globalUtils.SetResultMessage(pfile, keys, []int {execRes},  from, to, path)
	if err != nil {
		HandleError(err, "Cannot set output and cannot run command :" + err.Error() + err.Error(), false)
	}
	return mess
	}


func getOutputPath(basePath string, output string) (string, bool) {
	fileExt := config.LocalConfig.FileExtensions[output]
	currentPath := strings.Split(basePath, ".")
	if !(currentPath[1] == "dwg" || currentPath[1] == "dxf") {
		return currentPath[0], false
	}
	outpath := currentPath[0] + fileExt
	return outpath, true
}
