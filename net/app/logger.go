package app

import (
	"os"

	"github.com/yz778899/vGate/net/app/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

func InitLogger(cfg *config.RootConfig) error {
	// 1. 解析日志级别
	var level zapcore.Level
	switch cfg.Logger.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 2. 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("01-02 15:04:05.000"), // yy-mm-dd 01-02 15:04:05
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 3. 选择编码格式
	var encoder zapcore.Encoder
	if cfg.Logger.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 4. 构建输出目标
	var writers []zapcore.WriteSyncer
	for _, path := range cfg.Logger.OutputPaths {
		if path == "stdout" {
			writers = append(writers, zapcore.AddSync(os.Stdout))
		} else if path == "stderr" {
			writers = append(writers, zapcore.AddSync(os.Stderr))
		} else {
			// 文件输出，使用 lumberjack 切割
			writers = append(writers, zapcore.AddSync(&lumberjack.Logger{
				Filename:   path,
				MaxSize:    cfg.Logger.Lumberjack.MaxSize,
				MaxBackups: cfg.Logger.Lumberjack.MaxBackups,
				MaxAge:     cfg.Logger.Lumberjack.MaxAge,
				Compress:   cfg.Logger.Lumberjack.Compress,
			}))
		}
	}

	// 5. 构建错误输出目标
	var errorWriters []zapcore.WriteSyncer
	for _, path := range cfg.Logger.ErrorOutputPaths {
		if path == "stderr" {
			errorWriters = append(errorWriters, zapcore.AddSync(os.Stderr))
		} else {
			errorWriters = append(errorWriters, zapcore.AddSync(&lumberjack.Logger{
				Filename:   path,
				MaxSize:    cfg.Logger.Lumberjack.MaxSize,
				MaxBackups: cfg.Logger.Lumberjack.MaxBackups,
				MaxAge:     cfg.Logger.Lumberjack.MaxAge,
				Compress:   cfg.Logger.Lumberjack.Compress,
			}))
		}
	}

	// 6. 创建 Core
	multiWriter := zapcore.NewMultiWriteSyncer(writers...)
	multiErrorWriter := zapcore.NewMultiWriteSyncer(errorWriters...)

	core := zapcore.NewCore(encoder, multiWriter, level)

	// 7. 错误级别单独输出
	if len(errorWriters) > 0 {
		errorCore := zapcore.NewCore(encoder, multiErrorWriter, zapcore.ErrorLevel)
		core = zapcore.NewTee(core, errorCore)
	}

	// 8. 创建 Logger
	opts := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(cfg.Logger.CallerSkip),
	}

	Log = zap.New(core, opts...)

	Log.Info("日志系统初始化成功",
		zap.String("level", cfg.Logger.Level),
		zap.String("encoding", cfg.Logger.Encoding),
	)

	return nil
}
