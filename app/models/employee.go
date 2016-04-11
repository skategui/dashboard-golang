package models

import (
	"fmt"
	"github.com/revel/revel"
	"regexp"
)


type Employee struct {
	EmployeeID         int
	Email              string
	Username, Password string
	HashedPassword     []byte
	BrandID		   int
}

func (u *Employee) String() string {
	return fmt.Sprintf("User(%s)", u.Username)
}


func (user *Employee) Validate(v *revel.Validation) {
	var userRegex = regexp.MustCompile("^\\w*$")
	v.Check(user.Username,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{4},
		revel.Match{userRegex},
	)

	ValidatePassword(v, user.Password).
	Key("user.Password")

	v.Check(user.Email,
		revel.Required{},
		revel.MaxSize{100},
	)
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{5},
	)
}
