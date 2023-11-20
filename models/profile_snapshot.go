package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type ProfileSnapshot struct {
	Id                         uint `orm:"pk"`
	Hub                        string
	ProfileUsername            string
	SnapshotDate               time.Time
	UserId                     int
	Description                string
	FullName                   string
	AuthType                   string
	ExpirationDate             string
	NumberOfLogins             uint
	OutgoingUnicastPackets     float64
	OutgoingUnicastTotalSize   float64
	OutgoingBroadcastPackets   float64
	OutgoingBroadcastTotalSize float64
	IncomingUnicastPackets     float64
	IncomingUnicastTotalSize   float64
	IncomingBroadcastPackets   float64
	IncomingBroadcastTotalSize float64
	CreatedOn                  time.Time
	UpdatedOn                  time.Time
	IncomingBytes              float64
	OutgoingBytes              float64
	IncomingPackets            float64
	OutgoingPackets            float64
}

func (ps *ProfileSnapshot) TableName() string {
	return "v_latest_profile"
}

// GetProfileSnapshotsWithProfileUsername gets the snapshot for a user given a profile Username
func GetProfileSnapshotsWithUserID(userID int) []ProfileSnapshot {
	var ret []ProfileSnapshot
	orm.NewOrm().QueryTable((&ProfileSnapshot{}).TableName()).Filter("user_id", userID).All(&ret)
	return ret
}
