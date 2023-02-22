package in

import (
	"context"
	"handler/function/internal/application/model"
)

type SearchByFilters interface {
	SearchByFilter(queryFilter *model.QueryFilter, ctx context.Context) (*model.CharacterSearchResult, error)
}
