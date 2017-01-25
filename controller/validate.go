package controller

import (
	"fmt"
	"skilldirectory/data"
	"skilldirectory/model"
)

// Return non-nil error if skillID doesn't point to valid record in database
func validateSkillID(skillID string, session data.DataAccess) error {
	err := session.Read("skills", skillID, &model.Skill{})
	if err != nil {
		return fmt.Errorf("failed to find record with this ID in database: %q", skillID)
	}
	return nil
}
