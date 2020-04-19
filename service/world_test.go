package service

import (
	"context"
	"fmt"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.core/structs"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	log = logger.Init("test logger", true, false, ioutil.Discard)
	log.Info("test logger")
	if path, err := filepath.Abs("../config"); err != nil {
		log.Fatal(err)
	} else {
		viper.AddConfigPath(path)
		viper.SetConfigType("yaml")

		viper.SetConfigName(".world.circleci")

		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
	initRedis()
	os.Exit(m.Run())
}

func TestWorldTime(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if rnc, err := worldTime(ctx); err != nil {
		t.Error(err)
	} else {
		if reflect.TypeOf(rnc) != reflect.TypeOf(structs.NcMiscGameTimeAck{}) {
			t.Errorf("expected nc struct of type: %v but instead got %v", reflect.TypeOf(rnc).String(), reflect.TypeOf(structs.NcMiscGameTimeAck{}).String())
		}
	}
}

func TestReturnToServerSelect(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if rnc, err := returnToServerSelect(ctx); err != nil {
		t.Error(err)
	} else {
		if reflect.TypeOf(rnc) != reflect.TypeOf(structs.NcUserWillWorldSelectAck{}) {
			t.Errorf("expected nc struct of type: %v but instead got %v", reflect.TypeOf(rnc).String(), reflect.TypeOf(structs.NcUserWillWorldSelectAck{}).String())
		}
	}
}
