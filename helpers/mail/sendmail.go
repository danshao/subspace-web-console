package mail

import (
	"bytes"
	"math/rand"
	tt "text/template"
	"time"

	"github.com/astaxie/beego"

	"gitlab.ecoworkinc.com/Subspace/web-console/models"

	"fmt"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"os"
	"strings"
)

const ATTACH_DIR = "/tmp"

type VPNProfileData struct {
	Host                     string
	PreSharedKey             string
	Username                 string
	Password                 string
	FileNamePrefix           string
	WindowsPBKContent        string
	AppleMobileConfigContent string
}

type userInfoTemplate struct {
	Username string
	Password string
}

type systemSettingsChangeInfoTemplate struct {
	DataType string
	Value    string
}

type mailServer struct {
	smtpserver string
	port       int
	sender     string
	password   string
}

func TestSMTP(emailTo string) bool {
	var (
		systemInfo, _ = models.GetSystemInfo()
		ms            = mailServer{smtpserver: systemInfo.SmtpHost, port: systemInfo.SmtpPort, sender: systemInfo.SmtpUsername, password: systemInfo.SmtpPassword}
		message       = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" xmlns=\"http://www.w3.org/1999/xhtml\"><head> <meta http-equiv=\"Content-Type\" content=\"text/html; charset=3DUTF-8\" /> <meta name=\"viewport\" content=\"width=3Ddevice-width, initial-scale=1.0\" /> <title>Subspace Email</title></head><body style=\"-ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; font: 13px/1.4 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Tahoma, sans-serif; margin: 0; padding: 10px 0 0; width: 100% !important;\" bgcolor=\"#C7C7C7\"> <style type=\"text/css\"> body { width: 100% !important; -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; margin: 0; padding: 0; font: 13px/1.4 \"Helvetica\", \"Lucida Grande\", \"Lucida Sans Unicode\", Tahoma, sans-serif; } img { max-width: 100%; outline: none; text-decoration: none; -ms-interpolation-mode: bicubic; } @media (max-width: 500px) { .email-header { display: block; } .email-header td { display: block; } .email-header tr { display: block; } .email-header thead { display: block; } .email-header tbody { display: block; } } </style> <table cellpadding=\"0\" cellspacing=\"0\" border=\"0\" align=\"center\" style=\"border: 2px solid #b5b5b5; border-collapse: collapse; margin: 0 auto; max-width: 600px; min-width: 320px; mso-table-lspace: 0pt; mso-table-rspace: 0pt;\" bgcolor=\"white\"> <tr> <td style=\"border-collapse: collapse; color: white; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px; padding: 20px 30px 10px;\" bgcolor=\"black\"> <table style=\"width: 100%;\" class=\"email-header\"> <tr> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;\"> <h1 style=\"color: white; font-family: Sans-Serif; font-size: 30px; line-height: 1.4; margin: 0 0 10px;\"> SMTP Setup Success </h1> <p style=\"color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;\"> \"Success is stumbling from failure to failure with no loss of enthusiasm.\" <br />─ Winston S. Churchill </p> </td> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px; padding-left: 20px; width: 50px;\" valign=\"top\"> <img src=\"https://subspace-assets.s3.amazonaws.com/subspace.png\" alt=\"Subspace\" style=\"-ms-interpolation-mode: bicubic; height: 50px; margin: 0; max-width: 100%; outline: none; text-decoration: none;\" /> </td> </tr> </table> </td> </tr> <tr> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px; padding: 10px 35px;\"> <h2 style=\"color: black; font-size: 15px; line-height: 1.4; margin: 1.5em 0 0.5em;\"> Looks like your SMTP settings are working. </h2> <p style=\"color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;\"> You've unlocked the ability to email in Subspace! </p> </td> </tr> <tr> <td style=\"border-collapse: collapse; color: #777; font: 11px/1.4 'Helvetica', 'Lucida Grande','Lucida Sans Unicode', Tahoma, sans-serif; padding: 10px 20px;\" align=\"center\" bgcolor=\"#eee\"> Need help with anything? Hit up <a style=\"color: #777; text-decoration: underline;\" href=\"http://support.ecowork.com\">support</a>. </td> </tr> </table></body></html>"
	)

	// compose the message
	m := gomail.NewMessage()
	m.SetHeader("Subject", "Subspace SMTP Test Email")
	m.SetBody("text/html", message)
	m.SetHeader("To", emailTo)
	if systemInfo.SmtpSenderName != "" && systemInfo.SmtpSenderEmail != "" {
		m.SetAddressHeader("From", systemInfo.SmtpSenderEmail, systemInfo.SmtpSenderName)
	} else if systemInfo.SmtpSenderName != "" && systemInfo.SmtpSenderEmail == "" {
		m.SetAddressHeader("From", ms.sender, systemInfo.SmtpSenderName)
	} else if systemInfo.SmtpSenderName == "" && systemInfo.SmtpSenderEmail != "" {
		m.SetAddressHeader("From", systemInfo.SmtpSenderEmail, "Subspace")
	} else {
		m.SetAddressHeader("From", ms.sender, "Subspace")
	}

	// send it
	if systemInfo.SmtpAuthentication {
		d := gomail.NewDialer(ms.smtpserver, ms.port, ms.sender, ms.password)
		if err := d.DialAndSend(m); err == nil {
			models.SMTPValid(true)
			return true
		}
	} else {
		d := gomail.Dialer{Host: ms.smtpserver, Port: ms.port}
		if err := d.DialAndSend(m); err == nil {
			models.SMTPValid(true)
			return true
		}
	}

	models.SMTPValid(false)
	return false
}

// SendSystemSettingsChangeEmail sends an email with system changes to all users
func SendSystemSettingsChangeEmail(dataType, value string) {
	var (
		systemInfo, _            = models.GetSystemInfo()
		ms                       = mailServer{smtpserver: systemInfo.SmtpHost, port: systemInfo.SmtpPort, sender: systemInfo.SmtpUsername, password: systemInfo.SmtpPassword}
		systemSettingsChangeInfo = systemSettingsChangeInfoTemplate{dataType, value}
		message                  = systemSettingsChangeTemplate(systemSettingsChangeInfo)
	)

	users := models.GetAllUsers()

	for _, v := range users {
		// send mail slowly over 3 minutes
		go func(recipient string) {
			time.Sleep(time.Duration(rand.Int31n(180)) * time.Second)
			// compose the message
			m := gomail.NewMessage()
			m.SetHeader("Subject", "Subspace System Configuration Change")
			m.SetBody("text/html", message)
			m.SetHeader("To", recipient)
			if systemInfo.SmtpSenderName != "" && systemInfo.SmtpSenderEmail != "" {
				m.SetAddressHeader("From", systemInfo.SmtpSenderEmail, systemInfo.SmtpSenderName)
			} else if systemInfo.SmtpSenderName != "" && systemInfo.SmtpSenderEmail == "" {
				m.SetAddressHeader("From", ms.sender, systemInfo.SmtpSenderName)
			} else if systemInfo.SmtpSenderName == "" && systemInfo.SmtpSenderEmail != "" {
				m.SetAddressHeader("From", systemInfo.SmtpSenderEmail, "Subspace")
			} else {
				m.SetAddressHeader("From", ms.sender, "Subspace")
			}

			// send it
			if systemInfo.SmtpAuthentication {
				d := gomail.NewDialer(ms.smtpserver, ms.port, ms.sender, ms.password)
				if err := d.DialAndSend(m); err == nil {
					beego.Debug("[SUCCESS] SendSystemSettingsChangeEmail to", recipient)
				} else {
					beego.Debug("[FAILED] SendSystemSettingsChangeEmail to", recipient)
				}
			} else {
				d := gomail.Dialer{Host: ms.smtpserver, Port: ms.port}
				if err := d.DialAndSend(m); err == nil {
					beego.Debug("[SUCCESS] SendSystemSettingsChangeEmail to", recipient)
				} else {
					beego.Debug("[FAILED] SendSystemSettingsChangeEmail to", recipient)
				}
			}
		}(v.Email)
	}

}

// SendUserEmail sends an email with a user's login information for the web-console
func SendUserEmail(username, password string) {
	var (
		systemInfo, _ = models.GetSystemInfo()
		ms            = mailServer{smtpserver: systemInfo.SmtpHost, port: systemInfo.SmtpPort, sender: systemInfo.SmtpUsername, password: systemInfo.SmtpPassword}
		userInfo      = userInfoTemplate{username, password}
		message       = userTemplate(userInfo)
	)

	// compose the message
	m := gomail.NewMessage()
	m.SetHeader("Subject", "Welcome to Subspace")
	m.SetBody("text/html", message)
	m.SetHeader("To", username)
	if systemInfo.SmtpSenderName != "" && systemInfo.SmtpSenderEmail != "" {
		m.SetAddressHeader("From", systemInfo.SmtpSenderEmail, systemInfo.SmtpSenderName)
	} else if systemInfo.SmtpSenderName != "" && systemInfo.SmtpSenderEmail == "" {
		m.SetAddressHeader("From", ms.sender, systemInfo.SmtpSenderName)
	} else if systemInfo.SmtpSenderName == "" && systemInfo.SmtpSenderEmail != "" {
		m.SetAddressHeader("From", systemInfo.SmtpSenderEmail, "Subspace")
	} else {
		m.SetAddressHeader("From", ms.sender, "Subspace")
	}

	// send it
	if systemInfo.SmtpAuthentication {
		d := gomail.NewDialer(ms.smtpserver, ms.port, ms.sender, ms.password)
		if err := d.DialAndSend(m); err == nil {
			beego.Debug("[SUCCESS] Sent user information to", username)
		} else {
			beego.Debug("[FAILED] Did not send user information to", username)
		}
	} else {
		d := gomail.Dialer{Host: ms.smtpserver, Port: ms.port}
		if err := d.DialAndSend(m); err == nil {
			beego.Debug("[SUCCESS] Sent user information to", username)
		} else {
			beego.Debug("[FAILED] Did not send user information to", username)
		}
	}
}

// SendProfileEmail sends an email with a user's vpn connection info
func SendProfileEmail(vpnProfileInfo VPNProfileData, emailTo string) {
	var (
		systemInfo, _ = models.GetSystemInfo()
		ms            = mailServer{smtpserver: systemInfo.SmtpHost, port: systemInfo.SmtpPort, sender: systemInfo.SmtpUsername, password: systemInfo.SmtpPassword}
		message       = vpnTemplate(vpnProfileInfo)
	)

	// compose the message
	m := gomail.NewMessage()
	m.SetHeader("Subject", "Subspace VPN Profile")
	m.SetBody("text/html", message)
	m.SetHeader("To", emailTo)
	if systemInfo.SmtpSenderName != "" && systemInfo.SmtpSenderEmail != "" {
		m.SetAddressHeader("From", systemInfo.SmtpSenderEmail, systemInfo.SmtpSenderName)
	} else if systemInfo.SmtpSenderName != "" && systemInfo.SmtpSenderEmail == "" {
		m.SetAddressHeader("From", ms.sender, systemInfo.SmtpSenderName)
	} else if systemInfo.SmtpSenderName == "" && systemInfo.SmtpSenderEmail != "" {
		m.SetAddressHeader("From", systemInfo.SmtpSenderEmail, "Subspace")
	} else {
		m.SetAddressHeader("From", ms.sender, "Subspace")
	}

	// Write VPN settings attachment to /tmp and delete it after send.
	if appleProfilePath, err := writeToTempFile(fmt.Sprintf("%s-apple.mobileconfig", vpnProfileInfo.FileNamePrefix), []byte(vpnProfileInfo.AppleMobileConfigContent)); nil == err {
		m.Attach(appleProfilePath)
		defer os.Remove(appleProfilePath)
	} else {
		beego.Warn("Write apple mobileconfig attachment fail.")
	}
	if windowsProfilePath, err := writeToTempFile(fmt.Sprintf("%s-windows.pbk", vpnProfileInfo.FileNamePrefix), []byte(vpnProfileInfo.WindowsPBKContent)); nil == err {
		m.Attach(windowsProfilePath)
		defer os.Remove(windowsProfilePath)
	} else {
		beego.Warn("Write windows pbk attachment fail.")
	}

	// send it
	if systemInfo.SmtpAuthentication {
		d := gomail.NewDialer(ms.smtpserver, ms.port, ms.sender, ms.password)
		if err := d.DialAndSend(m); err == nil {
			beego.Debug("[SUCCESS] Sent VPN profile information to", emailTo)
		} else {
			beego.Debug("[FAILED] Did not send VPN profile information to", emailTo)
		}
	} else {
		d := gomail.Dialer{Host: ms.smtpserver, Port: ms.port}
		if err := d.DialAndSend(m); err == nil {
			beego.Debug("[SUCCESS] Sent VPN profile information to", emailTo)
		} else {
			beego.Debug("[FAILED] Did not send VPN profile information to", emailTo)
		}
	}
}

func writeToTempFile(fileName string, content []byte) (absolutePath string, e error) {
	// Try to make dir
	_ = os.Mkdir(ATTACH_DIR, 777)
	absolutePath = strings.Join([]string{ATTACH_DIR, fileName}, "/")
	e = ioutil.WriteFile(absolutePath, content, 0644)
	return absolutePath, e
}

func userTemplate(userInfo userInfoTemplate) (template string) {
	const UserInfoTemplate = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" xmlns=\"http://www.w3.org/1999/xhtml\"><head> <meta http-equiv=\"Content-Type\" content=\"text/html; charset=3DUTF-8\" /> <meta name=\"viewport\" content=\"width=3Ddevice-width, initial-scale=1.0\" /> <title>Subspace Email</title></head><body style=\"-ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; font: 13px/1.4 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Tahoma, sans-serif; margin: 0; padding: 10px 0 0; width: 100% !important;\" bgcolor=\"#C7C7C7\"> <style type=\"text/css\"> body { width: 100% !important; -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; margin: 0; padding: 0; font: 13px/1.4 \"Helvetica\", \"Lucida Grande\", \"Lucida Sans Unicode\", Tahoma, sans-serif; } img { max-width: 100%; outline: none; text-decoration: none; -ms-interpolation-mode: bicubic; } @media (max-width: 500px) { .email-header { display: block; } .email-header td { display: block; } .email-header tr { display: block; } .email-header thead { display: block; } .email-header tbody { display: block; } } </style> <table cellpadding=\"0\" cellspacing=\"0\" border=\"0\" align=\"center\" style=\"border: 2px solid #b5b5b5; border-collapse: collapse; margin: 0 auto; max-width: 600px; min-width: 320px; mso-table-lspace: 0pt; mso-table-rspace: 0pt;\" bgcolor=\"white\"> <tr> <td style=\"border-collapse: collapse; color: white; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px; padding: 20px 30px 10px;\" bgcolor=\"black\"> <table style=\"width: 100%;\" class=\"email-header\"> <tr> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;\"> <h1 style=\"color: white; font-family: Sans-Serif; font-size: 30px; line-height: 1.4; margin: 0 0 10px;\"> Welcome to Subspace </h1> <p style=\"color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;\"> \"To be yourself in a world that is constantly trying to make you something else is the greatest accomplishment.\" <br />─ Ralph Waldo Emerson </p> </td> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px; padding-left: 20px; width: 50px;\" valign=\"top\"> <img src=\"https://subspace-assets.s3.amazonaws.com/subspace.png\" alt=\"Subspace\" style=\"-ms-interpolation-mode: bicubic; height: 50px; margin: 0; max-width: 100%; outline: none; text-decoration: none;\" /> </td> </tr> </table> </td> </tr> <tr> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px; padding: 10px 35px;\"> <h2 style=\"color: black; font-size: 15px; line-height: 1.4; margin: 1.5em 0 0.5em;\"> It's great to have you aboard. </h2> <p style=\"color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;\"> Here's some information to get you started. </p> </td> </tr> <tr> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px; padding: 10px 35px;\"> <table style=\"width: 100%;\" class=\"email-header\"> <tr> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;\"> <h2 style=\"color: black; font-size: 15px; line-height: 1.4; margin: 1.5em 0 0.5em;\"> Web Console </h2> </td> </tr> <tr> <td style=\"font-weight: bold; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;padding: 5px;text-align: left;\">Username</td> <td style=\"font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;padding: 5px;text-align: left;\">{{.Username}}</td> </tr> <tr> <td style=\"font-weight: bold; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;padding: 5px;text-align: left;\">Password</td> <td style=\"font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;padding: 5px;text-align: left;\">{{.Password}}</td> </tr> </table> </td> </tr> <tr> <td style=\"border-collapse: collapse; color: #777; font: 11px/1.4 'Helvetica', 'Lucida Grande','Lucida Sans Unicode', Tahoma, sans-serif; padding: 10px 20px;\" align=\"center\" bgcolor=\"#eee\"> Need help with anything? Hit up <a style=\"color: #777; text-decoration: underline;\" href=\"http://support.ecowork.com\">support</a>. </td> </tr> </table></body></html>"

	tmpl, err := tt.New("template.name").Parse(UserInfoTemplate)
	var doc bytes.Buffer
	err = tmpl.Execute(&doc, userInfo)
	if err != nil {
		panic(err)
	}

	return doc.String()
}

func vpnTemplate(vpnProfileInfo VPNProfileData) (template string) {
	const ProfileTemplate = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" xmlns=\"http://www.w3.org/1999/xhtml\"><head> <meta http-equiv=\"Content-Type\" content=\"text/html; charset=3DUTF-8\" /> <meta name=\"viewport\" content=\"width=3Ddevice-width, initial-scale=1.0\" /> <title>Subspace Email</title></head><body style=\"-ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; font: 13px/1.4 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Tahoma, sans-serif; margin: 0; padding: 10px 0 0; width: 100% !important;\" bgcolor=\"#C7C7C7\"> <style type=\"text/css\"> body { width: 100% !important; -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; margin: 0; padding: 0; font: 13px/1.4 \"Helvetica\", \"Lucida Grande\", \"Lucida Sans Unicode\", Tahoma, sans-serif; } img { max-width: 100%; outline: none; text-decoration: none; -ms-interpolation-mode: bicubic; } @media (max-width: 500px) { .email-header { display: block; } .email-header td { display: block; } .email-header tr { display: block; } .email-header thead { display: block; } .email-header tbody { display: block; } } </style> <table cellpadding=\"0\" cellspacing=\"0\" border=\"0\" align=\"center\" style=\"border: 2px solid #b5b5b5; border-collapse: collapse; margin: 0 auto; max-width: 600px; min-width: 320px; mso-table-lspace: 0pt; mso-table-rspace: 0pt;\" bgcolor=\"white\"> <tr> <td style=\"border-collapse: collapse; color: white; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px; padding: 20px 30px 10px;\" bgcolor=\"black\"> <table style=\"width: 100%;\" class=\"email-header\"> <tr> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;\"> <h1 style=\"color: white; font-family: Sans-Serif; font-size: 30px; line-height: 1.4; margin: 0 0 10px;\"> Subspace VPN Profile Information </h1> <p style=\"color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;\"> \"You may say I'm a dreamer, but I'm not the only one. I hope someday you'll join us. And the world will live as one.\" <br />─ John Lennon </p> </td> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px; padding-left: 20px; width: 50px;\" valign=\"top\"> <img src=\"https://subspace-assets.s3.amazonaws.com/subspace.png\" alt=\"Subspace\" style=\"-ms-interpolation-mode: bicubic; height: 50px; margin: 0; max-width: 100%; outline: none; text-decoration: none;\" /> </td> </tr> </table> </td> </tr> <tr> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px; padding: 10px 35px;\"> <h2 style=\"color: black; font-size: 15px; line-height: 1.4; margin: 1.5em 0 0.5em;\"> A Subspace VPN profile has been created for you. </h2> <p style=\"color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;\"> Subspace supports all major platforms. Here are a few quick links for you: </p> <ul style=\"padding: 0 0 0 10px;\"> <li> <a href=\"https://subspace.zendesk.com/hc/en-us/articles/115007748207\" style=\"color: #30A9F4; text-decoration: none;\">VPN Profile usage on macOS/iOS</a> </li> <li> <a href=\"https://support.ecowork.com/hc/en-us/articles/115007906988\" style=\"color: #30A9F4; text-decoration: none;\">VPN Profile usage on Windows</a> </li> <li> <a href=\"https://subspace.zendesk.com/hc/en-us/articles/115007748267\" style=\"color: #30A9F4; text-decoration: none;\">VPN Profile usage on Android</a> </li> </ul> </td> </tr> <tr> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px; padding: 10px 35px;\"> <table style=\"width: 100%;\" class=\"email-header\"> <tr> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;\"> <h2 style=\"color: black; font-size: 15px; line-height: 1.4; margin: 1.5em 0 0.5em;\"> VPN profile credentials </h2> </td> </tr> <tr> <td style=\"font-weight: bold; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;padding: 5px;text-align: left;\">Username</td> <td style=\"font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;padding: 5px;text-align: left;\">{{.Username}}</td> </tr> <tr> <td style=\"font-weight: bold; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;padding: 5px;text-align: left;\">Password</td> <td style=\"font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;padding: 5px;text-align: left;\">{{.Password}}</td> </tr> <tr> <td style=\"font-weight: bold; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;padding: 5px;text-align: left;\">Host Address</td> <td style=\"font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;padding: 5px;text-align: left;\">{{.Host}}</td> </tr> <tr> <td style=\"font-weight: bold; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;padding: 5px;text-align: left;\">Pre-Shared Key</td> <td style=\"font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;padding: 5px;text-align: left;\">{{.PreSharedKey}}</td> </tr> </table> </td> </tr> <tr> <td style=\"border-collapse: collapse; color: #777; font: 11px/1.4 'Helvetica', 'Lucida Grande','Lucida Sans Unicode', Tahoma, sans-serif; padding: 10px 20px;\" align=\"center\" bgcolor=\"#eee\"> Need help with anything? Hit up <a style=\"color: #777; text-decoration: underline;\" href=\"http://support.ecowork.com\">support</a>. </td> </tr> </table></body></html>"

	tmpl, err := tt.New("template.name").Parse(ProfileTemplate)
	var doc bytes.Buffer
	err = tmpl.Execute(&doc, vpnProfileInfo)
	if err != nil {
		panic(err)
	}

	return doc.String()
}

func systemSettingsChangeTemplate(systemSettingsChangeInfo systemSettingsChangeInfoTemplate) (template string) {
	const SystemSettingsChangeTemplate = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" xmlns=\"http://www.w3.org/1999/xhtml\"><head> <meta http-equiv=\"Content-Type\" content=\"text/html; charset=3DUTF-8\" /> <meta name=\"viewport\" content=\"width=3Ddevice-width, initial-scale=1.0\" /> <title>Subspace Email</title></head><body style=\"-ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; font: 13px/1.4 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Tahoma, sans-serif; margin: 0; padding: 10px 0 0; width: 100% !important;\" bgcolor=\"#C7C7C7\"> <style type=\"text/css\"> body { width: 100% !important; -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; margin: 0; padding: 0; font: 13px/1.4 \"Helvetica\", \"Lucida Grande\", \"Lucida Sans Unicode\", Tahoma, sans-serif; } img { max-width: 100%; outline: none; text-decoration: none; -ms-interpolation-mode: bicubic; } @media (max-width: 500px) { .email-header { display: block; } .email-header td { display: block; } .email-header tr { display: block; } .email-header thead { display: block; } .email-header tbody { display: block; } } </style> <table cellpadding=\"0\" cellspacing=\"0\" border=\"0\" align=\"center\" style=\"border: 2px solid #b5b5b5; border-collapse: collapse; margin: 0 auto; max-width: 600px; min-width: 320px; mso-table-lspace: 0pt; mso-table-rspace: 0pt;\" bgcolor=\"white\"> <tr> <td style=\"border-collapse: collapse; color: white; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px; padding: 20px 30px 10px;\" bgcolor=\"black\"> <table style=\"width: 100%;\" class=\"email-header\"> <tr> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;\"> <h1 style=\"color: white; font-family: Sans-Serif; font-size: 30px; line-height: 1.4; margin: 0 0 10px;\"> System Settings Updated </h1> <p style=\"color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;\"> \"The measure of intelligence is the ability to change.\" <br />─ Albert Einstein </p> </td> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px; padding-left: 20px; width: 50px;\" valign=\"top\"> <img src=\"https://subspace-assets.s3.amazonaws.com/subspace.png\" alt=\"Subspace\" style=\"-ms-interpolation-mode: bicubic; height: 50px; margin: 0; max-width: 100%; outline: none; text-decoration: none;\" /> </td> </tr> </table> </td> </tr> <tr> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px; padding: 10px 35px;\"> <h2 style=\"color: black; font-size: 15px; line-height: 1.4; margin: 1.5em 0 0.5em;\"> An administrator has changed the system configuration of Subspace. </h2> <span style=\"color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;\"> You may need to update the configuration of your Subspace VPN profile on your device. </span> </td> </tr> <tr> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px; padding: 10px 35px;\"> <table style=\"width: 100%;\" class=\"email-header\"> <tr> <td style=\"border-collapse: collapse; color: #878787; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;\"> <h2 style=\"color: black; font-size: 15px; line-height: 1.4; margin: 1.5em 0 0.5em;\"> Subspace System Changes </h2> </td> </tr> <tr> <td style=\"font-weight: bold; font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;padding: 5px;text-align: left;\">{{.DataType}}</td> <td style=\"font-family: 'Helvetica', 'Lucida Grande', 'Lucida Sans Unicode', Verdana, sans-serif; font-size: 13px;padding: 5px;text-align: left;\">{{.Value}}</td> </tr> </table> </td> </tr> <tr> <td style=\"border-collapse: collapse; color: #777; font: 11px/1.4 'Helvetica', 'Lucida Grande','Lucida Sans Unicode', Tahoma, sans-serif; padding: 10px 20px;\" align=\"center\" bgcolor=\"#eee\"> Need help with anything? Hit up <a style=\"color: #777; text-decoration: underline;\" href=\"http://support.ecowork.com\">support</a>. </td> </tr> </table></body></html>"

	tmpl, err := tt.New("template.name").Parse(SystemSettingsChangeTemplate)
	var doc bytes.Buffer
	err = tmpl.Execute(&doc, systemSettingsChangeInfo)
	if err != nil {
		panic(err)
	}

	return doc.String()
}
