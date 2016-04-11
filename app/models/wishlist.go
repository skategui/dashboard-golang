package models

import "github.com/revel/revel"

type WishList struct {
	WishListID      int
	UserID      	int
	ProductID    	int

	//transient
	User         	*User
	Product         *Product
}


func (wishlist *WishList) Validate(v *revel.Validation) {
	v.Check(wishlist.UserID,
		revel.Required{},
	).Message("UserID missing")
	v.Check(wishlist.ProductID,
		revel.Required{},
	).Message("productID missing")
}
