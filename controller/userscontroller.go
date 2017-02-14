package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

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
	return fmt.Errorf("GET requests not currently supported.")
}

func (c UsersController) Post() error {
	return c.authenticateUser()
}

func (c UsersController) Delete() error {
	return fmt.Errorf("DELETE requests not currently supported.")
}

func (c UsersController) Put() error {
	return fmt.Errorf("PUT requests not currently supported.")
}

func (c UsersController) Options() error {
	return c.handleOptionsRequest()
}

func (c *UsersController) authenticateUser() error {
	body, err := ioutil.ReadAll(c.r.Body)
	if err != nil {
		return err
	}
	c.Println(string(body))
	// Create a new User instance and unmarshal the request data into it
	credentials := model.AuthCredentials{}
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		return errors.MarshalingError(err)
	}

	if err = c.validatePOSTBody(&credentials); err != nil {
		return err
	}
	// Will execute the request to Github
	reqBody := url.Values{}
	reqBody.Set("client_id", credentials.Id)
	reqBody.Set("client_secret", credentials.Secret)
	reqBody.Set("code", credentials.Code)

	client := http.Client{}
	req, err := http.NewRequest("POST", "https://www.github.com/", bytes.NewBufferString(reqBody.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(reqBody.Encode())))

	response, err := client.Do(req)
	if err != nil {
		return err
	}

	tokenBody, _ := ioutil.ReadAll(response.Body)
	tokenResp := model.TokenResponse{}
	err = json.Unmarshal(tokenBody, &tokenResp)
	if err != nil {
		return err
	}

	c.Printf("Got token: %s", tokenResp.Token)
	b, err := json.Marshal(tokenResp)
	if err != nil {
		return err
	}

	c.w.Write(b)

	return nil
}

func (c *UsersController) handleOptionsRequest() error {
	c.w.Header().Set("Access-Control-Allow-Origin", "*")
	c.w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	c.w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, access-control-allow-methods")
	c.w.Write([]byte(""))
	return nil
}

func (c *UsersController) validatePOSTBody(credentials *model.AuthCredentials) error {
	if credentials.Id == "" || credentials.Code == "" || credentials.Secret == "" {
		return errors.IncompletePOSTBodyError(fmt.Errorf(
			"%q and %q fields must be non-empty", "Login", "Password"))
	}
	return nil
}
