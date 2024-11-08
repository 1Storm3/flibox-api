package helper

import (
	"net/url"
	"strings"
)

func ExtractS3Key(photoURL string) (string, error) {
	parsedURL, err := url.Parse(photoURL)
	if err != nil {
		return "", err
	}
	segments := strings.Split(parsedURL.Path, "/")
	return segments[len(segments)-1], nil
}
