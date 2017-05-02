package utils

import (
	"regexp"
	"testing"
)

func Test_GenerateUUID4(t *testing.T) {
	uuid4 := GenerateUUID4()
	validUUID4 := regexp.MustCompile(`[a-f0-9]{8}-?[a-f0-9]{4}-?4[a-f0-9]{3}-?[89ab][a-f0-9]{3}-?[a-f0-9]{12}`)
	if !validUUID4.MatchString(uuid4) {
		t.Errorf("Generated UUIDv4 is not valid: %s", uuid4)
	}
}

func Test_GenerateUUID4_notEmpty(t *testing.T) {
	uuid4 := GenerateUUID4()
	if uuid4 == "" {
		t.Error("Generated UUIDv4 is empty")
	}
}

func Test_GeneratedUUID4_unique(t *testing.T) {
	uuid4a := GenerateUUID4()
	uuid4b := GenerateUUID4()

	if uuid4a == uuid4b {
		t.Errorf("Generated UUIDv4's are identical: %s == %s", uuid4a, uuid4b)
	}
}
