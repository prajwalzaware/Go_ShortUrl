package utils

import "net/url"

func IsValidURL(inputURL string) bool {
	_, err := url.ParseRequestURI(inputURL)
	return err == nil
}
