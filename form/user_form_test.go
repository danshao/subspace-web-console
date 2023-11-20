package form

import (
	"testing"
)

func TestCreateAdminForm(t *testing.T) {
	f1 := CreateAdminForm{
		Email: "carterlin@ecoworkinc.com",
		Alias: "CarterLin",
		Password: "Abcd1234",
		ConfirmPassword: "Abcd1234",
	}
	if err := f1.Validate(); nil != err {
		t.Error(err)
	}

	f2 := CreateAdminForm{
		Email: "carterlin@ecoworkinc.com",
		Alias: "CarterLin",
		Password: "Abcd1234",
		ConfirmPassword: "abcd1234",
	}
	if err := f2.Validate(); nil == err {
		t.Error("Abcd1234 not match abcd1234")
	}
}

func TestUserCreateForm(t *testing.T) {
	f1 := UserCreateForm{
		Email: "carterlin@ecoworkinc.com",
		Password: "",
		ConfirmPassword: "",
		GeneratePassword: true,
		Alias: "CarterLin",
		Role: "admin",
		VerifyEmail: false,
		CreateVPNProfile: true,
		EmailDelivery: false,
	}
	if err := f1.Validate(); nil != err {
		t.Error(err)
	}
}

func TestUserUpdateForm(t *testing.T) {
	f1 := UserUpdateForm{
		Password: "Abcd1234",
		ConfirmPassword: "Abcd1234",
		Alias: "CarterLin",
		Role: "admin",
	}
	if err := f1.Validate(); nil != err {
		t.Error(err)
	}

	f2 := UserUpdateForm{
		Password: "Abcd1234",
		ConfirmPassword: "abcd1234",
		Alias: "CarterLin",
		Role: "admin",
	}
	if err := f2.Validate(); nil == err {
		t.Error("Abcd1234 not match abcd1234")
	}

	f3 := UserUpdateForm{
		Password: "",
		ConfirmPassword: "",
		Alias: "CarterLin",
		Role: "admin",
	}
	if err := f3.Validate(); nil != err {
		t.Error("Password can keep the same")
	}

	f4 := UserUpdateForm{
		Password: "",
		ConfirmPassword: "",
		Alias: "CarterLin",
		Role: "member",
	}
	if err := f4.Validate(); nil == err {
		t.Error("User role do not accept member")
	}
}