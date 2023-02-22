package out

import (
	"context"
	"handler/function/internal/application/model"
)

type StarwarRepository interface {
	SearchCharacterByFilters(queryFilter *model.QueryFilter, ctx context.Context) (*model.CharacterSearchResult, error)
	GetFilms(localId int, url string, ctx context.Context) (*model.SearchFilmResult, error)
}
