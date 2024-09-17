package transformation

import (
	"context"
	"errors"

	"github.com/blues/jsonata-go"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	ErrNotAMap = errors.New("interface{} is not a struct")
)

type JsonataTransformer struct {
	Expression string `yaml:"expression"`
}

func (jt *JsonataTransformer) Transform(ctx context.Context, data *structpb.Struct) (*structpb.Struct, error) {
	expr, err := jsonata.Compile(jt.Expression)
	if err != nil {
		return nil, err
	}
	raw, err := expr.Eval(data.AsMap())
	if err != nil {
		if errors.Is(err, jsonata.ErrUndefined) {
			return &structpb.Struct{}, nil
		}
		return nil, err
	}
	newM, ok := raw.(map[string]any)
	if !ok {
		return nil, ErrNotAMap
	}
	return structpb.NewStruct(newM)
}
