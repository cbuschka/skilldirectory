package model

// SkillIcon models a small (1 MB max) graphical image for a Skill
type SkillIcon struct {
	SkillID string `json:"id"`
	URL     string `json:"url"`
}

// NewSkillIcon returns a new SkillIcon with specified params
func NewSkillIcon(skillID, url string) SkillIcon {
	return SkillIcon{
		SkillID: skillID,
		URL:     url,
	}
}

// GetType satisfies data.ReadAllInterface
func (i SkillIcon) GetType() interface{} {
	return SkillIcon{}
}
