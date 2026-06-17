package credentials

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Credentials struct {
	MailLogin string
	PassAPI   string
}

func ResolvePath(flagPath string) string {
	if flagPath != "" {
		return flagPath
	}
	if env := os.Getenv("CC_CREDENTIALS_FILE"); env != "" {
		return env
	}
	return "password_api.txt"
}

func Load(path string) (Credentials, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Credentials{}, fmt.Errorf("read credentials file %s: %w", path, err)
	}

	creds := Credentials{}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(strings.ToLower(key))
		value = strings.TrimSpace(value)
		switch key {
		case "user", "maillogin", "login", "email":
			creds.MailLogin = value
		case "password", "passapi", "pass":
			creds.PassAPI = value
		}
	}
	if err := scanner.Err(); err != nil {
		return Credentials{}, fmt.Errorf("parse credentials file %s: %w", path, err)
	}
	if creds.MailLogin == "" {
		return Credentials{}, fmt.Errorf("credentials file %s: missing user/login", path)
	}
	if creds.PassAPI == "" {
		return Credentials{}, fmt.Errorf("credentials file %s: missing password", path)
	}
	return creds, nil
}
