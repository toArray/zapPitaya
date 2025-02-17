package zapPitaya

import (
	"context"
	"fmt"
	"github.com/topfreegames/pitaya/v2"
	"github.com/topfreegames/pitaya/v2/constants"
	"github.com/topfreegames/pitaya/v2/logger/interfaces"
	"go.uber.org/zap"
)

// ZapCtx 定义用户区分数据的管理
// 如果key是zapCtx则在WithField和WithFields时取出pitaya的上下文数据并提取
const ZapCtx = "zap-ctx" // 特殊上下文标识

// CustomSessionData 自定义数据格式
type CustomSessionData map[string]any

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
		if k == ZapCtx {
			ctx, ok := v.(context.Context)
			if !ok {
				continue
			}

			// 获得所有需要记录的自定义日志列
			res := z.getFieldsList(ctx)
			if len(res) > 0 {
				fieldsList = append(fieldsList, res...)
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
	if key == ZapCtx {
		ctx, ok := value.(context.Context)
		if !ok {
			return z
		}

		// 获得所有需要记录的自定义日志列
		res := z.getFieldsList(ctx)
		newZap.logger = z.logger.With(res...)
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

// getFieldsList 获得日志列
func (z *ZapLog) getFieldsList(ctx context.Context) (fieldsList []zap.Field) {

	fieldsList = make([]zap.Field, 0)
	fieldsList = append(fieldsList, zap.Any(constants.PeerIDKey, pitaya.GetFromPropagateCtx(ctx, constants.PeerIDKey)))
	fieldsList = append(fieldsList, zap.Any(constants.RouteKey, pitaya.GetFromPropagateCtx(ctx, constants.RouteKey)))
	fieldsList = append(fieldsList, zap.Any(constants.PeerServiceKey, pitaya.GetFromPropagateCtx(ctx, constants.PeerServiceKey)))
	fieldsList = append(fieldsList, zap.Any(constants.RequestIDKey, pitaya.GetFromPropagateCtx(ctx, constants.RequestIDKey)))
	fieldsList = append(fieldsList, zap.Any(constants.StartTimeKey, pitaya.GetFromPropagateCtx(ctx, constants.StartTimeKey)))

	// 获取数据
	data := pitaya.GetFromPropagateCtx(ctx, ZapCustomSessionData)
	if data == nil {
		return
	}

	res, ok := data.(CustomSessionData)
	if !ok {
		return
	}

	// 自定义字段
	fieldsList = append(fieldsList, zap.Any(ZapCustomSessionData, res))
	return
}
