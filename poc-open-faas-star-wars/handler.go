package function

import (
	"fmt"
	"handler/function/internal/adapter/controller"
	"handler/function/internal/adapter/rest"
	"handler/function/internal/application/usecase/starwar"
	"log"
	"net/http"
	"strings"
)

var routes = map[string]func(http.ResponseWriter, *http.Request){}

func init() {
	restAdapter, err := rest.NewStarwarRestAdapter()
	if err != nil {
		log.Println("Error initializing rest adapter")
		return
	}
	searchByFilter := starwar.NewGetByName(restAdapter)
	starwarsController := controller.NewStarWarController(searchByFilter)
	routes["/starwars"] = starwarsController.SearchStarWarCharacter

}

func Handle(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/v1/starwars") {
		routes["/starwars"](w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(fmt.Sprintf("Url: %s %s", r.URL.Path, http.StatusText(http.StatusNotFound))))
}
