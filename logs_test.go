package zapPitaya

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"testing"
)

func TestGetPitayaLogger(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	m := make(map[string]interface{})
	m["test"] = "test"
	m["num"] = 100

	log := GetPitayaLogger(logger)

	log.Debug("this is Debug", "21212")
	log.Debugf("this is Debugf, %s, %d, %+v", "string", 1000, m)
	log.Debugln("this is Debugln")

	log.Error(fmt.Sprintf("122222"))
	log.Errorf(" %s, %d, %v", "string", 1000, m)
	log.Errorln("this is Errorln")

	log.Info("this is Info")
	log.Infof("this is Infof, %s, %d, %v", "string", 1000, m)
	log.Infoln("this is Infoln")

	log.Warn("this is Warn")
	log.Warnf("this is Warnf, %s, %d, %v", "string", 1000, m)
	log.Warnln("this is Warnln")

	log.Panic("this is Panic")
	log.Panicf("this is Panicf, %s, %d, %v", "string", 1000, m)
	log.Panicln("this is Panicln")

	log.Fatal("this is Fatal")
	log.Fatalf("this is Fatalf, %s, %d, %v", "string", 1000, m)
	log.Fatalln("this is Fatalln")
}

func TestWithFile(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	m := make(map[string]interface{})
	m["test"] = "test"
	m["num"] = 100
	
	log := GetPitayaLogger(logger)
	log = log.WithField("test", "test")
	log = log.WithError(fmt.Errorf("test err"))
	log = log.WithField(ZapCtx, context.TODO())
	log = log.WithField("err", fmt.Errorf("test err2"))
	log.Error("this is Info")
	log.Errorf("this is Infof, %s, %d, %+v", "string", 1000, m)
	log.Errorln("this is Infoln ")
}
