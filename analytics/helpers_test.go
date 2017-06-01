package analytics

import (
	"reflect"
	"sort"
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

// GockBootstrap returns a UserSession ready for use in testing the GA service.
// Note: This endpoint IRL ALWAYS returns 200 OK, and sends back a small 1px gif image.
// This means we must rely on whether our hit is put together properly.
// Validate at https://ga-dev-tools.appspot.com/hit-builder/
func MockUserSession() *UserSession {
	u := *NewUserSession("test-uid")
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
