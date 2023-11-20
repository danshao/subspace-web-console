package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"gitlab.ecoworkinc.com/Subspace/subspace-utility/subspace/vpn"
	"strings"
)

type Profile struct {
	Id             int       `orm:"pk"`                          // auto-increment
	Hub            string                                        // required
	Username       string                                        // required
	FullName       string    `orm:"null"`                        //
	Description    string    `orm:"null"`                        //
	UserId         int                                           // required
	VpnHost        string                                        // required
	PreSharedKey   string                                        // required
	Enabled        bool      `orm:"default(1)"`                  //
	LoginCount     int       `orm:"default(0)"`                  //
	IncomingBytes  float64   `orm:"default(0)"`                  //
	OutgoingBytes  float64   `orm:"default(0)"`                  //
	PasswordHash   string                                        // required
	NtlmSecureHash string                                        // required
	RevokedDate    time.Time `orm:"null;type(datetime)"`         //
	LastLoginDate  time.Time `orm:"null;type(datetime)"`         //
	CreatedDate    time.Time `orm:"auto_now_add;type(datetime)"` //
	UpdatedDate    time.Time `orm:"auto_now_add;type(datetime)"` //
}

func (p *Profile) TableName() string {
	return "profiles"
}

// CreateProfile creates a vpn profile
func CreateProfile(userID int, hub, username, description, email, password, host, psk string) (int, error) {
	//Generate password hash follow softether rule.
	passwordHash, err := vpn.GenerateSoftetherPasswordHash(password + strings.ToUpper(username))
	if nil != err {
		return -1, err
	}

	profile := Profile{
		Hub:            hub,
		Username:       username,
		FullName:       email,
		Description:    description,
		UserId:         userID,
		VpnHost:        host,
		PreSharedKey:   psk,
		Enabled:        true,
		PasswordHash:   passwordHash,
		NtlmSecureHash: vpn.GenerateNtLmPasswordHash(password),
	}
	o := orm.NewOrm()
	id, err := o.Insert(&profile)
	if err != nil {
		beego.Error("error occured while inserting into profile: ", err)
		return 0, err
	}
	return int(id), nil
}

// GetAllProfiles gets all of the SoftEther users
func GetAllProfiles(limit ...int) []Profile {
	var (
		limitRows []int
		users     []Profile
	)
	if limitRows = limit; len(limitRows) == 0 {
		limitRows = append(limitRows, 1000)
	}
	o := orm.NewOrm()
	_, err := o.QueryTable((&Profile{}).TableName()).Limit(limitRows[0]).All(&users)
	if err != nil {
		beego.Error("error occured while reading all profiles: ", err)
	}

	return users
}

//TODO Return nil, error when something wrong.
// GetProfileWithID gets a specific profile given an id
func GetProfileWithID(profileID int) Profile {
	o := orm.NewOrm()
	profile := Profile{Id: profileID}
	err := o.Read(&profile)
	if err == orm.ErrNoRows {
		beego.Warn("No result found.")
	} else if err == orm.ErrMissPK {
		beego.Warn("No primary key found.")
	} else if err != nil {
		beego.Error("error occured while reading profile of user: ", err)
	}
	return profile
}

// GetProfileWithUserID gets a specific profile given an user id
func GetProfilesWithUserID(userID int) []Profile {
	var (
		profiles []Profile
	)
	o := orm.NewOrm()
	_, err := o.QueryTable((&Profile{}).TableName()).Filter("user_id", userID).All(&profiles)
	if err != nil {
		beego.Error("Error occured GetProfileWithUserID: ", err)
	}

	return profiles
}

func UpdateProfileWithID(id int, description string) (int, error) {
	var (
		numAffect int64
		err       error
	)
	o := orm.NewOrm()
	profile := Profile{Id: id}
	if o.Read(&profile) == nil {
		profile.Description = description
		numAffect, err = o.Update(&profile)
		if err != nil {
			beego.Error("Error occured while updating profile "+profile.Username+": ", err)
		}
	}
	return int(numAffect), err
}

func EnableProfileWithID(profileID int) (int64, error) {
	o := orm.NewOrm()
	result, _ := o.Raw("UPDATE `profiles` SET `enabled` = TRUE, `revoked_date` = NULL WHERE `id` = ?", profileID).Exec()
	return result.RowsAffected()
}

func DisableProfileWithID(profileID int) (int64, error) {
	o := orm.NewOrm()
	result, _ := o.Raw("UPDATE `profiles` SET `enabled` = FALSE, `revoked_date` = CURRENT_TIMESTAMP WHERE `id` = ?", profileID).Exec()
	return result.RowsAffected()
}

// EnableProfilesWithUserID disables a specific profile given a USER_ID
func EnableProfilesWithUserID(userID int, revokedDate time.Time, enabled bool) (int, error) {
	o := orm.NewOrm()
	var (
		numRows int64
		err     error
	)

	if revokedDate.IsZero() {
		res, err := o.Raw("UPDATE "+(&Profile{}).TableName()+" SET enabled = ? WHERE user_id = ? ", enabled, userID).Exec()
		if err == nil {
			numRows, _ = res.RowsAffected()
		}
	} else {
		res, err := o.Raw("UPDATE "+(&Profile{}).TableName()+" SET enabled = ?, revoked_date = ? WHERE user_id = ? ", enabled, revokedDate, userID).Exec()
		if err == nil {
			numRows, _ = res.RowsAffected()
		}
	}

	return int(numRows), err
}

// DeleteProfilesWithUserID deletes a specific profile given a USER_ID
func DeleteProfilesWithUserID(userID int) (numRowsModified int, err error) {
	o := orm.NewOrm()

	if numRowsModified, err := o.QueryTable((&Profile{}).TableName()).Filter("user_id", userID).Delete(); err == nil {
		return int(numRowsModified), err
	}

	beego.Error("Error occured while DeleteProfilesWithUserID: ", userID, ", error: ", err)
	return 0, err
}

// DeleteProfileWithID deletes a profile given an id
func DeleteProfileWithID(profileID int) (int, error) {
	o := orm.NewOrm()
	profile := Profile{Id: profileID}
	num, err := o.Delete(&profile)
	if err != nil {
		beego.Error("Error occured while deleting profile: ", profileID, ", error: ", err)
	}
	return int(num), err
}
