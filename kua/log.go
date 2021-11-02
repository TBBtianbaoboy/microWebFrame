package kua

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"time"
)


const (
    LevelTrace = iota//最高等级
    LevelDebug
    LevelInfo
    LevelWarning
    LevelError
    LevelCritical
)

var level = 1

type logInfo struct {
	Level string
	Time string
	Err   string
	Message string
}

// 获取日志等级
func GetLevel() int {
    return level
}

//手动修改日志等级
func SetLevel(l int) {
    level = l
}

//将数据转换为JSON格式
func logJson(le int,e error) []byte {
	SetLevel(ConfigStore.LogLevel)
	logInfoForJson := logInfo{
		Level: "",
		Time: time.Now().UTC().String(),
		Err: "",
		Message: "failed",
	}

	switch le {
		case LevelDebug:
			logInfoForJson.Err = debug(e)
			logInfoForJson.Level = "Debug"
		case LevelWarning:
			logInfoForJson.Err = warn(e)
			logInfoForJson.Level = "Warming"
		case LevelError:
			logInfoForJson.Err = rror(e)
			logInfoForJson.Level = "Error"
		case LevelCritical:
			logInfoForJson.Err = critical(e)
			logInfoForJson.Level = "Critical"
		default:
			logInfoForJson.Err = trace(e)
			logInfoForJson.Level = "Trace"
	}

	data,err := json.Marshal(logInfoForJson)
	if err != nil {
		log.Fatal(err)
	}
	if logInfoForJson.Err=="NO" {
		data = append(data, 'n')
	}
	return data
}


//写入日志
func writeLog(message []byte){
	file, err := os.OpenFile(ConfigStore.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
		log.Fatal(err)
    }
    defer file.Close()
    writer := bufio.NewWriter(file)
	_,err = writer.Write(message)
	if err != nil {
		log.Fatal(err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

//真正被外界调用的写日志函数
func WriteErrToLogFile(le int,e error) {
	data := logJson(le,e)
	if data[len(data)-1] != 'n' {
		writeLog(data)
	}
}

//以下五个为等级分控
func trace(e error) string {
    if level <= LevelTrace {
        return e.Error()
    }
	return "NO"
}

func debug(e error) string{
    if level <= LevelDebug {
        return e.Error()
    }
	return "NO"
}

func info(e error) string {
    if level <= LevelInfo {
		return e.Error()
    }
	return "NO"
}

func warn(e error) string{
    if level <= LevelWarning {
		return e.Error()
    }
	return "NO"
}

func rror(e error) string{
    if level <= LevelError {
	    return e.Error()
    }
	return "NO"
}

func critical(e error) string{
    if level <= LevelCritical {
        return e.Error()
    }
	return "NO"
}
