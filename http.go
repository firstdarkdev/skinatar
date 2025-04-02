package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"image/png"
	"net/http"
	"strconv"
	"strings"
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
	w.Header().Set("Content-Type", "image/png")
	err = png.Encode(w, img)
	if err != nil {
		http.Error(w, "Failed to encode image", http.StatusInternalServerError)
		return
	}
}
