package rss

import (
	"os"
	"testing"
)

func TestMediaGroup(t *testing.T) {
	data, err := os.ReadFile("testdata/media_group_test")
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

	if item.Title.String() != "title!" {
		t.Errorf("Expected title 'title!', got %s", item.Title)
	}

	if item.Link != "https://www.video.com/example-video" {
		t.Errorf("Expected link 'https://www.video.com/example-video', got %s", item.Link)
	}

	if len(item.Metadata.MediaGroups) != 1 {
		t.Errorf("Expected 1 content, got %d", len(item.Metadata.MediaGroups))
	}

	media := item.Metadata.MediaGroups[0]

	mediaLink := media.Content[0]

	if mediaLink.URL != "https://www.videos.com/video.mp4" {
		t.Errorf("Expected URL 'https://www.videos.com/video.mp4', got %s", mediaLink.URL)
	}

	if mediaLink.Type != "video/mp4" {
		t.Errorf("Expected type 'video/mp4', got %s", mediaLink.Type)
	}

	if mediaLink.Width != "640" {
		t.Errorf("Expected width 640, got %s", mediaLink.Width)
	}

	if mediaLink.Height != "390" {
		t.Errorf("Expected height 390, got %s", mediaLink.Height)
	}

	thumbnail := media.Thumbnail

	if thumbnail.URL != "https://www.videos.com/video_thumb.jpg" {
		t.Errorf("Expected URL 'https://www.videos.com/video_thumb.jpg', got %s", thumbnail.URL)
	}

	if thumbnail.Width != "400" {
		t.Errorf("Expected width 400, got %s", thumbnail.Width)
	}

	if thumbnail.Height != "300" {
		t.Errorf("Expected height 300, got %s", thumbnail.Height)
	}

	if media.Description == nil || media.Description.Content != "description!" {
		t.Errorf("Expected description 'description!', got %v", media.Description)
	}
}

func TestMediaContent(t *testing.T) {
	data, err := os.ReadFile("testdata/media_content_test")
	if err != nil {
		t.Fatalf("Reading file: %v", err)
	}

	feed, err := Parse(data, ParseOptions{})
	if err != nil {
		t.Fatalf("Parsing %v", err)
	}

	if len(feed.Items) != 2 {
		t.Fatalf("Expected 2 item, got %d", len(feed.Items))
	}

	content := feed.Items[0].MediaContents()
	if content != "" {
		t.Errorf("item 1: content was wrong: %s", content)
	}

	content = feed.Items[1].MediaContents()
	if content != "media:description" {
		t.Errorf("item 2: content was wrong: '%s'", content)
	}

}
