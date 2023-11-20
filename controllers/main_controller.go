package controllers

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	geoip2 "github.com/oschwald/geoip2-golang"
	"gitlab.ecoworkinc.com/Subspace/server-status-api/serverstatus"
	"gitlab.ecoworkinc.com/Subspace/subspace-utility/subspace/repository"
	"gitlab.ecoworkinc.com/Subspace/web-console/definition"
	"gitlab.ecoworkinc.com/Subspace/web-console/form"
	"gitlab.ecoworkinc.com/Subspace/web-console/helpers"
	"gitlab.ecoworkinc.com/Subspace/web-console/helpers/system"
	"gitlab.ecoworkinc.com/Subspace/web-console/models"
	"gitlab.ecoworkinc.com/Subspace/web-console/business"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Prepare() {
	// navbar user info
	if c.GetSession("auth") != nil {
		if c.GetSession("auth").(bool) {
			c.Data["Username"] = c.GetSession("username")
			c.Data["UserLink"] = c.GetSession("userlink")
			if c.GetSession("role") != nil {
				c.Data["Role"] = "(" + strings.Title(c.GetSession("role").(string)) + ")"
			} else {
				c.Data["Role"] = "()"
			}
		}
	}
}

// SignIn defines the default sign in page
func (c *MainController) SignIn() {
	currentUserID := c.GetSession("userID")
	if nil != currentUserID {
		currentUser := models.GetUserWithID(currentUserID.(int))
		if _, err := business.CanSignIn(currentUser); nil == err {
			c.Redirect(c.URLFor("MainController.Index"), 302)
			return
		}
	}

	c.Data["Title"] = DefaultTitle + "Sign In"
	c.TplName = "content/main/sign_in.tpl"
	c.Layout = LayoutPath["init"]["Layout"]
	c.LayoutSections = make(map[string]string)
	c.LayoutSections = LayoutPath["init"]

	beego.Debug("origin uri from session: ", c.GetSession("origin-uri"))

	var (
		flash   = beego.ReadFromRequest(&c.Controller)
		_       = flash
		success = new(bool)
	)

	*success = false

	if strings.ToLower(c.Ctx.Request.Method) == "post" {
		p := form.UserSignInForm{}

		if err := c.ParseForm(&p); err != nil {
			msg := "Parse form input data error: " + err.Error()
			beego.Error(msg)
			ShowErrorMessage(c.Controller, msg)
			return
		}

		c.Data[FORM_DATA] = p

		if _, err := IsFormValid(&p); nil != err {
			msg := "Form validate fail: " + err.Error()
			beego.Error(msg)
			ShowErrorMessage(c.Controller, msg)
			return
		}

		user := models.AuthenticateUser(p.Email, p.Password)
		if nil == user {
			msg := "Username or password is incorrect."
			beego.Debug(msg)
			ShowErrorMessage(c.Controller, msg)
			models.WriteLog("user_sign_in", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), 0, 0)
			return
		}

		if _, err := business.CanSignIn(user); nil != err {
			msg := err.Error()
			beego.Debug(msg)
			ShowErrorMessage(c.Controller, msg)
			models.WriteLog("user_sign_in", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), 0, 0)
			return
		}

		// Can sign in
		c.SetSession("username", user.Email)
		c.SetSession(definition.SESSION_USER_ID, user.Id)
		c.SetSession("role", user.Role)
		c.SetSession("userlink", "/users/"+strconv.Itoa(user.Id))
		c.SetSession("auth", true)

		*success = true
		models.WriteLog("user_sign_in", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), 0, 0)

		originURI := c.GetSession("origin-uri")
		if originURI != nil && !strings.Contains(originURI.(string), ".js") {
			if url := c.GetSession("origin-uri").(string); url != "" {
				c.Redirect(url, 302)
				c.DelSession("origin-uri")
			}
		} else {
			c.Redirect(c.URLFor("MainController.Index"), 302)
			return
		}
	} else {
		flash := beego.ReadFromRequest(&c.Controller)
		if n, ok := flash.Data["notice"]; ok {
			beego.Debug("flash notice: ", n)
		}
		if n, ok := flash.Data["error"]; ok {
			beego.Debug("flash error: ", n)
		}
	}
}

func (c *MainController) SignOut() {
	c.DelSession("auth")
	c.DelSession("username")
	c.DelSession(definition.SESSION_USER_ID)
	c.DelSession("role")
	c.DelSession("userlink")
	success := new(bool)
	*success = true
	models.WriteLog("user_sign_out", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), 0, 0)

	c.Redirect(c.URLFor("InitController.SignIn"), 302)
}

// Index defines the main dashboard page
func (c *MainController) Index() {
	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("Dashboard", "main", "index.tpl")
	c.LayoutSections["Scripts_Google_Map"] = "layout/main/_scripts_google_map.tpl"

	// Server Status
	var (
		redisHost     = beego.AppConfig.String("host")
		redisPort     = "6379"
		redisServer   = redisHost + ":" + redisPort
		password      = ""
		db            = 0
		timeSliceHour []string
		m, s, t       string
	)
	client := serverstatus.GetClient(redisServer, password, db)

	// Uptime
	upTime := serverstatus.GetVPNRunningTime(client)
	// Parse duration: HHhMMmxSSs
	duration, _ := time.ParseDuration(strconv.Itoa(upTime) + "s")
	// Transfer to DDdHHhMMm
	timeSlice := strings.Split(duration.String(), "m")
	if len(timeSlice) > 1 { // must have minute, but not sure about hour
		timeSliceHour = strings.Split(timeSlice[0], "h")
		if len(timeSliceHour) > 1 { // must have hour
			dInt, _ := strconv.Atoi(timeSliceHour[0])
			d := strconv.Itoa(dInt / 24)
			hInt, _ := strconv.Atoi(timeSliceHour[0])
			h := strconv.Itoa(hInt % 24)
			m = timeSliceHour[1]
			s = timeSlice[1]
			if d != "0" {
				t = d + "d" + h + "h" + m + "m" + s
			} else {
				t = h + "h" + m + "m" + s
			}
		} else { // no hour
			m = timeSlice[0]
			s = timeSlice[1]
			t = m + "m" + s
		}
	} else { // only second exists
		t = duration.String()
	}
	c.Data["UpTime"] = t
	beego.Debug("Server running time (seconds):", strconv.Itoa(upTime))

	// Status
	serverStatus := serverstatus.GetServerStatus(client)
	c.Data["Availability"] = serverStatus
	beego.Debug("Server Status:", serverStatus)

	// Total traffic
	totalIncoming := serverstatus.GetVPNIncoming(client)
	totalOutgoing := serverstatus.GetVPNOutgoing(client)
	beego.Debug("Total Incoming Traffic (bytes):", totalIncoming, "| Total Outgoing Traffic (bytes):", totalOutgoing)
	c.Data["TotalTrafficIncoming"] = helpers.ByteSizeFmt(uint64(totalIncoming))
	c.Data["TotalTrafficOutgoing"] = helpers.ByteSizeFmt(uint64(totalOutgoing))

	// Session Status
	var (
		serverIP  = beego.AppConfig.String("host")
		dataList  []map[string]string
		dataChunk = make(map[string]string)
		// softetherPassword = beego.AppConfig.String("softehter_management_password")
		// userID         string
		// DataList       = make(map[string]interface{})
	)

	sessionRepository := repository.InitSessionRepositoryWithHost(serverIP)
	sessions, err := sessionRepository.GetAllSessions()
	if err == nil {
		var userIDSlice []int
		for _, session := range sessions {
			if session.SessionName == "SID-SECURENAT-1" {
				continue
			}
			userID, _ := strconv.Atoi(strings.Split(session.UserNameAuthentication, "_")[0])
			userIDSlice = append(userIDSlice, userID)
		}

		users := models.GetUsersWithIDs(userIDSlice)
		if users != nil {
			userMap := make(map[int]map[string]string)
			for _, u := range users {
				userMap[u.Id] = map[string]string{
					"Email": u.Email,
					"Alias": u.Alias,
				}
			}
			for _, session := range sessions {
				if session.SessionName == "SID-SECURENAT-1" {
					continue
				}
				incomingUint, _ := strconv.ParseUint(session.IncomingDataSize, 10, 64)
				incoming := helpers.ByteSizeFmt(incomingUint)
				outgoingUint, _ := strconv.ParseUint(session.OutgoingDataSize, 10, 64)
				outgoing := helpers.ByteSizeFmt(outgoingUint)
				userID, _ := strconv.Atoi(strings.Split(session.UserNameAuthentication, "_")[0])
				dataChunk = map[string]string{
					"UserID":       strings.Split(session.UserNameAuthentication, "_")[0],
					"UserEmail":    userMap[userID]["Email"],
					"UserAlias":    userMap[userID]["Alias"],
					"SessionName":  session.SessionName,
					"SourceIP":     session.ClientIPAddress,
					"SessionStart": helpers.LocalTimeFmt(session.ConnectionStartedAt),
					"IncomingByte": incoming,
					"OutgoingByte": outgoing,
				}
				dataList = append(dataList, dataChunk)
			}
		}
	}
	c.Data["DataList"] = dataList

	//
	// GeoIP
	//
	type location struct {
		Lat  float64
		Lng  float64
		Info string
	}
	var locS []location

	// Open GeoLite2-City DB
	geodb, err := geoip2.Open("./static/resource/GeoLite2-City.mmdb")
	if err != nil {
		beego.Error(err)
	}
	defer geodb.Close()

	/*
	**  Parse IPs
	 */

	// Server IP
	if instanceIpV4, err := system.RefreshCurrentServerIp(); nil == err {
		checkIPValid := net.ParseIP(instanceIpV4)
		record, err := geodb.City(checkIPValid)
		if err != nil {
			beego.Warn(err)
		}

		loc := location{record.Location.Latitude, record.Location.Longitude, "server"}
		locS = append(locS, loc)
	}

	// Client IP
	for _, v := range dataList {
		checkIPValid := net.ParseIP(v["SourceIP"])
		record, err := geodb.City(checkIPValid)
		if err != nil {
			beego.Warn(err)
		}

		loc := location{record.Location.Latitude, record.Location.Longitude, v["UserEmail"]}
		locS = append(locS, loc)
	}

	c.Data["location"] = locS
}

func (c *MainController) ForgotPassword() {
	c.Data["Title"] = DefaultTitle + "Password Recovery"
	c.TplName = "content/main/forgot_password.tpl"
	c.Layout = LayoutPath["init"]["Layout"]
	c.LayoutSections = make(map[string]string)
	c.LayoutSections = LayoutPath["init"]
	c.LayoutSections["Scripts_Custom"] = "layout/init/_scripts_custom.tpl"

	var (
		flash      = beego.ReadFromRequest(&c.Controller)
		_          = flash
		success    = new(bool)
		instanceID string
	)

	*success = false

	if runmode == "dev" {
		instanceID = "12345678"
	} else if runmode == "prod" {
		c := ec2metadata.New(session.New())
		instanceID, _ = c.GetMetadata("instance-id")
	}

	if strings.ToLower(c.Ctx.Request.Method) == "post" {
		p := form.PasswordRecoveryForm{}
		if err := c.ParseForm(&p); err != nil {
			msg := "Parse form input data error: " + err.Error()
			beego.Error(msg)
			return
		}

		c.Data[FORM_DATA] = p

		if _, err := IsFormValid(&p); nil != err {
			msg := err.Error()
			beego.Debug(msg)
			ShowErrorMessage(c.Controller, msg)
			return
		}

		user, _ := models.GetUserWithEmail(p.Email)
		if user.Role == "admin" && (p.InstanceID == instanceID || p.InstanceID == instanceID[2:]) { // Instance ID match and user is admin
			models.UpdateUserWithID(user.Id, p.Password, user.Alias, user.Role)

			*success = true
			models.WriteLog("user_password_recovery", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), 0, 0)

			msg := "Successfully changed password for " + p.Email + ". Please login with your new credentials."
			beego.Debug(msg)
			ShowNoticeMessage(c.Controller, msg)

			c.Redirect(c.URLFor("MainController.SignIn"), 302)
		} else {
			models.WriteLog("user_password_recovery", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), 0, 0)

			msg := "Email or Instance ID is incorrect."
			beego.Debug(msg)
			ShowErrorMessage(c.Controller, msg)
		}
	}
}

func (c *MainController) Maintenance() {
	c.Data["Title"] = DefaultTitle + "Dashboard"
	// c.Data["json"] = "This page is temporarily down for maintenance."
	c.Data["json"] = c.Ctx.Input.Param(":param")
	// c.Abort("500")

	c.ServeJSON()
}
