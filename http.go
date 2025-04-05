package main

import (
	"fmt"
	"github.com/chai2010/webp"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// Init the HTTP (or rest api) server
func initHttp() {
	r := mux.NewRouter()
	r.HandleFunc("/{type}/{id}", handleAvatar).Methods("GET")

	http.ListenAndServe(":8080", r)
}

// Handle the incoming avatar request.
func handleAvatar(w http.ResponseWriter, r *http.Request) {
	// Ratelimiter
	ip := getIP(r)
	limiter := getRateLimiter(ip)

	if !limiter.Allow() {
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}

	// Default values passed from URL
	vars := mux.Vars(r)
	identifier := vars["id"]
	mode := vars["type"]
	scaleStr := r.URL.Query().Get("scale")

	// Double check that the requested render is supported
	if mode != "isometric" && mode != "body" && mode != "avatar" && mode != "head" {
		http.Error(w, "Unsupported render mode. Valid render modes are isometric, head, body & avatar", http.StatusBadRequest)
		return
	}

	// Initialize default scale
	scale := 512

	// Parse scale from string to int, if possible
	if scaleStr != "" {
		if parsedScale, err := strconv.Atoi(scaleStr); err == nil {
			scale = parsedScale
		} else {
			fmt.Println("Invalid scale value, using default")
		}
	}

	// To prevent overload on the server, we enforce some limits on scaling.
	// Not smaller than 16, so that it's visible, and not more than 1024px
	if scale < 16 {
		scale = 16
	} else if scale > 1024 {
		scale = 1024
	}

	var uuid string
	var err error

	// Check if the supplied ID is a valid UUID
	if isValidUUID(identifier) {
		// We don't need the - in the UUID, so we remove it
		uuid = strings.ReplaceAll(identifier, "-", "")

		// Check if the supplied ID is a texture hash
	} else if isValidSHA256Hash(identifier) {
		uuid = identifier
	} else {

		// Supplied ID was likely a username. So we try to resolve it
		uuid, err = getUUID(identifier)
		if err != nil || uuid == "" {
			uuid = strings.ReplaceAll(generateOfflineUUID(identifier).String(), "-", "")
		}
	}

	cachePath := path.Join(renderDir, fmt.Sprintf("%s_%s_%s.web", mode, uuid, strconv.Itoa(scale)))
	cachedFile, err := os.Open(cachePath)
	if err == nil {
		w.Header().Set("Content-Type", "image/webp")
		_, err = io.Copy(w, cachedFile)
		if err != nil {
			http.Error(w, "Failed to read image: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer cachedFile.Close()

	// Request the skin from the MOJANG servers
	skinPath, err := fetchSkin(uuid)
	if err != nil {
		skinPath = "fallback.png"
	}

	// Render the skin for the API
	img, err := renderSkin(skinPath, mode, scale, uuid, true)
	if err != nil {
		http.Error(w, "Failed to render skin", http.StatusInternalServerError)
		return
	}

	// Encode the image ready for browser rendering
	f, _ := os.Create(cachePath)
	err = webp.Encode(f, img, &webp.Options{
		Lossless: true,
		Quality:  100,
		Exact:    true,
	})

	if err != nil {
		http.Error(w, "Failed to encode image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/webp")
	http.ServeContent(w, r, cachePath, time.Now(), f)
}
