package testutil

import (
	"google.golang.org/protobuf/types/known/structpb"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

var eventData = map[string]interface{}{
	"Image": map[string]any{
		"Width":  800,
		"Height": 600,
		"Title":  "View from 15th Floor",
		"Thumbnail": map[string]any{
			"Url":    "http://www.example.com/image/481989943",
			"Height": 125,
			"Width":  100,
		},
		"Animated": false,
		"IDs":      []any{116, 943, 234, 38793},
	},
}

func BuildTestEvent() *optimusv1.LogEvent {
	data, err := structpb.NewStruct(eventData)
	if err != nil {
		return nil
	}
	return &optimusv1.LogEvent{
		Id:     "test",
		Source: "testing",
		Data:   data,
	}
}
