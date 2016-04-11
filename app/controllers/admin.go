package controllers

import (
"github.com/revel/revel"
"golang.org/x/crypto/bcrypt"
"regexp"
"dashboard/app/routes"
"dashboard/app/models"
)

type Admin struct {
	Application
}


func (c Admin) Index() revel.Result {
	var user = c.connected()
	if user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Admin.Login())
	}
	var userName = user.Username
	return c.Render(userName)
}

func (c Admin) Registration() revel.Result {
	return c.Render()
}

func (c Admin) Login() revel.Result {
	return c.Render()
}



func (c Admin) TryLogin(username, password string, remember bool) revel.Result {

	c.Validation.Required(username).Message("Username can't be empty")
	c.Validation.Required(password).Message("Password can't be empty")
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		return c.Redirect(routes.Admin.Login())
	}
	user := c.getUser(username)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
		if err == nil {
			c.Session["user"] = username
			if remember {
				c.Session.SetDefaultExpiration()
			} else {
				c.Session.SetNoExpiration()
			}
			c.Flash.Success("Welcome, " + username)
			return c.Redirect(routes.Admin.Index())
		}
	}
	c.Flash.Out["username"] = username
	c.Flash.Error("Wrong credential")
	return c.Redirect(routes.Admin.Login())
}

func (c Admin) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.Admin.Index())
}


func (c Admin) Register(username, email, password, verifyPassword string) revel.Result {

	models.ValidatePassword(c.Validation, password).Message("Please verify your password")
	c.Validation.Required(verifyPassword).Message("Please verify your password")
	c.Validation.Required(verifyPassword == password).Message("Your password doesn't match")

	c.Validation.Required(username).Message("Username can't be empty")
	c.Validation.MinSize(username, 4).Message("Username too short.")
	c.Validation.Required(email).Message("Email addr can't be empty")


	var EMAIL_PATTERN = "^[_A-Za-z0-9-\\+]+(\\.[_A-Za-z0-9-]+)*@[A-Za-z0-9-]+(\\.[A-Za-z0-9]+)*(\\.[A-Za-z]{2,})$";

	c.Validation.Match(email, regexp.MustCompile(EMAIL_PATTERN)).Message("Invalid email addr")



	if c.Validation.HasErrors() {
		c.Validation.Keep()
		return c.Redirect(routes.Admin.Registration())
	}

	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	newUser := models.Employee{Email: email, Username:username, Password:password, HashedPassword: bcryptPassword, BrandID: 1} // should be dynamic
	err := c.Txn.Insert(&newUser)
	if err != nil {
		panic(err)
	}

	c.Session["user"] = username
	c.Session.SetDefaultExpiration()
	c.Flash.Success("Welcome, " + username)
	return c.Redirect(routes.Admin.Index())

}