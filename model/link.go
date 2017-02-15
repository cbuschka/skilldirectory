package model

import "github.com/jinzhu/gorm"
// Link has a many-to-one relationship to a Skill
type Link struct {
	gorm.Model
	Name     string `json:"name"`
	URL      string `json:"url"`
	LinkType string `json:"link_type"`
}

const (
	BlogLinkType     = "blog"     //BlogLinkType is a blog enum
	TutorialLinkType = "tutorial" //TutorialLinkType is a tutorial enum
	WebpageLinkType  = "webpage"  //WebpageLinkType is a webpage enum
)

// NewLink is a Link constructor
func NewLink(name, url, linkType string) Link {
	return Link{
		Name:     name,
		URL:      url,
		LinkType: linkType,
	}
}

// IsValidLinkType is a switch that validates a give linkType string
func IsValidLinkType(linkType string) bool {
	switch linkType {
	case
		BlogLinkType,
		TutorialLinkType,
		WebpageLinkType:
		return true
	}
	return false
}

//GetType returns the implemented type
func (l Link) GetType() interface{} {
	return Link{}
}
