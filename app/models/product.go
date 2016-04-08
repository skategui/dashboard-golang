package models

import (
"fmt"
"github.com/revel/revel"
)

type Product struct {
	ProductID           int64 `db:"productID" json:"productID"`
	Name               string `db:"name" json:"name"`
	Description	   string `db:"description" json:"description"`
	Quantity	      int `db:"quantity" json:"quantity"`
}

func (product *Product) String() string {
	return fmt.Sprintf("Product ID: %d, name :  %s, desc: %s, quantity: %d", product.ProductID, product.Name, product.Description, product.Quantity)
}


func (product *Product) Validate(v *revel.Validation) {
	v.Check(product.Name,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{4},
	)
	v.Check(product.Description,
		revel.Required{},
		revel.MinSize{4},
		revel.MaxSize{100},
	)
}

