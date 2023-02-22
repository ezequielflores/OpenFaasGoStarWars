package starwar

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"handler/function/internal/application/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

var modelSearchResultOk = model.CharacterSearchResult{
	Count:    1,
	Next:     "",
	Previous: "",
	Character: []model.Character{{
		Name:      "Test POC Datafaas",
		Height:    "202",
		Mass:      "136",
		HairColor: "none",
		SkinColor: "white",
		EyeColor:  "yellow",
		BirthYear: "41.9BBY",
		Gender:    "male",
		Homeworld: "https://swapi.dev/api/planets/1/",
		Films: []string{
			"https://swapi.dev/api/films/1/",
		},
		FilmsDetail: []model.Film{
			{
				Title:        "",
				EpisodeId:    1,
				OpeningCrawl: "It is a period of civil war.\r\nRebel spaceships, striking\r\nfrom a hidden base, have won\r\ntheir first victory against\r\nthe evil Galactic Empire.\r\n\r\nDuring the battle, Rebel\r\nspies managed to steal secret\r\nplans to the Empire's\r\nultimate weapon, the DEATH\r\nSTAR, an armored space\r\nstation with enough power\r\nto destroy an entire planet.\r\n\r\nPursued by the Empire's\r\nsinister agents, Princess\r\nLeia races home aboard her\r\nstarship, custodian of the\r\nstolen plans that can save her\r\npeople and restore\r\nfreedom to the galaxy....",
				Director:     "George Lucas",
				Producer:     "Gary Kurtz, Rick McCallum",
				ReleaseDate:  "1977-05-25",
			},
		},
		Species: []string{
			"https://swapi.dev/api/species/1/",
			"https://swapi.dev/api/species/2/",
			"https://swapi.dev/api/species/3/",
			"https://swapi.dev/api/species/4/",
			"https://swapi.dev/api/species/5/",
		},
		Vehicles: []string{
			"https://swapi.dev/api/vehicles/4/",
			"https://swapi.dev/api/vehicles/6/",
			"https://swapi.dev/api/vehicles/7/",
			"https://swapi.dev/api/vehicles/8/",
		},
		Starships: []string{
			"https://swapi.dev/api/starships/13/",
		},

		Created: "2014-12-10T15:18:20.704000Z",
		Edited:  "2014-12-20T21:17:50.313000Z",
		Url:     "https://swapi.dev/api/people/4/",
	}},
}

var filmsModelResult = model.Film{
	Title:        "",
	EpisodeId:    1,
	OpeningCrawl: "It is a period of civil war.\r\nRebel spaceships, striking\r\nfrom a hidden base, have won\r\ntheir first victory against\r\nthe evil Galactic Empire.\r\n\r\nDuring the battle, Rebel\r\nspies managed to steal secret\r\nplans to the Empire's\r\nultimate weapon, the DEATH\r\nSTAR, an armored space\r\nstation with enough power\r\nto destroy an entire planet.\r\n\r\nPursued by the Empire's\r\nsinister agents, Princess\r\nLeia races home aboard her\r\nstarship, custodian of the\r\nstolen plans that can save her\r\npeople and restore\r\nfreedom to the galaxy....",
	Director:     "George Lucas",
	Producer:     "Gary Kurtz, Rick McCallum",
	ReleaseDate:  "1977-05-25",
}

var searchFilmModelResul = model.SearchFilmResult{
	Film:        &filmsModelResult,
	LocalFilmId: 0,
}

type MockSearchByFilter struct {
	mock.Mock
}

func (m *MockSearchByFilter) GetFilms(
	localId int,
	url string,
	ctx context.Context) (*model.SearchFilmResult, error) {

	args := m.Called(localId, url, ctx)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	firstParameter := args.Get(0)
	films := firstParameter.(*model.SearchFilmResult)

	if localId == 0 {
		return films, nil
	}

	return nil, nil
}

func (m *MockSearchByFilter) SearchCharacterByFilters(
	queryFilter *model.QueryFilter,
	ctx context.Context) (*model.CharacterSearchResult, error) {

	args := m.Called(queryFilter, ctx)

	firstArg := args.Get(0)

	if firstArg != nil && queryFilter.Page == "1" {
		return firstArg.(*model.CharacterSearchResult), nil
	}

	return nil, args.Error(1)
}

func TestNewGetByName(t *testing.T) {

	mockSearchByFilter := new(MockSearchByFilter)
	newGetByName := NewGetByName(mockSearchByFilter)

	assert.NotNil(t, newGetByName)
}

func TestSearchByFilter_GetFilmsInfo(t *testing.T) {

	mockSearchByFilter := new(MockSearchByFilter)
	newGetByName := NewGetByName(mockSearchByFilter)
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&p=1",
		nil,
	)

	filmUrls := []string{
		"https://swapi.dev/api/films/1/",
	}

	ctx := req.Context()
	mockSearchByFilter.On("GetFilms", 0, filmUrls[0], ctx).
		Return(&searchFilmModelResul, nil)

	getFilmsInfo := newGetByName.GetFilmsInfo(filmUrls, ctx)
	assert.Equal(t, filmsModelResult, getFilmsInfo[0])
}

func TestSearchByFilter_GetFilmsInfo_Error(t *testing.T) {

	mockSearchByFilter := new(MockSearchByFilter)
	newGetByName := NewGetByName(mockSearchByFilter)
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&p=1",
		nil,
	)

	filmUrls := []string{
		"https://swapi.dev/api/films/1/",
	}

	ctx := req.Context()
	mockSearchByFilter.On("GetFilms", 0, filmUrls[0], ctx).
		Return(nil, fmt.Errorf("generic error"))

	getFilmsInfo := newGetByName.GetFilmsInfo(filmUrls, ctx)
	assert.Equal(t, model.Film{}, getFilmsInfo[0])
}

func TestSearchByFilter(t *testing.T) {

	mockSearchByFilter := new(MockSearchByFilter)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&p=1",
		nil,
	)

	ctx := req.Context()
	queryFilter := model.QueryFilter{
		Search: "darth",
		Page:   "1",
	}

	filmUrls := modelSearchResultOk.Character[0].Films

	mockSearchByFilter.On("SearchCharacterByFilters",
		mock.IsType(&model.QueryFilter{}), ctx).
		Return(&modelSearchResultOk, nil)

	mockSearchByFilter.On("GetFilms", 0, filmUrls[0], ctx).
		Return(&searchFilmModelResul, nil)

	newGetByName := NewGetByName(mockSearchByFilter)
	searchCharacterResult, err := newGetByName.SearchByFilter(&queryFilter, ctx)

	assert.Nil(t, err)
	assert.Equal(t, &modelSearchResultOk, searchCharacterResult)

}

func TestSearchByFilterErro(t *testing.T) {

	mockSearchByFilter := new(MockSearchByFilter)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&p=1",
		nil,
	)

	ctx := req.Context()
	queryFilter := model.QueryFilter{
		Search: "darth",
		Page:   "1",
	}

	mockSearchByFilter.On("SearchCharacterByFilters",
		mock.IsType(&model.QueryFilter{}), ctx).
		Return(nil, fmt.Errorf("generic error"))

	newGetByName := NewGetByName(mockSearchByFilter)
	searchCharacterResult, err := newGetByName.SearchByFilter(&queryFilter, ctx)

	assert.Nil(t, searchCharacterResult)
	assert.Equal(t, "generic error", err.Error())

}
