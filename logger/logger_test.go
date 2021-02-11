package logger

import (
	"testing"
	"time"
)

func TestFatal(t *testing.T) {
	defer Fatal(500, "Fatal 1").
		Trace("SELECT * FROM `app` LIMIT 20").
		TimeCost(time.Now())
	defer Fatal(500, "Fatal 2").
		Trace("SELECT * FROM `user` LIMIT 20").
		Write()
	defer Fatal(500, "Fatal 3").Write()
	defer Fatal(500, "").Write()
}

func TestError(t *testing.T) {
	defer Error(500, "Error 1").Trace("SELECT * FROM `app` LIMIT 20").TimeCost(time.Now())
	defer Error(500, "Error 2").Trace("SELECT * FROM `user` LIMIT 20").Write()
	defer Error(500, "Error 3").Write()
	defer Error(500, "").Write()
}

func TestInfo(t *testing.T) {
	defer Info(RETRIEVE, "SELECT * FROM `user` LIMIT 20").TimeCost(time.Now())
	defer Info(UPDATE, "UPDATE `user` SET A=B").TimeCost(time.Now())
	defer Info(CREATE, "CREATE `user` SET A=B").TimeCost(time.Now())
	defer Info(DELETE, "DELETE `user` SET A=B").TimeCost(time.Now())
}

func TestDebug(t *testing.T) {
	defer Debug(RETRIEVE, "SELECT * FROM `user` LIMIT 20").TimeCost(time.Now())
	defer Debug(UPDATE, "UPDATE `user` SET A=B").TimeCost(time.Now())
	defer Debug(CREATE, "CREATE `user` SET A=B").TimeCost(time.Now())
	defer Debug(DELETE, "DELETE `user` SET A=B").TimeCost(time.Now())
}

func TestLevelFatal(t *testing.T) {
	DefaultLogger.SetLevel(LevelFatal)
	DefaultErrorLogger.SetLevel(LevelFatal)
	defer Debug(RETRIEVE, "SELECT * FROM `user` LIMIT 20").TimeCost(time.Now())
	defer Info(RETRIEVE, "SELECT * FROM `user` LIMIT 20").TimeCost(time.Now())
	defer Error(500, "Error 4").Write()
	defer Fatal(500, "Fatal 4").Write()
}

func TestLevelError(t *testing.T) {
	DefaultLogger.SetLevel(LevelError)
	DefaultErrorLogger.SetLevel(LevelError)
	defer Debug(RETRIEVE, "SELECT * FROM `user` LIMIT 20").TimeCost(time.Now())
	defer Info(RETRIEVE, "SELECT * FROM `user` LIMIT 20").TimeCost(time.Now())
	defer Error(500, "Error 5").Write()
	defer Fatal(500, "Fatal 5").Write()
}

func TestLevelInfo(t *testing.T) {
	DefaultLogger.SetLevel(LevelInfo)
	defer Debug(RETRIEVE, "SELECT * FROM `user` LIMIT 20").TimeCost(time.Now())
	defer Info(RETRIEVE, "SELECT * FROM `user` LIMIT 20").TimeCost(time.Now())
	defer Error(500, "Error 6").Write()
	defer Fatal(500, "Fatal 6").Write()
}

func TestLevelDebug(t *testing.T) {
	DefaultLogger.SetLevel(LevelDebug)
	defer Debug(RETRIEVE, "SELECT * FROM `user` LIMIT 20").TimeCost(time.Now())
	defer Info(RETRIEVE, "SELECT * FROM `user` LIMIT 20").TimeCost(time.Now())
	defer Error(500, "Error 7").Write()
	defer Fatal(500, "Fatal 7").Write()
}
