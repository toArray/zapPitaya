package zapPitaya

import (
	"context"
	"fmt"
	"github.com/topfreegames/pitaya/v2"
	"github.com/topfreegames/pitaya/v2/constants"
	"github.com/topfreegames/pitaya/v2/logger/interfaces"
	"go.uber.org/zap"
)

// CTX 定义用户区分数据的管理
// 如果key是zapCtx则在WithField和WithFields时取出pitaya的上下文数据并提取
const CTX = "zapCtx"
const SESSION_UID = "session_uid"

// GetPitayaLogger 获得pitaya的logger
// 这里只实现pitaya接口
func GetPitayaLogger(zap *zap.Logger) interfaces.Logger {
	return &ZapLog{logger: zap}
}

// NewDefaultPitayaLogger 使用zap example 实现
func NewDefaultPitayaLogger() interfaces.Logger {
	return &ZapLog{logger: zap.NewExample()}
}

// ZapLog zap log
type ZapLog struct {
	logger *zap.Logger
}

func (z *ZapLog) Fatal(format ...interface{}) {
	z.logger.Sugar().Fatal(format)
}

func (z *ZapLog) Fatalf(format string, args ...interface{}) {
	z.logger.Sugar().Fatalf(fmt.Sprintf(format, args...))
}

func (z *ZapLog) Fatalln(args ...interface{}) {
	z.logger.Sugar().Fatalln(args)
}

func (z *ZapLog) Debug(args ...interface{}) {
	z.logger.Sugar().Debug(args)
}

func (z *ZapLog) Debugf(format string, args ...interface{}) {
	z.logger.Sugar().Debugf(fmt.Sprintf(format, args...))
}

func (z *ZapLog) Debugln(args ...interface{}) {
	z.logger.Sugar().Debugln(args)
}

func (z *ZapLog) Error(args ...interface{}) {
	z.logger.Sugar().Error(args)
}

func (z *ZapLog) Errorf(format string, args ...interface{}) {
	z.logger.Sugar().Errorf(fmt.Sprintf(format, args...))
}

func (z *ZapLog) Errorln(args ...interface{}) {
	z.logger.Sugar().Errorln(args)
}

func (z *ZapLog) Info(args ...interface{}) {
	z.logger.Sugar().Info(args)
}

func (z *ZapLog) Infof(format string, args ...interface{}) {
	z.logger.Sugar().Infof(fmt.Sprintf(format, args...))
}

func (z *ZapLog) Infoln(args ...interface{}) {
	z.logger.Sugar().Infoln(args)
}

func (z *ZapLog) Warn(args ...interface{}) {
	z.logger.Sugar().Warn(args)
}
func (z *ZapLog) Warnf(format string, args ...interface{}) {
	z.logger.Sugar().Warnf(fmt.Sprintf(format, args...))
}
func (z *ZapLog) Warnln(args ...interface{}) {
	z.logger.Sugar().Warnln(args)
}

func (z *ZapLog) Panic(args ...interface{}) {
	z.logger.Sugar().Panic(args)
}
func (z *ZapLog) Panicf(format string, args ...interface{}) {
	z.logger.Sugar().Panicf(fmt.Sprintf(format, args...))
}
func (z *ZapLog) Panicln(args ...interface{}) {
	z.logger.Sugar().Panicln(args)
}

// WithFields 日志添加列
// 如果是自定义上下文key则特殊处理
func (z *ZapLog) WithFields(fields map[string]interface{}) interfaces.Logger {
	fieldsList := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		if k == CTX {
			ctx, ok := v.(context.Context)
			if !ok {
				continue
			}

			// ready
			peerIDKey := pitaya.GetFromPropagateCtx(ctx, constants.PeerIDKey)
			routeKey := pitaya.GetFromPropagateCtx(ctx, constants.RouteKey)
			peerServiceKey := pitaya.GetFromPropagateCtx(ctx, constants.PeerServiceKey)
			requestIDKey := pitaya.GetFromPropagateCtx(ctx, constants.RequestIDKey)
			startTimeKey := pitaya.GetFromPropagateCtx(ctx, constants.StartTimeKey)

			// append
			fieldsList = append(fieldsList, zap.Any(constants.PeerIDKey, peerIDKey))
			fieldsList = append(fieldsList, zap.Any(constants.RouteKey, routeKey))
			fieldsList = append(fieldsList, zap.Any(constants.PeerServiceKey, peerServiceKey))
			fieldsList = append(fieldsList, zap.Any(constants.RequestIDKey, requestIDKey))
			fieldsList = append(fieldsList, zap.Any(constants.StartTimeKey, startTimeKey))

			// uid
			session := pitaya.GetSessionFromCtx(ctx)
			if session != nil {
				fieldsList = append(fieldsList, zap.Any(SESSION_UID, session.UID()))
			} else {
				fieldsList = append(fieldsList, zap.Any(SESSION_UID, pitaya.GetFromPropagateCtx(ctx, SESSION_UID)))
			}
		} else {
			fieldsList = append(fieldsList, zap.Any(k, v))
		}
	}

	logger := z.logger.With(fieldsList...)
	return &ZapLog{logger: logger}
}

// WithField 日志添加列
// 如果是自定义上下文key则特殊处理
func (z *ZapLog) WithField(key string, value interface{}) interfaces.Logger {

	newZap := &ZapLog{}
	if key == CTX {
		ctx, ok := value.(context.Context)
		if !ok {
			return z
		}

		// 上下文内数据
		peerIDKey := pitaya.GetFromPropagateCtx(ctx, constants.PeerIDKey)
		routeKey := pitaya.GetFromPropagateCtx(ctx, constants.RouteKey)
		peerServiceKey := pitaya.GetFromPropagateCtx(ctx, constants.PeerServiceKey)
		requestIDKey := pitaya.GetFromPropagateCtx(ctx, constants.RequestIDKey)
		startTimeKey := pitaya.GetFromPropagateCtx(ctx, constants.StartTimeKey)
		logger := z.logger.With(
			zap.Any(constants.PeerIDKey, peerIDKey),
			zap.Any(constants.RouteKey, routeKey),
			zap.Any(constants.PeerServiceKey, peerServiceKey),
			zap.Any(constants.RequestIDKey, requestIDKey),
			zap.Any(constants.StartTimeKey, startTimeKey),
		)

		// uid
		session := pitaya.GetSessionFromCtx(ctx)
		if session != nil {
			logger = logger.With(zap.Any(SESSION_UID, session.UID()))
		} else {
			logger = logger.With(zap.Any(SESSION_UID, pitaya.GetFromPropagateCtx(ctx, SESSION_UID)))
		}

		newZap.logger = logger
	} else {
		newZap.logger = z.logger.With(zap.Any(key, value))
	}

	return newZap
}

func (z *ZapLog) WithError(err error) interfaces.Logger {
	return &ZapLog{logger: z.logger.With(zap.Error(err))}
}

func (z *ZapLog) GetInternalLogger() any {
	return z.logger
}
