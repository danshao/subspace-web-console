package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type System struct {
	Restriction                     string    `orm:"pk"`
	InstanceId                      string    `orm:"null"`                        // default null
	SubspaceVersion                 string    `orm:"null"`                        // default null
	SubspaceBuildNumber             uint      `orm:"null"`                        // default null
	VpnServerVersion                string    `orm:"null"`                        // default null
	VpnServerBuildNumber            uint      `orm:"null"`                        // default null
	VpnServerAdministrationPassword string    `orm:"default(subspace)"`           //
	VpnServerAdministrationPort     uint      `orm:"default(992)"`                //
	VpnHubName                      string    `orm:"default(subspace)"`           //
	Ip                              string    `orm:"null"`                        // default null
	IpUpdatedDate                   time.Time `orm:"null"`                        // default null
	Host                            string    `orm:"null"`                        // default null
	HostUpdatedDate                 time.Time `orm:"null"`                        // default null
	PreSharedKey                    string    `orm:"null"`                        // default null
	PreSharedKeyUpdatedDate         time.Time `orm:"null"`                        // default null
	Uuid                            string    `orm:"null"`                        // default null
	UuidUpdatedDate                 time.Time `orm:"null"`                        // default null
	SmtpHost                        string    `orm:"null"`                        //
	SmtpPort                        int       `orm:"default(587)"`                // default null
	SmtpAuthentication              bool      `orm:"default(0)"`                  //
	SmtpUsername                    string    `orm:"null"`                        // default null
	SmtpPassword                    string    `orm:"null"`                        // default null
	SmtpSenderName                  string    `orm:"null"`                        //default null
	SmtpSenderEmail                 string    `orm:"null"`                        // default null
	SmtpValid                       bool      `orm:"null"`                        // default null
	LatestBackupDate                time.Time `orm:"null"`                        // default null
	UserSchemaVersion               uint      `orm:"null"`                        // default null
	ProfileSchemaVersion            uint      `orm:"null"`                        // default null
	ConfigSchemaVersion             uint      `orm:"null"`                        // default null
	UpdatedDate                     time.Time `orm:"auto_now_add;type(datetime)"` // not null, default current_timestamp
	CreatedAt                       time.Time `orm:"auto_now_add;type(datetime)"` // not null, default current_timestamp
}

func (s *System) TableName() string {
	return "system"
}

// CreateProfile creates a vpn profile
func InitSystemTable() error {
	s := System{
		SmtpPort:    587,
		SmtpValid:   false,
		UpdatedDate: time.Now(),
		CreatedAt:   time.Now(),
	}
	o := orm.NewOrm()
	_, err := o.Insert(&s)

	if err != nil {
		beego.Debug("Error occured while initializing system table:", err)
	}

	return err
}

// GetSystemInfo
func GetSystemInfo() (systemInfo System, err error) {
	o := orm.NewOrm()
	if err := o.QueryTable("system").One(&systemInfo); err != nil {
		beego.Error("SELECT * from system failed with error: ", err)
	}
	//XXX 不做讀取第二次，因為目前只有備份還原會需要用到 IP，因此不透過 Raw query 做，回傳的 IP 皆為空字串。
	//var maps []orm.Params
	//o.Raw("SELECT INET6_NTOA(ip) AS ip FROM " + (&System{}).TableName()).Values(&maps)
	//systemInfo.Ip = maps[0]["ip"].(string)
	return systemInfo, err
}

// UpdateSystemUUID
func UpdateSystemUUID(uuid string) (int, error) {
	o := orm.NewOrm()
	var (
		numRows int64
		err     error
	)

	if res, err := o.Raw("UPDATE "+(&System{}).TableName()+" SET uuid = ?, uuid_updated_date = ?, updated_date = ?", uuid, time.Now(), time.Now()).Exec(); err == nil {
		numRows, _ = res.RowsAffected()
		return int(numRows), err
	}

	return 0, err
}

// GetSystemUUID
func GetSystemUUID() (uuid string, err error) {
	var (
		systemInfo System
	)
	o := orm.NewOrm()
	if err := o.QueryTable("system").One(&systemInfo); err != nil {
		beego.Error("SELECT * from system failed with error: ", err)
	}
	return systemInfo.Uuid, err
}

func UpdateInstanceId(instanceId string) (int, error) {
	o := orm.NewOrm()
	var (
		numRows int64
		err     error
	)

	if res, err := o.Raw("UPDATE "+(&System{}).TableName()+" SET instance_id = ?, updated_date = ?", instanceId, time.Now()).Exec(); err == nil {
		numRows, _ = res.RowsAffected()
		return int(numRows), err
	}

	return 0, err
}

// UpdatePreSharedKey
func UpdatePreSharedKey(psk string) (int, error) {
	o := orm.NewOrm()
	var (
		numRows int64
		err     error
	)

	if res, err := o.Raw("UPDATE "+(&System{}).TableName()+" SET pre_shared_key = ?, pre_shared_key_updated_date = ?, updated_date = ?", psk, time.Now(), time.Now()).Exec(); err == nil {
		numRows, _ = res.RowsAffected()
		return int(numRows), err
	}

	return 0, err
}

// GetPreSharedKey
func GetPreSharedKey() (psk string, err error) {
	var (
		systemInfo System
	)
	o := orm.NewOrm()
	if err := o.QueryTable("system").One(&systemInfo); err != nil {
		beego.Error("SELECT * from system failed with error: ", err)
	}
	return systemInfo.PreSharedKey, err
}

// UpdateSystemHost
func UpdateSystemHost(host string) (int, error) {
	o := orm.NewOrm()
	var (
		numRows int64
		err     error
	)

	if res, err := o.Raw("UPDATE "+(&System{}).TableName()+" SET host = ?, host_updated_date = ?, updated_date = ?", host, time.Now(), time.Now()).Exec(); err == nil {
		numRows, _ = res.RowsAffected()
		return int(numRows), err
	}

	return 0, err
}

// UpdateSystemIP
func UpdateSystemIP(ip string) (int, error) {
	o := orm.NewOrm()
	var (
		numRows int64
		err     error
	)

	var maps []orm.Params
	o.Raw("SELECT INET6_NTOA(ip) AS ip FROM " + (&System{}).TableName()).Values(&maps)

	originalIp := ""
	if 0 < len(maps) && nil != maps[0]["ip"] {
		originalIp = maps[0]["ip"].(string)
	}
	if ip != originalIp {
		if res, err := o.Raw("UPDATE "+(&System{}).TableName()+" SET ip = INET6_ATON(?), ip_updated_date = ?, updated_date = ?", ip, time.Now(), time.Now()).Exec(); err == nil {
			numRows, _ = res.RowsAffected()
			return int(numRows), err
		}
	}

	return 0, err
}

// UpdateSMTP
func UpdateSMTP(host, port, username, password, senderName, senderEmail string, authentication bool) (int, error) {
	o := orm.NewOrm()
	var (
		numRows int64
		err     error
	)

	if "" == port {
		port = "0"
	}

	if authentication {
		if res, err := o.Raw("UPDATE "+(&System{}).TableName()+" SET smtp_host = ?, smtp_port = ?, smtp_authentication = ?, smtp_username = ?, smtp_password = ?, smtp_sender_name = ?, smtp_sender_email = ?, updated_date = ?", host, port, authentication, username, password, senderName, senderEmail, time.Now()).Exec(); err == nil {
			numRows, _ = res.RowsAffected()
			return int(numRows), err
		}
	} else {
		if res, err := o.Raw("UPDATE "+(&System{}).TableName()+" SET smtp_host = ?, smtp_port = ?, smtp_authentication = ?, smtp_sender_name = ?, smtp_sender_email = ?, updated_date = ?", host, port, authentication, senderName, senderEmail, time.Now()).Exec(); err == nil {
			numRows, _ = res.RowsAffected()
			return int(numRows), err
		}
	}

	return 0, err
}

// SMTPValid
func SMTPValid(valid bool) (int, error) {
	o := orm.NewOrm()
	var (
		numRows int64
		err     error
	)

	if res, err := o.Raw("UPDATE "+(&System{}).TableName()+" SET smtp_valid = ?, updated_date = ?", valid, time.Now()).Exec(); err == nil {
		numRows, _ = res.RowsAffected()
		return int(numRows), err
	}

	return 0, err
}

func UpdateSubspaceInformation(subspaceVersion string, subspaceBuild uint) (int64, error) {
	o := orm.NewOrm()
	if res, err := o.Raw("UPDATE "+(&System{}).TableName()+" SET subspace_version = ?, subspace_build_number = ?, updated_date = ?", subspaceVersion, subspaceBuild, time.Now()).Exec(); err == nil {
		numRows, _ := res.RowsAffected()
		return numRows, err
	} else {
		return 0, err
	}
}

func UpdateVpnServerInformation(vpnServerVersion string, vpnServerBuild uint, vpnServerPassword string, vpnServerCmdPort uint, vpnServerDefaultHub string) (int64, error) {
	o := orm.NewOrm()
	if res, err := o.Raw("UPDATE "+(&System{}).TableName()+" SET "+
		"vpn_server_version = ?, "+
		"vpn_server_build_number = ?, "+
		"vpn_server_administration_password = ?, "+
		"vpn_server_administration_port = ?, "+
		"vpn_hub_name = ?, "+
		"updated_date = ?", vpnServerVersion, vpnServerBuild, vpnServerPassword, vpnServerCmdPort, vpnServerDefaultHub, time.Now()).Exec(); err == nil {
		numRows, _ := res.RowsAffected()
		return numRows, err
	} else {
		return 0, err
	}
}

func UpdateConfigSchemaVersion(version uint) (int64, error) {
	o := orm.NewOrm()
	if res, err := o.Raw("UPDATE "+(&System{}).TableName()+" SET config_schema_version = ?, updated_date = ?", version, time.Now()).Exec(); err == nil {
		numRows, _ := res.RowsAffected()
		return numRows, err
	} else {
		return 0, err
	}
}

func UpdateUserSchemaVersion(version uint) (int64, error) {
	o := orm.NewOrm()
	if res, err := o.Raw("UPDATE "+(&System{}).TableName()+" SET user_schema_version = ?, updated_date = ?", version, time.Now()).Exec(); err == nil {
		numRows, _ := res.RowsAffected()
		return numRows, err
	} else {
		return 0, err
	}
}

func UpdateProfileSchemaVersion(version uint) (int64, error) {
	o := orm.NewOrm()
	if res, err := o.Raw("UPDATE "+(&System{}).TableName()+" SET profile_schema_version = ?, updated_date = ?", version, time.Now()).Exec(); err == nil {
		numRows, _ := res.RowsAffected()
		return numRows, err
	} else {
		return 0, err
	}
}

func UpdateBackupTime(t time.Time) (int64, error) {
	o := orm.NewOrm()
	if res, err := o.Raw("UPDATE "+(&System{}).TableName()+" SET latest_backup_date = ?, updated_date = ?", t, time.Now()).Exec(); err == nil {
		numRows, _ := res.RowsAffected()
		return numRows, err
	} else {
		return 0, err
	}
}
