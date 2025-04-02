package main

import (
	"os"
)

// Storage directory for skin cache
const skinCacheDir = "cached_skins"

// Main entry point
func main() {
	os.MkdirAll(skinCacheDir, os.ModePerm)
	cleanupCache()
	initRedis()
	initHttp()

	go startCleanupRoutine()
}
