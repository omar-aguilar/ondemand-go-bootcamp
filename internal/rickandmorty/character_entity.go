package rickandmorty

type CharacterID = int

type Character struct {
	ID      CharacterID `json:"id"`
	Name    string      `json:"name"`
	Species string      `json:"species"`
	Type    string      `json:"type"`
	Gender  string      `json:"gender"`
	Image   string      `json:"image"`
	Url     string      `json:"url"`
	Created string      `json:"created"`
}

type CharacterList = []Character

type API struct {
	Results CharacterList
}
