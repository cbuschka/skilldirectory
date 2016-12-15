package model

type Link struct {
	Name     string
	URL      string
	SkillID  string
	LinkType string
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
