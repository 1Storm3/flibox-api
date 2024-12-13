package helper

import (
	"bytes"
	_ "embed"
	"log"
	"net/url"
	"strings"
	"text/template"
)

//go:embed template/email.html
var emailTemplate string

func ExtractS3Key(photoURL string) (string, error) {
	parsedURL, err := url.Parse(photoURL)
	if err != nil {
		return "", err
	}
	segments := strings.Split(parsedURL.Path, "/")
	return segments[len(segments)-1], nil
}

func TakeHTMLTemplate(appUrl, verificationToken string) (string, error) {
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		log.Printf("Ошибка при создании шаблона: %v", err)
		return "", err
	}

	data := struct {
		AppUrl string
		Token  string
	}{
		AppUrl: appUrl,
		Token:  verificationToken,
	}

	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, data); err != nil {
		log.Printf("Ошибка при выполнении шаблона: %v", err)
		return "", err
	}
	return emailBody.String(), nil
}
