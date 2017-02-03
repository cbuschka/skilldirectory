package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"skilldirectory/errors"
	"skilldirectory/model"
)

type UsersController struct {
	*BaseController
}

func (c UsersController) Base() *BaseController {
	return c.BaseController
}

func (c UsersController) Get() error {
	return fmt.Errorf("GET requests nor currently supported.")
}

func (c UsersController) Post() error {
	return c.authenticateUser()
}

func (c UsersController) Delete() error {
	return fmt.Errorf("DELETE requests nor currently supported.")
}

func (c UsersController) Put() error {
	return fmt.Errorf("PUT requests nor currently supported.")
}

func (c *UsersController) authenticateUser() error {
	body, _ := ioutil.ReadAll(c.r.Body)

	// Create a new User instance and unmarshal the request data into it
	user := model.User{}
	err := json.Unmarshal(body, &user)
	if err != nil {
		return errors.MarshalingError(err)
	}

	if err = c.validatePOSTBody(&user); err != nil {
		return err
	}

	acc, err := c.getUserAccount(&user)
	if err != nil {
		return err
	}

	b, err := json.Marshal(acc)
	if err != nil {
		return errors.MarshalingError(err)
	}

	c.w.Write(b)

	return nil
}

func (c *UsersController) getUserAccount(user *model.User) (model.UserAccount, error) {
	if user.Login != "test" || user.Password != "test" {
		return model.UserAccount{}, errors.InvalidLoginData(fmt.Errorf("Invalid login data provided"))
	}
	return model.UserAccount{Login: "test", DisplayName: "Foo Bar"}, nil
}

func (c *UsersController) validatePOSTBody(user *model.User) error {
	if user.Login == "" || user.Password == "" {
		return errors.IncompletePOSTBodyError(fmt.Errorf(
			"%q and %q fields must be non-empty", "Login", "Password"))
	}
	return nil
}
