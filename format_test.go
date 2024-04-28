package env

import "testing"

func Test__formatName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"DBName", "DB_NAME"},
		{"DB", "DB"},
		{"Name", "NAME"},
		{"AllowedOrigins", "ALLOWED_ORIGINS"},
		{"TestA", "TEST_A"},
		{"UserID", "USER_ID"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := formatName(test.name)
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}
