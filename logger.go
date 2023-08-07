package ginz

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func LoadLogger(level string) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	lv, err := logrus.ParseLevel(level)
	if err != nil {
		log.Panic("InvalidLogLevel", err)
	}
	logrus.SetLevel(lv)

	// 显示文件和位置
	logrus.SetReportCaller(true)
	// 设置为json格式
	// logrus.SetFormatter(&logrus.JSONFormatter{
	// 	TimestampFormat: "2006-01-02 15:04:05",
	// })
	// 设置为自定义格式
	logrus.SetFormatter(&MyFormatter{})

	// logfile, _ := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)

	// 日志分割
	rotateLogger := &lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    10, // mb
		MaxBackups: 20,
		MaxAge:     30,   //days
		Compress:   true, // disabled by default
	}
	logrus.SetOutput(io.MultiWriter(os.Stdout, rotateLogger)) // 控制台+文件
	// logrus.SetOutput(logfile) // 只写到文件

	// 不需要颜色
	// gin.DisableConsoleColor()
	// gin的日志写到控制台和日志文件中, 默认只有os.Stdout
	// gin.DefaultWriter = io.MultiWriter(os.Stdout, rotateLogger)
}

type MyFormatter struct {
}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog, fileVal string

	if entry.HasCaller() {
		fileVal = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
	}

	newLog = fmt.Sprintf("%s || %s || %s || %s\n", timestamp, entry.Level, fileVal, entry.Message)

	b.WriteString(newLog)
	return b.Bytes(), nil
}
