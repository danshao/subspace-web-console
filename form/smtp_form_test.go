package form

import (
	"testing"
)

func TestSmtpSettingsForm(t *testing.T) {
	f1 := SmtpSettingsForm{
		SMTPHost: "smtp.gmail.com",
		SMTPPort: "587",
		Authentication: true,
		Username: "ecowork.subspace@gmail.com",
		Password: "1234",
		SenderName: "Subspace Admin",
		SenderEmail: "ecowork.subspace@gmail.com",
	}
	if err := f1.Validate(); nil != err {
		t.Error(err)
	}

	f2 := SmtpSettingsForm{
		SMTPHost: "smtp.gmail.com",
		SMTPPort: "-1",
		Authentication: true,
		Username: "ecowork.subspace@gmail.com",
		Password: "1234",
		SenderName: "Subspace Admin",
		SenderEmail: "ecowork.subspace@gmail.com",
	}
	if err := f2.Validate(); nil == err {
		t.Error(err)
	}

	f3 := SmtpSettingsForm{
		SMTPHost: "smtp.gmail.com",
		SMTPPort: "65536",
		Authentication: true,
		Username: "ecowork.subspace@gmail.com",
		Password: "1234",
		SenderName: "Subspace Admin",
		SenderEmail: "ecowork.subspace@gmail.com",
	}
	if err := f3.Validate(); nil == err {
		t.Error(err)
	}

	f4 := SmtpSettingsForm{
		SMTPHost: "",
		SMTPPort: "587",
		Authentication: true,
		Username: "ecowork.subspace@gmail.com",
		Password: "1234",
		SenderName: "Subspace Admin",
		SenderEmail: "ecowork.subspace@gmail.com",
	}
	if err := f4.Validate(); nil != err {
		t.Error(err)
	}

	f5 := SmtpSettingsForm{
		SMTPHost: "smtp.gmail.com",
		SMTPPort: "587",
		Authentication: true,
		Username: "ecowork.subspace@gmail.com",
		Password: "1234",
		SenderName: "Subspace Admin",
		SenderEmail: "ecowork.subspace",
	}
	if err := f5.Validate(); nil == err {
		t.Error("ecowork.subspace not a valid email")
	}

	f6 := SmtpSettingsForm{
		SMTPHost: "!!!",
		SMTPPort: "587",
		Authentication: true,
		Username: "ecowork.subspace@gmail.com",
		Password: "1234",
		SenderName: "Subspace Admin",
		SenderEmail: "ecowork.subspace@gmail.com",
	}
	if err := f6.Validate(); nil == err {
		t.Error("Not a valid smtp host.")
	}

	f7 := SmtpSettingsForm{
		SMTPHost: "",
		SMTPPort: "",
		Authentication: true,
		Username: "",
		Password: "",
		SenderName: "",
		SenderEmail: "",
	}
	if err := f7.Validate(); nil != err {
		t.Error("All empty Should be a valid smtp host.")
	}

	f8 := SmtpSettingsForm{
		SMTPHost: "",
		SMTPPort: "0",
		Authentication: true,
		Username: "",
		Password: "",
		SenderName: "",
		SenderEmail: "",
	}
	if err := f8.Validate(); nil != err {
		t.Error("Port 0 Should be a valid smtp host.")
	}
}