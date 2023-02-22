package rest

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"handler/function/internal/application/model"
	"handler/function/pkg"
	"net/http"
	"net/http/httptest"
	"testing"
)

var jsonSearchResponse = "{\n\"count\":1,\n\"next\":null,\n\"previous\":null,\n\"results\":[\n{\n\"name\":\"Darth Vader\",\n\"height\":\"202\",\n\"mass\":\"136\",\n\"hair_color\":\"none\",\n\"skin_color\":\"white\",\n\"eye_color\":\"yellow\",\n\"birth_year\":\"41.9BBY\",\n\"gender\":\"male\",\n\"homeworld\":\"https://swapi.dev/api/planets/1/\",\n\"films\":[\n\"https://swapi.dev/api/films/1/\",\n\"https://swapi.dev/api/films/2/\",\n\"https://swapi.dev/api/films/3/\",\n\"https://swapi.dev/api/films/6/\"\n],\n\"species\":[\n\n],\n\"vehicles\":[\n\n],\n\"starships\":[\n\"https://swapi.dev/api/starships/13/\"\n],\n\"created\":\"2014-12-10T15:18:20.704000Z\",\n\"edited\":\"2014-12-20T21:17:50.313000Z\",\n\"url\":\"https://swapi.dev/api/people/4/\"\n}\n]\n}"

var structFilmResponse = model.Film{
	Title:        "A New Hope",
	EpisodeId:    1,
	OpeningCrawl: "It is a period of civil war.\r\nRebel spaceships, striking\r\nfrom a hidden base, have won\r\ntheir first victory against\r\nthe evil Galactic Empire.\r\n\r\nDuring the battle, Rebel\r\nspies managed to steal secret\r\nplans to the Empire's\r\nultimate weapon, the DEATH\r\nSTAR, an armored space\r\nstation with enough power\r\nto destroy an entire planet.\r\n\r\nPursued by the Empire's\r\nsinister agents, Princess\r\nLeia races home aboard her\r\nstarship, custodian of the\r\nstolen plans that can save her\r\npeople and restore\r\nfreedom to the galaxy....",
	Director:     "George Lucas",
	Producer:     "Gary Kurtz, Rick McCallum",
	ReleaseDate:  "1977-05-25",
}

var structSearchFilmResponse = model.SearchFilmResult{
	Film:        &structFilmResponse,
	LocalFilmId: 0,
}

var jsonFilmResponse = "{\n\"title\": \"A New Hope\",\n\"episode_id\": 1,\n\"opening_crawl\": \"It is a period of civil war.\\r\\nRebel spaceships, striking\\r\\nfrom a hidden base, have won\\r\\ntheir first victory against\\r\\nthe evil Galactic Empire.\\r\\n\\r\\nDuring the battle, Rebel\\r\\nspies managed to steal secret\\r\\nplans to the Empire's\\r\\nultimate weapon, the DEATH\\r\\nSTAR, an armored space\\r\\nstation with enough power\\r\\nto destroy an entire planet.\\r\\n\\r\\nPursued by the Empire's\\r\\nsinister agents, Princess\\r\\nLeia races home aboard her\\r\\nstarship, custodian of the\\r\\nstolen plans that can save her\\r\\npeople and restore\\r\\nfreedom to the galaxy....\",\n\"director\": \"George Lucas\",\n\"producer\": \"Gary Kurtz, Rick McCallum\",\n\"release_date\": \"1977-05-25\"\n}"

var structResponse = model.CharacterSearchResult{
	Count:    1,
	Next:     "",
	Previous: "",
	Character: []model.Character{{
		Name:      "Darth Vader",
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
		FilmsDetail: []model.Film{},
		Species:     []string{},
		Vehicles:    []string{},
		Starships: []string{
			"https://swapi.dev/api/starships/13/",
		},

		Created: "2014-12-10T15:18:20.704000Z",
		Edited:  "2014-12-20T21:17:50.313000Z",
		Url:     "https://swapi.dev/api/people/4/",
	}},
}

func TestNewStarwarRestAdapter(t *testing.T) {
	t.Setenv("character_search_timeout", "31000")
	t.Setenv("character_search_url", "http://url.test")
	restAdapter, err := NewStarwarRestAdapter()

	assert.Nil(t, err)
	assert.NotNil(t, restAdapter)
}

func TestNewStarwarRestAdapterError(t *testing.T) {

	restAdapter, err := NewStarwarRestAdapter()

	assert.Nil(t, restAdapter)
	assert.NotNil(t, err)
}

func TestSearchCharacterByFilters(t *testing.T) {
	t.Setenv("character_search_timeout", "31000")
	t.Setenv("character_search_url", "http://url.test")

	restAdapter, _ := NewStarwarRestAdapter()

	httpmock.ActivateNonDefault(restAdapter.client.GetClient())
	defer httpmock.DeactivateAndReset()

	restMockResponder := httpmock.NewStringResponder(200, jsonSearchResponse)
	httpmock.RegisterResponder("GET", "http://url.test?page=1&search=darth", restMockResponder)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&p=1",
		nil,
	)

	queryFilter := model.QueryFilter{
		Search: "darth",
		Page:   "1",
	}
	searchResult, err := restAdapter.SearchCharacterByFilters(&queryFilter, req.Context())

	assert.Nil(t, err)
	assert.Equal(t, &structResponse, searchResult)
}

func TestSearchCharacterByFiltersError(t *testing.T) {
	t.Setenv("character_search_timeout", "31000")
	t.Setenv("character_search_url", "http://url.test")

	restAdapter, _ := NewStarwarRestAdapter()

	httpmock.ActivateNonDefault(restAdapter.client.GetClient())
	defer httpmock.DeactivateAndReset()

	restMockResponder := httpmock.NewErrorResponder(fmt.Errorf("error from api"))
	httpmock.RegisterResponder("GET", "http://url.test?page=1&search=darth", restMockResponder)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&p=1",
		nil,
	)

	queryFilter := model.QueryFilter{
		Search: "darth",
		Page:   "1",
	}
	searchResult, err := restAdapter.SearchCharacterByFilters(&queryFilter, req.Context())

	assert.Nil(t, searchResult)
	assert.NotNil(t, err)
	assert.Equal(t, "Get \"http://url.test?page=1&search=darth\": error from api", err.Error())
}

func TestStarwarRestAdapterHttsCode404(t *testing.T) {

	t.Setenv("character_search_timeout", "31000")
	t.Setenv("character_search_url", "http://url.test")

	restAdapter, _ := NewStarwarRestAdapter()

	httpmock.ActivateNonDefault(restAdapter.client.GetClient())
	defer httpmock.DeactivateAndReset()

	restMockResponder := httpmock.NewStringResponder(404, `{"mesg":"Not found"}`)
	httpmock.RegisterResponder("GET", "http://url.test?page=1&search=darth", restMockResponder)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&p=1",
		nil,
	)

	queryFilter := model.QueryFilter{
		Search: "darth",
		Page:   "1",
	}

	response, err := restAdapter.SearchCharacterByFilters(&queryFilter, req.Context())
	assert.NotNil(t, err)
	assert.IsTypef(t, pkg.GenericException{}, err, "It's not same type")
	exception := err.(pkg.GenericException)
	assert.Nil(t, response)
	assert.Equal(t, "Error searching character http status: 404", exception.Msj)

}

func TestStarwarRestAdapterGetFilmsOK(t *testing.T) {

	t.Setenv("character_search_timeout", "31000")
	t.Setenv("character_search_url", "http://url.test")

	restAdapter, _ := NewStarwarRestAdapter()

	httpmock.ActivateNonDefault(restAdapter.client.GetClient())
	defer httpmock.DeactivateAndReset()

	restMockResponder := httpmock.NewStringResponder(200, jsonFilmResponse)
	httpmock.RegisterResponder("GET", "https://swapi.dev/api/films/1/", restMockResponder)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&p=1",
		nil,
	)

	response, err := restAdapter.GetFilms(0, "https://swapi.dev/api/films/1/", req.Context())

	assert.Nil(t, err)
	assert.Equal(t, &structSearchFilmResponse, response)

}

func TestStarwarRestAdapterGetFilmsError(t *testing.T) {

	t.Setenv("character_search_timeout", "31000")
	t.Setenv("character_search_url", "http://url.test")

	restAdapter, _ := NewStarwarRestAdapter()

	httpmock.ActivateNonDefault(restAdapter.client.GetClient())
	defer httpmock.DeactivateAndReset()

	restMockResponder := httpmock.NewErrorResponder(fmt.Errorf("error from api"))
	httpmock.RegisterResponder("GET", "https://swapi.dev/api/films/1/", restMockResponder)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&p=1",
		nil,
	)

	response, err := restAdapter.GetFilms(0, "https://swapi.dev/api/films/1/", req.Context())

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "error getting films: Get \"https://swapi.dev/api/films/1/\": error from api\n", err.Error())

}

func TestStarwarRestAdapterGetFilmsHttp404(t *testing.T) {

	t.Setenv("character_search_timeout", "31000")
	t.Setenv("character_search_url", "http://url.test")

	restAdapter, _ := NewStarwarRestAdapter()

	httpmock.ActivateNonDefault(restAdapter.client.GetClient())
	defer httpmock.DeactivateAndReset()

	restMockResponder := httpmock.NewStringResponder(404, `{"mesg":"Not found"}`)
	httpmock.RegisterResponder("GET", "https://swapi.dev/api/films/1/", restMockResponder)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&p=1",
		nil,
	)

	response, err := restAdapter.GetFilms(0, "https://swapi.dev/api/films/1/", req.Context())

	assert.NotNil(t, err)
	assert.IsTypef(t, pkg.GenericException{}, err, "It's not same type")
	exception := err.(pkg.GenericException)
	assert.Nil(t, response)
	assert.Equal(t, "Error searching films http status: 404", exception.Msj)

}
