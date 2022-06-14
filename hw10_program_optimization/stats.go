package hw10programoptimization

import (
	"bufio"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"io"
	"strings"
)

type User struct {
	ID       int    `json:"-"`
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string `json:"Email"`
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
}

type DomainStat map[string]int

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	reader := bufio.NewReader(r)
	result := make(DomainStat)

	isEOF := false
	for {
		lineContent, _, err := reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				isEOF = true
			} else {
				return result, err
			}
		}

		var user User
		if err := json.Unmarshal(lineContent, &user); err != nil && !isEOF {
			return nil, err
		}

		matched := strings.Contains(user.Email, domain)
		if matched {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}

		if isEOF {
			break
		}
	}

	return result, nil
}
