package databases

import (
	"errors"
	"os"
)

func LoadConfig() (env, db_host, db_user, db_pass, db_name, db_ssl, db_port string, err error) {
	env, set := os.LookupEnv("env")
	if !set {
		return "", "", "", "", "", "", "", errors.New("env environment variable not set")
	}
	db_host, set = os.LookupEnv("db_host")
	if !set {
		return "", "", "", "", "", "", "", errors.New("db_host environment variable not set")
	}
	db_user, set = os.LookupEnv("db_user")
	if !set {
		return "", "", "", "", "", "", "", errors.New("db_user environment variable not set")
	}
	db_pass, set = os.LookupEnv("db_pass")
	if !set {
		return "", "", "", "", "", "", "", errors.New("db_pass environment variable not set")
	}
	db_name, set = os.LookupEnv("db_name")
	if !set {
		return "", "", "", "", "", "", "", errors.New("db_name environment variable not set")
	}
	db_ssl, set = os.LookupEnv("db_ssl")
	if !set {
		return "", "", "", "", "", "", "", errors.New("db_ssl environment variable not set")
	}
	db_port, set = os.LookupEnv("db_port")
	if !set {
		return "", "", "", "", "", "", "", errors.New("db_port environment variable not set")
	}
	return env, db_host, db_user, db_pass, db_name, db_ssl, db_port, nil
}
