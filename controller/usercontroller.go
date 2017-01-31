package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"skilldirectory/errors"
	"skilldirectory/model"
)

type UserController struct {
	*BaseController
}

func (c UserController) Base() *BaseController {
	return c.BaseController
}

func (c UserController) Get() error {
	return fmt.Errorf("GET requests nor currently supported.")
}

func (c UserController) Post() error {
	return c.authenticateUser()
}

func (c UserController) Delete() error {
	return fmt.Errorf("DELETE requests nor currently supported.")
}

func (c UserController) Put() error {
	return fmt.Errorf("PUT requests nor currently supported.")
}

func (c *UserController) authenticateUser() error {
	body, _ := ioutil.ReadAll(c.r.Body)
	user := model.User{}
	err := json.Unmarshal(body, &user)
	if err != nil {
		return errors.MarshalingError(err)
	}

	c.Printf("user: %q; pass: %q", user.Login, user.Password)

	err = c.validatePOSTBody(&user)
	if err != nil {
		return err
	}

	b, err := json.Marshal(user)
	if err != nil {
		return errors.MarshalingError(err)
	}
	c.w.Write(b)

	return nil
}

func (c *UserController) validatePOSTBody(user *model.User) error {
	if user.Login == "" || user.Password == "" {
		return errors.IncompletePOSTBodyError(fmt.Errorf(
			"%q and %q fields must be non-empty", "Login", "Password"))
	}
	return nil
}
