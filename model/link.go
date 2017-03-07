package model

import "github.com/jinzhu/gorm"

// Link has a many-to-one relationship to a Skill
type Link struct {
	gorm.Model
	Name     string `json:"name"`
	URL      string `json:"url"`
	SkillID  uint   `gorm:"index" json:"skill_id"`
	LinkType string `json:"link_type"`
}

const (
	BlogLinkType          = "blog"           // BlogLinkType is a blog enum
	TutorialLinkType      = "tutorial"       // TutorialLinkType is a tutorial enum
	WebpageLinkType       = "webpage"        // WebpageLinkType is a webpage enum
	DeveloperToolLinkType = "developer-tool" // DeveloperToolLinkType is a developer-tool enum
)

// NewLink is a Link constructor
func NewLink(id, skillID uint, name, url, linkType string) Link {
	link := Link{
		Name:     name,
		URL:      url,
		SkillID:  skillID,
		LinkType: linkType,
	}
	link.ID = id
	return link
}

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

func (l Link) GetID() uint {
	return l.ID
}

//GetType returns the implemented type
func (l Link) GetType() interface{} {
	return Link{}
}

func QueryLink(id uint) Link {
	var link Link
	link.ID = id
	return link
}
