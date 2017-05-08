package analytics

import (
	"sort"
	"testing"

	"reflect"

	"net/url"

	gock "gopkg.in/h2non/gock.v1"
)

// GockBootstrap returns a UserSession ready for use in testing the GA service.
// Note: This endpoint IRL ALWAYS returns 200 OK, and sends back a small 1px gif image.
// This means we must rely on whether our hit is put together properly.
// Validate at https://ga-dev-tools.appspot.com/hit-builder/
func MockUserSession() *UserSession {
	u := *NewUserSession()
	gock.New("https://www.google-analytics.com").
		Post("/collect").
		Reply(200).
		SetHeader("Content-Type", "image/gif").
		BodyString("GIF89a�����,D;")
	gock.InterceptClient(u.HTTPClient)
	return &u
}

func TeardownMockUserSession(t *testing.T) {
	defer gock.Off()

	if gock.IsDone() != true {
		t.Errorf("Pending mock requests are still in queue")
	}
}

// Test Payload Hit Validation: https://goo.gl/4E4ISD
func TestScreenViewPostMeasurement(t *testing.T) {
	validScreenViewHitURL, _ := url.Parse("https://www.google-analytics.com/collect?v=1&t=screenview&tid=UA-12345678-9&cid=testid&an=lumogon&av=testversion&cd=testscreen")
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
	validEventHitURL, _ := url.Parse("https://www.google-analytics.com/collect?v=1&t=event&tid=UA-12345678-9&cid=testid&an=lumogon&av=testversion&ec=testcategory&ea=testaction")
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

func compareKeys(a []reflect.Value, b []reflect.Value) bool {
	x := []string{}
	y := []string{}

	for _, e := range a {
		x = append(x, e.String())
	}

	for _, e := range b {
		y = append(y, e.String())
	}
	sort.Strings(x)
	sort.Strings(y)

	return reflect.DeepEqual(x, y)
}

func TestDisableTransmitPostMeasurement(t *testing.T) {
	defer gock.Off()

	t.Skipf("Skipping temporarily")

}
