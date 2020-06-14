package core

import (
	"fmt"
	"time"
)

type Media struct {
	Title     string
	FilePath  string
	Uploader  string
	URL       string
	Thumbnail string
	Duration  time.Duration
}

func NewMedia(title, filePath, uploader, url, thumbnail string, durationInSeconds int) *Media {
	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", durationInSeconds))
	return &Media{
		Title:     title,
		FilePath:  filePath,
		Uploader:  uploader,
		URL:       url,
		Thumbnail: thumbnail,
		Duration:  duration,
	}
}
