package obs

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// Level defines the level of the log
type Level int

const (
	LEVEL_OFF   Level = 500
	LEVEL_ERROR Level = 400
	LEVEL_WARN  Level = 300
	LEVEL_INFO  Level = 200
	LEVEL_DEBUG Level = 100
)

var logLevelMap = map[Level]string{
	LEVEL_OFF:   "[OFF]: ",
	LEVEL_ERROR: "[ERROR]: ",
	LEVEL_WARN:  "[WARN]: ",
	LEVEL_INFO:  "[INFO]: ",
	LEVEL_DEBUG: "[DEBUG]: ",
}

type logConfType struct {
	level        Level
	logToConsole bool
	logFullPath  string
	maxLogSize   int64
	backups      int
}

func getDefaultLogConf() logConfType {
	return logConfType{
		level:        LEVEL_WARN,
		logToConsole: false,
		logFullPath:  "",
		maxLogSize:   1024 * 1024 * 30, // 30MB
		backups:      10,
	}
}

var logConf logConfType

type loggerWrapper struct {
	fullPath   string
	fd         *os.File
	ch         chan string
	wg         sync.WaitGroup
	queue      []string
	logger     *log.Logger
	index      int
	cacheCount int
	closed     bool
}

func (lw *loggerWrapper) doInit() {
	lw.queue = make([]string, 0, lw.cacheCount)
	lw.logger = log.New(lw.fd, "", 0)
	lw.ch = make(chan string, lw.cacheCount)
	lw.wg.Add(1)
	go lw.doWrite()
}

func (lw *loggerWrapper) rotate() {
	stat, err := lw.fd.Stat()
	if err != nil {
		_err := lw.fd.Close()
		if _err != nil {
			doLog(LEVEL_WARN, "Failed to close file with reason: %v", _err)
		}
		panic(err)
	}
	if stat.Size() >= logConf.maxLogSize {
		_err := lw.fd.Sync()
		if _err != nil {
			panic(err)
		}
		_err = lw.fd.Close()
		if _err != nil {
			doLog(LEVEL_WARN, "Failed to close file with reason: %v", _err)
		}
		if lw.index > logConf.backups {
			lw.index = 1
		}
		_err = os.Rename(lw.fullPath, lw.fullPath+"."+IntToString(lw.index))
		if _err != nil {
			panic(err)
		}
		lw.index++

		fd, err := os.OpenFile(lw.fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			panic(err)
		}
		lw.fd = fd
		lw.logger.SetOutput(lw.fd)
	}
}

func (lw *loggerWrapper) doFlush() {
	lw.rotate()
	for _, m := range lw.queue {
		lw.logger.Println(m)
	}
	err := lw.fd.Sync()
	if err != nil {
		panic(err)
	}
}

func (lw *loggerWrapper) doClose() {
	lw.closed = true
	close(lw.ch)
	lw.wg.Wait()
}

func (lw *loggerWrapper) doWrite() {
	defer lw.wg.Done()
	for {
		msg, ok := <-lw.ch
		if !ok {
			lw.doFlush()
			_err := lw.fd.Close()
			if _err != nil {
				doLog(LEVEL_WARN, "Failed to close file with reason: %v", _err)
			}
			break
		}
		if len(lw.queue) >= lw.cacheCount {
			lw.doFlush()
			lw.queue = make([]string, 0, lw.cacheCount)
		}
		lw.queue = append(lw.queue, msg)
	}

}

func (lw *loggerWrapper) Printf(format string, v ...interface{}) {
	if !lw.closed {
		msg := fmt.Sprintf(format, v...)
		lw.ch <- msg
	}
}

var consoleLogger *log.Logger
var fileLogger *loggerWrapper
var lock = new(sync.RWMutex)

func isDebugLogEnabled() bool {
	return logConf.level <= LEVEL_DEBUG
}

func isErrorLogEnabled() bool {
	return logConf.level <= LEVEL_ERROR
}

func isWarnLogEnabled() bool {
	return logConf.level <= LEVEL_WARN
}

func isInfoLogEnabled() bool {
	return logConf.level <= LEVEL_INFO
}

func reset() {
	if fileLogger != nil {
		fileLogger.doClose()
		fileLogger = nil
	}
	consoleLogger = nil
	logConf = getDefaultLogConf()
}

// InitLog enable logging function with default cacheCnt
func InitLog(logFullPath string, maxLogSize int64, backups int, level Level, logToConsole bool) error {
	return InitLogWithCacheCnt(logFullPath, maxLogSize, backups, level, logToConsole, 50)
}

// InitLogWithCacheCnt enable logging function
func InitLogWithCacheCnt(logFullPath string, maxLogSize int64, backups int, level Level, logToConsole bool, cacheCnt int) error {
	lock.Lock()
	defer lock.Unlock()
	if cacheCnt <= 0 {
		cacheCnt = 50
	}
	reset()
	if fullPath := strings.TrimSpace(logFullPath); fullPath != "" {
		_fullPath, err := filepath.Abs(fullPath)
		if err != nil {
			return err
		}

		if !strings.HasSuffix(_fullPath, ".log") {
			_fullPath += ".log"
		}

		stat, fd, err := initLogFile(_fullPath)
		if err != nil {
			return err
		}

		prefix := stat.Name() + "."
		index := 1
		var timeIndex int64 = 0
		walkFunc := func(path string, info os.FileInfo, err error) error {
			if err == nil {
				if name := info.Name(); strings.HasPrefix(name, prefix) {
					if i := StringToInt(name[len(prefix):], 0); i >= index && info.ModTime().Unix() >= timeIndex {
						timeIndex = info.ModTime().Unix()
						index = i + 1
					}
				}
			}
			return err
		}

		if err = filepath.Walk(filepath.Dir(_fullPath), walkFunc); err != nil {
			_err := fd.Close()
			if _err != nil {
				doLog(LEVEL_WARN, "Failed to close file with reason: %v", _err)
			}
			return err
		}

		fileLogger = &loggerWrapper{
			fullPath: _fullPath, fd: fd, index: index, cacheCount: cacheCnt, closed: false,
		}
		fileLogger.doInit()
	}
	if maxLogSize > 0 {
		logConf.maxLogSize = maxLogSize
	}
	if backups > 0 {
		logConf.backups = backups
	}
	logConf.level = level
	if logToConsole {
		consoleLogger = log.New(os.Stdout, "", log.LstdFlags)
	}
	return nil
}

func initLogFile(_fullPath string) (os.FileInfo, *os.File, error) {
	stat, err := os.Stat(_fullPath)
	if err == nil && stat.IsDir() {
		return nil, nil, fmt.Errorf("logFullPath:[%s] is a directory", _fullPath)
	} else if err = os.MkdirAll(filepath.Dir(_fullPath), os.ModePerm); err != nil {
		return nil, nil, err
	}

	fd, err := os.OpenFile(_fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		_err := fd.Close()
		if _err != nil {
			doLog(LEVEL_WARN, "Failed to close file with reason: %v", _err)
		}
		return nil, nil, err
	}

	if stat == nil {
		stat, err = os.Stat(_fullPath)
		if err != nil {
			_err := fd.Close()
			if _err != nil {
				doLog(LEVEL_WARN, "Failed to close file with reason: %v", _err)
			}
			return nil, nil, err
		}
	}

	return stat, fd, nil
}

// CloseLog disable logging and synchronize cache data to log files
func CloseLog() {
	if logEnabled() {
		lock.Lock()
		defer lock.Unlock()
		reset()
	}
}

func logEnabled() bool {
	return consoleLogger != nil || fileLogger != nil
}

// DoLog writes log messages to the logger
func DoLog(level Level, format string, v ...interface{}) {
	doLog(level, format, v...)
}

func doLog(level Level, format string, v ...interface{}) {
	if logEnabled() && logConf.level <= level {
		msg := fmt.Sprintf(format, v...)
		if _, file, line, ok := runtime.Caller(1); ok {
			index := strings.LastIndex(file, "/")
			if index >= 0 {
				file = file[index+1:]
			}
			msg = fmt.Sprintf("%s:%d|%s", file, line, msg)
		}
		prefix := logLevelMap[level]
		defer func() {
			_ = recover()
			// ignore ch closed error
		}()
		if consoleLogger != nil {
			consoleLogger.Printf("%s%s", prefix, msg)
		}
		if fileLogger != nil {
			nowDate := FormatUtcNow("2006-01-02T15:04:05Z")
			fileLogger.Printf("%s %s%s", nowDate, prefix, msg)
		}
	}
}

func checkAndLogErr(err error, level Level, format string, v ...interface{}) {
	if err != nil {
		doLog(level, format, v...)
	}
}

func logResponseHeader(respHeader http.Header) string {
	resp := make([]string, 0, len(respHeader)+1)
	for key, value := range respHeader {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if strings.HasPrefix(key, HEADER_PREFIX) || strings.HasPrefix(key, HEADER_PREFIX_OBS) {
			key = key[len(HEADER_PREFIX):]
		}
		_key := strings.ToLower(key)
		if _, ok := allowedLogResponseHTTPHeaderNames[_key]; ok {
			resp = append(resp, fmt.Sprintf("%s: [%s]", key, value[0]))
		}
		if _key == HEADER_REQUEST_ID {
			resp = append(resp, fmt.Sprintf("%s: [%s]", key, value[0]))
		}
	}
	return strings.Join(resp, " ")
}

func logRequestHeader(reqHeader http.Header) string {
	resp := make([]string, 0, len(reqHeader)+1)
	for key, value := range reqHeader {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		_key := strings.ToLower(key)
		if _, ok := allowedRequestHTTPHeaderMetadataNames[_key]; ok {
			resp = append(resp, fmt.Sprintf("%s: [%s]", key, value[0]))
		}
	}
	return strings.Join(resp, " ")
}
