package logger

type RollingInterval uint8

const (
	RollingIntervalDay RollingInterval = 0 + iota
	RollingIntervalWeek
	RollingIntervalMonth
)

type config struct {
	limitLevel     Level
	writeToConsole writeToConsole
	writeToFile    WriteToFile
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
	}
}

func (config *config) SetLimitLevel(level Level) *config {
	config.limitLevel = level
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

	if option.RollingInterval <= RollingIntervalWeek && config.writeToFile.RollingInterval != option.RollingInterval {
		config.writeToFile.RollingInterval = option.RollingInterval
	}

	if option.RollingFileSize > 0 && config.writeToFile.RollingFileSize != option.RollingFileSize {
		config.writeToFile.RollingFileSize = option.RollingFileSize
	}

	return config
}

func (config *config) CreateLogger() ILogger {
	return CreateLogger(*config)
}
