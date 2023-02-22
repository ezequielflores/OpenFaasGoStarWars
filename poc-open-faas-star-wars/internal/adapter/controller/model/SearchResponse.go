package model

import (
	domainModel "handler/function/internal/application/model"
)

type QueryFilter struct {
	Search string
	Page   string
}

type ControllerFilm struct {
	Title        string `json:"title"`
	EpisodeId    int    `json:"episode_id"`
	OpeningCrawl string `json:"opening_crawl"`
	Director     string `json:"director"`
	Producer     string `json:"producer"`
	ReleaseDate  string `json:"release_date"`
}

type ControllerCharacter struct {
	Name        string           `json:"name"`
	Height      string           `json:"height"`
	Mass        string           `json:"mass"`
	HairColor   string           `json:"hair_color"`
	SkinColor   string           `json:"skin_color"`
	EyeColor    string           `json:"eye_color"`
	BirthYear   string           `json:"birth_year"`
	Gender      string           `json:"gender"`
	Homeworld   string           `json:"homeworld"`
	Films       []string         `json:"films"`
	FilmsDetail []ControllerFilm `json:"films_details"`
	Species     []string         `json:"species"`
	Vehicles    []string         `json:"vehicles"`
	Starships   []string         `json:"starships"`
	Created     string           `json:"created"`
	Edited      string           `json:"edited"`
	Url         string           `json:"url"`
}

type ControllerCharacterSearchResult struct {
	Count     int                   `json:"count"`
	Next      string                `json:"next"`
	Previous  string                `json:"previous"`
	Character []ControllerCharacter `json:"results"`
}

type ControllerSearchFilmResult struct {
	Film        *ControllerFilm
	LocalFilmId int
}

type ControllerFilmResult struct {
	Title        string `json:"title"`
	EpisodeId    int    `json:"episode_id"`
	OpeningCrawl string `json:"opening_crawl"`
	Director     string `json:"director"`
	Producer     string `json:"producer"`
	ReleaseDate  string `json:"release_date"`
}

func FromDomain(modelResponse *domainModel.CharacterSearchResult) *ControllerCharacterSearchResult {

	characters := make([]ControllerCharacter, len(modelResponse.Character))

	for i, ch := range modelResponse.Character {
		filmsDetail := make([]ControllerFilm, len(ch.FilmsDetail))
		for i, fd := range ch.FilmsDetail {
			filmsDetail[i] = ControllerFilm{
				Title:        fd.Title,
				EpisodeId:    fd.EpisodeId,
				OpeningCrawl: fd.OpeningCrawl,
				Director:     fd.Director,
				Producer:     fd.Producer,
				ReleaseDate:  fd.ReleaseDate,
			}
		}

		characters[i] = ControllerCharacter{
			Name:        ch.Name,
			Height:      ch.Height,
			Mass:        ch.Mass,
			HairColor:   ch.HairColor,
			SkinColor:   ch.SkinColor,
			EyeColor:    ch.EyeColor,
			BirthYear:   ch.BirthYear,
			Gender:      ch.Gender,
			Homeworld:   ch.Homeworld,
			Films:       ch.Films,
			FilmsDetail: filmsDetail,
			Species:     ch.Species,
			Vehicles:    ch.Vehicles,
			Starships:   ch.Starships,
			Created:     ch.Created,
			Edited:      ch.Edited,
			Url:         ch.Url,
		}
	}

	return &ControllerCharacterSearchResult{
		Count:     modelResponse.Count,
		Next:      modelResponse.Next,
		Previous:  modelResponse.Previous,
		Character: characters,
	}
}
