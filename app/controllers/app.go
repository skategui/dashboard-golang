package controllers

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/revel/revel"
	"dashboard/app/models"
	"dashboard/app/routes"
)

type Application struct {
	GorpController
}

func (c Application) AddUser() revel.Result {
	if user := c.connected(); user != nil {
		c.RenderArgs["user"] = user
	}
	return nil
}

func (c Application) connected() *models.Employee {
	if c.RenderArgs["user"] != nil {
		return c.RenderArgs["user"].(*models.Employee)
	}
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username)
	}
	return nil
}

func (c Application) getUser(username string) *models.Employee {
	users, err := c.Txn.Select(models.Employee{}, `select * from Employee where Username = ?`, username)
	if err != nil {
		panic(err)
	}
	if len(users) == 0 {
		return nil
	}
	return users[0].(*models.Employee)
}

func (c Application) Index() revel.Result {
	if c.connected() != nil {
		return c.Redirect(routes.Hotels.Index())
	}
	c.Flash.Error("Please log in first")
	return c.Redirect(routes.Admin.Login())
}

func (c Application) Register() revel.Result {
	return c.Render()
}

func (c Application) SaveUser(user models.Employee, verifyPassword string) revel.Result {
	c.Validation.Required(verifyPassword)
	c.Validation.Required(verifyPassword == user.Password).
		Message("Password does not match")
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Application.Register())
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost)
	err := c.Txn.Insert(&user)
	if err != nil {
		panic(err)
	}

	c.Session["user"] = user.Username
	c.Flash.Success("Welcome, " + user.Username)
	return c.Redirect(routes.Hotels.Index())
}
