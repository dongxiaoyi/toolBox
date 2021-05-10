package configs

import (
	"fmt"
	"github.com/dongxiaoyi/toolBox/pkg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"time"
)

var Logger *zap.SugaredLogger

// 日志器
func LogLevel() map[string]zapcore.Level {
	level := make(map[string]zapcore.Level)
	level["debug"] = zap.DebugLevel
	level["info"] = zap.InfoLevel
	level["warn"] = zap.WarnLevel
	level["error"] = zap.ErrorLevel
	level["dpanic"] = zap.DPanicLevel
	level["panic"] = zap.PanicLevel
	level["fatal"] = zap.FatalLevel
	return level
}

// 初始化日志
func NewLogger(isRotate, rotateConsole, disableStacktrace, disableCeller, isConsole bool, encodingType string) *zap.SugaredLogger {
	curPath := pkg.AbsPath()

	logLevelOpt := "DEBUG" // 日志级别
	levelMap := LogLevel()
	logLevel, _ := levelMap[logLevelOpt]
	atomicLevel := zap.NewAtomicLevelAt(logLevel)

	encodingConfig := zapcore.EncoderConfig{
		TimeKey: "Time",
		LevelKey: "Level",
		NameKey: "Log",
		CallerKey: "Celler",
		MessageKey: "Message",
		StacktraceKey: "Stacktrace",
		LineEnding: zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("[2006-01-02 15:04:05]"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller: zapcore.FullCallerEncoder,
	}

	// 初始化字段
	filedKey := "service"
	fieldValue := "toolBox"
	filed := zap.Fields(zap.String(filedKey, fieldValue))

	outputPath := []string{
		"stdout",
		//path.Join(curPath, "check.log"),
	}
	errorPath := []string{
		"stderr",
		//path.Join(curPath, "check.log"),
	}

	// 是否开启日志滚动
	if isRotate {
		rotateFile := path.Join(curPath, "check.log")
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   rotateFile,
			MaxSize:    50,  // 默认是100MB
			MaxBackups: 14,  // 保存14份
			MaxAge:     30,  // 保留30天
			Compress:   true, // 日志压缩
		})
		var core zapcore.Core
		// 是否在日志切割开启的情况下开启console打印
		if rotateConsole {
			consoleDebugging := zapcore.Lock(os.Stdout)
			// 多个core使用NewTee
			if isConsole {
				core = zapcore.NewTee(
					zapcore.NewCore(
						zapcore.NewJSONEncoder(encodingConfig),
						writer,
						logLevel,
					),
					zapcore.NewCore(
						zapcore.NewJSONEncoder(encodingConfig),
						consoleDebugging,
						logLevel,
					),
				)
			} else {
				core = zapcore.NewTee(
					zapcore.NewCore(
						zapcore.NewJSONEncoder(encodingConfig),
						writer,
						logLevel,
					),
				)
			}

		} else {
			core = zapcore.NewCore(
				zapcore.NewJSONEncoder(encodingConfig),
				writer,
				logLevel,
				)
		}

		if !disableCeller {
			caller := zap.AddCaller()

			if !disableStacktrace {
				stacktrace := zap.AddStacktrace(logLevel)
				Logger = zap.New(core, caller, stacktrace, filed).Sugar()
			} else {
				Logger = zap.New(core, caller, filed).Sugar()
			}
		} else {
			if !disableStacktrace {
				stacktrace := zap.AddStacktrace(logLevel)
				Logger = zap.New(core, stacktrace, filed).Sugar()
			} else {
				Logger = zap.New(core, filed).Sugar()
			}
		}
	} else {
		// TODO: 传参确定日志器的参数（非配置文件：下级节点不保留配置文件）
		logCfg := zap.Config{
			Level: atomicLevel,
			Development: true,
			DisableCaller: disableCeller,
			DisableStacktrace: disableStacktrace,
			Encoding: encodingType,
			EncoderConfig: encodingConfig, // console or json
			InitialFields: map[string]interface{}{filedKey: fieldValue},
			OutputPaths: outputPath,
			ErrorOutputPaths: errorPath,
		}

		logger, err := logCfg.Build()
		if err != nil {
			panic(fmt.Sprintf("Loggger初始化失败: %v", err))
		}

		Logger = logger.Sugar()
	}

	return Logger
}

