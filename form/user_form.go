package form

import (
	"strings"
	"errors"
	"bytes"
	"github.com/asaskevich/govalidator"
	"gitlab.ecoworkinc.com/Subspace/subspace-utility/subspace/utils/validator"
)

// Init Create admin form
type CreateAdminForm struct {
	Alias           string `form:"alias" valid:"-"`
	Email           string `form:"email" valid:"email,required"`
	Password        string `form:"password" valid:"subspacePassword,required"`
	ConfirmPassword string `form:"confirm_password" valid:"confirmPassword,required"`
}

func (form *CreateAdminForm) beforeValidate() {
	form.Alias = strings.TrimSpace(form.Alias)
	form.Email = strings.TrimSpace(form.Email)
}

func (form *CreateAdminForm) IsValid() (bool, error) {
	err := form.Validate()
	return nil == err, err
}

func (form *CreateAdminForm) Validate() (err error) {
	form.beforeValidate()
	_, err = govalidator.ValidateStruct(form)
	return err
}


type UserCreateForm struct {
	Email            string `form:"email" valid:"email,required"`
	Password         string `form:"password"`
	ConfirmPassword  string `form:"confirm_password"`
	GeneratePassword bool   `form:"autogenPassword"`
	Alias            string `form:"alias"`
	Role             string `form:"role" valid:"in(admin|user)"`
	VerifyEmail      bool   `form:"verifyEmail"`
	CreateVPNProfile bool   `form:"createVPNProfile"`
	EmailDelivery    bool   `form:"delivery"`
}

func (form *UserCreateForm) beforeValidate() {
	form.Email = strings.TrimSpace(form.Email)
	form.Role = strings.TrimSpace(form.Role)
}

func (form *UserCreateForm) IsValid() (bool, error) {
	err := form.Validate()
	return nil == err, err
}

func (form *UserCreateForm) Validate() (err error) {
	form.beforeValidate()

	if _, err := govalidator.ValidateStruct(form); nil != err {
		return err
	}

	var buffer bytes.Buffer
	if !validator.IsValidEmail(form.Email) {
		buffer.WriteString("Email is not valid;")
	}

	if !form.GeneratePassword {
		if !validator.IsValidPassword(form.Password) {
			buffer.WriteString("Password Validation Failed: Please match the requested format. Password must have a minimum of 8 characters and at least 1 upper case letter, 1 lower case letter and 1 number;")
		}
		if form.Password != form.ConfirmPassword {
			buffer.WriteString("Error: Passwords do not match;")
		}
	}

	if "" != buffer.String() {
		err = errors.New(buffer.String())
	}
	return err
}



type UserUpdateForm struct {
	Alias           string `form:"alias"`
	Password        string `form:"password"`
	ConfirmPassword string `form:"confirm_password"`
	Role            string `form:"roleUpdate" valid:"in(admin|user)"`
}

func (form *UserUpdateForm) beforeValidate() {
	form.Role = strings.TrimSpace(form.Role)
}

func (form *UserUpdateForm) IsValid() (bool, error) {
	err := form.Validate()
	return nil == err, err
}

func (form *UserUpdateForm) Validate() (err error) {
	form.beforeValidate()
	if _, err := govalidator.ValidateStruct(form); nil != err {
		return err
	}

	var buffer bytes.Buffer
	if "" != form.Password || "" != form.ConfirmPassword {
		if !validator.IsValidPassword(form.Password) {
			buffer.WriteString("Password Validation Failed: Please match the requested format. Password must have a minimum of 8 characters and at least 1 upper case letter, 1 lower case letter and 1 number;")
		}
		if form.Password != form.ConfirmPassword {
			buffer.WriteString("Error: Passwords do not match;")
		}
	}

	if "" != buffer.String() {
		err = errors.New(buffer.String())
	}
	return err
}