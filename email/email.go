package email

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var From string
var Key string

func SenEmail(subject, content string, tos []string) (string, error) {
	from := mail.NewEmail("coderkl", From)
	to := mail.NewEmail("", tos[0])

	plainTextContent := ""
	htmlContent := "<pre>" + content + "</pre>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	for index, value := range tos {
		if index > 0 {
			to := mail.NewEmail("", value)
			message.Personalizations[0].AddTos(to)
		}
	}

	client := sendgrid.NewSendClient(Key)
	response, err := client.Send(message)

	responseMsg := fmt.Sprintf("responseCode:%d responseMsg:%s", response.StatusCode, response.Body)
	return responseMsg, err
}

func GenerateHtml(content string) string {
	head := `<!DOCTYPE html>
<html>
<head><meta charset='utf-8'></head>
<body>
<pre>`

	foot := `</pre>
</body>
</html>
`
	return head + content + foot
}
