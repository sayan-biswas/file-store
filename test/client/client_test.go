package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockStoreService struct {
	mock.Mock
}

func (m *MockStoreService) Add(ctx context.Context, data []byte) {}

func TestAdd(t *testing.T) {
}
