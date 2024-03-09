package logger

import (
	"context"
	"gonetlib/util"
	"path/filepath"
	"strings"
	"time"
)

type RollingInterval uint8

const (
	RollingIntervalInfinite RollingInterval = 0 + iota
	RollingIntervalYear
	RollingIntervalMonth
	RollingIntervalDay
	RollingIntervalHour
	RollingIntervalMinute
)

type config struct {
	limitLevel     Level
	writeToConsole writeToConsole
	writeToFile    WriteToFile
	tickDuration   time.Duration
}

type writeTo struct {
	enable bool
}

type writeToConsole struct {
	writeTo
}

type WriteToFile struct {
	writeTo
	Filepath        string
	RollingInterval RollingInterval
	RollingFileSize uint32
}

func (wtf *WriteToFile) makeRollingFilepath() string {
	filename := util.GetFileNameWithoutExt(wtf.Filepath)
	ext := filepath.Ext(wtf.Filepath)
	dir := filepath.Dir(wtf.Filepath)

	var sb strings.Builder
	sb.WriteString(dir)
	sb.WriteString("/")
	sb.WriteString(filename)

	now := time.Now()

	switch wtf.RollingInterval {
	case RollingIntervalYear:
		sb.WriteString("_")
		sb.WriteString(now.Format("2006"))
		break
	case RollingIntervalMonth:
		sb.WriteString("_")
		sb.WriteString(now.Format("2006_01"))
		break
	case RollingIntervalDay:
		sb.WriteString("_")
		sb.WriteString(now.Format("2006_01_02"))
		break
	case RollingIntervalHour:
		sb.WriteString("_")
		sb.WriteString(now.Format("2006_01_02_15"))
		break
	case RollingIntervalMinute:
		sb.WriteString("_")
		sb.WriteString(now.Format("2006_01_02_1504"))
		break
	default:
		break
	}

	sb.WriteString(ext)

	return sb.String()
}

func CreateLoggerConfig() *config {
	return &config{
		limitLevel: DebugLevel,
		writeToConsole: writeToConsole{
			writeTo: writeTo{enable: false},
		},
		writeToFile: WriteToFile{
			writeTo:         writeTo{enable: false},
			Filepath:        "",
			RollingInterval: RollingIntervalDay,
			RollingFileSize: 1024 * 1024 * 10, //10mb
		},
		tickDuration: time.Second, //1초
	}
}

func (config *config) MinimumLevel(level Level) *config {
	config.limitLevel = level
	return config
}

func (config *config) TickDuration(ms time.Duration) *config {
	config.tickDuration = ms
	return config
}

func (config *config) WriteToConsole() *config {
	config.writeToConsole.enable = true
	return config
}

func (config *config) WriteToFile(option WriteToFile) *config {
	config.writeToFile.enable = true

	if option.Filepath != "" && option.Filepath != config.writeToFile.Filepath {
		config.writeToFile.Filepath = option.Filepath
	}

	if option.RollingInterval <= RollingIntervalMonth && config.writeToFile.RollingInterval != option.RollingInterval {
		config.writeToFile.RollingInterval = option.RollingInterval
	}

	if option.RollingFileSize > 0 && config.writeToFile.RollingFileSize != option.RollingFileSize {
		config.writeToFile.RollingFileSize = option.RollingFileSize
	}

	return config
}

func (config *config) CreateLogger(ctx context.Context) ILogger {
	return CreateLogger(*config, ctx)
}
