package short

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const maxKeyLen = 16

type link struct {
	uri      string
	shortKey string
}

func checkLink(ln link) error {
	if err := checkShortKey(ln.shortKey); err != nil {
		return err
	}
	u, err := url.ParseRequestURI(ln.uri)

	if err != nil {
		return err
	}
	if u.Host == "" {
		return errors.New("empty host")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("scheme must be http or https")
	}
	return nil
}

func checkShortKey(k string) error {
	if strings.TrimSpace(k) == "" {
		return errors.New("empty key")
	}
	if len(k) > maxKeyLen {
		return fmt.Errorf("key too long (max %d)", maxKeyLen)
	}
	return nil
}
