package main

import (
	"github.com/LukeHandle/minecraft"
)

const (
	// Get the skin size in bytes. Stored as a []uint8, one byte each,
	// plus bounces. So 64 * 64 bytes and we'll throw in an extra 16
	// bytes of overhead.
	SKIN_SIZE = (64 * 64) + 16

	// Define a 64 MB cache size.
	CACHE_SIZE = 2 << 25

	// Based off those, calculate the maximum number of skins we'll store
	// in memory.
	SKIN_NUMBER = CACHE_SIZE / SKIN_SIZE
)

// Cache object that stores skins in memory.
type CacheMemory struct {
	// Map of usernames to minecraft skins. Lookups here are O(1), so that
	// makes my happy.
	Skins map[string]minecraft.Skin
	// Additionally keep a *slice* of usernames which we can update
	Usernames []string
}

// Find the position of a string in a slice. Returns -1 on failure.
func indexOf(str string, list []string) int {
	for index, value := range list {
		if value == str {
			return index
		}
	}

	return -1
}

func (c *CacheMemory) setup() error {
	c.Skins = map[string]minecraft.Skin{}
	c.Usernames = []string{}

	log.Info("Loaded Memory cache")
	return nil
}

// Returns whether the item exists in the cache.
func (c *CacheMemory) has(username string) bool {
	if _, exists := c.Skins[username]; exists {
		return true
	} else {
		return false
	}
}

// Retrieves the item from the cache. We'll promote it to the "top" of the
// cache, effectively updating its expiry time.
func (c *CacheMemory) pull(username string) minecraft.Skin {
	index := indexOf(username, c.Usernames)
	c.Usernames = append(c.Usernames, username)
	c.Usernames = append(c.Usernames[:index], c.Usernames[index+1:]...)

	return c.Skins[username]
}

// Adds the skin to the cache, remove the oldest, expired skin if the cache
// list is full.
func (c *CacheMemory) add(username string, skin minecraft.Skin) {
	if len(c.Usernames) >= SKIN_NUMBER {
		first := c.Usernames[0]
		delete(c.Skins, first)
		c.Usernames = append(c.Usernames[1:], username)
	} else {
		c.Usernames = append(c.Usernames, username)
	}

	c.Skins[username] = skin
}

// The byte size of the cache. Fairly rough... don't really want to venture
// into the land of manual memory management, because there be dragons.
func (c *CacheMemory) memory() uint64 {
	return uint64(len(c.Usernames) * SKIN_SIZE)
}
