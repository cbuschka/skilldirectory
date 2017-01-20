package model

// SkillIcon models a small (1 MB max) graphical image for a Skill
type SkillIcon struct {
	ID      string `json:"id'`
	SkillID string `json:"skill_id"`
	Icon    []byte `json:"icon"` // contains image's constituent bytes
}

// NewSkillIcon returns a new SkillIcon with specified params
func NewSkillIcon(id, skillID string, icon []byte) SkillIcon {
	return SkillIcon{
		ID:      id,
		SkillID: skillID,
		Icon:    icon,
	}
}

// isValidSize returns true if number of bytes in receiver's
// SkillIcon field is <= 1 MB, returns false otherwise.
func (i *SkillIcon) isValidSize() bool {
	if len(i.Icon) <= 1000000 {
		return true
	}
	return false
}

// GetType satisfies data.ReadAllInterface
func (i *SkillIcon) GetType() interface{} {
	return SkillIcon{}
}
