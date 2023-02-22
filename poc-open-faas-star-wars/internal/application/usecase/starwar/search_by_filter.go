package starwar

import (
	"context"
	"handler/function/internal/application/model"
	"handler/function/internal/application/port/in"
	"handler/function/internal/application/port/out"
	"log"
)

var _ in.SearchByFilters = (*SearchByFilter)(nil)

type SearchByFilter struct {
	starwarRepository out.StarwarRepository
}

func NewGetByName(starwarRepository out.StarwarRepository) *SearchByFilter {
	return &SearchByFilter{
		starwarRepository: starwarRepository,
	}
}

func (c *SearchByFilter) GetFilmsInfo(urlFilms []string, ctx context.Context) []model.Film {

	channel := make(chan *model.SearchFilmResult, len(urlFilms))
	defer close(channel)
	filmsByCharacter := make([]model.Film, len(urlFilms))
	for localId, movie := range urlFilms {
		go func(localId int, movie string, localContext context.Context) {
			response, err := c.starwarRepository.GetFilms(localId, movie, localContext)
			if err != nil {
				log.Printf("Error getting films method %s\n", err.Error())
				channel <- &model.SearchFilmResult{Film: &model.Film{}}
				return
			}
			channel <- response
		}(localId, movie, ctx)
	}

	for range urlFilms {
		response := <-channel
		log.Printf("Response films %+v\n", response.Film)
		filmsByCharacter[response.LocalFilmId] = *response.Film
	}

	return filmsByCharacter
}

func (c *SearchByFilter) SearchByFilter(queryFilter *model.QueryFilter, ctx context.Context) (*model.CharacterSearchResult, error) {

	characterResult, err := c.starwarRepository.SearchCharacterByFilters(queryFilter, ctx)

	if err != nil {
		return nil, err
	}

	cantCharacter := len(characterResult.Character)

	for i := 0; i < cantCharacter; i++ {
		filmsInfo := c.GetFilmsInfo(characterResult.Character[i].Films, ctx)
		characterResult.Character[i].FilmsDetail = filmsInfo
	}

	return characterResult, nil
}
