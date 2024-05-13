package logger

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/wj-dominic/gonetlib/util"
)

type RollingInterval uint8

const (
	RollingIntervalInvalid RollingInterval = 0 + iota
	RollingIntervalInfinite
	RollingIntervalYear
	RollingIntervalMonth
	RollingIntervalDay
	RollingIntervalHour
	RollingIntervalMinute
)

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
	RollingFileSize int64
}

func (wtf *WriteToFile) makeRollingFilepath() string {
	filename := util.GetFileNameWithoutExt(wtf.Filepath)
	ext := filepath.Ext(wtf.Filepath)
	dir := filepath.Dir(wtf.Filepath)

	var sb strings.Builder
	sb.WriteString(dir)
	sb.WriteString("/")

	var realFileName strings.Builder
	realFileName.WriteString(filename)

	//rolling interval for date
	now := time.Now()
	switch wtf.RollingInterval {
	case RollingIntervalYear:
		realFileName.WriteString("_")
		realFileName.WriteString(now.Format("2006"))
		break
	case RollingIntervalMonth:
		realFileName.WriteString("_")
		realFileName.WriteString(now.Format("2006_01"))
		break
	case RollingIntervalDay:
		realFileName.WriteString("_")
		realFileName.WriteString(now.Format("2006_01_02"))
		break
	case RollingIntervalHour:
		realFileName.WriteString("_")
		realFileName.WriteString(now.Format("2006_01_02_15"))
		break
	case RollingIntervalMinute:
		realFileName.WriteString("_")
		realFileName.WriteString(now.Format("2006_01_02_1504"))
		break
	default:
		break
	}

	rollingNumber := wtf.getRollingNumber(dir, realFileName.String())
	if rollingNumber > 0 {
		realFileName.WriteString("(")
		realFileName.WriteString(strconv.Itoa(int(rollingNumber)))
		realFileName.WriteString(")")
	}

	sb.WriteString(realFileName.String())
	sb.WriteString(ext)

	return sb.String()
}

func (wtf *WriteToFile) getRollingNumber(dir string, filename string) uint32 {
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	//마지막 파일 찾기
	lastRollingNumber := uint32(0)
	lastFileSize := int64(0)
	hasRollingFile := bool(false)
	for _, entry := range entries {
		if entry.IsDir() == true {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			panic(err)
		}

		if strings.Contains(info.Name(), filename) == false {
			continue
		}

		if strings.Contains(info.Name(), "(") == true {
			begin := strings.Index(info.Name(), "(")
			end := strings.Index(info.Name(), ")")

			numberStr := info.Name()[begin+1 : end]
			number, err := strconv.Atoi(numberStr)
			if err != nil {
				panic(err)
			}

			if lastRollingNumber > uint32(number) {
				continue
			}

			lastRollingNumber = uint32(number)
			hasRollingFile = true
			lastFileSize = info.Size()
		}

		//숫자가 없는 파일이 덮어 쓰는 것을 방지
		if hasRollingFile == false {
			lastFileSize = info.Size()
		}
	}

	if lastFileSize >= wtf.RollingFileSize {
		lastRollingNumber++
	}

	return lastRollingNumber
}

type config struct {
	limitLevel     Level
	writeToConsole writeToConsole
	writeToFile    WriteToFile
	tickDuration   time.Duration
}

func NewLoggerConfig() *config {
	return &config{
		limitLevel: DebugLevel,
		writeToConsole: writeToConsole{
			writeTo: writeTo{enable: true},
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

	if option.RollingInterval <= RollingIntervalMinute &&
		option.RollingInterval > RollingIntervalInvalid &&
		config.writeToFile.RollingInterval != option.RollingInterval {
		config.writeToFile.RollingInterval = option.RollingInterval
	}

	if option.RollingFileSize > 0 && config.writeToFile.RollingFileSize != option.RollingFileSize {
		config.writeToFile.RollingFileSize = option.RollingFileSize
	}

	return config
}
func (config *config) CreateLogger() Logger {
	return CreateLogger(*config)
}
