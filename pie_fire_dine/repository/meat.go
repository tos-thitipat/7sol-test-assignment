package repository

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"strings"
)

type Meat struct {
	Name     string
	Category string
}

type MeatRepository interface {
	GetAllByCategory(context.Context, string) ([]Meat, error)
}

type meatRepository struct {
	filepath string
}

func NewMeatRepository(filepath string) MeatRepository {
	return &meatRepository{
		filepath: filepath,
	}
}

func (m meatRepository) GetAllByCategory(ctx context.Context, category string) ([]Meat, error) {
	var meats []Meat
	file, err := os.Open(m.filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	headerMap := make(map[string]int)
	for i, colName := range header {
		headerMap[colName] = i
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if strings.EqualFold(row[headerMap["Category"]], category) {
			meats = append(meats, Meat{
				Name:     row[headerMap["Meat"]],
				Category: row[headerMap["Category"]],
			})
		}
	}

	return meats, nil
}
