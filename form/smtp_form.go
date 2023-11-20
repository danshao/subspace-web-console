package form

import (
	"strings"
	"github.com/asaskevich/govalidator"
)

type SmtpSettingsForm struct {
	SMTPHost       string `form:"smtpserver" valid:"host"`
	SMTPPort       string `form:"smtpport" valid:"port"`
	Authentication bool   `form:"authentication"`
	Username       string `form:"username"`
	Password       string `form:"password"`
	SenderName     string `form:"sendername"`
	SenderEmail    string `form:"senderemail" valid:"email"`
}

func (form *SmtpSettingsForm) beforeValidate() {
	form.SMTPHost = strings.TrimSpace(form.SMTPHost)
	form.SMTPPort = strings.TrimSpace(form.SMTPPort)
	if "0" == form.SMTPPort {
		form.SMTPPort = ""
	}
	form.Username = strings.TrimSpace(form.Username)
	form.SenderEmail = strings.TrimSpace(form.SenderEmail)
}

func (form *SmtpSettingsForm) IsValid() (bool, error) {
	err := form.Validate()
	return nil == err, err
}

func (form *SmtpSettingsForm) Validate() (err error) {
	form.beforeValidate()
	_, err = govalidator.ValidateStruct(form)

	return err
}

func (form *SmtpSettingsForm) IsEmpty() bool {
	return "" == form.SMTPHost &&
					( "0" == form.SMTPPort || "" == form.SMTPHost ) &&
					(!form.Authentication || ( form.Authentication && "" == form.Username && "" == form.Password)) &&
					"" == form.SenderName &&
					"" == form.SenderEmail
}
