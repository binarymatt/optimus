package transformation

import (
	"context"

	"github.com/jmespath/go-jmespath"
	"google.golang.org/protobuf/types/known/structpb"
)

type JmesTransformer struct {
	Expression string `yaml:"expression"`
	path       *jmespath.JMESPath
}

func (jt *JmesTransformer) Initialize() (err error) {
	jt.path, err = jmespath.Compile(jt.Expression)
	return
}
func (jt *JmesTransformer) Transform(ctx context.Context, data *structpb.Struct) (*structpb.Struct, error) {
	result, err := jt.path.Search(data.AsMap())
	if err != nil {
		return nil, err
	}
	newM, ok := result.(map[string]any)
	if !ok {
		return nil, ErrNotAMap
	}
	return structpb.NewStruct(newM)
}
