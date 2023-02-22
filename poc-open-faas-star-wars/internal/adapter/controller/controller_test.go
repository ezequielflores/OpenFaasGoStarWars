package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	controllerModel "handler/function/internal/adapter/controller/model"
	"handler/function/internal/application/model"
	"handler/function/pkg"
	"io"
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
			"https://swapi.dev/api/films/2/",
			"https://swapi.dev/api/films/3/",
			"https://swapi.dev/api/films/6/",
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

var controllerSearchResultOk = controllerModel.ControllerCharacterSearchResult{
	Count:    1,
	Next:     "",
	Previous: "",
	Character: []controllerModel.ControllerCharacter{{
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
			"https://swapi.dev/api/films/2/",
			"https://swapi.dev/api/films/3/",
			"https://swapi.dev/api/films/6/",
		},
		FilmsDetail: []controllerModel.ControllerFilm{{
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
	},
	},
}

type NewMockedSearchByFilterOk struct {
	mock.Mock
}

func (m *NewMockedSearchByFilterOk) SearchByFilter(
	queryFilter *model.QueryFilter,
	ctx context.Context) (*model.CharacterSearchResult, error) {

	args := m.Called(queryFilter, ctx)
	firstParameter := args.Get(0)
	err := args.Error(1)
	if firstParameter != nil {
		response := firstParameter.(model.CharacterSearchResult)
		return &response, err
	}

	return nil, err
}

func TestNewStarWarController(t *testing.T) {

	mockObject := new(NewMockedSearchByFilterOk)
	controller := NewStarWarController(mockObject)

	assert.NotNil(t, controller)
	assert.IsType(t, &StarWarController{}, controller)
}

func TestBuildQueryParametersOk(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&p=10",
		nil,
	)

	queryFilter, err := BuildQueryParameters(req)
	expectedQueryFilter := &model.QueryFilter{
		Search: "data",
		Page:   "10",
	}

	assert.Nil(t, err)
	assert.Equal(t, expectedQueryFilter, queryFilter)
}

func TestBuildQueryParametersLetterPageParam(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&p=AAAAA",
		nil,
	)

	queryFilter, err := BuildQueryParameters(req)
	expectedError := pkg.GenericException{
		Msj:        "Invalid page parameter",
		StatusCode: http.StatusBadRequest,
	}

	assert.Nil(t, queryFilter)
	assert.Equal(t, &expectedError, err)
}

func TestBuildQueryParametersNegativePageNumber(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&p=-1",
		nil,
	)

	queryFilter, err := BuildQueryParameters(req)
	expectedError := &pkg.GenericException{
		Msj:        "Invalid page parameter",
		StatusCode: http.StatusBadRequest,
	}

	assert.Nil(t, queryFilter)
	assert.Equal(t, expectedError, err)
}

func TestBuildQueryParametersWithoutPageNumber(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data",
		nil,
	)

	queryFilter, err := BuildQueryParameters(req)
	expectedQueryFilter := &model.QueryFilter{
		Search: "data",
	}

	assert.Nil(t, err)
	assert.Equal(t, expectedQueryFilter, queryFilter)
}

func TestBuildQueryParametersWithoutQueryParam(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar",
		nil,
	)

	queryFilter, err := BuildQueryParameters(req)
	expectedQueryFilter := &model.QueryFilter{}

	assert.Nil(t, err)
	assert.Equal(t, expectedQueryFilter, queryFilter)
}

func TestSearchStarWarCharacterOk(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&p=1",
		nil,
	)
	rec := httptest.NewRecorder()

	mockObject := new(NewMockedSearchByFilterOk)
	mockObject.On("SearchByFilter",
		mock.IsType(&model.QueryFilter{}),
		req.Context()).
		Return(modelSearchResultOk, nil)

	controller := NewStarWarController(mockObject)
	controller.SearchStarWarCharacter(rec, req)

	jsonFullResponse, _ := json.Marshal(controllerSearchResultOk)
	body, _ := io.ReadAll(rec.Result().Body)

	assert.EqualValues(t, rec.Result().StatusCode, http.StatusOK)
	assert.JSONEqf(t, string(jsonFullResponse), string(body), "")
}

func TestSearchStarWarCharacterWithoutPageOk(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar",
		nil,
	)
	rec := httptest.NewRecorder()

	mockObject := new(NewMockedSearchByFilterOk)
	mockObject.On("SearchByFilter",
		mock.IsType(&model.QueryFilter{}),
		req.Context()).
		Return(modelSearchResultOk, nil)

	controller := NewStarWarController(mockObject)

	controller.SearchStarWarCharacter(rec, req)

	assert.EqualValues(t, rec.Result().StatusCode, http.StatusOK)
}

func TestSearchStarWarCharacterPageWithLetter(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&p=a",
		nil,
	)

	rec := httptest.NewRecorder()
	mockObject := new(NewMockedSearchByFilterOk)
	controller := NewStarWarController(mockObject)
	controller.SearchStarWarCharacter(rec, req)
	body, _ := io.ReadAll(rec.Result().Body)

	assert.EqualValues(t, rec.Result().StatusCode, http.StatusBadRequest)
	assert.EqualValues(t, "Invalid page parameter", string(body))
}
func TestSearchStarWarCharacterPageNegative(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&p=-1",
		nil,
	)

	rec := httptest.NewRecorder()
	mockObject := new(NewMockedSearchByFilterOk)
	controller := NewStarWarController(mockObject)
	controller.SearchStarWarCharacter(rec, req)
	body, _ := io.ReadAll(rec.Result().Body)

	assert.EqualValues(t, rec.Result().StatusCode, http.StatusBadRequest)
	assert.EqualValues(t, "Invalid page parameter", string(body))
}

func TestSearchStarWarCharacterWithUseCaseError(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&p=1",
		nil,
	)

	rec := httptest.NewRecorder()

	mockObject := new(NewMockedSearchByFilterOk)
	mockObject.On("SearchByFilter",
		mock.IsType(&model.QueryFilter{}),
		req.Context()).
		Return(nil, pkg.GenericException{
			Msj:        fmt.Sprint("Error searching character http status: 502"),
			StatusCode: 502,
		})

	controller := NewStarWarController(mockObject)
	controller.SearchStarWarCharacter(rec, req)
	body, _ := io.ReadAll(rec.Result().Body)

	assert.EqualValues(t, http.StatusBadGateway, rec.Result().StatusCode)
	assert.EqualValues(t, "Error searching character http status: 502", string(body))
}

func TestSearchStarWarCharacterWithoutPageOks(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar",
		nil,
	)
	rec := httptest.NewRecorder()
	mockObject := new(NewMockedSearchByFilterOk)
	mockObject.On("SearchByFilter",
		mock.IsType(&model.QueryFilter{}),
		req.Context()).
		Return(modelSearchResultOk, nil)

	controller := NewStarWarController(mockObject)

	controller.SearchStarWarCharacter(rec, req)

	assert.EqualValues(t, rec.Result().StatusCode, http.StatusOK)
}
