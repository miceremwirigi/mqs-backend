package utils

import (
	"errors"
	"os"
	"strconv"
)

func LoadSMTPConfig() (smtp_host string, smtp_port int, smtp_user, smtp_pass string, err error) {
	smtp_host, set := os.LookupEnv("smtp_host")
	if !set {
		return "", 0, "", "", errors.New("smtp_host environment variable not set")
	}
	smtp_port_string, set := os.LookupEnv("smtp_port")
	if !set {
		return "", 0, "", "", errors.New("smtp_port environment variable not set")
	}
	smtp_port, err = strconv.Atoi(smtp_port_string)
	if err != nil {
		return "", 0, "", "", errors.New("invalid smtp_port value: " + smtp_port_string)
	}
	smtp_user, set = os.LookupEnv("smtp_user")
	if !set {
		return "", 0, "", "", errors.New("smtp_user environment variable not set")
	}
	smtp_pass, set = os.LookupEnv("smtp_pass")
	if !set {
		return "", 0, "", "", errors.New("smtp_pass environment variable not set")
	}
	return smtp_host, smtp_port, smtp_user, smtp_pass, nil
}
