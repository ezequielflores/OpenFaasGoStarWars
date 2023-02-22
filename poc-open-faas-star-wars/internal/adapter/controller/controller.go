package controller

import (
	"encoding/json"
	controllerModel "handler/function/internal/adapter/controller/model"
	"handler/function/internal/application/model"
	"handler/function/internal/application/port/in"
	"handler/function/pkg"
	"log"
	"net/http"
	"strconv"
)

type StarWarController struct {
	searchByFilters in.SearchByFilters
}

func NewStarWarController(
	searchByFilters in.SearchByFilters,
) *StarWarController {
	return &StarWarController{
		searchByFilters: searchByFilters,
	}
}

func BuildQueryParameters(r *http.Request) (*model.QueryFilter, *pkg.GenericException) {

	queryFilter := model.QueryFilter{
		Search: r.URL.Query().Get("q"),
		Page:   "",
	}
	page := r.URL.Query().Get("p")

	if page != "" {
		if pageValid, errorInt := strconv.Atoi(page); errorInt == nil && pageValid > 0 {
			log.Printf("page %d\n", pageValid)
			queryFilter.Page = page
		} else {
			return nil, &pkg.GenericException{
				Msj:        "Invalid page parameter",
				StatusCode: http.StatusBadRequest,
			}
		}
	}

	return &queryFilter, nil
}

func (c *StarWarController) SearchStarWarCharacter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	queryParameters, queryError := BuildQueryParameters(r)

	if queryError != nil {
		w.WriteHeader(queryError.StatusCode)
		w.Write([]byte(queryError.Msj))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	searchResult, err := c.searchByFilters.SearchByFilter(queryParameters, ctx)

	if err != nil {
		statusCode, errorMsj := pkg.GetErrorDetail(err)
		w.WriteHeader(statusCode)
		w.Write([]byte(errorMsj))
		return
	}

	jsonResult, jsonError := json.Marshal(controllerModel.FromDomain(searchResult))
	if jsonError != nil {
		log.Printf("Error converting to json %s", jsonError.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(jsonError.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResult)
}
