package messaging_test

import (
	"journie/pkg/messaging"
	"reflect"
	"regexp"
	"testing"
)

// TestGetPlatformUserIdValid calls messaging.GetPlatformUserId with an id, checking
// for a valid return value.
func TestGetPlatformUserIdValid(t *testing.T) {
	userId := "123456"
	want := regexp.MustCompile(`\b` + `telegram-123456` + `\b`)
	str, err := messaging.GetPlatformUserId(userId)
	if !want.MatchString(str) || err != nil {
		t.Fatalf(`GetPlatformUserId("123456") = %q, %v, want match for %#q, nil`, str, err, want)
	}
}

// TestPlatformUserIdEmpty calls messaging.GetPlatformUserId with an empty string,
// checking for an error.
func TestPlatformUserIdEmpty(t *testing.T) {
	str, err := messaging.GetPlatformUserId("")
	if str != "" || err == nil {
		t.Fatalf(`Hello("") = %q, %v, want "", error`, str, err)
	}
}

// TestParsePlatformUserIdValid calls messaging.ParsePlatformUserId with an id, checking
// for a valid return value.
func TestParsePlatformUserIdValid(t *testing.T) {
	platformUserId := "telegram-123456"
	want := &messaging.UserModel{
		Platform: "telegram",
		UserId:   "123456",
	}

	obj, err := messaging.ParsePlatformUserId(platformUserId)
	if !reflect.DeepEqual(want, obj) || err != nil {
		t.Fatalf(`ParsePlatformUserId("telegram-123456") = %q, %v, want match for %#q, nil`, obj, err, want)
	}
}

// TestParsePlatformUserIdInvalid calls messaging.ParsePlatformUserId with wrong format
// checking for an error.
func TestParsePlatformUserIdInvalid(t *testing.T) {
	platformUserId := "123456"
	obj, err := messaging.ParsePlatformUserId(platformUserId)

	if obj != nil || err == nil {
		t.Fatalf(`ParsePlatformUserId("123456") = %q, %v, want "", error`, obj, err)
	}
}
