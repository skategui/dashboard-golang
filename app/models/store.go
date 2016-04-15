package models

import (
"fmt"
"github.com/revel/revel"
"github.com/kellydunn/golang-geo"
	"encoding/json"
)

type Store struct {
	StoreID           int
	Name              string
	BrandID           int
	Description	  string
	Address		  string
	Longitude	  float64
	Latitude	  float64

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

func  (store *Store) GetPositionByAddrName() {
	geo.SetOpenCageAPIKey("5788cf4e67ac29311b77bf14d6e96d60")
	geocoder := new (geo.OpenCageGeocoder)
	fmt.Println(store.Address)
	p,err := geocoder.Geocode(store.Address)
	if (err != nil) {
		panic(err)
	}
	store.Latitude = p.Lat()
	store.Longitude = p.Lng()
	fltB, _ := json.Marshal(p)
	fmt.Println(string(fltB))
}

