package utils

import (
	"errors"

	"google.golang.org/protobuf/proto"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

var ErrNilEvent = errors.New("nil event")

func deepCopy(src *optimusv1.LogEvent) (*optimusv1.LogEvent, error) {
	var dst optimusv1.LogEvent
	if src == nil {
		return nil, ErrNilEvent
	}
	bytes, err := proto.Marshal(src)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(bytes, &dst)
	return &dst, err
}
func CopyLogEvent(event *optimusv1.LogEvent) (*optimusv1.LogEvent, error) {
	return deepCopy(event)
}
