package main

// Data structures for data returned from mojang
type MojangProfile struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Properties []SkinProp `json:"properties"`
}

type SkinProp struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SkinTextures struct {
	Timestamp   int64  `json:"timestamp"`
	ProfileID   string `json:"profileId"`
	ProfileName string `json:"profileName"`
	Textures    struct {
		Skin struct {
			URL string `json:"url"`
		} `json:"SKIN"`
	} `json:"textures"`
}
