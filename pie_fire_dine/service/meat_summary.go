package service

import (
	"context"
	"fmt"
	"pie_fire_dine/config"
	"pie_fire_dine/contract"
	"pie_fire_dine/errs"
	"pie_fire_dine/external"
	"pie_fire_dine/repository"
	"regexp"
	"strings"

	"log/slog"
)

type MeatSummary interface {
	GetMeatSummary(context.Context, string) (map[string]contract.MeatCategoryCountMap, error)
}

type meatSummary struct {
	meatRepo   repository.MeatRepository
	httpClient external.ExternalApiRequester
}

func NewMeatSummaryService(meatRepo repository.MeatRepository, httpClient external.ExternalApiRequester) MeatSummary {
	return meatSummary{
		meatRepo:   meatRepo,
		httpClient: httpClient,
	}
}

func (m meatSummary) GetMeatSummary(ctx context.Context, category string) (map[string]contract.MeatCategoryCountMap, error) {
	logAttrs := slog.Group("get_meat_summary",
		slog.String("category", category),
	)

	var meatTypeCount contract.MeatCategoryCountMap = make(contract.MeatCategoryCountMap)
	var meatTypes []string

	meats, err := m.meatRepo.GetAllByCategory(ctx, category)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Error getting meats by category, err: %v", err), logAttrs)
		return nil, errs.NewUnexpectedError()
	}
	if len(meats) == 0 {
		slog.ErrorContext(ctx, fmt.Sprintf("No meats found for category: %s", category), logAttrs)
		return nil, errs.NewBadRequest()
	}
	for _, meat := range meats {
		meatTypes = append(meatTypes, meat.Name)
	}

	regexFromCategories, err := createRegexFromMeatTypes(meatTypes)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Error creating regex from meats, err: %v", err), logAttrs)
		return nil, errs.NewUnexpectedError()
	}

	statusCode, res, err := m.httpClient.GetSourceText(ctx, config.GetDefaultSourceTextURL())
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Error getting source text, err: %v", err), logAttrs)
		return nil, errs.NewUnexpectedError()
	}
	if statusCode != 200 {
		slog.ErrorContext(ctx, fmt.Sprintf("Error getting source text, status code: %d", statusCode), logAttrs)
		return nil, errs.NewUnexpectedError()
	}

	meatTypeMatched := regexFromCategories.FindAllString(string(res), -1)

	for _, meatType := range meatTypeMatched {
		lowerMeatType := strings.ToLower(meatType)
		if _, ok := meatTypeCount[lowerMeatType]; !ok {
			meatTypeCount[lowerMeatType] = 1
			continue
		}
		meatTypeCount[lowerMeatType]++
	}

	result := map[string]contract.MeatCategoryCountMap{
		category: meatTypeCount,
	}

	return result, nil
}

func createRegexFromMeatTypes(meatTypes []string) (*regexp.Regexp, error) {
	var escapedMeatType []string
	for _, meatType := range meatTypes {
		escaped := regexp.QuoteMeta(meatType)
		escapedMeatType = append(escapedMeatType, escaped)
	}
	pattern := strings.Join(escapedMeatType, "|")

	regex, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %w", err)
	}

	return regex, nil
}
