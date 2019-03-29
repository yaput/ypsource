package dialogflow

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	df "cloud.google.com/go/dialogflow/apiv2"
	"google.golang.org/api/option"
)

func TestClientInit(t *testing.T) {
	type typeDummyContextsClient func(ctx context.Context, opts ...option.ClientOption) (*df.ContextsClient, error)
	type typeDummySessionsClient func(ctx context.Context, opts ...option.ClientOption) (*df.SessionsClient, error)
	type args struct {
		projectID    string
		authFilePath string
		ctx          context.Context
	}
	tests := []struct {
		name                string
		args                args
		dummyContextsClient typeDummyContextsClient
		dummySessionsClient typeDummySessionsClient
		want                *DFClient
	}{
		{
			name: "Test Init Success",
			args: args{
				authFilePath: "",
				projectID:    "123-321",
			},
			want: &DFClient{
				projectID:        "123-321",
				authJSONFilePath: "",
			},
			dummyContextsClient: func(ctx context.Context, opts ...option.ClientOption) (*df.ContextsClient, error) {
				return nil, nil
			},
			dummySessionsClient: func(ctx context.Context, opts ...option.ClientOption) (*df.SessionsClient, error) {
				return nil, nil
			},
		},
		{
			name: "Test Failed Init Contexts",
			args: args{
				authFilePath: "",
				projectID:    "123-321",
			},
			want: nil,
			dummyContextsClient: func(ctx context.Context, opts ...option.ClientOption) (*df.ContextsClient, error) {
				return nil, fmt.Errorf("Error Test")
			},
			dummySessionsClient: func(ctx context.Context, opts ...option.ClientOption) (*df.SessionsClient, error) {
				return nil, nil
			},
		},
		{
			name: "Test Failed Init Contexts",
			args: args{
				authFilePath: "",
				projectID:    "123-321",
			},
			want: nil,
			dummySessionsClient: func(ctx context.Context, opts ...option.ClientOption) (*df.SessionsClient, error) {
				return nil, fmt.Errorf("Error Test")
			},
		},
	}
	for _, tt := range tests {
		dfNewContextsClient = tt.dummyContextsClient
		dfNewSessionClient = tt.dummySessionsClient
		t.Run(tt.name, func(t *testing.T) {
			if got := ClientInit(tt.args.projectID, tt.args.authFilePath, tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientInit() = %v, want %v", got, tt.want)
			}
		})
	}
}
