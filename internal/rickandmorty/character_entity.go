package rickandmorty

import (
	"strconv"
	"strings"
)

type Location struct {
	Name string
	Url  string
}

type Origin struct {
	Name string
	Url  string
}

type EpisodeURL = string

type CharacterID = int

type Character struct {
	ID       CharacterID
	Name     string
	Species  string
	Type     string
	Gender   string
	Image    string
	Url      string
	Created  string
	Origin   Origin       `json:"-"`
	Location Location     `json:"-"`
	Episode  []EpisodeURL `json:"-"`
}

type Info struct {
	Count int
	Pages int
	Next  *string
	Prev  *string
}

type API struct {
	Info    Info
	Results []Character
}

func (c Character) ToCSVEntry() string {
	entry := []string{
		strconv.Itoa(c.ID),
		c.Name,
		c.Species,
		c.Type,
		c.Gender,
		c.Image,
		c.Url,
		c.Created,
	}
	return strings.Join(entry, ",")
}
