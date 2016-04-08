package models

type AssoBrandProduct struct {
	AssoBrandProductID           int
	BrandID      		     int
	ProductID    		     int

	//transient
	Brand         	*Brand
	Product         *Product
}
