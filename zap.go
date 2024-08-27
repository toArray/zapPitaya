package zapPitaya

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/topfreegames/pitaya/v2/logger/interfaces"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"sync"
	"time"
)

// 共享的日志实例
var (
	log  *zap.Logger
	once sync.Once
)

// InitZap 初始化日志实例
func InitZap(vipConfig *viper.Viper) interfaces.Logger {
	once.Do(func() {

		level := vipConfig.GetString("zap.level")        // 日志级别
		outType := vipConfig.GetString("zap.outType")    // 日志格式
		filePath := vipConfig.GetString("zap.filePath")  // 输出路径
		filePath = strings.TrimSuffix(filePath, "/")     // 输出路径
		maxSize := vipConfig.GetInt("zap.maxSize")       // 每个日志文件的最大大小 (MB)
		maxBackups := vipConfig.GetInt("zap.maxBackups") // 保留旧日志文件的最大数量
		maxAge := vipConfig.GetInt("zap.maxAge")         // 旧日志文件的最大保存天数
		compress := vipConfig.GetBool("zap.compress")    // 是否压缩旧日志

		// 日志等级
		logLevel, err := zapcore.ParseLevel(level)
		if err != nil {
			panic(err)
		}

		// 区分日志类型
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderEncoder := zapcore.NewJSONEncoder(encoderConfig)
		var core zapcore.Core

		if outType == "file" {
			// 根据时间生成日志文件名、生成info文件日志
			infoFilename := fmt.Sprintf("%s/%s-info.log", filePath, time.Now().Format("2006-01-02"))
			infoWriter := &lumberjack.Logger{
				Filename:   infoFilename, // 日志文件名
				MaxSize:    maxSize,      // 每个日志文件的最大大小 (MB)
				MaxBackups: maxBackups,   // 保留旧日志文件的最大数量
				MaxAge:     maxAge,       // 旧日志文件的最大保存天数
				Compress:   compress,     // 是否压缩旧日志
			}

			// 根据时间生成日志文件名、生成error文件日志
			errorFilename := fmt.Sprintf("%s/%s-error.log", filePath, time.Now().Format("2006-01-02"))
			errorWriter := &lumberjack.Logger{
				Filename:   errorFilename, // 日志文件名
				MaxSize:    maxSize,       // 每个日志文件的最大大小 (MB)
				MaxBackups: maxBackups,    // 保留旧日志文件的最大数量
				MaxAge:     maxAge,        // 旧日志文件的最大保存天数
				Compress:   compress,      // 是否压缩旧日志
			}

			// 输出到文件
			core = zapcore.NewTee(
				zapcore.NewCore(encoderEncoder, zapcore.AddSync(infoWriter), zapcore.InfoLevel),
				zapcore.NewCore(encoderEncoder, zapcore.AddSync(errorWriter), zapcore.ErrorLevel))
		} else {
			// 输出到控制台
			core = zapcore.NewTee(zapcore.NewCore(encoderEncoder, zapcore.Lock(os.Stdout), logLevel))
		}

		// 初始化
		log = zap.New(core)
	})

	// 设置日志格式
	return GetPitayaLogger(log)
}

// GetLogger 获取共享的日志实例
func GetLogger() *zap.Logger {
	return log
}

// Close 关闭日志实例
func Close() {
	if log != nil {
		log.Sync()
	}
}
