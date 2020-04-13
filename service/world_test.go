package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/structs"
	"reflect"
	"testing"
)

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
