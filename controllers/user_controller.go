package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"gitlab.ecoworkinc.com/Subspace/softetherlib/softether"
	"gitlab.ecoworkinc.com/Subspace/subspace-utility/subspace/repository"
	"gitlab.ecoworkinc.com/Subspace/vpn-profile-generator/vpnprofile"
	"gitlab.ecoworkinc.com/Subspace/web-console/form"
	"gitlab.ecoworkinc.com/Subspace/web-console/helpers"
	"gitlab.ecoworkinc.com/Subspace/web-console/helpers/mail"
	"gitlab.ecoworkinc.com/Subspace/web-console/helpers/system"
	"gitlab.ecoworkinc.com/Subspace/web-console/models"
)

// UserController beego.Controller for user management
type UserController struct {
	beego.Controller
}

func (c *UserController) Prepare() {
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

// ListUsers displays a list of current users
func (c *UserController) ListUsers() {
	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("User Management", "user", "index.tpl")

	var (
		flash = beego.ReadFromRequest(&c.Controller)
		_     = flash

		qs = models.GetAllUsers()
	)

	c.Data["UserList"] = qs
}

// UserInfo gets the information of a specific user
func (c *UserController) UserInfo() {
	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("User Management", "user", "user_get.tpl")
	c.LayoutSections["Scripts_Custom"] = "layout/user/_scripts_custom.tpl"

	type sessionDetail struct {
		SessionName         string
		ClientIPAddress     string
		ConnectionStartedAt string
		IncomingDataSize    string
		OutgoingDataSize    string
	}

	type profileDetail struct {
		Id            int              `json:"id"`
		Username      string           `json:"username"`
		Description   string           `json:"description"`
		Enabled       bool             `json:"enabled"`
		LoginCount    int              `json:"login_count"`
		IncomingBytes string           `json:"incoming_bytes"`
		OutgoingBytes string           `json:"outgoing_bytes"`
		RevokedDate   time.Time        `json:"revoked_date"`
		LastLoginDate string           `json:"last_login_data"`
		SessionList   []*sessionDetail `json:"session_list"`
		UserEnabled   bool             `json:"user_enabled"`
	}

	var (
		profileList   []profileDetail
		totalUpload   float64
		totalDownload float64

		userID, _        = strconv.Atoi(c.Ctx.Input.Param(":id"))
		user             = models.GetUserWithID(userID)
		profiles         = models.GetProfilesWithUserID(userID)
		profileSnapshots = models.GetProfileSnapshotsWithUserID(userID)

		flash = beego.ReadFromRequest(&c.Controller)
		_     = flash
	)

	// If user is not exist, redirect to user list.
	if "" == user.Email {
		c.Redirect(c.URLFor("UserController.ListUsers"), 302)
		return
	}

	// Profile Transferred Data
	for _, v := range profileSnapshots {
		totalUpload += v.OutgoingBytes
		totalDownload += v.IncomingBytes
	}

	// Session Status
	var (
		serverIP = beego.AppConfig.String("host")
	)
	sessionRepository := repository.InitSessionRepositoryWithHost(serverIP)
	sessions, err := sessionRepository.GetSessionsByUserId(userID)
	if err != nil {
		beego.Error("Error occurred while getting session details by user ID: ", err)
	}

	for _, profile := range profiles {
		var (
			sessionList []*sessionDetail
		)

		for _, session := range sessions {
			if session.UserNameAuthentication == profile.Username {
				incoming, _ := strconv.Atoi(session.IncomingDataSize)
				outgoing, _ := strconv.Atoi(session.OutgoingDataSize)
				s := &sessionDetail{
					SessionName:         session.SessionName,                               // string
					ClientIPAddress:     session.ClientIPAddress,                           // string
					ConnectionStartedAt: helpers.LocalTimeFmt(session.ConnectionStartedAt), // string
					IncomingDataSize:    helpers.ByteSizeFmt(uint64(incoming)),             // string
					OutgoingDataSize:    helpers.ByteSizeFmt(uint64(outgoing)),             // string
				}
				sessionList = append(sessionList, s)
			} else {
				continue
			}
		}

		p := profileDetail{
			Id:            profile.Id,                                         // int
			Username:      profile.Username,                                   // string
			Description:   profile.Description,                                // string
			Enabled:       profile.Enabled,                                    // bool
			LoginCount:    profile.LoginCount,                                 // int
			IncomingBytes: helpers.ByteSizeFmt(uint64(profile.IncomingBytes)), // float64
			OutgoingBytes: helpers.ByteSizeFmt(uint64(profile.OutgoingBytes)), // float64
			RevokedDate:   profile.RevokedDate,                                // time.Time
			LastLoginDate: helpers.LocalTimeFmt(profile.LastLoginDate),        // string
			SessionList:   sessionList,
			UserEnabled:   user.Enabled,
		}
		profileList = append(profileList, p)
	}

	beego.Debug("profile list: ", profileList)
	c.Data["User"] = user
	c.Data["ProfileList"] = profileList
	c.Data["TotalUpload"] = helpers.ByteSizeFmt(uint64(totalUpload))
	c.Data["TotalDownload"] = helpers.ByteSizeFmt(uint64(totalDownload))
}

// UserCreate creates a new user (Rollback needed)
func (c *UserController) UserCreate() {
	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("User Management", "user", "user_create.tpl")
	c.LayoutSections["Scripts_Custom"] = "layout/user/_scripts_custom.tpl"

	var (
		flash   = beego.ReadFromRequest(&c.Controller)
		_       = flash
		success = new(bool)
	)

	systemInfo, _ := models.GetSystemInfo()
	c.Data["CanSendEmail"] = systemInfo.SmtpValid

	*success = false
	if strings.ToLower(c.Ctx.Request.Method) == "post" {
		p := form.UserCreateForm{}

		if err := c.ParseForm(&p); err != nil { // Check if form is parsed correctly
			msg := "Parse form error: " + err.Error()
			beego.Error(msg)
			ShowErrorMessage(c.Controller, msg)
			c.Redirect(c.URLFor("UserController.UserCreate"), 302)
			return
		}

		if valid, err := IsFormValid(&p); !valid {
			msg := "Validate form error: " + err.Error()
			beego.Debug(msg)
			ShowErrorMessage(c.Controller, msg)
			c.Redirect(c.URLFor("UserController.UserCreate"), 302)
			return
		}

		var (
			email            = p.Email
			password         = p.Password
			role             = p.Role
			alias            = p.Alias
			generatePassword = p.GeneratePassword
			createVPNProfile = p.CreateVPNProfile
			emailDelivery    = p.EmailDelivery
			emailVerified    = false
			enabled          = true
		)

		//Get VPN host to connect.
		vpnHostForClient, err := system.GetVpnHostForClient()
		if nil != err {
			beego.Debug("Cannot fetch EC2 instance public IP.", err)

			flash := beego.NewFlash()
			flash.Error("Error: Get VPN host fail.")
			flash.Store(&c.Controller)

			c.Redirect(c.URLFor("UserController.UserCreate"), 302)
			return
		}

		// Create User Flow
		if userExists := models.UserEmailExists(email); userExists { // first check if user exists
			beego.Debug("User with email", email, "already exists.")

			flash := beego.NewFlash()
			flash.Error("Error: User with email " + email + " already exists.")
			flash.Store(&c.Controller)

			c.Redirect(c.URLFor("UserController.UserCreate"), 302)
			return
		} else { // if not, let's create the user
			if generatePassword { // if AUTO GENERATE PASSWORD checkbox is checked
				password = helpers.GeneratePassword(10)
				beego.Debug("UserCreate generated password for new user ", email, ": ", password)
			}

			// User creating sequence:
			// 1. Create Subspace user in user DB
			// 2. Create SofeEther user in SoftEther
			// 3. Create Subspace profile in profile DB
			// 4. Set password for newly created SoftEther user
			// 5. Generate VPN profile
			if userID, err := models.CreateUser(email, password, alias, role, emailVerified, enabled); err == nil {
				if createVPNProfile { // if AUTOGENERATE VPN PROFILE checkbox is checked
					var (
						timeStamp          = strconv.FormatInt(time.Now().UnixNano(), 10)
						profileUsername    = strings.Join([]string{strconv.Itoa(userID), timeStamp}, "_")
						profilePassword    = helpers.GeneratePassword(10)
						profileDescription = "default"
						vpnPsk, _          = models.GetPreSharedKey()
						se                 = softether.SoftEther{IP: VPN_HOST, Password: SOFTETHER_MANGEMENT_PASSWORD, Hub: SOFTETHER_HUB}
					)

					if returnCode := se.CreateUser(profileUsername, email, profileDescription); returnCode == 0 { // 1. Add user to SoftEther
						if returnCode = se.SetUserPassword(profileUsername, profilePassword); returnCode == 0 { // 2. Set password for newly created SoftEther user
							if profileID, err := models.CreateProfile(userID, SOFTETHER_HUB, profileUsername, profileDescription, email, profilePassword, vpnHostForClient, vpnPsk); err == nil { // 3. Insert new Profile into DB
								*success = true

								accountRepo := repository.InitVpnAccountRepositoryWithHost(beego.AppConfig.String("host"))
								err := accountRepo.SetAccountCache(profileUsername, profilePassword)
								if nil != err {
									beego.Error("Cannot set account data into redis.")
								}

								// Generate VPN profile
								systemUUID, _ := models.GetSystemUUID()

								vpnServer := vpnprofile.Server{Host: vpnHostForClient, PreSharedKey: vpnPsk}
								vpnUser := vpnprofile.User{Username: profileUsername, Password: profilePassword}
								vpnMeta := vpnprofile.Metadata{Identifier: vpnprofile.FormatMobileConfigIdentifier(systemUUID, SOFTETHER_HUB, userID, profileID), Description: profileDescription}
								vpnProfilePrefix := strings.Join([]string{strconv.Itoa(profileID), profileDescription, timeStamp}, "-")

								vpnProfile := map[string]string{
									"username":           profileUsername,
									"password":           profilePassword,
									"description":        profileDescription,
									"host":               vpnHostForClient,
									"key":                vpnPsk,
									"appleProfileName":   vpnProfilePrefix + "-apple.mobileconfig",
									"windowsProfileName": vpnProfilePrefix + "-windows.pbk",
									"appleProfilePath":   fmt.Sprintf("/users/%d/profiles/%d/download/apple", userID, profileID),
									"windowsProfilePath": fmt.Sprintf("/users/%d/profiles/%d/download/windows", userID, profileID),
									"ttlInMinutes":       fmt.Sprintf("%d", repository.VPN_ACCOUNT_TTL/time.Minute),
								}

								if generatePassword {
									flash := beego.NewFlash()
									flash.Notice("User with email " + email + " successfully created with auto-generated password '" + password + "'")
									flash.Store(&c.Controller)
								} else {
									flash := beego.NewFlash()
									flash.Notice("User with email " + email + " successfully created.")
									flash.Store(&c.Controller)
								}

								models.WriteLog("user_create", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))
								models.WriteLog("profile_create", helpers.NewSubspaceRawLog(c.Ctx, success).String(), profileUsername, "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))

								if emailDelivery {
									windowsProfile := vpnServer.GenerateProfile(vpnprofile.WINDOWS, vpnUser, vpnMeta)
									appleProfile := vpnServer.GenerateProfile(vpnprofile.APPLE, vpnUser, vpnMeta)

									mail.SendUserEmail(email, password)
									mail.SendProfileEmail(mail.VPNProfileData{
										Host:                     vpnHostForClient,
										PreSharedKey:             vpnPsk,
										Username:                 profileUsername,
										Password:                 profilePassword,
										FileNamePrefix:           vpnProfilePrefix,
										WindowsPBKContent:        windowsProfile,
										AppleMobileConfigContent: appleProfile,
									}, email)
									c.Redirect(c.URLFor("UserController.ListUsers"), 302)
									return
								} else {
									c.SetSession("vpnProfile", vpnProfile)
									c.Redirect(c.URLFor("UserController.ProfileInfo", ":id", userID, ":profile_id", profileID), 302)
									return
								}
							} else {
								beego.Error("Error while adding profile to Profiles table: ", err)

								flash := beego.NewFlash()
								flash.Error("Failed to create user. Please try again.")
								flash.Store(&c.Controller)

								// Roll back
								se.DeleteUser(profileUsername)
								beego.Warn("Roll back, delete SoftEther user: ", profileUsername)
								models.DeleteUserWithUserID(userID)
								beego.Warn("Roll back, delete Subspace user: ", strconv.Itoa(userID))
							}
						} else {
							beego.Debug("Error while setting softether user's password: ", softether.Strerror(returnCode))

							flash := beego.NewFlash()
							flash.Error("Failed to create user. Please try again.")
							flash.Store(&c.Controller)

							// Roll back
							se.DeleteUser(profileUsername)
							beego.Warn("Roll back, delete SoftEther user: ", profileUsername)
							models.DeleteUserWithUserID(userID)
							beego.Warn("Roll back, delete Subspace user: ", strconv.Itoa(userID))
						}
					} else {
						beego.Debug("Error while creating softether user: ", softether.Strerror(returnCode))

						flash := beego.NewFlash()
						flash.Error("Failed to create user. Please try again.")
						flash.Store(&c.Controller)

						// Roll back
						models.DeleteUserWithUserID(userID)
						beego.Warn("Roll back, delete Subspace user: ", strconv.Itoa(userID))
					}

					// CASE: PROFILE CREATE FAILED
					models.WriteLog("profile_create", helpers.NewSubspaceRawLog(c.Ctx, success).String(), profileUsername, "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))
					models.WriteLog("user_create", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))

					c.Redirect(c.URLFor("UserController.UserCreate"), 302)
					return
				}

				// CASE: USER CREATION SUCCESSFUL && AUTOGENERATE VPN PROFILE UNCHECKED
				*success = true
				if generatePassword {
					flash := beego.NewFlash()
					flash.Notice("User with email " + email + " successfully created with auto-generated password '" + password + "'")
					flash.Store(&c.Controller)
				} else {
					flash := beego.NewFlash()
					flash.Notice("User with email " + email + " successfully created.")
					flash.Store(&c.Controller)
				}

				models.WriteLog("user_create", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))

				if emailDelivery {
					mail.SendUserEmail(email, password)
				}
				c.Redirect(c.URLFor("UserController.ListUsers"), 302)
				return
			} else { // CASE: USER CREATION FAILED && AUTOGENERATE VPN PROFILE UNCHECKED
				beego.Debug("Failed to create user. Error", err)

				flash := beego.NewFlash()
				flash.Error("User was not successfully created. Error: " + err.Error())
				flash.Store(&c.Controller)

				models.WriteLog("user_create", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))

				c.Redirect(c.URLFor("UserController.UserCreate"), 302)
			}
		}
	}
}

// UserUpdate updates a specific user's information
func (c *UserController) UserUpdate() {
	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("User Management", "user", "user_update.tpl")
	c.LayoutSections["Scripts_Custom"] = "layout/user/_scripts_custom.tpl"

	var (
		flash = beego.ReadFromRequest(&c.Controller)
		_     = flash

		userID, _     = strconv.Atoi(c.Ctx.Input.Param(":id"))
		currentUserID = c.GetSession("userID")
		success       = new(bool)
		user          = models.GetUserWithID(userID)
	)

	c.Data["User"] = user

	// Check if user is modifying self.
	if currentUserID == userID {
		c.Data["EnableAdvancedOptions"] = false
	} else {
		c.Data["EnableAdvancedOptions"] = true
	}

	*success = false
	if strings.ToLower(c.Ctx.Request.Method) == "post" {
		p := form.UserUpdateForm{}
		if err := c.ParseForm(&p); err != nil {
			msg := "Parse form error: " + err.Error()
			beego.Error(msg)
			ShowErrorMessage(c.Controller, msg)
			c.Redirect(c.URLFor("UserController.UserUpdate", ":id", userID), 302)
			return
		}

		_, err := IsFormValid(&p)
		if nil != err {
			msg := "Parse form input data error: " + err.Error()
			beego.Debug(msg)
			ShowErrorMessage(c.Controller, msg)
			c.Redirect(c.URLFor("UserController.UserUpdate", ":id", userID), 302)
			return
		}

		// Update User
		var (
			password = p.Password
			role     = p.Role
			alias    = p.Alias
		)

		// User should not change self role
		if currentUserID == userID {
			role = ""
		}

		if _, err := models.UpdateUserWithID(userID, password, alias, role); err != nil {
			beego.Debug("Update User", user.Email, " failed with error: ", err)
			ShowErrorMessage(c.Controller, "Error: User update was unsuccessful. Please try again.")
		} else {
			*success = true

			msg := "User " + user.Email + " was successfully updated."
			beego.Debug(msg)
			ShowNoticeMessage(c.Controller, msg)
		}
		models.WriteLog("user_edit", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))

		c.Redirect(c.URLFor("UserController.UserInfo", ":id", userID), 302)
	}
}

// UserEnable enables a specific user
func (c *UserController) UserEnable() {
	var (
		userID, _     = strconv.Atoi(c.Ctx.Input.Param(ARG_USER_ID))
		currentUserID = c.GetSession("userID")

		success = new(bool)
	)

	*success = false

	//User cannot enable self
	if currentUserID == userID {
		msg := fmt.Sprintf("Error enabling user %d. Error: user cannot enable self.", userID)
		beego.Error(msg)
		ShowErrorMessage(c.Controller, msg)
		models.WriteLog("user_enable", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))
		c.Redirect(c.URLFor("UserController.UserUpdate", ":id", userID), 302)
		return
	}

	user := models.GetUserWithID(userID)
	if _, err := models.EnableUserWithID(userID); nil == err {
		*success = true
		ShowNoticeMessage(c.Controller, fmt.Sprintf("User %s was successfully enabled.", user.Email))
	} else {
		beego.Error("Error enabling user ", userID, ". Error: ", err)
		ShowErrorMessage(c.Controller, fmt.Sprintf("User %s was NOT successfully enabled.", user.Email))
	}

	models.WriteLog("user_enable", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))
	c.Redirect(c.URLFor("UserController.UserUpdate", ":id", userID), 302)
}

// UserDisable disables a specific user
func (c *UserController) UserDisable() {
	var (
		userID, _     = strconv.Atoi(c.Ctx.Input.Param(ARG_USER_ID))
		user          = models.GetUserWithID(userID)
		currentUserID = c.GetSession("userID")
		success       = new(bool)

		se = softether.SoftEther{IP: VPN_HOST, Password: SOFTETHER_MANGEMENT_PASSWORD, Hub: SOFTETHER_HUB}
	)

	*success = false

	//User cannot disable self
	if currentUserID == userID {
		beego.Error("Error disabling user ", userID, ". Error: user cannot disable self.")

		ShowErrorMessage(c.Controller, "User cannot change self status.")

		models.WriteLog("user_disable", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))
		c.Redirect(c.URLFor("UserController.UserUpdate", ":id", userID), 302)
		return
	}

	revokeDate := time.Now()

	if profiles := models.GetProfilesWithUserID(userID); len(profiles) == 0 {
		beego.Debug("No profiles found for user id", userID)
		*success = true
	} else {
		for _, profile := range profiles { // Loop through profiles
			// User disabling sequence:
			// 1. Disable SoftEther user
			// 2. Disable profile in Profiles table
			if err := disconnectSessionsWithProfileUsername(profile.Username, se, c); err == nil { // 1. Disconnect sessions
				if returnCode := se.SetUserEnabled(profile.Username, false); returnCode == 0 { // 2. Disable SoftEther user
					if _, err := models.EnableProfilesWithUserID(userID, revokeDate, false); err == nil { // 3. Disable profile in Profiles table
						*success = true
					} else {
						beego.Debug("Error while disabling profile", profile.Id, "from Profiles table. Error:", err)

						flash := beego.NewFlash()
						flash.Notice("User " + user.Email + " was NOT successfully disabled.")
						flash.Store(&c.Controller)

						// Roll back
						se.SetUserEnabled(profile.Username, true)
						beego.Warn("Roll back, enable user profile: ", profile.Username)
					}
				} else {
					beego.Debug("Error while disabling profile", profile.Username, "from SoftEther. Error:", softether.Strerror(returnCode))

					flash := beego.NewFlash()
					flash.Notice("User " + user.Email + " was NOT successfully disabled.")
					flash.Store(&c.Controller)
				}
			} else {
				msg := "User " + user.Email + " was NOT successfully disabled."
				beego.Debug("Error while disconnecting session. Error:", err)

				ShowNoticeMessage(c.Controller, msg)
			}

			models.WriteLog("profile_disable", helpers.NewSubspaceRawLog(c.Ctx, success).String(), profile.Username, "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))
		}
	}

	if *success == true { // If profiles were successfully disabled, disable the user
		if _, err := models.DisableUserWithID(userID); err == nil {
			msg := "User " + user.Email + " was successfully disabled."
			beego.Debug(msg)
			ShowNoticeMessage(c.Controller, msg)
		} else {
			*success = false

			beego.Debug("Error while disabling user ", userID, " from User DB. Error: ", err)
			ShowErrorMessage(c.Controller, "User "+user.Email+" was not successfully disabled.")
		}
	}

	models.WriteLog("user_disable", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))
	c.Redirect(c.URLFor("UserController.UserUpdate", ":id", userID), 302)
}

// UserDelete deletes a user and associated profiles (TODO: Rollback)
func (c *UserController) UserDelete() {
	var (
		userID, _         = strconv.Atoi(c.Ctx.Input.Param(":id"))
		userEmail         = models.GetUserWithID(userID).Email
		vpnServerProfiles = models.GetProfilesWithUserID(userID)

		se = softether.SoftEther{IP: VPN_HOST, Password: SOFTETHER_MANGEMENT_PASSWORD, Hub: SOFTETHER_HUB}

		success = new(bool)
	)

	*success = false
	if adminCount, err := models.GetAdminCount(); err == nil && adminCount >= 1 { // 1. Check if there is at least 1 admin
		if len(vpnServerProfiles) == 0 {
			beego.Debug("No profiles found for user id", userID)
			*success = true
		} else {
			for _, profile := range vpnServerProfiles { // 2. Loop through all profiles for the user
				// User delete sequence:
				// 1. Remove the profile from SoftEther
				// 2. Delete profile from DB
				if err := disconnectSessionsWithProfileUsername(profile.Username, se, c); err == nil { // 3. Disconnect all sessions
					if returnCode := se.DeleteUser(profile.Username); returnCode == 0 { // 4. Remove the profile from SoftEther
						if _, err := models.DeleteProfileWithID(profile.Id); err == nil { // 5. Delete profile from DB
							*success = true
						} else {
							beego.Debug("Error deleting user ", userEmail, ". Unable to delete profile ", strconv.Itoa(profile.Id), " from Profiles table. Error: ", err)

							flash := beego.NewFlash()
							flash.Error("Unable to delete user " + userEmail + ". Please try again. Error: " + err.Error())
							flash.Store(&c.Controller)
						}
					} else {
						beego.Debug("Error deleting user ", userEmail, ". Unable to delete profile ", strconv.Itoa(profile.Id), " from SoftEther. Error: ", softether.Strerror(returnCode))

						flash := beego.NewFlash()
						flash.Error("Unable to delete user " + userEmail + ". Please try again. Error: " + softether.Strerror(returnCode))
						flash.Store(&c.Controller)
					}
				} else {
					beego.Debug("Error deleting user ", userEmail, ". Error: ", err)

					flash := beego.NewFlash()
					flash.Error("Unable to delete user " + userEmail + ". Please try again. Error: " + err.Error())
					flash.Store(&c.Controller)

				}

				models.WriteLog("profile_delete", helpers.NewSubspaceRawLog(c.Ctx, success).String(), profile.Username, "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))
			}
		}

		if *success == true { // delete user
			if _, err := models.DeleteUserWithUserID(userID); err == nil {
				flash := beego.NewFlash()
				flash.Notice("User " + userEmail + " and associated profiles have been deleted.")
				flash.Store(&c.Controller)
			} else {
				*success = false

				beego.Debug("Error deleting user ", userID, " from User DB. Error: ", err)

				flash := beego.NewFlash()
				flash.Error("Unable to delete user: "+userEmail+". Please try again. Error: ", err)
				flash.Store(&c.Controller)
			}
		}
	} else {
		beego.Debug("Error deleting user ", userID, "Not allowed to delete last admin.")

		flash := beego.NewFlash()
		flash.Error("Unable to delete user: " + userEmail + " because this is the last user.")
		flash.Store(&c.Controller)
	}
	models.WriteLog("user_delete", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))

	if *success == true {
		c.Redirect(c.URLFor("UserController.ListUsers"), 302)
	} else {
		c.Redirect(c.URLFor("UserController.UserUpdate", ":id", userID), 302)
	}
}

// ProfileCreate creates a profile (Rollback needed)
func (c *UserController) ProfileCreate() {
	var (
		userID, _ = strconv.Atoi(c.Ctx.Input.Param(":id"))
		user      = models.GetUserWithID(userID)

		success = new(bool)
	)

	if !user.Enabled { // check if user is disabled
		flash := beego.NewFlash()
		flash.Error("Error: Cannot add profile to a disabled user.")
		flash.Store(&c.Controller)

		c.Redirect(c.URLFor("UserController.UserInfo", ":id", userID), 302)
		return
	}

	c.Data["userEmail"] = user.Email
	c.Data["userID"] = userID

	systemInfo, _ := models.GetSystemInfo()
	c.Data["CanSendEmail"] = systemInfo.SmtpValid

	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("User Management", "user", "profile_create.tpl")
	c.LayoutSections["Scripts_Custom"] = "layout/user/_scripts_custom.tpl"

	*success = false
	if strings.ToLower(c.Ctx.Request.Method) == "post" {
		var (
			timeStamp          = strconv.FormatInt(time.Now().UnixNano(), 10)
			profileUsername    = strings.Join([]string{strconv.Itoa(userID), timeStamp}, "_")
			profilePassword    = helpers.GeneratePassword(10)
			profileDescription = c.GetString("description")
			profileDelivery, _ = c.GetBool("delivery")
			se                 = softether.SoftEther{IP: VPN_HOST, Password: SOFTETHER_MANGEMENT_PASSWORD, Hub: SOFTETHER_HUB}
			vpnPsk, _          = models.GetPreSharedKey()
		)

		//Get VPN host to connect.
		vpnHostForClient, err := system.GetVpnHostForClient()
		if nil != err {
			beego.Debug("Cannot get VPN host for client.", err)

			flash := beego.NewFlash()
			flash.Error("Error: Get VPN host fail.")
			flash.Store(&c.Controller)

			c.Redirect(c.URLFor("UserController.UserInfo"), 302)
			return
		}

		// Profile creating sequence:
		// 1. Add user to SoftEther
		// 2. Add Profile to DB
		// 3. Set password for SoftEther user
		userEmail := models.GetUserWithID(userID).Email
		if returnCode := se.CreateUser(profileUsername, userEmail, profileDescription); returnCode == 0 { // 1. Add user to SoftEther
			if profileID, err := models.CreateProfile(userID, SOFTETHER_HUB, profileUsername, profileDescription, userEmail, profilePassword, vpnHostForClient, vpnPsk); err == nil { // 2. Add Profile to DB
				if returnCode := se.SetUserPassword(profileUsername, profilePassword); returnCode == 0 { // 3. Set password for SoftEther user
					*success = true

					accountRepo := repository.InitVpnAccountRepositoryWithHost(beego.AppConfig.String("host"))
					err := accountRepo.SetAccountCache(profileUsername, profilePassword)
					if nil != err {
						beego.Error("Cannot set account data into redis.")
					}

					// Generate VPN profile
					systemUUID, _ := models.GetSystemUUID()
					vpnServer := vpnprofile.Server{Host: vpnHostForClient, PreSharedKey: vpnPsk}
					vpnUser := vpnprofile.User{Username: profileUsername, Password: profilePassword}
					vpnMeta := vpnprofile.Metadata{Identifier: vpnprofile.FormatMobileConfigIdentifier(systemUUID, SOFTETHER_HUB, userID, profileID), Description: profileDescription}
					vpnProfilePrefix := strings.Join([]string{strconv.Itoa(profileID), profileDescription, timeStamp}, "-")

					vpnProfile := map[string]string{
						"username":           profileUsername,
						"password":           profilePassword,
						"description":        profileDescription,
						"host":               vpnHostForClient,
						"key":                vpnPsk,
						"appleProfileName":   vpnProfilePrefix + "-apple.mobileconfig",
						"windowsProfileName": vpnProfilePrefix + "-windows.pbk",
						"appleProfilePath":   fmt.Sprintf("/users/%d/profiles/%d/download/apple", userID, profileID),
						"windowsProfilePath": fmt.Sprintf("/users/%d/profiles/%d/download/windows", userID, profileID),
						"ttlInMinutes":       fmt.Sprintf("%d", repository.VPN_ACCOUNT_TTL/time.Minute),
					}

					models.WriteLog("profile_create", helpers.NewSubspaceRawLog(c.Ctx, success).String(), profileUsername, "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))

					flash := beego.NewFlash()
					flash.Notice("Successfully created new profile " + profileUsername + " for user " + user.Email + ".")
					flash.Store(&c.Controller)

					if profileDelivery {
						windowsProfile := vpnServer.GenerateProfile(vpnprofile.WINDOWS, vpnUser, vpnMeta)
						appleProfile := vpnServer.GenerateProfile(vpnprofile.APPLE, vpnUser, vpnMeta)

						mail.SendProfileEmail(mail.VPNProfileData{
							Host:                     vpnHostForClient,
							PreSharedKey:             vpnPsk,
							Username:                 profileUsername,
							Password:                 profilePassword,
							FileNamePrefix:           vpnProfilePrefix,
							WindowsPBKContent:        windowsProfile,
							AppleMobileConfigContent: appleProfile,
						}, user.Email)
						c.Redirect(c.URLFor("UserController.UserInfo", ":id", userID), 302)
						return
					} else {
						c.SetSession("vpnProfile", vpnProfile)
						c.Redirect(c.URLFor("UserController.ProfileInfo", ":id", userID, ":profile_id", profileID), 302)
						return
					}
				} else {
					beego.Debug("Error while creating user in SoftEther: ", softether.Strerror(returnCode))

					flash := beego.NewFlash()
					flash.Error("Error creating new profile. Error: " + softether.Strerror(returnCode))
					flash.Store(&c.Controller)

					// Roll back
					models.DeleteProfileWithID(profileID)
					beego.Warn("Roll back, delete Subspace profile: ", profileID)
					se.DeleteUser(profileUsername)
					beego.Warn("Roll back, delete SoftEther user: ", profileUsername)
				}
			} else {
				beego.Debug("Error while inserting profile in Profiles table. Error: ", err)

				flash := beego.NewFlash()
				flash.Error("Error creating new profile. Error: " + err.Error())
				flash.Store(&c.Controller)

				// Roll back
				se.DeleteUser(profileUsername)
				beego.Warn("Roll back, delete SoftEther user: ", profileUsername)
			}
		} else {
			beego.Debug("Error while creating user in SoftEther: ", softether.Strerror(returnCode))

			flash := beego.NewFlash()
			flash.Error("Error creating new profile. Error: " + softether.Strerror(returnCode))
			flash.Store(&c.Controller)
		}

		models.WriteLog("profile_create", helpers.NewSubspaceRawLog(c.Ctx, success).String(), profileUsername, "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))

		c.Redirect(c.URLFor("UserController.UserInfo", ":id", userID), 302)
	}
}

// ProfileInfo gets the information of a specific profile
func (c *UserController) ProfileInfo() {
	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("User Management", "user", "profile_get.tpl")
	c.LayoutSections["VPNInfo"] = "content/init/setup_complete.tpl"

	var (
		userID, _ = strconv.Atoi(c.Ctx.Input.Param(":id"))
		user      = models.GetUserWithID(userID)

		flash = beego.ReadFromRequest(&c.Controller)
		_     = flash
	)

	c.Data["User"] = user

	if info := c.GetSession("vpnProfile"); info != nil {
		c.Data["vpnProfile"] = info
		c.DelSession("vpnProfile")
	} else {
		c.Redirect(c.URLFor("UserController.ListUsers"), 302)
	}
}

// ProfileUpdate updates information of a specific profile (Rollback needed)
func (c *UserController) ProfileUpdate() {
	var (
		userID, _    = strconv.Atoi(c.Ctx.Input.Param(":id"))
		profileID, _ = strconv.Atoi(c.Ctx.Input.Param(":profile_id"))
		description  = c.GetString("description")
		profile      = models.GetProfileWithID(profileID)
		se           = softether.SoftEther{IP: VPN_HOST, Password: SOFTETHER_MANGEMENT_PASSWORD, Hub: SOFTETHER_HUB}
		success      = new(bool)
	)

	*success = false
	// Profile updating sequence:
	// 1. Update profile in SoftEther
	// 2. Update profile description in Profiles tables
	if returnCode := se.SetUserInfo(profile.Username, "", description); returnCode == 0 { // 1. Update profile in SoftEther
		if _, err := models.UpdateProfileWithID(profileID, description); err == nil { // 2. Update profile description in Profiles tables
			*success = true

			flash := beego.NewFlash()
			flash.Notice("The profile description for " + profile.Username + " successfully updated.")
			flash.Store(&c.Controller)
		} else {
			beego.Debug("[FAILED] Error while updating user, ", profile.Username, " in MySQL. Error: ", err)

			flash := beego.NewFlash()
			flash.Error("The profile description for " + profile.Username + " did not successfully update.")
			flash.Store(&c.Controller)

			// Roll back
			se.SetUserInfo(profile.Username, "", profile.Description)
			beego.Warn("Roll back, reset profile description to the old one")
		}
	} else {
		beego.Debug("[FAILED] Error while updating user, ", profile.Username, " in SoftEther. Error: ", softether.Strerror(returnCode))

		flash := beego.NewFlash()
		flash.Error("The profile description for " + profile.Username + " did not successfully update.")
		flash.Store(&c.Controller)
	}

	models.WriteLog("profile_edit", helpers.NewSubspaceRawLog(c.Ctx, success).String(), profile.Username, "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))

	c.Redirect(c.URLFor("UserController.UserInfo", ":id", userID), 302)
}

// ProfileEnable enables a profile (Rollback needed)
func (c *UserController) ProfileEnable() {
	var (
		userID, _    = strconv.Atoi(c.Ctx.Input.Param(ARG_USER_ID))
		profileID, _ = strconv.Atoi(c.Ctx.Input.Param(ARG_PROFILE_ID))
		profile      = models.GetProfileWithID(profileID)

		success = new(bool)

		se = softether.SoftEther{IP: VPN_HOST, Password: SOFTETHER_MANGEMENT_PASSWORD, Hub: SOFTETHER_HUB}
	)

	*success = false
	// Profile NOT exist in database
	if "" == profile.Username {
		beego.Debug("Profile", profileID, "does not exist in database.")
		c.Redirect(c.URLFor("UserController.UserInfo", ":id", userID), 302)
	} else {
		// Profile enabling sequence:
		// 1. Enable SoftEther profile
		// 2. Enable the profile in Profiles table
		if returnCode := se.SetUserEnabled(profile.Username, true); returnCode == 0 { // 1. Enable SoftEther profile
			if numRows, _ := models.EnableProfileWithID(profileID); numRows > 0 { // 2. Enable the profile in Profiles table
				*success = true

				flash := beego.NewFlash()
				flash.Notice("Profile " + profile.Username + " was successfully enabled.")
				flash.Store(&c.Controller)
			} else {
				beego.Debug("[FAILED] No rows modified in Profiles table when Enabling Profile for profile: ", profile.Username)

				flash := beego.NewFlash()
				flash.Error("Profile " + profile.Username + " was not successfully enabled.")
				flash.Store(&c.Controller)

				// Roll back
				se.SetUserEnabled(profile.Username, false)
				beego.Warn("Roll back, disable user profile: ", profile.Username)
			}
		} else {
			beego.Debug("[FAILED] Enabling profile in SoftEther. Error: ", softether.Strerror(returnCode))

			flash := beego.NewFlash()
			flash.Error("Profile " + profile.Username + " was not successfully enabled.")
			flash.Store(&c.Controller)
		}
	}

	models.WriteLog("profile_enable", helpers.NewSubspaceRawLog(c.Ctx, success).String(), profile.Username, "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))
	c.Redirect(c.URLFor("UserController.UserInfo", ":id", userID), 302)
}

// ProfileDisable disables a profile (Rollback needed)
func (c *UserController) ProfileDisable() {
	var (
		userID, _    = strconv.Atoi(c.Ctx.Input.Param(ARG_USER_ID))
		profileID, _ = strconv.Atoi(c.Ctx.Input.Param(ARG_PROFILE_ID))
		profile      = models.GetProfileWithID(profileID)

		success = new(bool)

		se = softether.SoftEther{IP: VPN_HOST, Password: SOFTETHER_MANGEMENT_PASSWORD, Hub: SOFTETHER_HUB}
	)

	*success = false
	if "" == profile.Username {
		beego.Debug("Profile", profileID, "is not exist in database.")
		c.Redirect(c.URLFor("UserController.UserInfo", ":id", userID), 302)
	} else {
		// Profile disabling sequence:
		// 1. Disable the profile in SoftEther
		// 2. Disable the profile in Profiles table
		if err := disconnectSessionsWithProfileUsername(profile.Username, se, c); err == nil { // 1. Disconnect all sessions for this profile from SoftEther
			if returnCode := se.SetUserEnabled(profile.Username, false); returnCode == 0 { // 2. Disable the profile in SoftEther
				if numRows, _ := models.DisableProfileWithID(profileID); numRows > 0 { // 3. Disable the profile in Profiles table
					*success = true

					flash := beego.NewFlash()
					flash.Notice("Profile " + profile.Username + " was successfully disabled.")
					flash.Store(&c.Controller)
				} else {
					beego.Debug("[FAILED] No rows modified in Profiles table when Disabling Profile for profile: ", profile.Username)

					flash := beego.NewFlash()
					flash.Error("Profile " + profile.Username + " was not successfully disabled.")
					flash.Store(&c.Controller)

					// Roll back
					se.SetUserEnabled(profile.Username, true)
					beego.Warn("Roll back, enable user profile: ", profile.Username)
				}
			} else {
				beego.Debug("[FAILED] Disabling profile in SoftEther. Error: ", softether.Strerror(returnCode))

				flash := beego.NewFlash()
				flash.Error("Profile " + profile.Username + " was not successfully disabled.")
				flash.Store(&c.Controller)
			}
		} else {
			beego.Debug("[FAILED] Disconnecting sessions associated to the profile: ", profile.Username, ". Error: ", err)

			flash := beego.NewFlash()
			flash.Error("Profile " + profile.Username + " was not successfully disabled.")
			flash.Store(&c.Controller)
		}
	}

	models.WriteLog("profile_disable", helpers.NewSubspaceRawLog(c.Ctx, success).String(), profile.Username, "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))
	c.Redirect(c.URLFor("UserController.UserInfo", ":id", userID), 302)
}

// ProfileDelete deletes a profile (TODO: Rollback)
func (c *UserController) ProfileDelete() {
	var (
		userID, _    = strconv.Atoi(c.Ctx.Input.Param(":id"))
		profileID, _ = strconv.Atoi(c.Ctx.Input.Param(":profile_id")) // integer
		profile      = models.GetProfileWithID(profileID)

		success = new(bool)

		se = softether.SoftEther{IP: VPN_HOST, Password: SOFTETHER_MANGEMENT_PASSWORD, Hub: SOFTETHER_HUB}
	)

	if profile.Enabled { // check if profile is enabled
		flash := beego.NewFlash()
		flash.Error("Error: Disable the profile before deleting.")
		flash.Store(&c.Controller)
		c.Redirect(c.URLFor("UserController.UserInfo", ":id", userID), 302)
		return
	}

	*success = false
	// Profile deleting sequence:
	// 1. Delete profile from SoftEther
	// 2. Delete profile from profiles table
	if err := disconnectSessionsWithProfileUsername(profile.Username, se, c); err == nil { // 1. Disconnect sessions associated with the profile
		if errCode := se.DeleteUser(profile.Username); errCode == 0 { // 2. Delete profile from SoftEther
			if _, err := models.DeleteProfileWithID(profileID); err == nil { // 3. Delete profile from profiles table
				*success = true

				flash := beego.NewFlash()
				flash.Notice("The profile: " + profile.Username + " has been deleted.")
				flash.Store(&c.Controller)
			} else {
				beego.Debug("[FAILED] Deleting profile from MySQL. Error: ", err)

				flash := beego.NewFlash()
				flash.Error("Error: Deleting the profile: " + profile.Username + " failed.")
				flash.Store(&c.Controller)
			}
		} else {
			beego.Debug("[FAILED] Deleting profile from SoftEther. SoftEther Error: ", softether.Strerror(errCode))

			flash := beego.NewFlash()
			flash.Error("Error: Deleting the profile: " + profile.Username + " failed.")
			flash.Store(&c.Controller)
		}
	} else {
		beego.Debug("[FAILED] Disconnecting Session. Error: ", err)

		flash := beego.NewFlash()
		flash.Error("Error: Deleting the profile: " + profile.Username + " failed.")
		flash.Store(&c.Controller)
	}

	models.WriteLog("profile_delete", helpers.NewSubspaceRawLog(c.Ctx, success).String(), profile.Username, "", c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))

	c.Redirect(c.URLFor("UserController.UserInfo", ":id", userID), 302)
}

func (c *UserController) ProfileDownload() {
	success := new(bool)

	// Get argument from path
	platform := c.Ctx.Input.Param(ARG_PROFILE_PLATFORM)
	userID, _ := strconv.Atoi(c.Ctx.Input.Param(ARG_USER_ID))
	profileID, _ := strconv.Atoi(c.Ctx.Input.Param(ARG_PROFILE_ID))

	// Check user who access this url is authorized.
	operatorId := c.GetSession("userID").(int)
	operator := models.GetUserWithID(operatorId)
	if models.USER_ROLE_ADMIN != operator.Role {
		// Return error
		beego.Debug("Operator", operatorId, "is not admin and try to download profile.")
		c.CustomAbort(401, "Only admin can download.")
	}

	// Prepare system data for VPN settings config.
	systemInfo, err := models.GetSystemInfo()
	if nil != err {
		beego.Debug("Cannot fetch system info.", err)
		c.CustomAbort(500, "Cannot fetch system info.")
	}
	systemUUID := systemInfo.Uuid
	vpnPsk := systemInfo.PreSharedKey
	vpnHostForClient, err := system.GetVpnHostForClient()
	if nil != err {
		c.CustomAbort(500, "Cannot get host for client")
	}

	// Get profile data in MySQL.
	//TODO Cannot check profile is exist in DB or not, no error return for now.
	profile := models.GetProfileWithID(profileID)
	username := profile.Username
	description := profile.Description
	if "" == username {
		beego.Debug("Profile may not exist in database.")
		c.CustomAbort(400, "Error when fetch profile data.")
	}

	// Get VPN account password cached in Redis.
	accountRepo := repository.InitVpnAccountRepositoryWithHost(beego.AppConfig.String("host"))
	password, err := accountRepo.GetAccountCache(username)
	if nil != err {
		beego.Debug("VPN account data may not exist is Redis.", err)
		c.CustomAbort(400, "Cannot fetch account data.")
	}

	var (
		configPlatform vpnprofile.Platform
		filenameFormat string
		logType        string
	)
	switch platform {
	case "apple":
		configPlatform = vpnprofile.APPLE
		filenameFormat = "%s-apple.mobileconfig"
		logType = "profile_download_apple"
	case "windows":
		configPlatform = vpnprofile.WINDOWS
		filenameFormat = "%s-windows.pbk"
		logType = "profile_download_windows"
	default:
		beego.Debug("Operator", operatorId, "access unsupported platform", platform)
		c.CustomAbort(400, "Unsupported platform.")
	}

	*success = true
	// Write log
	models.WriteLog(logType, helpers.NewSubspaceRawLog(c.Ctx, success).String(), username, "", c.Ctx.Input.IP(), userID, operatorId)

	// Generate VPN profile for download
	vpnServer := vpnprofile.Server{Host: vpnHostForClient, PreSharedKey: vpnPsk}
	vpnUser := vpnprofile.User{Username: username, Password: password}
	vpnMeta := vpnprofile.Metadata{Identifier: vpnprofile.FormatMobileConfigIdentifier(systemUUID, SOFTETHER_HUB, userID, profileID), Description: description}
	vpnProfilePrefix := strings.Join([]string{strconv.Itoa(profileID), description, username}, "-")
	fileName := fmt.Sprintf(filenameFormat, vpnProfilePrefix)
	fileContent := vpnServer.GenerateProfile(configPlatform, vpnUser, vpnMeta)

	// Response to client
	buf := &bytes.Buffer{}
	buf.WriteString(fileContent)

	c.Ctx.Output.Header("Content-Type", "application/x-apple-aspen-config")
	c.Ctx.Output.Header("Content-Disposition", "attachment;filename="+fileName)
	c.Ctx.Output.Header("Expires", "0")
	c.Ctx.Output.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	c.Ctx.Output.Body(buf.Bytes())
}

// SessionDisconnect disconnects a session
func (c *UserController) SessionDisconnect() {
	var (
		userID, _   = strconv.Atoi(c.Ctx.Input.Param(":id"))
		sessionName = c.Ctx.Input.Param(":session_name")

		success = new(bool)

		se = softether.SoftEther{IP: VPN_HOST, Password: SOFTETHER_MANGEMENT_PASSWORD, Hub: SOFTETHER_HUB}
	)

	*success = false
	if err := se.DisconnectSession(sessionName); err == 0 { // Disconnect session from SoftEther server
		*success = true

		flash := beego.NewFlash()
		flash.Notice("Successfully disconnected the session: " + sessionName + ".")
		flash.Store(&c.Controller)
	} else { // Disconnect failed
		beego.Debug("Session Disconnect failed. Error: ", softether.Strerror(err), " Session: ", sessionName)

		flash := beego.NewFlash()
		flash.Error("Error: disconnecting the session: " + sessionName + " failed.")
		flash.Store(&c.Controller)
	}

	models.WriteLog("session_kicked", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", sessionName, c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))

	c.Redirect(c.URLFor("UserController.UserInfo", ":id", userID), 302)
}

func disconnectSessionsWithProfileUsername(profileUsername string, se softether.SoftEther, c *UserController) (err error) {
	// 1a. Get all sessions for the specified profileUsername from Redis
	var (
		serverIP  = beego.AppConfig.String("host")
		success   = new(bool)
		userID, _ = strconv.Atoi(c.Ctx.Input.Param(":id"))
	)
	sessionRepository := repository.InitSessionRepositoryWithHost(serverIP)
	sessions, err := sessionRepository.GetSessionsByProfileUserName(profileUsername)
	if err != nil {
		beego.Error("Error occurred while getting session details by user ID: ", err)
		return
	}

	*success = true
	// 1b. Loop through all sessions and disconnect
	if len(sessions) == 0 {
		return
	}
	for _, session := range sessions {
		if returnCode := se.DisconnectSession(session.SessionName); returnCode != 0 {
			*success = false
			err = errors.New("Error disconnecting session " + session.SessionName + ". Error: " + softether.Strerror(returnCode))
		}
		models.WriteLog("session_kicked", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", session.SessionName, c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))
	}

	return
}
