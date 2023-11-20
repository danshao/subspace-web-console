package aws

import (
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

//TODO 將 beego 與 Cloud platform 抽離，Web console 不需要知道執行在哪個平台執行才對。

func GetInstancePublicIpV4() (string, error) {
	c := ec2metadata.New(session.New())
	return c.GetMetadata("public-ipv4")
}

func GetEc2InstanceId() (string, error) {
	c := ec2metadata.New(session.New())
	return c.GetMetadata("instance-id")
}