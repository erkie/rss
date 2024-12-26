package rss

import "strings"

/*
Example:
<!-- inline -->
<media:content url="https://example.com/image.jpeg" type="image/jpeg" fileSize="362243"
	medium="image">
	<media:rating scheme="urn:simple">nonadult</media:rating>
	<media:description type="plain">media:description</media:description>
</media:content>

<!-- grouped -->
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

type MediaContent struct {
	URL    string `xml:"url,attr"`
	Type   string `xml:"type,attr"`
	Width  string `xml:"width,attr"`
	Height string `xml:"height,attr"`
	Medium string `xml:"medium,attr"`

	Rating      *MediaRating      `xml:"rating"`
	Description *MediaDescription `xml:"description"`
}

func (m MediaContent) IsImage() bool {
	return strings.HasPrefix("image/", strings.ToLower(m.Type)) || strings.ToLower(m.Medium) == "image"
}

type MediaGroup struct {
	Content     []*MediaContent   `xml:"content"`
	Thumbnail   *MediaContent     `xml:"thumbnail"`
	Description *MediaDescription `xml:"description"`
	Community   *MediaCommunity   `xml:"community"`
}

type MediaDescription struct {
	Type    string `xml:"string,attr"`
	Content string `xml:",chardata"`
}

type MediaThumbnail struct {
	URL    string `xml:"url,attr"`
	Type   string `xml:"type,attr"`
	Width  string `xml:"width,attr"`
	Height string `xml:"height,attr"`
}

type MediaCommunity struct {
	StarRating *MediaStarRating `xml:"starRating"`
	Statistics *MediaStatistics `xml:"statistics"`
}

type MediaRating struct {
	Scheme string `xml:"scheme,attr"`
}

type MediaStarRating struct {
	Count   string `xml:"count,attr"`
	Average string `xml:"average,attr"`
	Min     string `xml:"min,attr"`
	Max     string `xml:"max,attr"`
}

type MediaStatistics struct {
	Views string `xml:"views,attr"`
}
