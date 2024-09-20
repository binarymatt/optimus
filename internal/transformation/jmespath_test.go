package transformation

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/shoenig/test/must"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/binarymatt/optimus/internal/testutil"
)

var simpleJmesData = `
{
  "locations": [
    {"name": "Seattle", "state": "WA"},
    {"name": "New York", "state": "NY"},
    {"name": "Bellevue", "state": "WA"},
    {"name": "Olympia", "state": "WA"}
  ]
}`

func TestJmesTransformer_HappyPath(t *testing.T) {
	jt := &JmesTransformer{
		Expression: `locations[?state == 'WA'].name | sort(@) | {cities: join(', ', @)}`,
	}
	must.NoError(t, jt.Initialize())
	var inputMap map[string]any
	must.NoError(t, json.Unmarshal([]byte(simpleJmesData), &inputMap))
	inputStruct, err := structpb.NewStruct(inputMap)
	must.NoError(t, err)
	expected := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"cities": structpb.NewStringValue("Bellevue, Olympia, Seattle"),
		},
	}
	outS, err := jt.Transform(context.Background(), inputStruct)
	must.NoError(t, err)
	must.Eq(t, expected, outS, testutil.CmpTransform)

}
func TestJmespath_NonStructOutput(t *testing.T) {
	jt := &JmesTransformer{
		Expression: "a",
	}
	must.NoError(t, jt.Initialize())
	inputStruct, err := structpb.NewStruct(map[string]any{
		"a": "foo", "b": "bar", "c": "baz",
	})
	must.NoError(t, err)
	out, err := jt.Transform(context.Background(), inputStruct)
	must.ErrorIs(t, err, ErrNotAMap)
	must.Nil(t, out)
}
