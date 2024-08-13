package utils

import (
	"google.golang.org/protobuf/proto"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

func DeepCopy(src *optimusv1.LogEvent) (*optimusv1.LogEvent, error) {
	var dst optimusv1.LogEvent
	bytes, err := proto.Marshal(src)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(bytes, &dst)
	return &dst, err
}
func CopyLogEvent(event *optimusv1.LogEvent) *optimusv1.LogEvent {
	dst, err := DeepCopy(event)
	if err != nil {
		return nil
	}
	return dst
}
