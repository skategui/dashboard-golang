package controllers

import (
"github.com/revel/revel"
"dashboard/app/models"
"net/http"
"encoding/json"
	"dashboard/app/routes"
)

type Product struct {
	Application
}


func (c Product) parseProduct() (models.Product, error) {
	product := models.Product{}
	err := json.NewDecoder(c.Request.Body).Decode(&product)
	return product, err
}

func (c Product) AddProduct() revel.Result {

	if product, err := c.parseProduct(); err != nil {
		return c.RenderText("Unable to parse the Product from JSON.")
	} else {
		product.Validate(c.Validation)

		if c.Validation.HasErrors() {
			c.Response.Status = http.StatusBadRequest
			return c.RenderJson(c.Validation.Errors)
		}

		err := c.Txn.Insert(&product)
		if err != nil {
			panic(err)
		}
		c.Response.Status = http.StatusOK
		return c.RenderJson(product)
	}
}


func (c Product) SubmitProduct(name, description string, quantity int) revel.Result{

	c.Validation.Required(name).Message("Name can't be empty")
	c.Validation.Required(description).Message("Description can't be empty")

	c.Validation.Required(quantity).Message("Address can't be empty")

	if c.Validation.HasErrors() {
		return c.Redirect(routes.Catalogue.Index())
	}
	var product = models.Product{Name:name, Description:description, Quantity: quantity}

	err := c.Txn.Insert(&product)
	if err != nil {
		panic(err)
	}

	var tmp = c.GetProductByName(name)
	var assoProduct = models.AssoBrandProduct{BrandID:c.connected().BrandID, ProductID : tmp.ProductID}

	addError := c.Txn.Insert(&assoProduct)
	if addError != nil {
		panic(addError)
	}
	return c.Redirect(routes.Catalogue.Index())
}


func (c Product) GetAllProducts() revel.Result {
	var products []*models.Product
	products = loadProducts(c.Txn.Select(models.Product{},
		`select * from Product`))
	return c.RenderJson(products)
}

func (c Product) GetProductByName(name string) models.Product {
	var product models.Product
	err := c.Txn.SelectOne(&product,
		"select * from Product where name = ? ", name)
	if err != nil {
		panic(err)
	}
	return product
}



func (c Product) GetProductByID(id int) revel.Result {
	h, err := c.Txn.Get(models.Product{}, id)
	if err != nil {
		panic(err)
	}
	if h == nil {
		return nil
	}
	return c.RenderJson(h)
}


func (c Product) DeleteProductByID(id int) revel.Result {
	success, err := c.Txn.Delete(&models.Product{ProductID:id})
	if err != nil || success == 0 {
		return c.RenderText("Failed to remove Product with id %v", id)
	}
	return c.RenderJson(success)
}

func loadProducts(results []interface{}, err error) []*models.Product {
	if err != nil {
		panic(err)
	}
	var products []*models.Product
	for _, r := range results {
		products = append(products, r.(*models.Product))
	}
	return products
}