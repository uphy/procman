package config

import (
	"os"
)

type (
	Env struct {
		prefix string
	}
)

func NewEnv() *Env {
	return &Env{"PROCMAN_"}
}

func (e *Env) get(key string, defaultValue string) string {
	v := os.Getenv(e.prefix + key)
	if v == "" {
		return defaultValue
	}
	return v
}

func (e *Env) User() string {
	return e.get("USER", "admin")
}

func (e *Env) Password() string {
	return e.get("PASSWORD", "password")
}

func (e *Env) JWTSecret() string {
	return e.get("JWTSECRET", "s3cr3t")
}
