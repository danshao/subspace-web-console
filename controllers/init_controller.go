package controllers

import (
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"gitlab.ecoworkinc.com/Subspace/softetherlib/softether"
	"gitlab.ecoworkinc.com/Subspace/subspace-utility/subspace/administration/restore"
	"gitlab.ecoworkinc.com/Subspace/vpn-profile-generator/vpnprofile"
	"gitlab.ecoworkinc.com/Subspace/web-console/helpers"
	"gitlab.ecoworkinc.com/Subspace/web-console/models"

	"fmt"

	"github.com/astaxie/beego"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	uuid "github.com/satori/go.uuid"
	"gitlab.ecoworkinc.com/Subspace/subspace-utility/subspace/repository"
	"gitlab.ecoworkinc.com/Subspace/web-console/form"
	"gitlab.ecoworkinc.com/Subspace/web-console/helpers/system"
)

var (
	VPN_HOST = beego.AppConfig.String("host")

	SOFTETHER_MANGEMENT_PASSWORD = beego.AppConfig.String("softether_management_password")
	SOFTETHER_HUB                = beego.AppConfig.String("softether_hub")

	ec2InstanceID   string
	defaultPassword = "12345678"

	runmode = beego.AppConfig.String("RunMode")
)

const (
	SUBSPACE_VERSION      = "1.1.4"
	SUBSPACE_BUILD_NUMBER = 18

	VPN_SERVER_VERSION      = "4.22"
	VPN_SERVER_BUILD_NUMBER = 9634
	VPN_SERVER_CMD_PORT     = 992

	CONFIG_SCHEMA_VERSION  = 1
	USER_SCHEMA_VERSION    = 1
	PROFILE_SCHEMA_VERSION = 1
)

type InitController struct {
	beego.Controller
}

type InitRestoreCallback struct {
	controller *InitController
}

func (c InitRestoreCallback) OnStart() {
	beego.AppConfig.Set("maintenance_mode", "true")
}
func (c InitRestoreCallback) OnCancel() {
	beego.AppConfig.Set("maintenance_mode", "false")
}
func (c InitRestoreCallback) OnSuccess(yaml string) {
	beego.AppConfig.Set("maintenance_mode", "false")
	success := new(bool)
	*success = true
	var operatorID int
	if c.controller.GetSession("userID") != nil {
		operatorID = c.controller.GetSession("userID").(int)
	}
	models.WriteLog("config_restore", helpers.NewSubspaceRawLog(c.controller.Ctx, success).String(), "", "", c.controller.Ctx.Input.IP(), 0, operatorID)

	// Log out all users.
	c.controller.DestroySession()
	cmd := exec.Command("redis-cli", "-h", beego.AppConfig.String("host"), "-n", beego.AppConfig.String("session_redis_db_number"), "flushall")
	err := cmd.Run()
	if err != nil {
		beego.Warn("redis flush db failed: ", err)
	}
}
func (c InitRestoreCallback) OnFail(e error) {
	beego.AppConfig.Set("maintenance_mode", "false")
}

func init() {

	beego.Info("Run Mode: ", runmode)

	// Determine instance IP and VPN Profile HOST
	instanceIpV4, err := system.RefreshCurrentServerIp()
	if nil != err {
		beego.Error("Get server IP fail.", err)
	}

	if runmode == "dev" {
		// Do nothing
	} else if runmode == "prod" {
		c := ec2metadata.New(session.New())

		// Initial Setup PASSWORD
		ec2InstanceID, _ = c.GetMetadata("instance-id")
		if ec2InstanceID != "" {
			defaultPassword = ec2InstanceID
			beego.Info("InstanceID: ", ec2InstanceID)
		}
	}

	if err := models.InitSystemTable(); err == nil {
		// creating UUID Version 4
		uuid := uuid.NewV4()

		// update system table
		models.UpdateSystemUUID(uuid.String())
		models.UpdateInstanceId(ec2InstanceID)
		models.UpdatePreSharedKey("subspace")
		models.UpdateSystemIP(instanceIpV4)
		models.UpdateSubspaceInformation(SUBSPACE_VERSION, SUBSPACE_BUILD_NUMBER)
		models.UpdateVpnServerInformation(VPN_SERVER_VERSION, VPN_SERVER_BUILD_NUMBER, SOFTETHER_MANGEMENT_PASSWORD, VPN_SERVER_CMD_PORT, SOFTETHER_HUB)
		models.UpdateConfigSchemaVersion(CONFIG_SCHEMA_VERSION)
		models.UpdateUserSchemaVersion(USER_SCHEMA_VERSION)
		models.UpdateProfileSchemaVersion(PROFILE_SCHEMA_VERSION)
	}
}

// SetupStart defines the first launch onboarding page page
func (c *InitController) SetupStart() {
	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("Setup", "init", "setup_start.tpl")

	var (
		flash = beego.ReadFromRequest(&c.Controller)
		_     = flash
	)

	if strings.ToLower(c.Ctx.Request.Method) == "post" {
		instanceID := strings.TrimSpace(c.GetString("instanceID"))
		switch runmode {
		case "prod":
			if instanceID == ec2InstanceID || instanceID == ec2InstanceID[2:] {
				c.SetSession("instance_id", instanceID)
				c.SetSession("init_process", true)
				c.Redirect(c.URLFor("InitController.CreateAdmin"), 302)
				return
			}
			fallthrough
		case "dev":
			if instanceID == defaultPassword && ec2InstanceID == "" {
				c.SetSession("instance_id", instanceID)
				c.SetSession("init_process", true)
				c.Redirect(c.URLFor("InitController.CreateAdmin"), 302)
				return
			}
			fallthrough
		default:
			beego.Debug("Instance ID is incorrect.")

			flash := beego.NewFlash()
			flash.Error("Error: Incorrect Instance ID")
			flash.Store(&c.Controller)

			c.Redirect(c.URLFor("InitController.SetupStart"), 302)
		}
	}
}

// CreateAdmin create admin account at first launch
func (c *InitController) CreateAdmin() {
	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("Setup", "init", "create_admin.tpl")
	c.LayoutSections["Scripts_Custom"] = "layout/init/_scripts_custom.tpl"
	var (
		flash = beego.ReadFromRequest(&c.Controller)
		_     = flash
	)

	if strings.ToLower(c.Ctx.Request.Method) == "post" {
		var (
			p       = form.CreateAdminForm{}
			success = new(bool)
		)

		*success = false

		// Parse form
		if err := c.ParseForm(&p); err != nil { // 1. Parse form data
			msg := "Parse form input data error: " + err.Error()
			beego.Error(msg)
			ShowErrorMessage(c.Controller, msg)
			return
		}

		c.Data[FORM_DATA] = p

		// Validate form
		if _, err := IsFormValid(&p); nil != err {
			msg := "Validate form input data error: " + err.Error()
			beego.Debug(msg)
			ShowErrorMessage(c.Controller, msg)
			return
		}

		// Start Create admin account!
		var (
			email              = p.Email
			password           = p.Password
			alias              = p.Alias
			profileDescription string
			role               = "admin"
			enabled            = true
			emailVerified      = false
			errMessage         = ""
		)

		if alias == "" {
			alias = "Default Admin Account"
			profileDescription = "Default Admin VPN Profile"
		} else {
			profileDescription = alias + " Default VPN Profile"
		}

		// 1. Create Subspace user in user DB
		if userID, err := models.CreateUser(email, password, alias, role, emailVerified, enabled); err == nil {
			// beego.Debug("Successfully create Subpsace user: ", userID)
			se := softether.SoftEther{IP: VPN_HOST, Password: SOFTETHER_MANGEMENT_PASSWORD, Hub: SOFTETHER_HUB}
			var (
				timeStamp       = strconv.FormatInt(time.Now().UnixNano(), 10)
				profileUsername = strings.Join([]string{strconv.Itoa(userID), timeStamp}, "_")
				profilePassword = helpers.GeneratePassword(10)
				vpnPsk, _       = models.GetPreSharedKey()
			)

			// Get instance host for client mobile config
			if vpnHostForClient, err := system.GetVpnHostForClient(); nil == err {
				// 2. Create SofeEther user in SoftEther
				if errCode := se.CreateUser(profileUsername, email, profileDescription); errCode == 0 {
					// beego.Debug("Successfully create SoftEther user: ", profileUsername)
					// 3. Create Subspace profile in profile DB
					if profileID, err := models.CreateProfile(userID, SOFTETHER_HUB, profileUsername, profileDescription, email, profilePassword, vpnHostForClient, vpnPsk); err == nil {
						// beego.Debug("Successfully create Subspace user's profile: ", profileUsername)
						// 4. Set password for newly created SoftEther user
						if errCode = se.SetUserPassword(profileUsername, profilePassword); errCode == 0 {
							// beego.Debug("Successfully set SoftEther password for profile: ", profileUsername)

							*success = true

							accountRepo := repository.InitVpnAccountRepositoryWithHost(beego.AppConfig.String("host"))
							err := accountRepo.SetAccountCache(profileUsername, profilePassword)
							if nil != err {
								beego.Error("Cannot set account data into redis.")
							}

							// 5. Generate VPN profile
							systemUUID, _ := models.GetSystemUUID()
							vpnServer := vpnprofile.Server{Host: vpnHostForClient, PreSharedKey: vpnPsk}
							vpnUser := vpnprofile.User{Username: profileUsername, Password: profilePassword}
							vpnMeta := vpnprofile.Metadata{Identifier: vpnprofile.FormatMobileConfigIdentifier(systemUUID, SOFTETHER_HUB, userID, profileID), Description: profileDescription}
							vpnProfilePrefix := strings.Join([]string{strconv.Itoa(profileID), profileDescription, timeStamp}, "-")

							windowsProfile := vpnServer.GenerateProfile(vpnprofile.WINDOWS, vpnUser, vpnMeta)
							data := []byte(windowsProfile)
							filename := "./public/" + vpnProfilePrefix + "-windows.pbk"
							_ = ioutil.WriteFile(filename, data, 0644)

							appleProfile := vpnServer.GenerateProfile(vpnprofile.APPLE, vpnUser, vpnMeta)
							data = []byte(appleProfile)
							filename = "./public/" + vpnProfilePrefix + "-apple.mobileconfig"
							_ = ioutil.WriteFile(filename, data, 0644)

							// 6. Record user info in session
							c.SetSession("username", p.Email)
							c.SetSession("userID", userID)
							c.SetSession("role", "admin")
							c.SetSession("userlink", "/users/"+strconv.Itoa(userID))
							c.SetSession("auth", true)
							beego.Debug("Set session: ", c.CruSession)

							vpnProfile := map[string]string{
								"username":           profileUsername,
								"password":           profilePassword,
								"host":               vpnHostForClient,
								"key":                vpnPsk,
								"appleProfileName":   vpnProfilePrefix + "-apple.mobileconfig",
								"windowsProfileName": vpnProfilePrefix + "-windows.pbk",
								"appleProfilePath":   fmt.Sprintf("/users/%d/profiles/%d/download/apple", userID, profileID),
								"windowsProfilePath": fmt.Sprintf("/users/%d/profiles/%d/download/windows", userID, profileID),
								"ttlInMinutes":       fmt.Sprintf("%d", repository.VPN_ACCOUNT_TTL/time.Minute),
							}
							c.SetSession("vpnProfile", vpnProfile)
							// beego.Debug("Successfully save user info in session: ", c.GetSession("vpnProfile"))

							models.WriteLog("user_create", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))
							models.WriteLog("profile_create", helpers.NewSubspaceRawLog(c.Ctx, success).String(), profileUsername, "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))

						} else { // [FAILED] 4. Set password for newly created SoftEther user
							beego.Error("Error while setting softether user's password: ", softether.Strerror(errCode))
							errMessage = "Error while setting SoftEther user."
							// Roll back step 3
							models.DeleteProfilesWithUserID(userID)
							beego.Warn("Roll back, delete profile of user: ", strconv.Itoa(userID))
							// Roll back step 2
							se.DeleteUser(profileUsername)
							beego.Warn("Roll back, delete SoftEther user: ", profileUsername)
							// Roll back step 1
							models.DeleteUserWithUserID(userID)
							beego.Warn("Roll back, delete Subspace user: ", strconv.Itoa(userID))
						}
					} else { // [FAILED] 3. Create Subspace profile in profile DB
						beego.Error("Error while creating profile for Subspace user: ", err)
						errMessage = "Error while creating Subspace VPN profile."
						// Roll back step 2
						se.DeleteUser(profileUsername)
						beego.Warn("Roll back, delete SoftEther user: ", profileUsername)
						// Roll back step 1
						models.DeleteUserWithUserID(userID)
						beego.Warn("Roll back, delete Subspace user: ", strconv.Itoa(userID))
					}
				} else { // [FAILED] 2. Create SofeEther user in SoftEther, roll back.
					beego.Error("Error while creating softether user: ", softether.Strerror(errCode))
					errMessage = "Error while creating SoftEther user."
					// Roll back step 1
					models.DeleteUserWithUserID(userID)
					beego.Warn("Roll back, delete Subspace user: ", strconv.Itoa(userID))
				}
			} else {
				beego.Error("Error while fetch VPN host for client.", err)
				errMessage = "Error while fetch VPN host for client."
				// Roll back step 1
				models.DeleteUserWithUserID(userID)
				beego.Warn("Roll back, delete Subspace user: ", strconv.Itoa(userID))
			}
		} else { // [FAILED] 1. Create Subspace user in user DB
			beego.Error("Error while creating Subspace user: ", err)
			errMessage = "Error while creating Subspace user."
		}

		if !*success {
			flash := beego.NewFlash()
			flash.Error(errMessage)
			flash.Store(&c.Controller)
			models.WriteLog("user_create", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), -1, -1)
			return
		}
		c.Redirect(c.URLFor("InitController.SetupComplete"), 302)
	}
}

// SetupComplete is the final page in the initial setup process
func (c *InitController) SetupComplete() {
	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("Setup", "init", "setup_complete.tpl")
	c.DelSession("init_process")
	if info := c.GetSession("vpnProfile"); info != nil {
		c.Data["vpnProfile"] = info
		c.DelSession("vpnProfile")
	} else {
		c.Redirect(c.URLFor("MainController.Index"), 302)
		return
	}
}

func (c *InitController) CleanInstall() {
	authenticator := false
	if runmode == "prod" && ec2InstanceID != "" && (c.GetSession("instance_id") == ec2InstanceID || c.GetSession("instance_id") == ec2InstanceID[2:]) {
		authenticator = true
	} else if runmode == "dev" && c.GetSession("instance_id") == defaultPassword {
		authenticator = true
	}
	// Start restore
	if authenticator {
		beego.Debug("receiving " + c.Ctx.Request.Method + " event from ajax")
		var (
			host   = beego.AppConfig.String("host")
			dbURI  = "subspace:subspace@tcp(" + host + ":3306)/subspace?charset=utf8&parseTime=True&loc=Local"
			folder = "./public/"
		)
		// 1. start backup / restore (POST)
		if strings.ToLower(c.Ctx.Request.Method) == "post" {
			option := c.GetString("backup_restore_option")
			beego.Debug("do: ", c.GetString("backup_restore_option"))

			if option == "restore" {
				_, fileHeader, err := c.GetFile("upload_config")
				start := false
				if err != nil {
					beego.Debug("Error: ", err)
				} else {
					beego.Debug("File name: ", fileHeader.Filename)
					beego.Debug("File header: ", fileHeader.Header)

					filePath := folder + "restore-" + fileHeader.Filename
					c.SaveToFile("upload_config", filePath)

					restoreController := restore.GetInstance()
					restoreController.SetDatabaseUri(dbURI)

					restoreController.SetCallback(InitRestoreCallback{c})
					start = restoreController.Start(filePath)
				}
				if start {
					c.Data["json"] = fileHeader.Filename
				} else {
					c.Data["json"] = "false"
				}
			} else {
				c.Data["json"] = "[Error] I'm not sure what you want, please be explicit. Simply and clearly state your purpose and try again."
			}
			c.ServeJSON()
			return
		}

		// 2. ask status (GET)
		if strings.ToLower(c.Ctx.Request.Method) == "get" {
			var option string
			c.Ctx.Input.Bind(&option, "backup_restore_option")
			beego.Debug("ask status: ", option)

			if option == "restore" {
				c.Data["json"] = "restore ok"
				restoreController := restore.GetInstance()
				status := restoreController.GetStatus()
				switch status.Step {
				case restore.IDLE:
					beego.Info("restore idle")
					c.Data["json"] = "restore idle"
				case restore.RUNNING:
					beego.Info("restore running")
					c.Data["json"] = "restore running"
				case restore.SUCCEED:
					beego.Info("restore succeeded")
					c.Data["json"] = "restore succeeded" // + status.Result
				case restore.FAILED:
					beego.Error("restore failed", status.Error)
					c.Data["json"] = "restore failed: " + status.Error.Error()
				case restore.CANCELING:
					beego.Warn("restore canceling")
					c.Data["json"] = "restore canceling"
				case restore.CANCELED:
					beego.Warn("restore canceled")
					c.Data["json"] = "restore canceled"
				case restore.UNKNOWN:
					beego.Error("something wrong")
					c.Data["json"] = "restore wrong"
				}
			} else {
				c.Data["json"] = "[Error] I'm not sure what you want, please be explicit. Simply and clearly state your purpose and try again."
			}
			c.ServeJSON()
		}
	} else {
		c.Data["json"] = "[Error] authentication failed."
		c.ServeJSON()
		return
	}
}

func (c *InitController) Maintenance() {
	c.Data["Title"] = DefaultTitle + "Setup"
	c.Data["json"] = "This page is temporarily down for maintenance."
	c.ServeJSON()
}
