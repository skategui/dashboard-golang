package models

import (
	"fmt"
	"github.com/revel/revel"
	"regexp"
)


type User struct {
	UserID             int
	Email              string
	Password  	   string
	HashedPassword     []byte
	Token		   string
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Email)
}


func (user *User) Validate(v *revel.Validation) {
	var EMAIL_PATTERN = "^[_A-Za-z0-9-\\+]+(\\.[_A-Za-z0-9-]+)*@[A-Za-z0-9-]+(\\.[A-Za-z0-9]+)*(\\.[A-Za-z]{2,})$";
	v.Check(user.Email,
		revel.Required{},
		revel.MaxSize{100},
		revel.MinSize{4},
		revel.Match{regexp.MustCompile(EMAIL_PATTERN)},
	).Message("Invalid email addr")

	v.Check(user.Password,
		revel.Required{},
		revel.MaxSize{30},
		revel.MinSize{5},
	).Message("Invalid password. Should be be between 5 and 30 characters")

	v.Check(user.Email,
		revel.Required{},
		revel.MaxSize{100},
	).Message("Email is required")

	v.Check(user.Password,
		revel.Required{},
		revel.MaxSize{100},
	).Message("Password is required")
}
