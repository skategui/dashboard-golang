package controllers

import (
"github.com/revel/revel"
"dashboard/app/models"
"fmt"
"dashboard/app/routes"
)

type Store struct {
	Application
}


func (c Store) SubmitStore(name, description, address, postcode, city, country string) revel.Result{

	c.Validation.Required(name).Message("Name can't be empty")
	c.Validation.Required(description).Message("Description can't be empty")

	c.Validation.Required(address).Message("Address can't be empty")

	if c.Validation.HasErrors() {
		return c.Redirect(routes.Store.Stores())
	}
	var addr = address + ", " + postcode + ", " + city + ", " + country
	var store = models.Store{Name:name, Description:description, Address:addr, BrandID:c.connected().BrandID}

	err := c.Txn.Insert(&store)
	if err != nil {
		panic(err)
	}
	return c.Redirect(routes.Store.Stores())
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

func (c Store) DeleteStoreByID(id int) revel.Result {
	success, err := c.Txn.Delete(&models.Store{StoreID:id})
	if err != nil || success == 0 {
		return c.RenderText("Failed to remove Store with id %v", id)
	}
	return c.RenderJson(success)
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
	store := c.GetStoreByID(id)
	store.GetPositionByAddrName()
	return c.Render(store)
}
