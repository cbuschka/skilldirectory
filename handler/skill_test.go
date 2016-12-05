package handler

import (
  "testing"
  "fmt"
  "skilldirectory/data"
)

type FakeDataAccessor struct {}
func (f FakeDataAccessor)Save(s string, i interface{}) error {return nil}
func (f FakeDataAccessor) Read(s string, i interface{}) error {
  fmt.Println("READ")
  return nil
}
func (f FakeDataAccessor) Delete(s string) error {return nil}

func TestLoadSkill(t *testing.T) {
  fmt.Println("PASS")
  skillsConnector = data.NewAccessor(FakeDataAccessor{})
  _, err := loadSkill("1234")
  if err != nil {
    t.Errorf("Load skill error %s", err.Error())
  }
}
