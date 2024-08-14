package zapPitaya

import (
	"context"
	"github.com/topfreegames/pitaya/v2"
	"github.com/topfreegames/pitaya/v2/constants"
	"github.com/topfreegames/pitaya/v2/logger/interfaces"
	"go.uber.org/zap"
)

// CTX 定义用户区分数据的管理
// 如果key是zapCtx则在WithField和WithFields时取出pitaya的上下文数据并提取
const CTX = "zapPitayaCtx"

// InitZapLogger 初始化zap
// 这里只实现pitaya接口
func SetZapLog(logger *zap.Logger) interfaces.Logger {
	return &ZapLog{logger: logger}
}

// ZapLog zap log
type ZapLog struct {
	logger *zap.Logger
}

func (z *ZapLog) Fatal(format ...interface{}) {
	z.logger.Sugar().Fatal(format)
}
func (z *ZapLog) Fatalf(format string, args ...interface{}) {
	z.logger.Sugar().Fatal(format)
}
func (z *ZapLog) Fatalln(args ...interface{}) {
	z.logger.Sugar().Fatalln(args)
}

func (z *ZapLog) Debug(args ...interface{}) {
	z.logger.Sugar().Debug(args)
}
func (z *ZapLog) Debugf(format string, args ...interface{}) {
	z.logger.Sugar().Debugf(format, args)
}
func (z *ZapLog) Debugln(args ...interface{}) {
	z.logger.Sugar().Debugln(args)
}

func (z *ZapLog) Error(args ...interface{}) {
	z.logger.Sugar().Error(args)
}
func (z *ZapLog) Errorf(format string, args ...interface{}) {
	z.logger.Sugar().Errorf(format, args)
}
func (z *ZapLog) Errorln(args ...interface{}) {
	z.logger.Sugar().Errorln(args)
}

func (z *ZapLog) Info(args ...interface{}) {
	z.logger.Sugar().Info(args)
}
func (z *ZapLog) Infof(format string, args ...interface{}) {
	z.logger.Sugar().Infof(format, args)
}
func (z *ZapLog) Infoln(args ...interface{}) {
	z.logger.Sugar().Infoln(args)
}

func (z *ZapLog) Warn(args ...interface{}) {
	z.logger.Sugar().Warn(args)
}
func (z *ZapLog) Warnf(format string, args ...interface{}) {
	z.logger.Sugar().Warnf(format, args)
}
func (z *ZapLog) Warnln(args ...interface{}) {
	z.logger.Sugar().Warnln(args)
}

func (z *ZapLog) Panic(args ...interface{}) {
	z.logger.Sugar().Panic(args)
}
func (z *ZapLog) Panicf(format string, args ...interface{}) {
	z.logger.Sugar().Panicf(format, args)
}
func (z *ZapLog) Panicln(args ...interface{}) {
	z.logger.Sugar().Panicln(args)
}

func (z *ZapLog) WithFields(fields map[string]interface{}) interfaces.Logger {

	return z
}
func (z *ZapLog) WithField(key string, value interface{}) interfaces.Logger {

	if key == CTX {
		ctx, ok := value.(context.Context)
		if !ok {
			return z
		}

		uid := pitaya.GetSessionFromCtx(ctx).UID()
		route := pitaya.GetFromPropagateCtx(ctx, constants.RouteKey)
		service := pitaya.GetFromPropagateCtx(ctx, constants.PeerServiceKey)
		requestID := pitaya.GetFromPropagateCtx(ctx, constants.RequestIDKey)
		logger := z.logger.With(zap.Any("route", route),
			zap.Any("service", service),
			zap.Any("requestID", requestID),
			zap.Any("uid", uid))
		return &ZapLog{
			logger: logger,
		}
	}

	logger := z.logger.With(zap.Any(key, value))
	return &ZapLog{
		logger: logger,
	}
}
func (z *ZapLog) WithError(err error) interfaces.Logger {
	return z
}

func (z *ZapLog) GetInternalLogger() any {
	return z.logger
}
