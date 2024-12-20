package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go" //nolint:depguard
)

var (
	json             = jsoniter.ConfigCompatibleWithStandardLibrary
	ErrInvalidDomain = fmt.Errorf("invalid domain")
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	domain = strings.Trim(domain, " ")
	if domain == "" {
		return result, ErrInvalidDomain
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		userEmail, err := unmarshalUserEmail(scanner.Bytes())
		if err != nil {
			return result, err
		}

		foundDomain := ExtractMatchedDomain(userEmail, domain)
		if foundDomain != "" {
			result[foundDomain]++
		}
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func unmarshalUserEmail(data []byte) (string, error) {
	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return "", fmt.Errorf("unmarshal user error: %w", err)
	}
	return user.Email, nil
}

func ExtractMatchedDomain(email, domain string) string {
	email = strings.ToLower(email)

	if !strings.HasSuffix(email, domain) {
		return ""
	}

	emailParts := strings.SplitN(email, "@", 2)
	if len(emailParts) != 2 {
		return ""
	}

	return emailParts[1]
}
