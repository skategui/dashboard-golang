package controllers

import (
"github.com/revel/revel"
"dashboard/app/models"
"net/http"
"fmt"
"dashboard/app/routes"
)

type Store struct {
	Application
}


func (c Store) SubmitStore(name, description, address string) revel.Result {

	c.Validation.Required(name).Message("Name can't be empty")
	c.Validation.Required(description).Message("Description can't be empty")

	c.Validation.Required(address).Message("Address can't be empty")

	if c.Validation.HasErrors() {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(c.Validation.Errors)
	}
	var store = models.Store{Name:name, Description:description, Address:address}

	err := c.Txn.Insert(&store)
	if err != nil {
		panic(err)
	}
	c.Response.Status = http.StatusOK
	return c.RenderJson(store)
}


func (c Store) GetStoreByID(id int) *models.Store  {
	store, err := c.Txn.Get(models.Store{}, id)
	if err != nil {
		panic(err)
	}
	if store == nil {
		panic(err)
		return nil
	}
	return store.(*models.Store)
}



func (c Store) GetAllStoreByBrandID(brandID int32)  []*models.Store {
	var stores []*models.Store
	stores = loadStores(c.Txn.Select(models.Store{},
		`select * from Store where brandID = ?`, brandID))
	return stores
}

func loadStores(results []interface{}, err error) []*models.Store {
	if err != nil {
		panic(err)
	}
	var stores []*models.Store
	for _, r := range results {
		stores = append(stores, r.(*models.Store))
	}
	return stores
}

func (c Store) Stores() revel.Result {
	var user = c.connected()
	if user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Admin.Login())
	}
	fmt.Println("user brandID : %i", user.BrandID)
	var stores = loadStores(c.Txn.Select(models.Store{},
		`select * from Store where BrandID = ?`, 1 ))
	return c.Render(stores)
}


func (c Store) Store(id int) revel.Result {
	if c.connected() == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Admin.Login())
	}
	return c.Render()
}
