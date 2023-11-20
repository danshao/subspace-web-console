package models

import (
	"time"

	"gitlab.ecoworkinc.com/Subspace/web-console/helpers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	USER_ROLE_ADMIN  = "admin"
	USER_ROLE_MEMBER = "member"
)

type User struct {
	Id               int       `orm:"pk"`     // auto-increment
	Email            string    `orm:"unique"` // required
	Alias            string    `orm:"null"`   //
	Role             string    // required, enum('admin', 'user')
	EmailVerified    bool      `orm:"default(0)"` // required
	Enabled          bool      `orm:"default(1)"` // required
	PasswordHash     string    // required
	SetPasswordToken string    `orm:"null"`                        //
	RevokedDate      time.Time `orm:"null;type(datetime)"`         //
	LastLoginDate    time.Time `orm:"null;type(datetime)"`         //
	CreatedDate      time.Time `orm:"auto_now_add;type(datetime)"` //
	UpdatedDate      time.Time `orm:"auto_now_add;type(datetime)"` //
}

func (u *User) TableName() string {
	return "users"
}

func (u User) GetUserId() int {
	return u.Id
}

func (u User) GetRole() string {
	return u.Role
}

func (u User) IsEnabled() bool {
	return u.Enabled
}


// CreateUser creates a user.
func CreateUser(email, password, alias, role string, emailVerified, enabled bool) (int, error) {
	passwordHash, _ := helpers.HashPassword(password)
	user := &User{
		Email:         email,
		Role:          role,
		Alias:         alias,
		EmailVerified: emailVerified,
		Enabled:       enabled,
		PasswordHash:  passwordHash,
	}
	o := orm.NewOrm()

	id, err := o.Insert(user)
	if err != nil {
		beego.Error("CreateUser DB INSERT failed with error: ", err)
		return 0, err
	}

	return int(id), nil
}

// GetAllUsers gets all the users from the database.
func GetAllUsers(limit ...int) []User {
	var (
		limitRows []int
		users     []User
	)
	o := orm.NewOrm()
	if limitRows = limit; len(limitRows) == 0 {
		limitRows = append(limitRows, 1000)
	}

	_, err := o.QueryTable("users").Limit(limitRows[0]).All(&users)
	if err != nil {
		beego.Error("GetAllUsers DB SELECT failed with error: ", err)
	}

	return users
}

// AuthenticateUser authenticates a user
// Validates email, password, and role = 'admin'
func AuthenticateUser(email, password string) *User {
	var user User
	o := orm.NewOrm()
	cond := orm.NewCondition()

	condAuth := cond.And("email", email)
	err := o.QueryTable((&User{}).TableName()).SetCond(condAuth).One(&user)
	if err == orm.ErrMultiRows {
		beego.Error("The email has multiple records: ", email)
		return nil
	}
	if err == orm.ErrNoRows {
		beego.Error("No email has been found: ", email)
		return nil
	}
	passwordHash := user.PasswordHash
	if !helpers.CheckPasswordHash(password, passwordHash) {
		beego.Error("Password not match for user: ", email)
		return nil
	}
	return &user
}

//TODO Return nil, error when something wrong.
// GetUserWithID gets the information of a specific user
func GetUserWithID(userID int) User {
	user := User{Id: userID}
	o := orm.NewOrm()

	err := o.Read(&user)
	if err == orm.ErrNoRows {
		beego.Debug("No result found.")
	} else if err == orm.ErrMissPK {
		beego.Debug("No primary key found.")
	} else if err != nil {
		beego.Error("GetUserWithID DB SELECT failed with error: ", err)
	}

	return user
}

func GetUserWithEmail(email string) (user User, err error) {
	o := orm.NewOrm()

	if err = o.QueryTable((&User{}).TableName()).Filter("email", email).One(&user); err != nil {
		beego.Error("Error occurred when query batch users: ", err)
	}

	return user, err
}

// GetAdminCount gets the number of current administrators
func GetAdminCount() (numRows int, err error) {
	o := orm.NewOrm()

	if numRows, err := o.QueryTable((&User{}).TableName()).Filter("role", "admin").Count(); err == nil {
		return int(numRows), err
	}
	beego.Error("Error occured while GetAdminCount: ", err)
	return 0, err
}

// GetUsersWithBatchID gets the information of a bunch of users using user's ID
func GetUsersWithIDs(userIDs []int) []User {
	if userIDs == nil {
		return nil
	}
	o := orm.NewOrm()
	var users []User

	// num, err := o.Raw("Select * From "+(&User{}).TableName()+" Where id IN(?)", userIDs).QueryRows(&users)
	num, err := o.QueryTable((&User{}).TableName()).Filter("id__in", userIDs).All(&users)
	if err != nil {
		beego.Error("Error occurred when query batch users: ", err)
	}
	beego.Debug("Query result: ", num)

	return users
}

// GetCountUserWithEmail returns a count of the number of rows with the user email address
func UserEmailExists(email string) bool {
	o := orm.NewOrm()
	numRows, _ := o.QueryTable((&User{}).TableName()).Filter("email", email).Count()
	if numRows == 0 {
		return false
	}
	return true
}

// UpdateUserWithID
func UpdateUserWithID(userID int, password, alias, role string) (num int64, err error) {
	o := orm.NewOrm()
	user := User{Id: userID}

	if o.Read(&user) == nil {
		if password != "" {
			user.PasswordHash, _ = helpers.HashPassword(password)
		}
		if role != "" {
			user.Role = role
		}
		user.Alias = alias

		num, err = o.Update(&user)
	}

	return
}

// DeleteUserWithUserID deletes a user given a specified USER_ID
func DeleteUserWithUserID(userID int) (numRows int, err error) {
	o := orm.NewOrm()

	if numRows, err := o.QueryTable((&User{}).TableName()).Filter("id", userID).Delete(); err == nil {
		return int(numRows), nil
	}

	return 0, err
}

// EnableUserWithID disables a specific user given a USER_ID
func EnableUserWithID(userID int) (int64, error) {
	o := orm.NewOrm()
	result, _ := o.Raw("UPDATE `users` SET `enabled` = TRUE, `revoked_date` = NULL WHERE `id` = ?", userID).Exec()
	return result.RowsAffected()
}

// DisableUserWithID disables a specific user given a USER_ID
func DisableUserWithID(userID int) (int64, error) {
	o := orm.NewOrm()
	result, _ := o.Raw("UPDATE `users` SET `enabled` = FALSE, `revoked_date` = CURRENT_TIMESTAMP WHERE `id` = ?", userID).Exec()
	return result.RowsAffected()
}
