package commands

import (
	"fmt"
	"testing"
	"time"

	"github.com/ScholarlyKiwi/pokedex/internal/pokecache"
)

func TestCleanInput(t *testing.T) {

	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " hello world too long",
			expected: []string{"hello", "world", "too", "long"},
		},
		{
			input:    " hello world test",
			expected: []string{"hello", "world", "test"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("Slice Length Mismatch: %d != %d", len(actual), len(c.expected))
		} else {
			for i := range actual {
				word := actual[i]
				expectedWord := c.expected[i]
				if word != expectedWord {
					t.Errorf("Unexpected Word return: \nExpected: %s\nActual: %s", expectedWord, word)
				}
			}
		}
	}
}

func TestGetLocation(t *testing.T) {

	var config CmdConfig
	config.Map_next = "https://pokeapi.co/api/v2/location-area/"

	newCache, err := pokecache.NewCache(time.Second * 5)
	if err != nil {
		fmt.Printf("Error creating cache: %v", err)
	}
	config.Cache = newCache

	result, err := getPokeLocationAreas("https://pokeapi.co/api/v2/location-area/", &config)

	if err != nil {
		fmt.Printf("getPokeLocationAreas returned Error: %s\n", err)
	}
	fmt.Printf("getPokeLocationAreas return data: %v\n", result)
	for key, value := range result.Results {
		fmt.Printf("getPokeLocationAreas return results: %v: %v, %v\n", key, value, result.Results[key].Url)
	}
}

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache, _ := pokecache.NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Second
	const waitTime = baseTime + (5 * time.Second)
	cache, _ := pokecache.NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestReapLoopNo(t *testing.T) {
	const baseTime = 5 * time.Second
	const waitTime = baseTime - 5*time.Millisecond
	cache, _ := pokecache.NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key still")
		return
	}
}

func TestReapLoopZero(t *testing.T) {
	const baseTime = 5 * time.Second
	const waitTime = baseTime - 0*time.Millisecond
	cache, _ := pokecache.NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key still")
		return
	}
}

func TestGetPokemon(t *testing.T) {
	var config CmdConfig

	cases := []struct {
		name string
	}{
		{
			name: "tentacool",
		},
		{
			name: "shellos",
		},
	}

	newCache, err := pokecache.NewCache(time.Second * 5)
	if err != nil {
		fmt.Printf("Error creating cache: %v", err)
	}
	config.Cache = newCache

	for _, name := range cases {

		result, err := getpokePokemon(name.name, &config)

		if err != nil {
			t.Errorf("getpokePokemon returned Error: %s\n", err)
		}
		result_name := result.Name
		if result_name != name.name {
			t.Errorf("PokeAPI returned Pokemon named %v instead of %v.", result_name, name)
		}
	}
}
