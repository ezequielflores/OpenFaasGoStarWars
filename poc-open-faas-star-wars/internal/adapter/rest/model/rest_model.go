package model

import "handler/function/internal/application/model"

type RestFilmResult struct {
	Title        string `json:"title"`
	EpisodeId    int    `json:"episode_id"`
	OpeningCrawl string `json:"opening_crawl"`
	Director     string `json:"director"`
	Producer     string `json:"producer"`
	ReleaseDate  string `json:"release_date"`
}

func (r RestFilmResult) ToDomain() *model.Film {
	return &model.Film{
		Title:        r.Title,
		EpisodeId:    r.EpisodeId,
		OpeningCrawl: r.OpeningCrawl,
		Director:     r.Director,
		Producer:     r.Producer,
		ReleaseDate:  r.ReleaseDate,
	}
}

type RestCharacter struct {
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
	FilmsDetail []RestFilmResult `json:"films_details"`
	Species     []string         `json:"species"`
	Vehicles    []string         `json:"vehicles"`
	Starships   []string         `json:"starships"`
	Created     string           `json:"created"`
	Edited      string           `json:"edited"`
	Url         string           `json:"url"`
}

type RestCharacterSearchResult struct {
	Count     int             `json:"count"`
	Next      string          `json:"next"`
	Previous  string          `json:"previous"`
	Character []RestCharacter `json:"results"`
}

func (c *RestCharacterSearchResult) ToDomain() (*model.CharacterSearchResult, error) {

	characters := make([]model.Character, len(c.Character))
	for i, c := range c.Character {
		characters[i] = model.Character{
			Name:        c.Name,
			Height:      c.Height,
			Mass:        c.Mass,
			HairColor:   c.HairColor,
			SkinColor:   c.SkinColor,
			EyeColor:    c.EyeColor,
			BirthYear:   c.BirthYear,
			Gender:      c.Gender,
			Homeworld:   c.Homeworld,
			Films:       c.Films,
			FilmsDetail: []model.Film{},
			Species:     c.Species,
			Vehicles:    c.Vehicles,
			Starships:   c.Starships,
			Created:     c.Created,
			Edited:      c.Edited,
			Url:         c.Url,
		}
	}

	return &model.CharacterSearchResult{
		Count:     c.Count,
		Next:      c.Next,
		Previous:  c.Previous,
		Character: characters,
	}, nil
}
