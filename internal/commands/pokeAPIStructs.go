package commands

type pokeResourceList struct {
	Count    int                    `json:"count"`
	Next     string                 `json:"next"`
	Previous string                 `json:"previous"`
	Results  []pokeNamedApiResource `json:"results"`
}

type pokeNamedApiResource struct {
	Name string `json:"name"`
	Url  string `json:"utl"`
}

type pokeEncounter struct {
	MinLevel        int                    `json:"min_level"`
	MaxLevel        int                    `json:"max_level"`
	ConditionValues []pokeNamedApiResource `json:"condition_values"`
	Chance          int                    `json:"chance"`
	Method          pokeNamedApiResource   `json:"method"`
}

type pokeVersionEnctounerDetail struct {
	Version           pokeNamedApiResource `json:"version"`
	MaxChance         int                  `json:"max_chance"`
	encounter_details []pokeEncounter
}

type pokePokemonEncounter struct {
	Pokemon        pokeNamedApiResource         `json:"pokemon"`
	VersionDetails []pokeVersionEnctounerDetail `json:"version_details"`
}

type pokeNames struct {
	Name     string               `json:"name"`
	Language pokeNamedApiResource `json:"language"`
}

type pokeEncounterMethodRates struct {
	EncounterMethod pokeNamedApiResource         `json:"encounter_method"`
	VersionDetails  []pokeVersionEnctounerDetail `json:"version_details"`
}

type pokeLocationArea struct {
	Id                   int                        `json:"id"`
	Name                 string                     `json:"name"`
	GameIndex            int                        `json:"game_index"`
	EncounterMethodRates []pokeEncounterMethodRates `json:"EncounterMethodRates"`
	Location             pokeNamedApiResource       `json:"location"`
	Names                []pokeNames                `json:"names"`
	PokemonEncounters    []pokePokemonEncounter     `json:"pokemon_encounters"`
}

type pokePokemonStat struct {
	Stat     pokeNamedApiResource `json:"stat"`
	Effort   int                  `json:"effort"`
	BaseStat int                  `json:"base_stat"`
}

type pokePokemonType struct {
	Slot int                  `json:"slot"`
	Type pokeNamedApiResource `json:"type"`
}

type pokePokemon struct {
	Id             int               `json:"id"`
	Name           string            `json:"name"`
	BaseExperience int               `json:"base_experience"`
	Height         int               `json:"height"`
	IsDefault      bool              `json:"is_default"`
	Order          int               `json:"order"`
	Weight         int               `json:"weight"`
	Stats          []pokePokemonStat `json:"stats"`
	Types          []pokePokemonType `json:"types"`
}
