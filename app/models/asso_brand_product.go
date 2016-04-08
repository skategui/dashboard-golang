package models

type AssoBrandProduct struct {
	AssoBrandProductID           int64
	BrandID      		     int64
	ProductID    		     int64

	//transient
	Brand         	*Brand
	Product         *Product
}
