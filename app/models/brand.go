package models

import (
"fmt"
"github.com/revel/revel"
)

type Brand struct {
	BrandID           int64		`db:"brandID" json:"brandID"`
	Name              string	`db:"name" json:"name"`
}

func (brand *Brand) String() string {
	return fmt.Sprintf("Brand ID: %d, name :  %s", brand.BrandID, brand.Name)
}


func (brand *Brand) Validate(v *revel.Validation) {
	v.Check(brand.Name,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{4},
	)
}

