package system

import (
	"gitlab.ecoworkinc.com/Subspace/web-console/models"
	"gitlab.ecoworkinc.com/Subspace/web-console/helpers/aws"
	"github.com/astaxie/beego"
)

//Determine VPN host in mobile config.
func GetVpnHostForClient() (string, error) {
	sys, err := models.GetSystemInfo()
	if nil != err {
		return "", err
	}

	if "" != sys.Host {
		return sys.Host, nil
	} else {
		if currentIp, err := RefreshCurrentServerIp(); nil == err {
			return currentIp, nil
		} else {
			return "", err
		}
	}
}

func RefreshCurrentServerIp() (string, error) {
	softetherHost := beego.AppConfig.String("host")
	runmode := beego.AppConfig.String("RunMode")

	instanceIpV4, err := aws.GetInstancePublicIpV4()
	if nil != err && "dev" == runmode {
		instanceIpV4 = softetherHost
	}
	if "" != instanceIpV4 {
		models.UpdateSystemIP(instanceIpV4)
		return instanceIpV4, nil
	} else {
		return "", err
	}
}