package rss

/*
Example:
<media:group>
   <media:content url="https://www.videos.com/video.mp4" type="application/mp4" width="300" height="200"/>
   <media:thumbnail url="https://www.videos.com/video.jpg" width="400" height="300"/>
   <media:description>This is a video</media:description>
   <media:community>
    <media:starRating count="12756" average="5.00" min="1" max="5"/>
    <media:statistics views="275138"/>
   </media:community>
</media:group>
*/

type Media struct {
	Groups []MediaGroup `xml:"group"`
}

func (m *Media) IsPresent() bool {
	return len(m.Groups) > 0
}

func (m *Media) Description() string {
	for _, media := range m.Groups {
		if media.Description != "" {
			return media.Description
		}
	}
	return ""
}

type MediaItem struct {
	URL    string `xml:"url,attr"`
	Type   string `xml:"type,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

type MediaGroup struct {
	Media       []*MediaItem    `xml:"content"`
	Thumbnail   *MediaItem      `xml:"thumbnail"`
	Description string          `xml:"description"`
	Community   *MediaCommunity `xml:"community"`
}

type MediaCommunity struct {
	StarRating *MediaStarRating `xml:"starRating"`
	Statistics *MediaStatistics `xml:"statistics"`
}

type MediaStarRating struct {
	Count   int     `xml:"count,attr"`
	Average float64 `xml:"average,attr"`
	Min     int     `xml:"min,attr"`
	Max     int     `xml:"max,attr"`
}

type MediaStatistics struct {
	Views int `xml:"views,attr"`
}
