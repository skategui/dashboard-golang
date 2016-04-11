package controllers

import (
	"github.com/revel/revel"
	"dashboard/app/models"
	"encoding/json"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"crypto/rand"
	"fmt"
)

type Api struct {
	Application
}

type SignupResponse struct {
	StatusCode int
	Token    string
}

type SimplesReponse struct {
	StatusCode int
	Message    string
}

type ErrorResponse struct {
	StatusCode int
	Error    string
}


func (c Api) parseLogin() (models.User, error) {
	user := models.User{}
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	return user, err
}



func (c Api) Signup() revel.Result {

	if user, err := c.parseLogin(); err != nil {
		return c.RenderText("Unable to parse the User from JSON.")
	} else {
		// Validate the model
		user.Validate(c.Validation)
		if c.Validation.HasErrors() {
			// Do something better here!
			fmt.Sprintln("error : %s", c.Validation.Errors)
			return c.RenderJson(c.Validation.Errors)
		} else {
			return c.insertNewUser(user)
		}
	}
}

func (c Api) insertNewUser(user models.User) revel.Result {

	_ ,userExist := c.isUserExist(user.Email, user.Password)
	if (userExist == true) {
		response := ErrorResponse{StatusCode:http.StatusOK, Error:"User already exist"}
		return c.RenderJson(response)
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Token = c.randToken()
	err := c.Txn.Insert(&user);
	if err != nil {
		panic(err)
		return c.RenderJson(err)
	} else {
		response := SignupResponse{StatusCode:http.StatusOK, Token:user.Token}
		return c.RenderJson(response)
	}
}


func (c Api) GetUserByEmail(email string) (models.User, bool) {
	var user models.User
	err := c.Txn.SelectOne(&user,
		"select * from User where Email = ? LIMIT 1 ", email)
	if err != nil {
		return user, false
	}
	return user, true
}

func (c Api) randToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}


func (c Api) isUserExist(email, password string) (models.User, bool) {
	user, found := c.GetUserByEmail(email)
	if (found == false) {
		return user , false;
	}
	err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err == nil {
		return user, true
	}
	return user, false
}


func (c Api) Login() revel.Result {

	if loginInfo, err := c.parseLogin(); err != nil {
		return c.RenderText("Unable to parse the User from JSON.")
	} else {
		loginInfo.Validate(c.Validation)
		if c.Validation.HasErrors() {
			// Do something better here!
			fmt.Sprintln("error : %s", c.Validation.Errors)
			return c.RenderJson(c.Validation.Errors)
		}
		user, userExist := c.isUserExist(loginInfo.Email, loginInfo.Password)
		if userExist == true {
			response := SignupResponse{StatusCode:http.StatusBadRequest, Token:user.Token}
			return c.RenderJson(response)
		}
		response := ErrorResponse{StatusCode:http.StatusBadRequest, Error:"Wrong credential"}
		return c.RenderJson(response)
	}
}