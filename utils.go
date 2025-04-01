package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

// Helper method to convert images into the format expected by the skin renderer
func convertToNRGBA(img image.Image) *image.NRGBA {
	if res, ok := img.(*image.NRGBA); ok {
		return res
	}

	res := image.NewNRGBA(img.Bounds())
	draw.Draw(res, img.Bounds(), img, image.Pt(0, 0), draw.Src)

	return res
}

// Download the skin from a URL into the cache, for faster results and rate limit preventions
func downloadAndCacheSkin(skinURL string, cachePath string, uuid string) (string, error) {
	// Download the image
	skinResp, err := http.Get(skinURL)
	if err != nil || skinResp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch skin image")
	}
	defer skinResp.Body.Close()

	// Create the local cache
	file, err := os.Create(cachePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Decode the image and re-encode to PNG format
	img, _, err := image.Decode(skinResp.Body)
	if err != nil {
		return "", err
	}
	png.Encode(file, img)

	cacheSkin(uuid, cachePath)
	return cachePath, nil
}

// Helper method to try and convert a username into a UUID
func getUUID(username string) (string, error) {
	cachedUUID, err := redisClient.Get(ctx, "username:"+username).Result()
	if err == nil {
		return cachedUUID, nil
	}

	mojangURL := fmt.Sprintf("https://api.mojang.com/users/profiles/minecraft/%s", username)
	resp, err := http.Get(mojangURL)
	if err != nil || resp.StatusCode != 200 {
		return "", err
	}
	defer resp.Body.Close()

	var profile MojangProfile
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &profile)

	cacheUsername(username, profile.ID)
	return profile.ID, nil
}

// Helper method to check if the supplied ID matches a UUID
func isValidUUID(identifier string) bool {
	if len(identifier) == 36 {
		re := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
		return re.MatchString(identifier)
	}
	if len(identifier) == 32 {
		re := regexp.MustCompile(`^[0-9a-fA-F]{32}$`)
		return re.MatchString(identifier)
	}
	return false
}

// Helper method to check if the supplied ID matches a SHA256 hash (or texture ID in this case)
func isValidSHA256Hash(identifier string) bool {
	if len(identifier) == 64 {
		re := regexp.MustCompile(`^[0-9a-f]{64}$`)
		return re.MatchString(identifier)
	}
	return false
}

// Fallback method to ensure steve skin can be returned.
func generateOfflineUUID(username string) uuid.UUID {
	offlineUUIDStr := "OfflinePlayer:" + username

	hash := md5.New()
	io.WriteString(hash, offlineUUIDStr)
	md5Hash := hash.Sum(nil)
	md5Hash[6] = (md5Hash[6] & 0x0f) | 0x30
	md5Hash[8] = (md5Hash[8] & 0x3f) | 0x80
	return uuid.UUID(md5Hash)
}
