package model

type Link struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	SkillID  string `json:"skill_id"`
	LinkType string `json:"link_type"`
}

const (
	BlogLinkType     = "blog"
	TutorialLinkType = "tutorial"
	WebpageLinkType  = "webpage"
)

func NewLink(name, url, skillID, linkType string) Link {
	return Link{
		Name:     name,
		URL:      url,
		SkillID:  skillID,
		LinkType: linkType,
	}
}

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

func (l Link) GetType() interface{} {
	return Link{}
}
