package main

import (
	"encoding/json"
	"os"
)

// Example
// {"url": "https://example.com/video/%s.ts", "folder": "Video"}
type video struct {
	Url    string `json:"url"`
	Folder string `json:"folder"`
}

func parseConfig(path string) (videoItems []video, err error) {
	coursesFile, err := os.ReadFile(*configFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(coursesFile, &videoItems)
	if err != nil {
		return
	}

	return
}
