package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/stretchr/testify/assert"
)

/*

// TODO: Come back to this to ensure we aren't sending duplicate dimensions to Cloudwatch

func ensureUniqueMap(m []map[string]string) error {
	tmpMap := make(map[string]string)

	fmt.Println(m)

	for _, val := range m {
		for k, v := range val {
			if _, ok := tmpMap[k]; ok {
				fmt.Printf("Found duplicate key %s\n", k)
				return errors.New("Duplicate key found")
			}
			tmpMap[k] = v
		}
	}

	return nil
}
*/
type mockCloudwatchClient struct{}

func (c mockCloudwatchClient) PutMetricData(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error) {
	// TODO add some checks here
	return &cloudwatch.PutMetricDataOutput{}, nil
}

func TestNewDroneClient(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer testServer.Close()

	os.Setenv("DRONE_SERVER", testServer.URL)
	os.Setenv("DRONE_TOKEN", "ci")

	c := newDroneClient()

	_, err := c.Self()

	assert.Equal(t, nil, err, "Encountered unexpected error")
}

func TestGetQueuedBuilds(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer testServer.Close()

	os.Setenv("DRONE_SERVER", testServer.URL)
	os.Setenv("DRONE_TOKEN", "ci")

	c := newDroneClient()

	r := getQueuedBuilds(c)

	assert.Len(t, r, 2, fmt.Sprintf("Length of builds was %d", len(r)))
}

func TestVerifyEnvVars(t *testing.T) {
	os.Setenv("DRONE_SERVER", "foo")
	os.Setenv("DRONE_TOKEN", "bar")
	os.Setenv("CLOUDWATCH_METRICS_NAMESPACE", "foobar")

	assert.Equal(t, nil, verifyEnvVars(), "Did not expect verifyEnvVars to return an error")

	os.Unsetenv("DRONE_SERVER")

	assert.NotEqual(t, nil, verifyEnvVars(), "Expected verifyEnvVars to return an error")
}

func TestReportBuilds(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer testServer.Close()

	cwc := mockCloudwatchClient{}

	os.Setenv("DRONE_SERVER", testServer.URL)
	os.Setenv("DRONE_TOKEN", "ci")

	c := newDroneClient()

	s := getQueuedBuilds(c)

	reportBuilds(c, cwc, s)
}

func mockHandler(w http.ResponseWriter, r *http.Request) {
	userBody := `{"id":1,"login":"ciuser","email":"","machine":false,"admin":true,"active":true,"avatar":"https://avatars.githubusercontent.com/u/1?v=4","syncing":false,"synced":1617653050,"created":1615672091,"updated":1615672091,"last_login":1617652785}`

	queueBody := `[
		{
		"id":123,"repo_id":2,"build_id":4,"number":1,"name":"ci",
		"kind":"pipeline","type":"docker","status":"running","errignore":false,
		"exit_code":0,"machine":"cimachine","os":"linux","arch":"amd64",
		"started":1617666144,"stopped":0,"created":1617666144,"updated":1617666144,"version":3,
		"on_success":true,"on_failure":false,
		"labels":{"class":"standard","os":"linux"}
		},
		{
		"id":123,"repo_id":2,"build_id":4,"number":1,"name":"ci",
		"kind":"pipeline","type":"docker","status":"running","errignore":false,
		"exit_code":0,"machine":"cimachine","os":"linux","arch":"amd64",
		"started":1617666144,"stopped":0,"created":1617666144,"updated":1617666144,"version":3,
		"on_success":true,"on_failure":false,
		"labels":{"class":"standard","os":"linux"}
		}
	]`

	routes := []struct {
		verb string
		path string
		body string
		code int
	}{
		{
			verb: "GET",
			path: "/api/user",
			body: userBody,
			code: 200,
		},
		{
			verb: "GET",
			path: "/api/queue",
			body: queueBody,
			code: 200,
		},
	}

	path := r.URL.Path
	verb := r.Method

	for _, route := range routes {
		if route.verb != verb {
			continue
		}
		if route.path != path {
			continue
		}
		if route.code == 204 {
			w.WriteHeader(204)
			return
		}

		w.WriteHeader(route.code)
		w.Write([]byte(route.body))
		return
	}
	w.WriteHeader(404)
}
