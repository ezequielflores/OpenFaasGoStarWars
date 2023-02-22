package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	restmodel "handler/function/internal/adapter/rest/model"
	"handler/function/internal/application/model"
	"handler/function/internal/application/port/out"
	"handler/function/pkg"
	"log"
	"os"
	"strconv"
	"time"
)

var _ out.StarwarRepository = (*StarwarRestAdapter)(nil)
var searchUrl string

type StarwarRestAdapter struct {
	client *resty.Client
}

func NewStarwarRestAdapter() (*StarwarRestAdapter, error) {

	searchUrl = os.Getenv("character_search_url")
	stringTimeout := os.Getenv("character_search_timeout")
	timeout, parserserr := strconv.Atoi(stringTimeout)

	if parserserr != nil {
		log.Printf("Error parsing tiemout parameter. Value: %s", stringTimeout)
		return nil, fmt.Errorf("error parsing timeout parameter %w", parserserr)
	}

	client := resty.
		New().
		SetTimeout(time.Duration(timeout) * time.Millisecond).
		SetRetryWaitTime(time.Duration(timeout) * time.Millisecond)

	return &StarwarRestAdapter{client: client}, nil
}

func searchBuilderParameter(queryFilter *model.QueryFilter) map[string]string {
	parameters := make(map[string]string)

	if queryFilter.Search != "" {
		parameters["search"] = queryFilter.Search
	}

	if queryFilter.Page != "" {
		parameters["page"] = queryFilter.Page
	}

	return parameters
}

func (a *StarwarRestAdapter) SearchCharacterByFilters(
	queryFilter *model.QueryFilter,
	ctx context.Context,
) (*model.CharacterSearchResult, error) {

	parameters := searchBuilderParameter(queryFilter)

	log.Printf("Call to searchUrl: %s with parameters: %s\n", searchUrl, parameters)

	response, err := a.client.
		R().
		SetContext(ctx).
		SetQueryParams(parameters).
		Get(searchUrl)

	if err != nil {
		log.Printf("Error %s\n", err.Error())
		return nil, err
	}

	log.Printf("Http status Returned %s\n", response.Status())
	log.Printf("Http Body Returned %s\n", string(response.Body()))

	if response.StatusCode() != 200 {
		fmt.Printf("Error searching star war character: %s\n", parameters)
		return nil, pkg.GenericException{
			Msj:        fmt.Sprintf("Error searching character http status: %s", response.Status()),
			StatusCode: response.StatusCode(),
		}
	}

	var responseObject = &restmodel.RestCharacterSearchResult{}
	if err := json.Unmarshal(response.Body(), responseObject); err != nil {
		log.Printf("Error parser response %s", err.Error())
		return nil, err
	}

	return responseObject.ToDomain()
}

func (a *StarwarRestAdapter) GetFilms(
	localId int,
	filmUrl string,
	ctx context.Context) (*model.SearchFilmResult, error) {

	log.Printf("Call to searchUrl: %s\n", filmUrl)

	response, err := a.client.
		R().
		SetContext(ctx).
		Get(filmUrl)

	if err != nil {
		log.Printf("Error getting films %s\n", err.Error())
		return nil, fmt.Errorf("error getting films: %s\n", err.Error())
	}

	log.Printf("Http status Returned %s\n", response.Status())
	log.Printf("Http Body Returned %s\n", response.Body())

	if response.StatusCode() != 200 {
		fmt.Printf("Error searching star war films with searchUrl: %s\n", filmUrl)
		return nil, pkg.GenericException{
			Msj:        fmt.Sprintf("Error searching films http status: %s", response.Status()),
			StatusCode: response.StatusCode(),
		}
	}

	searchFilmResult := &restmodel.RestFilmResult{}
	if err := json.Unmarshal(response.Body(), searchFilmResult); err != nil {
		log.Printf("Error parsing films response %s\n", err.Error())
		return nil, fmt.Errorf("error parsing films response: %s\n", err.Error())
	}

	return &model.SearchFilmResult{
		Film:        searchFilmResult.ToDomain(),
		LocalFilmId: localId,
	}, nil
}
