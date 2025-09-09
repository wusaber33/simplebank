package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string, //主体
		content string, //内容
		to []string, //电子邮件发送道德电子邮件地址列表
		cc []string, //抄送
		bcc []string, //密件抄送收件人
		attachFiles []string, //附件
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name, fromEmailAddress, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(subject, content string, to, cc, bcc []string, attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>",sender.name,sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To= to
	e.Cc = cc
	e.Bcc = bcc


	for _,f:=range attachFiles {
		_,err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w",f,err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress,sender.fromEmailPassword,smtpAuthAddress)
	return e.Send(smtpServerAddress,smtpAuth)
}