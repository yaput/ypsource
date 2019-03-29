package ypsource

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/genproto/protobuf/field_mask"

	structpb "github.com/golang/protobuf/ptypes/struct"

	"google.golang.org/api/option"

	df "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type DFClient struct {
	projectID        string
	authJSONFilePath string
	sessionClient    *df.SessionsClient
	contextClient    *df.ContextsClient
	ctx              context.Context
}

type Agent struct {
	projectID string
	// languageCode  string
	sessionClient *df.SessionsClient
	contextClient *df.ContextsClient
	ctx           context.Context
	parent        string
}

var dfNewSessionClient = df.NewSessionsClient
var dfNewContextsClient = df.NewContextsClient

func ClientInit(projectID, authFilePath string, ctx context.Context) *DFClient {
	sessionClient, err := dfNewSessionClient(ctx, option.WithCredentialsFile(authFilePath))
	if err != nil {
		log.Printf("%+v", err)
		return nil
	}

	contextClient, err := dfNewContextsClient(ctx, option.WithCredentialsFile(authFilePath))
	if err != nil {
		return nil
	}

	return &DFClient{
		projectID:        projectID,
		authJSONFilePath: authFilePath,
		sessionClient:    sessionClient,
		contextClient:    contextClient,
		ctx:              ctx,
	}
}

func (c *DFClient) NewAgent() *Agent {
	return &Agent{
		sessionClient: c.sessionClient,
		contextClient: c.contextClient,
		projectID:     c.projectID,
		ctx:           c.ctx,
		parent:        fmt.Sprintf("projects/%s/agent/sessions/", c.projectID),
	}
}

func (a *Agent) QueryText(text, languageCode, sessionID string) *dialogflowpb.QueryResult {
	request := dialogflowpb.DetectIntentRequest{
		Session: a.parent + sessionID,
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Text{
				Text: &dialogflowpb.TextInput{
					LanguageCode: languageCode,
					Text:         text,
				},
			},
		},
	}

	resp, err := a.sessionClient.DetectIntent(a.ctx, &request)
	if err != nil {
		log.Printf("%+v", err)
		return nil
	}

	return resp.GetQueryResult()
}

func (a *Agent) InvokeEvent(text string, param *structpb.Struct, languageCode, sessionID string) *dialogflowpb.QueryResult {
	request := dialogflowpb.DetectIntentRequest{
		Session: a.parent + sessionID,
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Event{
				Event: &dialogflowpb.EventInput{
					LanguageCode: languageCode,
					Name:         text,
					Parameters:   param,
				},
			},
		},
	}

	resp, err := a.sessionClient.DetectIntent(a.ctx, &request)
	if err != nil {
		log.Printf("%+v", err)
		return nil
	}

	return resp.GetQueryResult()
}

func (a *Agent) CreateContext(name, sessionID string, param *structpb.Struct, ttl int32) *dialogflowpb.Context {
	ctx, err := a.contextClient.CreateContext(a.ctx, &dialogflowpb.CreateContextRequest{
		Parent: a.parent + sessionID,
		Context: &dialogflowpb.Context{
			LifespanCount: ttl,
			Name:          a.parent + sessionID + "/contexts/" + name,
			Parameters:    param,
		},
	})
	if err != nil {
		log.Printf("%+v", err)
		return nil
	}

	return ctx
}

func (a *Agent) GetListContexts(sessionID string) {
	listCtx := a.contextClient.ListContexts(a.ctx, &dialogflowpb.ListContextsRequest{
		Parent: a.parent + sessionID,
	})

	log.Printf("%+v", listCtx)
	return
}

func (a *Agent) GetContext(name, sessionID string) *dialogflowpb.Context {
	context, err := a.contextClient.GetContext(a.ctx, &dialogflowpb.GetContextRequest{
		Name: a.parent + sessionID + "/contexts/" + name,
	})
	if err != nil {
		log.Printf("%+v", err)
		return nil
	}
	return context
}

func (a *Agent) UpdateContext(inputContext *dialogflowpb.Context, updatedField *field_mask.FieldMask) *dialogflowpb.Context {
	context, err := a.contextClient.UpdateContext(a.ctx, &dialogflowpb.UpdateContextRequest{
		Context:    inputContext,
		UpdateMask: updatedField,
	})
	if err != nil {
		log.Printf("[Error] Update context got: %+v", err)
		return nil
	}
	return context
}
