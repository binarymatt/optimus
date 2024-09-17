package transformation

import (
	"context"

	"github.com/jmespath/go-jmespath"
	"google.golang.org/protobuf/types/known/structpb"
)

type JmesTransformer struct {
	Expression string `yaml:"expression"`
}

func (jt *JmesTransformer) Transform(ctx context.Context, data *structpb.Struct) (*structpb.Struct, error) {
	result, err := jmespath.Search(jt.Expression, data.AsMap())
	if err != nil {
		return nil, err
	}
	newM, ok := result.(map[string]any)
	if !ok {
		return nil, ErrNotAMap
	}
	return structpb.NewStruct(newM)
}
