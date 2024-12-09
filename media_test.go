package rss

import (
	"os"
	"testing"
)

func TestMediaGroup(t *testing.T) {

	data, err := os.ReadFile("testdata/media_test")
	if err != nil {
		t.Fatalf("Reading file: %v", err)
	}

	feed, err := Parse(data, ParseOptions{})
	if err != nil {
		t.Fatalf("Parsing %v", err)
	}

	if len(feed.Items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(feed.Items))
	}

	item := feed.Items[0]

	if item.Title != "title!" {
		t.Errorf("Expected title 'title!', got %s", item.Title)
	}

	if item.Link != "https://www.video.com/example-video" {
		t.Errorf("Expected link 'https://www.video.com/example-video', got %s", item.Link)
	}

	if len(item.Media.Groups) != 1 {
		t.Errorf("Expected 1 content, got %d", len(item.Media.Groups))
	}

	media := item.Media.Groups[0]

	mediaLink := media.Media[0]

	if mediaLink.URL != "https://www.videos.com/video.mp4" {
		t.Errorf("Expected URL 'https://www.videos.com/video.mp4', got %s", mediaLink.URL)
	}

	if mediaLink.Type != "video/mp4" {
		t.Errorf("Expected type 'video/mp4', got %s", mediaLink.Type)
	}

	if mediaLink.Width != 640 {
		t.Errorf("Expected width 640, got %d", mediaLink.Width)
	}

	if mediaLink.Height != 390 {
		t.Errorf("Expected height 390, got %d", mediaLink.Height)
	}

	thumbnail := media.Thumbnail

	if thumbnail.URL != "https://www.videos.com/video_thumb.jpg" {
		t.Errorf("Expected URL 'https://www.videos.com/video_thumb.jpg', got %s", thumbnail.URL)
	}

	if thumbnail.Width != 400 {
		t.Errorf("Expected width 400, got %d", thumbnail.Width)
	}

	if thumbnail.Height != 300 {
		t.Errorf("Expected height 300, got %d", thumbnail.Height)
	}

	if media.Description != "description!" {
		t.Errorf("Expected description 'description!', got %s", media.Description)
	}

}
