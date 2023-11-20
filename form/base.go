package form

import (
	"fmt"
	"reflect"
	"github.com/asaskevich/govalidator"
	"gitlab.ecoworkinc.com/Subspace/subspace-utility/subspace/utils/validator"
)

// Call when this go file loaded.
func init() {
	govalidator.CustomTypeTagMap.Set("subspacePassword", govalidator.CustomTypeValidator(func(i interface{}, context interface{}) bool {
		password := fmt.Sprintf("%v", i)
		return validator.IsValidPassword(password)
	}))

	govalidator.CustomTypeTagMap.Set("confirmPassword", govalidator.CustomTypeValidator(func(i interface{}, context interface{}) bool {
		//TODO For now, hard code field name "Password" was set, change to something like confirmPassword(FIELD_NAME_TO_CONFIRM)
		//TODO If confirmPassword is empty string, it will not trigger this validation function
		confirmPassword := fmt.Sprintf("%v", i)
		password := reflect.ValueOf(context).FieldByName("Password").String()
		fmt.Println(password)
		return password == confirmPassword
	}))
}

type IForm interface {
	beforeValidate()
	IsValid() (bool, error)
	Validate() (errors error)
}