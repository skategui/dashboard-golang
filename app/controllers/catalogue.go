package controllers

import (
"github.com/revel/revel"
"dashboard/app/routes"
)

type Catalogue struct {
	Application
}


func (c Catalogue) Index() revel.Result {
	if c.connected() == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Admin.Login())
	}
	return c.Render()
}

func (c Catalogue) AddProduct() revel.Result {
	return c.Render()
}
