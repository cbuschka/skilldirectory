package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	body, _ := ioutil.ReadAll(c.r.Body)

	// Get the access code and client_id from the request
	credentials := model.AuthCredentials{}
	err := json.Unmarshal(body, &credentials)
	if err != nil {
		return errors.MarshalingError(err)
	}

	if err = c.validatePOSTBody(&credentials); err != nil {
		return err
	}

	// Check that the supplied client_id matches the one we have

	githubClientID, isSet := os.LookupEnv("GITHUB_CLIENT_ID")
	if !isSet {
		return errors.MissingCredentialsError(fmt.Errorf(
			"Missing client ID credential"))
	}
	if credentials.Id != githubClientID {
		return errors.InvalidPOSTBodyError(fmt.Errorf(
			"Invalid client_id supplied"))
	}

	// Get the Github client secret
	githubClientSecret, isSet := os.LookupEnv("GITHUB_CLIENT_SECRET")
	if !isSet {
		return errors.MissingCredentialsError(fmt.Errorf(
			"Missing client secret credential"))
	}
	credentials.Secret = githubClientSecret

	// Get the access token from Github
	response, err := c.getAccessToken(&credentials)
	if err != nil {
		return err
	}

	// Read the response body, and write it to the response
	tokenBody, _ := ioutil.ReadAll(response.Body)
	c.w.Write(tokenBody)

	return nil
}

func (c *UsersController) handleOptionsRequest() error {
	c.w.Header().Set("Access-Control-Allow-Origin", "*")
	c.w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	c.w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, access-control-allow-methods")
	c.w.Write([]byte(""))
	return nil
}

func (c *UsersController) getAccessToken(credentials *model.AuthCredentials) (*http.Response, error) {
	body, err := json.Marshal(credentials)
	if err != nil {
		return nil, errors.MarshalingError(err)
	}
	req, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token/", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(body)))

	client := http.Client{}
	// Return the response to the Client's request
	return client.Do(req)
}

func (c *UsersController) validatePOSTBody(credentials *model.AuthCredentials) error {
	if credentials.Id == "" || credentials.Code == "" {
		return errors.IncompletePOSTBodyError(fmt.Errorf(
			"%q and %q fields must be non-empty", "Id", "Code"))
	}
	return nil
}

// stripFileContents takes a []byte of a file's contents, and will convert it to a string,
// and then trim that string and return it
// Needed as some editors will automatically enter a newline at the end, which will cause
// the inequality check for the client_id to fail
func (c *UsersController) stripFileContents(b []byte) string {
	return strings.TrimSpace(string(b))
}
