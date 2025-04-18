package service

import (
	"context"
	"pie_fire_dine/contract"

	"github.com/stretchr/testify/mock"
)

type meatSummaryMock struct {
	mock.Mock
}

func NewHttpRequestMock() *meatSummaryMock {
	return &meatSummaryMock{}
}

func (m *meatSummaryMock) GetMeatSummary(ctx context.Context, category string) (map[string]contract.MeatCategoryCountMap, error) {
	args := m.Called(ctx, category)
	return args.Get(0).(map[string]contract.MeatCategoryCountMap), args.Error(1)
}
