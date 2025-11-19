package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
} 


func NewLogger(logLevel string) (*zap.SugaredLogger, func() error, error ){
	lvl := zap.NewAtomicLevel()
	
	if err := lvl.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, nil, fmt.Errorf("unmarshal log level: %w", err)
	}
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, nil, fmt.Errorf("mkdir log folder: %w", err)
	}
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	
	logFilePath := filepath.Join("logs", fmt.Sprintf("%s.log", timestamp))

	LogFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return  nil, nil, fmt.Errorf("open log file: %w",err)
	}
	
	
	cfg := zap.NewDevelopmentEncoderConfig()
	
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04-05.000000")

	encoder := zapcore.NewConsoleEncoder(cfg)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout),lvl),
		zapcore.NewCore(encoder, zapcore.AddSync(LogFile), lvl),
	)

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	return logger.Sugar(), LogFile.Close, nil
}
