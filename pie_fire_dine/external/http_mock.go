package external

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type httpRequestMock struct {
	mock.Mock
}

func NewHttpRequestMock() *httpRequestMock {
	return &httpRequestMock{}
}

func (m *httpRequestMock) GetSourceText(context context.Context, url string) (int, []byte, error) {
	args := m.Called(context, url)
	return args.Int(0), args.Get(1).([]byte), args.Error(2)
}
