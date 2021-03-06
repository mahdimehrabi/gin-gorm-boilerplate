package infrastructure

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

//Email -> Email Struct
type Email struct {
	logger Logger
	env    Env
}

//NewEmail -> return new Email struct
func NewEmail(
	logger Logger,
	env Env,
) Email {
	return Email{
		logger: logger,
		env:    env,
	}
}

//main method for sending email ,
//we mocked this method in suite_test.go to make sure email sending work , feel free to use that
func (e Email) SendEmail(ch chan error, to string, subject string, htmlFilePath string, templateData interface{}) {
	if e.env.Environment == "test" {
		ch <- nil
		return
	}
	m := gomail.NewMessage()
	m.SetHeader("From", e.env.SiteEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	htmlContent, err := ReadHTMLFile(htmlFilePath)
	if err != nil {
		ch <- err
	}
	tmpl := template.New("email")
	tmpl.Parse(string(htmlContent))
	var processed bytes.Buffer
	tmpl.ExecuteTemplate(&processed, "email", templateData)
	m.SetBody("text/html", processed.String())
	port, err := strconv.Atoi(e.env.MailPort)
	if err != nil {
		ch <- err
	}
	d := gomail.NewDialer(e.env.MailHost, port, e.env.MailHostUser, e.env.MailHostPassword)

	if err := d.DialAndSend(m); err != nil {
		ch <- err
	}
	ch <- nil
}

func ReadHTMLFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}
