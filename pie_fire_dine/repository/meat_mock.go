package repository

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type meatRepositoryMock struct {
	mock.Mock
}

func NewMeatRepositoryMock() *meatRepositoryMock {
	return &meatRepositoryMock{}
}

func (m *meatRepositoryMock) GetAllByCategory(ctx context.Context, category string) ([]Meat, error) {
	args := m.Called(ctx, category)
	return args.Get(0).([]Meat), args.Error(1)
}
