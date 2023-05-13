package mail

import (
	"testing"

	"github.com/julysNICK/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	config, err := utils.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Test subject"

	content := `
		<html>
			<body>
				<h1 style="color:red;">Test content</h1>
			</body>
		</html>
	`

	to := []string{"julysmartins54@gmail.com"}

	attachFile := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFile)

	require.NoError(t, err)

}
