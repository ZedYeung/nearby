package model

import "sync"

type Post struct {
	Mu       sync.Mutex
	User     string   `json:"user"`
	Message  string   `json:"message"`
	Location Location `json:"location"`
	URLs     []string `json:"url"`
}
