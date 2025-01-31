package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

type imageBuilderCtxKeyType string

var imageBuilderCtxKey imageBuilderCtxKeyType = "image-builder-interface"

type ImageBuilderClientStub struct{}

func init() {
	clients.GetImageBuilderClient = getImageBuilderClientStub
}

type contextReadError struct{}

func (m *contextReadError) Error() string {
	return "failed to find or convert dao stored in testing context"
}

func WithImageBuilderClient(parent context.Context) context.Context {
	ctx := context.WithValue(parent, imageBuilderCtxKey, &ImageBuilderClientStub{})
	return ctx
}

func getImageBuilderClientStub(ctx context.Context) (si clients.ImageBuilder, err error) {
	var ok bool
	if si, ok = ctx.Value(imageBuilderCtxKey).(*ImageBuilderClientStub); !ok {
		err = &contextReadError{}
	}
	return si, err
}

func (*ImageBuilderClientStub) Ready(ctx context.Context) error {
	return nil
}

func (mock *ImageBuilderClientStub) GetAWSAmi(ctx context.Context, composeID string) (string, error) {
	return "ami-0c830793775595d4b-test", nil
}

func (mock *ImageBuilderClientStub) GetGCPImageName(ctx context.Context, composeID string) (string, error) {
	return "projects/red-hat-image-builder/global/images/composer-api-871fa36d-0b5b-4001-8c95-a11f751a4d66-test", nil
}
