package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// globalSugar 是我们就为了方便调用而持有的全局对象
	globalSugar *zap.SugaredLogger
)

func init() {
	// 1. 初始化 (通常在 main 开头只做一次)
	InitLogger(&Config{
		Level:    "debug",
		Filename: "./logs/app.log",
		Console:  true,
	})
}

type Config struct {
	Level      string // debug, info, warn, error
	Filename   string // 日志文件路径
	MaxSize    int    // 单个文件最大尺寸 (MB)
	MaxBackups int    // 保留旧文件最大数量
	MaxAge     int    // 保留旧文件最大天数
	Compress   bool   // 是否压缩
	Console    bool   // 是否输出到控制台
}

func InitLogger(cfg *Config) {
	// 1. 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	atomicLevel.SetLevel(level)

	// 2. 配置 Encoder (格式化)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") // 人类可读时间

	// 控制台专用的 Encoder (带颜色)
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 开启颜色
	consoleEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	var cores []zapcore.Core

	// --- Core A: 文件输出 (JSON, 无色) ---
	if cfg.Filename != "" {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		})
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig), // 文件里用 JSON
			fileWriter,
			atomicLevel,
		))
	}

	// --- Core B: 控制台输出 (Text, 有色) ---
	if cfg.Console {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewConsoleEncoder(consoleEncoderConfig), // 控制台用 Text + Color
			zapcore.Lock(os.Stdout),
			atomicLevel,
		))
	}

	core := zapcore.NewTee(cores...)

	// 3. 创建 Logger
	// AddCaller: 显示行号
	// AddCallerSkip(1): 关键！因为我们封装了一层函数，所以要跳过1层调用栈，
	// 否则日志显示的行号永远是 logger.go 里的位置，而不是你 main.go 里的位置。
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// 4. 获取 SugaredLogger 并赋值给全局变量
	globalSugar = logger.Sugar()
}

// -------------------------------------------------------
// 以下是包级别的快捷函数，支持 Printf 风格的占位符
// -------------------------------------------------------

func Debug(template string, args ...interface{}) {
	globalSugar.Debugf(template, args...)
}

func Info(template string, args ...interface{}) {
	globalSugar.Infof(template, args...)
}

func Warn(template string, args ...interface{}) {
	globalSugar.Warnf(template, args...)
}

// Error 支持占位符，例如: logger.Error("用户 %s id %d 查询失败", name, id)
func Error(template string, args ...interface{}) {
	globalSugar.Errorf(template, args...)
}

// Fatal 会打印日志并调用 os.Exit(1)
func Fatal(template string, args ...interface{}) {
	globalSugar.Fatalf(template, args...)
}
