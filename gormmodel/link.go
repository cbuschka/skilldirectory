package gormmodel

import "github.com/jinzhu/gorm"

// Link has a many-to-one relationship to a Skill
type Link struct {
	gorm.Model
	Name     string `json:"name"`
	URL      string `json:"url"`
	SkillID  uint   `gorm:"index"`
	LinkType string `json:"link_type"`
}

const (
	BlogLinkType          = "blog"           // BlogLinkType is a blog enum
	TutorialLinkType      = "tutorial"       // TutorialLinkType is a tutorial enum
	WebpageLinkType       = "webpage"        // WebpageLinkType is a webpage enum
	DeveloperToolLinkType = "developer-tool" // DeveloperToolLinkType is a developer-tool enum
)

// NewLink is a Link constructor
// func NewLink(id, name, url, skillID, linkType string) Link {
// 	return Link{
// 		ID:       id,
// 		Name:     name,
// 		URL:      url,
// 		SkillID:  skillID,
// 		LinkType: linkType,
// 	}
// }

// IsValidLinkType is a switch that validates a give linkType string
func IsValidLinkType(linkType string) bool {
	switch linkType {
	case
		BlogLinkType,
		TutorialLinkType,
		WebpageLinkType,
		DeveloperToolLinkType:
		return true
	}
	return false
}

//GetType returns the implemented type
func (l Link) GetType() interface{} {
	return Link{}
}
