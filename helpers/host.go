package helpers

import (
	"net"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

/*
  檢察 DNS 是否指向這台機器
  使用者可以手動測試系統 host 或是 profile 的 host.
  usage:
  1. In hostname setting page, validate hostname by user.
  2. profiles list in user info, validate hostname by user.
*/

// DNSMatch compares ip associated to the hostname with local instance ipv4
func DNSMatch(hostString string) bool {
	addrs, _ := net.LookupIP(hostString)
	c := ec2metadata.New(session.New())
	originIPString, _ := c.GetMetadata("public-ipv4")
	originIPStringLocal, _ := c.GetMetadata("local-ipv4")
	for _, addr := range addrs {
		if addr.String() == originIPString || addr.String() == originIPStringLocal {
			return true
		}
	}
	return false
}
