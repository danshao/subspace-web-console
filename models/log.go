package models

import (
	"time"
	// "github.com/astaxie/beego"
	"database/sql"

	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Log struct {
	Id              int            `orm:"pk;auto" json:"id"` // auto-increment
	Type            int            `json:"type"`             // required
	Source          int            `json:"source"`           // required
	RawLog          string         `json:"raw_message"`      // required
	UserId          sql.NullInt64  `orm:"null" json:"user_id"`
	Hub             sql.NullString `orm:"null" json:"hub"`
	ProfileUsername sql.NullString `orm:"null" json:"profile_username"`
	SessionName     sql.NullString `orm:"null" json:"session_name"`
	OperatorId      sql.NullInt64  `orm:"null" json:"operator_id"`
	ClientIp        sql.NullString `orm:"null" json:"client_ip"`
	LogTime         time.Time      `orm:"type(datetime)" json:"log_time"` // required
	CreatedDate     time.Time      `orm:"auto_now_add;type(datetime)"`    // required
}

// var logSource = map[int]string{
// 	1: "SoftEther server log",
// 	2: "SoftEther security log",
// 	3: "SoftEther packet log",
// 	4: "Subspace console log",
// 	5: "Status check monit log",
// }

var logTypeMap = map[string]int{
	"undefined":                -1,
	"user_sign_in":             0,
	"user_sign_out":            1,
	"user_session_timeout":     2,
	"user_create":              10,
	"user_edit":                11,
	"user_disable":             12,
	"user_enable":              13,
	"user_delete":              14,
	"profile_create":           20,
	"profile_edit":             21,
	"profile_disable":          22,
	"profile_enable":           23,
	"profile_delete":           24,
	"profile_download_apple":   25,
	"profile_download_windows": 26,
	"sessions_connect":         30,
	"session_disconnect":       31,
	"session_kicked":           32,
	"session_auth_fail":        33,
	"hostname_update":          40,
	"uuid_update":              41,
	"preshared_key_update":     42,
	"subspace_start":           100,
	"subspace_stop":            101,
	"vpn_server_start":         110,
	"vpn_server_stop":          111,
	"config_backup":            120,
	"config_restore":           121,
}

func (l *Log) TableName() string {
	return "logs"
}

//TODO Handle client_ip with MySQL INET6_NTOA
func ScanLog(limit ...int) []Log {
	var (
		limitRows []int
		logs      []Log
	)
	if limitRows = limit; len(limitRows) == 0 {
		limitRows = append(limitRows, 1000)
	}
	o := orm.NewOrm()
	if _, err := o.QueryTable((&Log{}).TableName()).Limit(limitRows[0]).OrderBy("-id").All(&logs); err != nil {
		beego.Debug("Scan log error: ", err)
	}
	return logs
}

//TODO Handle client_ip with MySQL INET6_NTOA
func GetLogWithConditions(logType string, userID int, limit ...int) []Log {
	var (
		limitRows []int
		logs      []Log
	)
	if limitRows = limit; len(limitRows) == 0 {
		limitRows = append(limitRows, 1000)
	}
	o := orm.NewOrm()
	cond := orm.NewCondition()
	if logType != "" {
		cond = cond.And("type", logTypeMap[logType])
	}
	if userID != -1 {
		cond = cond.And("user_id", userID)
	}

	_, err := o.QueryTable((&Log{}).TableName()).SetCond(cond).Limit(limitRows[0]).All(&logs)
	if err != nil {
		beego.Debug("Scan log error: ", err)
	}

	return logs
}

//TODO change to func WriteLog1(log Log) (int, error) {}
// We have to write client_ip via mysql INET6_ATON, so we using raw query.
func WriteLog(logType, rawMsg, profileUsername, sessionName, clientIP string, userID, operatorID int) (int, error) {

	o := orm.NewOrm()
	log := Log{
		Source:     4, // log source type for web console
		Type:       logTypeMap[logType],
		RawLog:     rawMsg,
		LogTime:    time.Now(),
		UserId:     sql.NullInt64{Int64: int64(userID), Valid: 0 < userID},
		OperatorId: sql.NullInt64{Int64: int64(operatorID), Valid: 0 < operatorID},
		//Hub: sql.NullString{String: hub, Valid: "" != hub},
		ProfileUsername: sql.NullString{String: profileUsername, Valid: "" != profileUsername},
		SessionName:     sql.NullString{String: sessionName, Valid: "" != sessionName},
		ClientIp:        sql.NullString{String: clientIP, Valid: "" != clientIP},
	}

	// Required columns: source, type, log_time, raw_log
	values := make([]string, 0)
	values = append(values,
		fmt.Sprintf("`source` = %d", log.Source),
		fmt.Sprintf("`type` = %d", log.Type),
		fmt.Sprintf("`log_time` = '%s'", log.LogTime.Format("2006-01-02 15:04:05")),
		fmt.Sprintf("`raw_log` = '%s'", log.RawLog),
	)

	// Optional columns: user_id, operator_id, hub, profile_username, session_name, client_ip
	if value, _ := log.UserId.Value(); nil != value {
		values = append(values, fmt.Sprintf("`user_id` = %d", userID))
	}
	if value, _ := log.OperatorId.Value(); nil != value {
		values = append(values, fmt.Sprintf("`operator_id` = %d", operatorID))
	}
	//TODO Hub name is missing in args
	//if value, _ := log.Hub.Value(); nil != value {
	//	values = append(values, fmt.Sprintf("`hub` = '%s'", hubName))
	//}
	if value, _ := log.ProfileUsername.Value(); nil != value {
		values = append(values, fmt.Sprintf("`profile_username` = '%s'", profileUsername))
	}
	if value, _ := log.SessionName.Value(); nil != value {
		values = append(values, fmt.Sprintf("`session_name` = '%s'", sessionName))
	}
	if value, _ := log.ClientIp.Value(); nil != value {
		values = append(values, fmt.Sprintf("`client_ip` = INET6_ATON('%s')", clientIP))
	}

	rawSql := fmt.Sprintf("INSERT INTO `logs` SET %s;", strings.Join(values, ", "))
	fmt.Println(rawSql)

	result, err := o.Raw(rawSql).Exec()
	if err != nil {
		beego.Debug("Error occurred while writing log: ", err)
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if nil != err {
		beego.Debug("Error occurred while writing log: ", err)
	}
	return int(lastId), err
}
