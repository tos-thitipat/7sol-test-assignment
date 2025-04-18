package service_test

import (
	"context"
	"errors"
	"pie_fire_dine/contract"
	"pie_fire_dine/errs"
	"pie_fire_dine/external"
	"pie_fire_dine/repository"
	"pie_fire_dine/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMeatSummary(t *testing.T) {
	type repoMockReturn struct {
		meats []repository.Meat
		err   error
	}
	type externalMockReturn struct {
		statusCode int
		body       []byte
		err        error
	}
	type expectedReturn struct {
		result map[string]contract.MeatCategoryCountMap
		err    error
	}
	testCases := []struct {
		name         string
		category     string
		repoMock     repoMockReturn
		externalMock externalMockReturn
		expected     expectedReturn
	}{
		{
			name:     "beef summary success",
			category: "beef",
			repoMock: repoMockReturn{
				meats: []repository.Meat{
					{Name: "tenderloin", Category: "beef"},
					{Name: "picanha", Category: "beef"},
				},
				err: nil,
			},
			externalMock: externalMockReturn{
				statusCode: 200,
				body:       []byte("tenderloin picanha"),
				err:        nil,
			},
			expected: expectedReturn{
				result: map[string]contract.MeatCategoryCountMap{
					"beef": {
						"tenderloin": 1,
						"picanha":    1,
					},
				},
				err: nil,
			},
		},
		{
			name:     "pork summary success",
			category: "pork",
			repoMock: repoMockReturn{
				meats: []repository.Meat{
					{Name: "bacon", Category: "pork"},
					{Name: "ham", Category: "pork"},
				},
				err: nil,
			},
			externalMock: externalMockReturn{
				statusCode: 200,
				body:       []byte("bacon ham ham"),
				err:        nil,
			},
			expected: expectedReturn{
				result: map[string]contract.MeatCategoryCountMap{
					"pork": {
						"bacon": 1,
						"ham":   2,
					},
				},
				err: nil,
			},
		},
		{
			name:     "failed GetAllByCategory",
			category: "pork",
			repoMock: repoMockReturn{
				meats: nil,
				err:   errors.New("failed to get meats"),
			},
			externalMock: externalMockReturn{
				statusCode: 200,
				body:       []byte("bacon ham ham"),
				err:        nil,
			},
			expected: expectedReturn{
				result: nil,
				err:    errs.NewUnexpectedError(),
			},
		},
		{
			name:     "failed GetAllByCategory invalid category",
			category: "pork",
			repoMock: repoMockReturn{
				meats: []repository.Meat{},
				err:   nil,
			},
			externalMock: externalMockReturn{
				statusCode: 200,
				body:       []byte("bacon ham ham"),
				err:        nil,
			},
			expected: expectedReturn{
				result: nil,
				err:    errs.NewBadRequest(),
			},
		},
		{
			name:     "failed regexFromCategories",
			category: "pork",
			repoMock: repoMockReturn{
				meats: []repository.Meat{},
				err:   nil,
			},
			externalMock: externalMockReturn{
				statusCode: 200,
				body:       []byte("bacon ham ham"),
				err:        nil,
			},
			expected: expectedReturn{
				result: nil,
				err:    errs.NewBadRequest(),
			},
		},
		{
			name:     "failed GetSourceText error",
			category: "pork",
			repoMock: repoMockReturn{
				meats: []repository.Meat{
					{Name: "bacon", Category: "pork"},
					{Name: "ham", Category: "pork"},
				},
				err: nil,
			},
			externalMock: externalMockReturn{
				statusCode: 500,
				body:       nil,
				err:        errors.New("failed to get source text"),
			},
			expected: expectedReturn{
				result: nil,
				err:    errs.NewUnexpectedError(),
			},
		},
		{
			name:     "failed GetSourceText status != 200",
			category: "pork",
			repoMock: repoMockReturn{
				meats: []repository.Meat{
					{Name: "bacon", Category: "pork"},
					{Name: "ham", Category: "pork"},
				},
				err: nil,
			},
			externalMock: externalMockReturn{
				statusCode: 402,
				body:       nil,
				err:        nil,
			},
			expected: expectedReturn{
				result: nil,
				err:    errs.NewUnexpectedError(),
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()

			// Create a mock repository
			meatRepository := repository.NewMeatRepositoryMock()
			meatRepository.On("GetAllByCategory", ctx, c.category).Return(c.repoMock.meats, c.repoMock.err)

			// Create a mock external API requester
			externalClient := external.NewHttpRequestMock()
			externalClient.On("GetSourceText", ctx, "").Return(c.externalMock.statusCode, c.externalMock.body, c.externalMock.err)

			meatSummaryService := service.NewMeatSummaryService(meatRepository, externalClient)
			result, err := meatSummaryService.GetMeatSummary(ctx, c.category)
			assert.Equal(t, err, c.expected.err)
			assert.Equal(t, c.expected.result, result)
		})
	}
}
