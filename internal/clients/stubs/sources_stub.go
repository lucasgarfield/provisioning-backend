package stubs

import (
	"context"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http/sources"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
)

type sourcesCtxKeyType string

var sourcesCtxKey sourcesCtxKeyType = "sources-interface"

type SourcesIntegrationStub struct {
	store           *[]sources.Source
	authentications *[]sources.AuthenticationRead
}
type SourcesClientStub struct{}

func init() {
	// We are currently using SourcesClientStub
	clients.GetSourcesClient = getSourcesClientStub
}

// SourcesClient
func WithSourcesClient(parent context.Context) context.Context {
	ctx := context.WithValue(parent, sourcesCtxKey, &SourcesClientStub{})
	return ctx
}

func getSourcesClientStub(ctx context.Context) (si clients.Sources, err error) {
	var ok bool
	if si, ok = ctx.Value(sourcesCtxKey).(*SourcesClientStub); !ok {
		err = &contextReadError{}
	}
	return si, err
}

func (*SourcesClientStub) Ready(ctx context.Context) error {
	return nil
}

func (mock *SourcesClientStub) GetAuthentication(ctx context.Context, sourceId sources.ID) (*clients.Authentication, error) {
	return clients.NewAuthentication("arn:aws:iam::230214684733:role/Test", models.ProviderTypeAWS), nil
}

func (mock *SourcesClientStub) GetProvisioningTypeId(ctx context.Context) (string, error) {
	return "11", nil
}

func (mock *SourcesClientStub) ListProvisioningSources(ctx context.Context) ([]*clients.Source, error) {
	TestSourceData := []*clients.Source{
		{
			Id:           ptr.To("1"),
			Name:         ptr.To("source1"),
			SourceTypeId: ptr.To("1"),
			Uid:          ptr.To("5eebe172-7baa-4280-823f-19e597d091e9"),
		},
		{
			Id:           ptr.To("2"),
			Name:         ptr.To("source2"),
			SourceTypeId: ptr.To("2"),
			Uid:          ptr.To("31b5338b-685d-4056-ba39-d00b4d7f19cc"),
		},
	}
	return TestSourceData, nil
}

// APIClient
func WithSourcesIntegration(parent context.Context, init_store *[]sources.Source) context.Context {
	ctx := context.WithValue(parent, sourcesCtxKey, &SourcesIntegrationStub{store: init_store})
	return ctx
}

func (mock *SourcesIntegrationStub) ShowSourceWithResponse(ctx context.Context, id sources.ID, reqEditors ...sources.RequestEditorFn) (*sources.ShowSourceResponse, error) {
	lst := *mock.store
	return &sources.ShowSourceResponse{
		JSON200: &lst[0],
		HTTPResponse: &http.Response{
			StatusCode: 200,
		},
	}, nil
}

func (mock *SourcesIntegrationStub) ListApplicationTypeSourcesWithResponse(ctx context.Context, appTypeId sources.ID, params *sources.ListApplicationTypeSourcesParams, reqEditors ...sources.RequestEditorFn) (*sources.ListApplicationTypeSourcesResponse, error) {
	return &sources.ListApplicationTypeSourcesResponse{
		JSON200: &sources.SourcesCollection{
			Data: mock.store,
		},
		HTTPResponse: &http.Response{
			StatusCode: 200,
		},
	}, nil
}

func (mock *SourcesIntegrationStub) ListSourceAuthenticationsWithResponse(ctx context.Context, sourceId sources.ID, params *sources.ListSourceAuthenticationsParams, reqEditors ...sources.RequestEditorFn) (*sources.ListSourceAuthenticationsResponse, error) {
	return &sources.ListSourceAuthenticationsResponse{
		JSON200: &sources.AuthenticationsCollection{
			Data: mock.authentications,
		},
		HTTPResponse: &http.Response{
			StatusCode: 200,
		},
	}, nil
}

func (mock *SourcesIntegrationStub) ShowApplicationWithResponse(ctx context.Context, appId sources.ID, reqEditors ...sources.RequestEditorFn) (*sources.ShowApplicationResponse, error) {
	return &sources.ShowApplicationResponse{
		JSON200: &sources.Application{},
		HTTPResponse: &http.Response{
			StatusCode: 200,
		},
	}, nil
}
