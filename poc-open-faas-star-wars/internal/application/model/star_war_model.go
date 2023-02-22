package model

type QueryFilter struct {
	Search string
	Page   string
}

type Film struct {
	Title        string
	EpisodeId    int
	OpeningCrawl string
	Director     string
	Producer     string
	ReleaseDate  string
}

type Character struct {
	Name        string
	Height      string
	Mass        string
	HairColor   string
	SkinColor   string
	EyeColor    string
	BirthYear   string
	Gender      string
	Homeworld   string
	Films       []string
	FilmsDetail []Film
	Species     []string
	Vehicles    []string
	Starships   []string
	Created     string
	Edited      string
	Url         string
}

type CharacterSearchResult struct {
	Count     int
	Next      string
	Previous  string
	Character []Character
}

type SearchFilmResult struct {
	Film        *Film
	LocalFilmId int
}
