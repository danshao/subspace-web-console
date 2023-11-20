package form

import (
	"strings"
	"github.com/asaskevich/govalidator"
)

type UserSignInForm struct {
	Email    string `form:"email" valid:"email,required"`
	Password string `form:"password"`
}

func (form *UserSignInForm) beforeValidate() {
	form.Email = strings.TrimSpace(form.Email)
}

func (form *UserSignInForm) IsValid() (bool, error) {
	err := form.Validate()
	return nil == err, err
}

func (form *UserSignInForm) Validate() (err error) {
	form.beforeValidate()
	_, err = govalidator.ValidateStruct(form)
	return err
}


type PasswordRecoveryForm struct {
	Email      string `form:"email" valid:"email,required"`
	InstanceID string `form:"instance_id" valid:"-"`
	Password   string `form:"password" valid:"subspacePassword,required"`
}

func (form *PasswordRecoveryForm) beforeValidate() {
	form.Email = strings.TrimSpace(form.Email)
	form.InstanceID = strings.TrimSpace(form.InstanceID)
}

func (form *PasswordRecoveryForm) IsValid() (bool, error) {
	err := form.Validate()
	return nil == err, err
}

func (form *PasswordRecoveryForm) Validate() (err error) {
	form.beforeValidate()
	_, err = govalidator.ValidateStruct(form)
	return err
}
