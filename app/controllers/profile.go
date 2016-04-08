package controllers

import (
"github.com/revel/revel"
"dashboard/app/routes"
)

type Profile struct {
	Application
}


func (c Profile) Index() revel.Result {

	var user = c.connected()
	if user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Admin.Login())
	}
	var userName = user.Username
	return c.Render(userName)
}

