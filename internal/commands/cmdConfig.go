package commands

import "github.com/ScholarlyKiwi/pokedex/internal/pokecache"

type CmdConfig struct {
	Input         []string
	Map_next      string
	Map_prev      string
	Cache         *pokecache.Cache
	CaughtPokemon map[string]any
}
