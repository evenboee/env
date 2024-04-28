package env

import (
	"os"
	"testing"
)

func Test__Bind(t *testing.T) {
	type Config struct {
		Host   string `env:"TEST,required"`
		Port   int    `env:",required"`
		UserID string `env:",required"`
	}

	os.Setenv("TEST", "localhost")
	os.Setenv("PORT", "8080")
	os.Setenv("USER_ID", "123")

	conf, err := Bind[Config]()
	if err != nil {
		t.Fatal("bind failed:", err)
	}

	if conf.Host != "localhost" {
		t.Errorf("expected %q, got %q", "localhost", conf.Host)
	}

	if conf.Port != 8080 {
		t.Errorf("expected %d, got %d", 8080, conf.Port)
	}

	if conf.UserID != "123" {
		t.Errorf("expected %q, got %q", "123", conf.UserID)
	}
}
