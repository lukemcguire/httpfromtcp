package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return Headers{}
}

// Parse
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	key, value, err := parseFieldLineFromString(string(data[:idx]))
	if err != nil {
		return 0, false, err
	}

	// if header key already exists extend its value by the new one
	if v, ok := h[key]; ok {
		value = strings.Join([]string{v, value}, ", ")
	}
	h.Set(key, value)

	return idx + 2, false, nil
}

// parseFieldLineFromString
func parseFieldLineFromString(str string) (key, value string, err error) {
	pattern := regexp.MustCompile(`^\s*(?P<key>[A-Za-z0-9!#$%&'*+.^_\x60|~-]+):\s*(?P<value>.+?)\s*$`)
	matches := pattern.FindStringSubmatch(str)
	if matches == nil {
		return "", "", fmt.Errorf("unable to parse field-line: %s", str)
	}
	keyIndex := pattern.SubexpIndex("key")
	key = strings.ToLower(matches[keyIndex])
	valueIndex := pattern.SubexpIndex("value")
	value = matches[valueIndex]

	return key, value, nil
}

func (h Headers) Set(key, value string) {
	h[key] = value
}

func (h Headers) Contains(key string) bool {
	_, ok := h[key]
	return ok
}
