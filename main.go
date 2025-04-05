package main

import (
	"os"
)

// Storage directory for skin cache
const skinCacheDir = "cached_skins"
const renderDir = "cached_skins/render"

// Main entry point
func main() {
	os.MkdirAll(skinCacheDir, os.ModePerm)
	os.MkdirAll(renderDir, os.ModePerm)
	cleanupCache()
	initRedis()
	initHttp()

	go startCleanupRoutine()
}
