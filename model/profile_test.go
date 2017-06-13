package model

import (
	"testing"
)

func TestProfile_Validate(t *testing.T) {
	config := &Profile{}
	valid := func(n int) {
		if err := config.Validate(); err == nil {
			t.Fatal("Test:", n, "this config need to be invalid")
		}
	}

	valid(0)

	config.Name = "x"
	valid(1)

	x := -1
	config.LinesToPreview = &x
	valid(2)

	*config.LinesToPreview = 10
	valid(3)

	config.Email = &Email{}
	valid(4)

	config.Email = &Email{
		ImapServer: "imap.example.com:143",
		Address:    "orders@example.com",
		Password:   "password",
	}
	valid(5)

	config.Mattermost = &Mattermost{}
	valid(6)

	config.Mattermost = &Mattermost{
		Server:   "https://mattermost.example.com",
		Team:     "team1",
		Channels: []string{"#orders"},
		User:     "mattermail@example.com",
		Password: "password",
	}

	if err := config.Validate(); err != nil {
		t.Fatal(err)
	}

	config.Filter = &Filter{}
	valid(7)

	config.Filter = nil

	if err := config.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestProfile_Fix(t *testing.T) {
	p := NewProfile()

	p.Fix()

	if *p.MailTemplate != defaultMailTemplate {
		t.Fatal("Expected MailTemplate:", defaultMailTemplate, " result:", *p.MailTemplate)
	}
	if *p.LinesToPreview != defaultLinesToPreview {
		t.Fatal("Expected LinesToPreview:", defaultLinesToPreview, " result:", *p.LinesToPreview)
	}
	if *p.RedirectChannel != defaultRedirectChannel {
		t.Fatal("Expected RedirectChannel:", defaultRedirectChannel, " result:", *p.RedirectChannel)
	}
	if *p.Attachment != defaultAttachment {
		t.Fatal("Expected Attachment:", defaultAttachment, " result:", *p.Attachment)
	}
	if *p.Disabled != defaultDisabled {
		t.Fatal("Expected Disabled:", defaultDisabled, " result:", *p.Disabled)
	}
}

func TestProfile_FormatMailTemplate(t *testing.T) {
	type args struct {
		from    string
		subject string
		message string
	}

	assert := func(n int, template, expected string, wantErr bool, args args) {
		profile := &Profile{
			MailTemplate: new(string),
		}
		*profile.MailTemplate = template

		result, err := profile.FormatMailTemplate(args.from, args.subject, args.message)
		gotErr := (err != nil)
		if gotErr != wantErr {
			if gotErr {
				t.Fatal("Test:", n, "Expected Err:", wantErr, " Result Err:", gotErr, " err:", err.Error())
			} else {
				t.Fatal("Test:", n, "Expected Err:", wantErr, " Result Err:", gotErr)
			}
		}
		if result != expected {
			t.Fatal("Test:", n, "Expected:", expected, " Result:", result)
		}
	}

	assert(0, "", "", false, args{from: "", subject: "", message: ""})
	assert(1, ">{{.From}}", ">test@test.com", false, args{from: "test@test.com", subject: "", message: ""})
	assert(2, "{{.From}}{{.Subject}}{{.Message}}", "", false, args{from: "", subject: "", message: ""})
	assert(3, ">{{.From}}, {{.Subject}}, {{.Message}}", ">test@test.com, subject, message", false, args{from: "test@test.com", subject: "subject", message: "message"})
	assert(4, ">{{.Nothing}}", "", true, args{from: "test@test.com", subject: "", message: ""})
	assert(5, ">{{.Noth%ing}}", "", true, args{from: "test@test.com", subject: "", message: ""})
}
