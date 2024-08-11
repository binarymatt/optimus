package filter

import (
	"context"
	"testing"

	"github.com/shoenig/test/must"
	"google.golang.org/protobuf/types/known/structpb"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

func TestBexprSetup_Err(t *testing.T) {
	bf := &BexprFilter{}
	err := bf.Setup()
	must.Error(t, err)
	must.Nil(t, bf.evaluator)
}

func TestBexprSetup(t *testing.T) {
	bf := &BexprFilter{
		Expression: "foo.bar == true",
	}
	err := bf.Setup()
	must.NoError(t, err)
	must.NotNil(t, bf.evaluator)
}

func TestBexprProcess(t *testing.T) {
	bf := &BexprFilter{
		Expression: "foo.bar == true",
	}
	must.NoError(t, bf.Setup())
	data, err := structpb.NewStruct(map[string]interface{}{
		"foo": map[string]any{
			"bar": true,
		},
	})
	must.NoError(t, err)

	event := &optimusv1.LogEvent{Data: data}
	ev, err := bf.Process(context.Background(), event)
	must.NoError(t, err)
	must.Eq(t, event, ev)

	bf2 := BexprFilter{
		Expression: "foo.bar == false",
	}
	must.NoError(t, bf2.Setup())
	ev2, err := bf2.Process(context.Background(), event)
	must.NoError(t, err)
	must.Nil(t, ev2)

	bf3 := BexprFilter{
		Expression: "foo == test",
	}
	must.NoError(t, bf3.Setup())
	ev3, err := bf3.Process(context.Background(), event)
	must.Error(t, err)
	must.Nil(t, ev3)
}
