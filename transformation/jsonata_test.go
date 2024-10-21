package transformation

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/shoenig/test/must"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/binarymatt/optimus/internal/testutil"
)

var simpleData = `
{
  "title": "test order",
  "orders": [
    {"price": 10, "quantity": 3},
    {"price": 0.5, "quantity": 10},
    {"price": 100, "quantity": 1}
  ]
}`
var simpleExpression = `
{
  "name": title,
  "total": $sum(orders.price),
  "count": $count(orders)
}`

func TestJsonNataTransformer_HappyPath(t *testing.T) {
	ctx := context.Background()
	expected := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"name":  structpb.NewStringValue("test order"),
			"total": structpb.NewNumberValue(110.5),
			"count": structpb.NewNumberValue(3),
		},
	}
	var i map[string]any
	err := json.Unmarshal([]byte(simpleData), &i)
	must.NoError(t, err)
	s, err := structpb.NewStruct(i)
	must.NoError(t, err)
	jt := &JsonataTransformer{
		Expression: simpleExpression,
	}
	must.NoError(t, jt.Initialize())
	newS, err := jt.Transform(ctx, s)
	must.NoError(t, err)
	must.Eq(t, expected, newS, testutil.CmpTransform)
}

func TestJsonNataTransfomer_NonStructOutput(t *testing.T) {
	ctx := context.Background()
	var i map[string]any
	err := json.Unmarshal([]byte(simpleData), &i)
	must.NoError(t, err)
	s, err := structpb.NewStruct(i)
	must.NoError(t, err)
	jt := &JsonataTransformer{
		Expression: `$sum(orders.(price * quantity))`,
	}
	must.NoError(t, jt.Initialize())

	newS, err := jt.Transform(ctx, s)
	must.ErrorIs(t, ErrNotAMap, err)
	must.Nil(t, newS)
}

func TestJsonNataTransformer_BadExpression(t *testing.T) {
	jt := &JsonataTransformer{
		Expression: `.Orders`,
	}
	must.Error(t, jt.Initialize())
}
func TestJsonNataTransformer_EmptyResult(t *testing.T) {

	ctx := context.Background()
	var i map[string]any
	err := json.Unmarshal([]byte(simpleData), &i)
	must.NoError(t, err)
	s, err := structpb.NewStruct(i)
	must.NoError(t, err)
	jt := &JsonataTransformer{
		Expression: `Orders`,
	}
	must.NoError(t, jt.Initialize())
	newS, err := jt.Transform(ctx, s)
	must.NoError(t, err)
	must.Nil(t, newS)
	// must.Eq(t, &structpb.Struct{}, newS, testutil.CmpTransform)

}
