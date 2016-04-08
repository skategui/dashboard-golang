package models

import (
"fmt"
"github.com/revel/revel"
)

type Store struct {
	StoreID           int
	Name              string
	BrandID           int
	Description	  string
	Address		  string
	Longitude	  float32
	Latitude	  float32

	//transient
	Brand         	*Brand
}

func (store *Store) String() string {
	return fmt.Sprintf("StoreID ID: %d, name :  %s, Address: %s", store.StoreID, store.Name, store.Address)
}


func (store *Store) Validate(v *revel.Validation) {
	v.Check(store.Name,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{4},
	)
	v.Check(store.Address,
		revel.Required{},
		revel.MinSize{4},
		revel.MaxSize{100},
	)
}

