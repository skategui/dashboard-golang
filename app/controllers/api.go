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

type SuccessResponse struct {
	Success bool
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


func (c Api) parseWishlist() (models.WishList, error) {
	wishlist := models.WishList{}
	err := json.NewDecoder(c.Request.Body).Decode(&wishlist)
	return wishlist, err
}

func (c Api) isAlreadyExist(productID, userID int) bool {
	wishlist := new(models.WishList)
	err := c.Txn.SelectOne(wishlist, `SELECT * FROM WishList WHERE UserID = ? AND ProductID = ?`, userID, productID)
	if err != nil {
		return false
	}
	return true
}

func (c Api) AddWishlist() revel.Result {
	if wishlist, err := c.parseWishlist(); err != nil {
		return c.RenderText("Unable to parse the User from JSON.")
	} else {
		wishlist.Validate(c.Validation)
		if c.Validation.HasErrors() {
			return c.RenderJson(c.Validation.Errors)
		}

		var alreadyExist = c.isAlreadyExist(wishlist.ProductID, wishlist.UserID)
		if (alreadyExist == true) {
			var response = ErrorResponse{StatusCode:http.StatusBadRequest, Error: "Product already in the wishlist"}
			return c.RenderJson(response)
		}
		err := c.Txn.Insert(&wishlist);
		if err != nil {
			panic(err)
			return c.RenderJson(err)
		}
		var response = SuccessResponse{Success : true}
		return c.RenderJson(response)
	}
}


func (c Api) RemoveWishlist() revel.Result {
	if wishlist, err := c.parseWishlist(); err != nil {
		return c.RenderText("Unable to parse the User from JSON.")
	} else {
		wishlist.Validate(c.Validation)
		if c.Validation.HasErrors() {
			return c.RenderJson(ErrorResponse{StatusCode:http.StatusBadRequest, Error: "Parameter missing"})
		}
		_, err := c.Txn.Exec("delete FROM WishList WHERE UserID = ? AND ProductID = ?", wishlist.UserID, wishlist.ProductID)
		if err != nil {
			return c.RenderJson(ErrorResponse{StatusCode : http.StatusInternalServerError, Error: "Can't remove from wishlist"})
		}
		return c.RenderJson(SuccessResponse{Success : true})
	}
}