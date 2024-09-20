package optimus

import (
	"context"
	"fmt"
	"testing"

	"github.com/shoenig/test/must"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/binarymatt/optimus/config"
	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/testutil"
)

var yamlStr = `---
inputs:
  httpInput:
    kind: http
filters:
  bFilter:
    kind: bexpr
    expression: action == "create"
    subscriptions:
      - httpInput
      - channelInput
transformations:
  jsonata:
    kind: jsonata
    expression: '{"user_email":principal.email,"path":path}'
    subscriptions:
      - bFilter
`

var inputs = []map[string]any{
	{
		"id":     "1",
		"topic":  "audit",
		"action": "delete",
		"system": "docs",
		"path":   "/docs/test_doc",
		"verb":   "DELETE",
		"principal": map[string]any{
			"id":    "1",
			"name":  "test user",
			"email": "test@example.com",
		},
	},
	{
		"id":     "2",
		"topic":  "audit",
		"system": "iam",
		"action": "create",
		"verb":   "POST",
		"path":   "/users/create",
		"principal": map[string]any{
			"id":    "2",
			"name":  "test user",
			"email": "test2@example.com",
		},
	},
	{
		"id":     "3",
		"topic":  "audit",
		"action": "create",
		"system": "docs",
		"path":   "/docs",
		"verb":   "POST",
		"principal": map[string]any{
			"id":    "3",
			"name":  "test user",
			"email": "test3@example.com",
		},
	},
}

func TestFlow(t *testing.T) {
	inputStruct1, err := structpb.NewStruct(inputs[0])
	must.NoError(t, err)
	inputStruct2, err := structpb.NewStruct(inputs[1])
	must.NoError(t, err)
	inputStruct3, err := structpb.NewStruct(inputs[2])
	must.NoError(t, err)
	inputEvents := []*optimusv1.LogEvent{
		{
			Id:   "1",
			Data: inputStruct1,
		},
		{
			Id:   "2",
			Data: inputStruct2,
		},
		{
			Id:   "3",
			Data: inputStruct3,
		},
	}
	outputEvents := []*optimusv1.LogEvent{
		{
			Id: "2",
			Data: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"path":       structpb.NewStringValue("/users/create"),
					"user_email": structpb.NewStringValue("test2@example.com"),
				},
			},
			Upstreams: []string{"channelInput", "bFilter", "jsonata"},
		},
		{
			Id: "3",
			Data: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"path":       structpb.NewStringValue("/docs"),
					"user_email": structpb.NewStringValue("test3@example.com"),
				},
			},
			Upstreams: []string{"channelInput", "bFilter", "jsonata"},
		},
	}
	inputChannel := make(chan *optimusv1.LogEvent)
	outputChannel := make(chan *optimusv1.LogEvent)
	cfg, err := config.NewWithYaml([]byte(yamlStr),
		config.WithChannelInput("channelInput", inputChannel),
		config.WithChannelOutput("channelOutput", outputChannel, []string{"jsonata"}),
	)
	must.NoError(t, err)
	o, err := New(cfg)
	must.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	// main routine
	eg.Go(func() error {
		return o.Run(ctx)
	})

	// input
	eg.Go(func() error {
		t.Logf("sending input events")
		inputChannel <- inputEvents[0]
		inputChannel <- inputEvents[1]
		inputChannel <- inputEvents[2]
		return nil
	})
	// output tests
	eg.Go(func() error {
		t.Logf("retrieving output events")
		defer cancel()
		first := <-outputChannel
		second := <-outputChannel

		t.Logf("validating output events")
		must.Eq(t, outputEvents[0], first, testutil.CmpTransform)
		must.Eq(t, outputEvents[1], second, testutil.CmpTransform)
		select {
		case msg := <-outputChannel:
			fmt.Println("received message", msg)
		default:
			fmt.Println("no message received")
			// cancel()
		}
		return nil
	})
	must.NoError(t, eg.Wait())
}
