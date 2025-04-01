package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func cleanupCache() {
	now := time.Now()

	err := filepath.Walk(skinCacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			age := now.Sub(info.ModTime())

			if age > time.Hour {
				err := os.Remove(path)
				if err != nil {
					fmt.Printf("Failed to remove file %s: %v\n", path, err)
				} else {
					fmt.Printf("Removed expired file: %s\n", path)
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error while cleaning up cache: %v\n", err)
	}
}

func startCleanupRoutine() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cleanupCache()
		}
	}
}
