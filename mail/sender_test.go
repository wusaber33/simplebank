package mail

import (
	"simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := util.LoadConfig("..")
	require.NoError(t,err)
	
	sender := NewGmailSender(config.EmailSenderName,config.EmailSenderAddress,config.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="http://techschool.guru">TechSchool</a></p>
	`
	to := []string{"wusaber33@outlook.com"}
	attachfiles := []string{"../dbml-error.log"}
	err = sender.SendEmail(subject, content, to, nil,nil,attachfiles)
	require.NoError(t, err)
}