package main

import (
	"os"
)

// Storage directory for skin cache
const skinCacheDir = "cached_skins"

// Main entry point
func main() {
	cleanupCache()
	os.MkdirAll(skinCacheDir, os.ModePerm)
	initRedis()
	initHttp()

	go startCleanupRoutine()
}
