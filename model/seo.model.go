package model

// Seo ...
type Seo struct {
	Title       string `json:"title" bson:"title,omitempty"`
	Description string `json:"description" bson:"description,omitempty"`
	Image       string `json:"image" bson:"image,omitempty"`
	Keywords    string `json:"keywords" bson:"keywords,omitempty"`
}
