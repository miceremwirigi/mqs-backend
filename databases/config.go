package databases

import (
	"errors"
	"os"
)

func LoadConfig() (db_name string, db_pass string, err error){
	env, set := os.LookupEnv("env")
	if !set {
		err := errors.New("environment variable not set")
		return "", "", err
	}
	if env != "test" {
		db_name, set = os.LookupEnv("db_name")
		if !set {
			err := errors.New("environment variable not set")
			return "", "", err
		}
		db_pass, set = os.LookupEnv("db_pass")
		if !set {
			err := errors.New("environment variable not set")
			return "", "", err
		}
	}
	return db_name, db_pass, nil
}
