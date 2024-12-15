package helper

import (
	"bytes"
	_ "embed"
	"net/url"
	"strings"
	"text/template"

	"github.com/1Storm3/flibox-api/internal/shared/logger"
	"go.uber.org/zap"
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
		logger.Info("Ошибка при создании шаблона: %v", zap.Error(err))
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
		logger.Info("Ошибка при выполнении шаблона: %v", zap.Error(err))
		return "", err
	}
	return emailBody.String(), nil
}
