package analytics

import (
	"testing"

	"reflect"

	"net/url"
)

// Test Payload Hit Validation: https://goo.gl/4E4ISD
func TestScreenViewPostMeasurement(t *testing.T) {
	validScreenViewHitURL, _ := url.Parse("https://www.google-analytics.com/collect?v=1&t=screenview&tid=UA-12345678-9&cid=testid&an=lumogon&av=testversion&cd=testscreen&cd1=1.27")
	u := *MockUserSession()
	u.Type = "screenview"
	u.ScreenName = "testscreen"

	res, err := u.PostMeasurement()
	if err != nil {
		t.Errorf("HTTP Client did not return 200. %s", err)
	}

	queryValues := reflect.ValueOf(res.Request.URL.Query()).MapKeys()
	requiredValues := reflect.ValueOf(validScreenViewHitURL.Query()).MapKeys()

	if compareKeys(queryValues, requiredValues) != true {
		t.Errorf("Event did not have required query parameters. Expected %s ::: Received %s", requiredValues, queryValues)
	}

	TeardownMockUserSession(t)
}

// Test Payload Hit Validation: https://goo.gl/7mC9TI
func TestEventPostMeasurement(t *testing.T) {
	validEventHitURL, _ := url.Parse("https://www.google-analytics.com/collect?v=1&t=event&tid=UA-12345678-9&cid=testid&an=lumogon&av=testversion&ec=testcategory&ea=testaction&cd1=1.27")
	u := *MockUserSession()
	u.Type = "event"
	u.Action = "testaction"
	u.Category = "testcategory"

	res, err := u.PostMeasurement()
	if err != nil {
		t.Errorf("HTTP Client did not return 200. %s", err)
	}

	queryValues := reflect.ValueOf(res.Request.URL.Query()).MapKeys()
	requiredValues := reflect.ValueOf(validEventHitURL.Query()).MapKeys()

	if compareKeys(queryValues, requiredValues) != true {
		t.Errorf("Event did not have required query parameters. Expected %s ::: Received %s", requiredValues, queryValues)
	}

	TeardownMockUserSession(t)
}

func TestDisableTransmitPostMeasurement(t *testing.T) {
	u := *MockUserSession()
	u.DisableTransmit = true

	res, _ := u.PostMeasurement()
	if res != nil {
		t.Errorf("User request to not post analytics failed")
	}
}
