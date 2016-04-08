package controllers

import (
"github.com/revel/revel"
"dashboard/app/routes"
	"dashboard/app/models"
)

type Catalogue struct {
	Application
}


func (c Catalogue) Index() revel.Result {
	if c.connected() == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Admin.Login())
	}
	var products = loadProducts(c.Txn.Select(models.Product{},
		`select Product.* from Product INNER JOIN AssoBrandProduct
		 ON AssoBrandProduct.ProductID = Product.productID AND
		 AssoBrandProduct.BrandID = ?`, c.connected().BrandID )) // should be a left join
	return c.Render(products)
}

func (c Catalogue) AddProduct() revel.Result {
	return c.Render()
}
