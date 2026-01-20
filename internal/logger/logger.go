package logger

import (
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogger() {
	lumberjackLogger := &lumberjack.Logger{
		Filename: "./log/app.log",
		MaxSize: 10,
		MaxBackups: 3,
		MaxAge: 15,
		Compress: true,
	}

	multiWriter := io.MultiWriter(os.Stdout, lumberjackLogger)

	log.SetOutput(multiWriter)
}