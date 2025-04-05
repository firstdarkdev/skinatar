package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/mineatar-io/skin-render"
	"image"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// Skin fetcher. Here, we try to fetch the skin from the cache, via mojangs servers, and finally, fall back to crafthead.
func fetchSkin(uuid string) (string, error) {
	cachePath := filepath.Join(skinCacheDir, uuid+".webp")
	if _, err := os.Stat(cachePath); err == nil {
		return cachePath, nil
	}

	// Texture hash was supplied, so we just pull it directly
	if isValidSHA256Hash(uuid) {
		return downloadAndCacheSkin(fmt.Sprintf("https://textures.minecraft.net/texture/%s", uuid), cachePath, uuid)
	}

	// Request the skin from mojang
	mojangURL := fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/profile/%s", uuid)
	resp, err := http.Get(mojangURL)
	if err != nil || resp.StatusCode != 200 {
		// Request failed. Either a fake account, or the servers are down, so we try with Crafthead
		return fetchSkinFromCraftHead(uuid, cachePath)
	}
	defer resp.Body.Close()

	// Decompile the mojang response, so that we can extract the texture data
	var profile MojangProfile
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &profile)

	if len(profile.Properties) == 0 {
		return "", fmt.Errorf("no properties found for UUID")
	}

	// Decode the BASE64 encoded texture url
	decoded, err := base64.StdEncoding.DecodeString(profile.Properties[0].Value)
	if err != nil {
		return "", fmt.Errorf("failed to decode Base64 skin data")
	}

	var skinData SkinTextures
	json.Unmarshal(decoded, &skinData)

	// Extract the skin URL and download it
	skinURL := skinData.Textures.Skin.URL
	if skinURL == "" {
		return "", fmt.Errorf("no skin URL found")
	}

	return downloadAndCacheSkin(skinURL, cachePath, uuid)
}

// Mojang failed us, so we try CraftHead instead
func fetchSkinFromCraftHead(uuid string, cachePath string) (string, error) {
	craftHeadURL := fmt.Sprintf("https://crafthead.net/skin/%s", uuid)
	return downloadAndCacheSkin(craftHeadURL, cachePath, uuid)
}

// Render the skin in the requested format
func renderSkin(skinPath string, mode string, scale int, uid string, overlay bool) (image.Image, error) {
	file, err := os.Open(skinPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := webp.Decode(file)
	if err != nil {
		return nil, err
	}

	switch mode {
	// Plain Avatar
	case "avatar":
		return skin.RenderFace(convertToNRGBA(img), skin.Options{
			Scale:   scale / 16,
			Slim:    skin.IsSlimFromUUID(uid),
			Overlay: overlay,
			Square:  true,
		}), nil

	// 3D head
	case "isometric", "head":
		return skin.RenderHead(convertToNRGBA(img), skin.Options{
			Scale:   scale / 16,
			Slim:    skin.IsSlimFromUUID(uid),
			Overlay: overlay,
			Square:  true,
		}), nil

	// Full Body
	case "body":
		return skin.RenderBody(convertToNRGBA(img), skin.Options{
			Scale:   scale / 39,
			Slim:    skin.IsSlimFromUUID(uid),
			Overlay: overlay,
			Square:  true,
		}), nil
	default:
		return img, nil
	}
}
