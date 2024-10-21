package transformation

import (
	"context"
	"errors"
	"log/slog"

	"github.com/blues/jsonata-go"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	ErrNotAMap = errors.New("interface{} is not a struct")
)

type JsonataTransformer struct {
	Expression string `hcl:"expression"`
	expr       *jsonata.Expr
}

func (jt *JsonataTransformer) Initialize() (err error) {
	jt.expr, err = jsonata.Compile(jt.Expression)
	return
}
func (jt *JsonataTransformer) Transform(ctx context.Context, data *structpb.Struct) (*structpb.Struct, error) {
	slog.Warn("starting jsonata transformation")
	raw, err := jt.expr.Eval(data.AsMap())
	if err != nil {
		if errors.Is(err, jsonata.ErrUndefined) {
			return nil, nil
		}
		return nil, err
	}
	newM, ok := raw.(map[string]any)
	if !ok {
		return nil, ErrNotAMap
	}
	slog.Warn("passing back new struct")
	return structpb.NewStruct(newM)
}
