package models

// Media represents a media element that can be either an image or a video.
// Only one of the fields (Image or Video) will be populated at a time.
type Media struct {
	// Image contains metadata for an image element
	Image *ImageMedia `json:"image,omitempty"`
	// Video contains metadata for a video element
	Video *VideoMedia `json:"video,omitempty"`
}

// ImageMedia contains the metadata and properties of an image.
type ImageMedia struct {
	// Key is the unique identifier for the image in the Articulate system
	Key string `json:"key"`
	// Type indicates the image format (jpg, png, etc.)
	Type string `json:"type"`
	// Width is the pixel width of the image
	Width int `json:"width,omitempty"`
	// Height is the pixel height of the image
	Height int `json:"height,omitempty"`
	// CrushedKey is the identifier for a compressed version of the image
	CrushedKey string `json:"crushedKey,omitempty"`
	// OriginalURL is the URL to the full-resolution image
	OriginalURL string `json:"originalUrl"`
	// UseCrushedKey indicates whether to use the compressed version
	UseCrushedKey bool `json:"useCrushedKey,omitempty"`
}

// VideoMedia contains the metadata and properties of a video.
type VideoMedia struct {
	// Key is the unique identifier for the video in the Articulate system
	Key string `json:"key"`
	// URL is the direct link to the video content
	URL string `json:"url"`
	// Type indicates the video format (mp4, webm, etc.)
	Type string `json:"type"`
	// Poster is the URL to the static thumbnail image for the video
	Poster string `json:"poster,omitempty"`
	// Duration is the length of the video in seconds
	Duration int `json:"duration,omitempty"`
	// InputKey is the original identifier for uploaded videos
	InputKey string `json:"inputKey,omitempty"`
	// Thumbnail is the URL to a smaller preview image
	Thumbnail string `json:"thumbnail,omitempty"`
	// OriginalURL is the URL to the source video file
	OriginalURL string `json:"originalUrl"`
}
