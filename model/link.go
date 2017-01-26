package model

type Link struct {
	ID       string `json:"id"`
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

func NewLink(id, name, url, skillID, linkType string) Link {
	return Link{
		ID:       id,
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
