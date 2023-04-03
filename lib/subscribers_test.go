package lib_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/novuhq/go-novu/lib"
	"github.com/novuhq/go-novu/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

const subscriberID = "62b51a44da1af31d109f5da7"

func TestSubscriberService_Identify_Success(t *testing.T) {
	var (
		subscriberPayload lib.SubscriberPayload
		receivedBody      lib.SubscriberPayload
		expectedRequest   lib.SubscriberPayload
		expectedResponse  lib.SubscriberResponse
	)

	subscriberService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if err := json.NewDecoder(req.Body).Decode(&receivedBody); err != nil {
			log.Printf("error in unmarshalling %+v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		t.Run("Header must contain ApiKey", func(t *testing.T) {
			authKey := req.Header.Get("Authorization")
			assert.True(t, strings.Contains(authKey, novuApiKey))
			assert.True(t, strings.HasPrefix(authKey, "ApiKey"))
		})

		t.Run("URL and request method is as expected", func(t *testing.T) {
			expectedURL := "/v1/subscribers"
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Equal(t, expectedURL, req.RequestURI)
		})

		t.Run("Request is as expected", func(t *testing.T) {
			fileToStruct(filepath.Join("../testdata", "identify_subscriber.json"), &expectedRequest)
			assert.Equal(t, expectedRequest, receivedBody)
		})

		var resp lib.SubscriberResponse
		fileToStruct(filepath.Join("../testdata", "subscriber_response.json"), &resp)

		w.WriteHeader(http.StatusOK)
		bb, _ := json.Marshal(resp)
		w.Write(bb)
	}))

	defer subscriberService.Close()

	ctx := context.Background()
	fileToStruct(filepath.Join("../testdata", "identify_subscriber.json"), &subscriberPayload)

	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: utils.MustParseURL(subscriberService.URL)})

	resp, err := c.SubscriberApi.Identify(ctx, subscriberID, subscriberPayload)
	require.Nil(t, err)
	assert.NotNil(t, resp)

	t.Run("Response is as expected", func(t *testing.T) {
		fileToStruct(filepath.Join("../testdata", "subscriber_response.json"), &expectedResponse)
		assert.Equal(t, expectedResponse, resp)
	})
}

func TestSubscriberService_Update_Success(t *testing.T) {
	var (
		updateSubscriber lib.SubscriberPayload
		receivedBody     lib.SubscriberPayload
		expectedRequest  lib.SubscriberPayload
		expectedResponse lib.SubscriberResponse
	)

	subscriberService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if err := json.NewDecoder(req.Body).Decode(&receivedBody); err != nil {
			log.Printf("error in unmarshalling %+v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		t.Run("Header must contain ApiKey", func(t *testing.T) {
			authKey := req.Header.Get("Authorization")
			assert.True(t, strings.Contains(authKey, novuApiKey))
			assert.True(t, strings.HasPrefix(authKey, "ApiKey"))
		})

		t.Run("URL and request method is as expected", func(t *testing.T) {
			expectedURL := "/v1/subscribers/" + subscriberID
			assert.Equal(t, http.MethodPut, req.Method)
			assert.Equal(t, expectedURL, req.RequestURI)
		})

		t.Run("Request is as expected", func(t *testing.T) {
			fileToStruct(filepath.Join("../testdata", "update_subscriber.json"), &expectedRequest)
			assert.Equal(t, expectedRequest, receivedBody)
		})

		var resp lib.SubscriberResponse
		fileToStruct(filepath.Join("../testdata", "subscriber_response.json"), &resp)

		w.WriteHeader(http.StatusOK)
		bb, _ := json.Marshal(resp)
		w.Write(bb)
	}))

	ctx := context.Background()
	fileToStruct(filepath.Join("../testdata", "update_subscriber.json"), &updateSubscriber)

	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: utils.MustParseURL(subscriberService.URL)})

	resp, err := c.SubscriberApi.Update(ctx, subscriberID, updateSubscriber)
	require.Nil(t, err)
	assert.NotNil(t, resp)

	t.Run("Response is as expected", func(t *testing.T) {
		fileToStruct(filepath.Join("../testdata", "subscriber_response.json"), &expectedResponse)
		assert.Equal(t, expectedResponse, resp)
	})
}

func TestSubscriberService_Delete_Success(t *testing.T) {
	var expectedResponse lib.SubscriberResponse

	ctx := context.Background()

	subscriberService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t.Run("Header must contain ApiKey", func(t *testing.T) {
			authKey := req.Header.Get("Authorization")
			assert.True(t, strings.Contains(authKey, novuApiKey))
			assert.True(t, strings.HasPrefix(authKey, "ApiKey"))
		})

		t.Run("URL and request method is as expected", func(t *testing.T) {
			expectedURL := "/v1/subscribers/" + subscriberID
			assert.Equal(t, http.MethodDelete, req.Method)
			assert.Equal(t, expectedURL, req.RequestURI)
		})

		var resp lib.SubscriberResponse
		fileToStruct(filepath.Join("../testdata", "subscriber_response.json"), &resp)

		w.WriteHeader(http.StatusOK)
		bb, _ := json.Marshal(resp)
		w.Write(bb)
	}))

	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: utils.MustParseURL(subscriberService.URL)})

	resp, err := c.SubscriberApi.Delete(ctx, subscriberID)
	require.Nil(t, err)
	assert.NotNil(t, resp)

	t.Run("Response is as expected", func(t *testing.T) {
		fileToStruct(filepath.Join("../testdata", "subscriber_response.json"), &expectedResponse)
		assert.Equal(t, expectedResponse, resp)
	})
}

func TestSubscriberService_GetSubscriber_Success(t *testing.T) {
	var expectedResponse lib.SubscriberResponse
	fileToStruct(filepath.Join("../testdata", "subscriber_response.json"), &expectedResponse)

	httpServer := createTestServer(t, TestServerOptions[io.Reader, *lib.SubscriberResponse]{
		expectedURLPath:    fmt.Sprintf("/v1/subscribers/%s", subscriberID),
		expectedSentMethod: http.MethodGet,
		expectedSentBody:   http.NoBody,
		responseStatusCode: http.StatusOK,
		responseBody:       &expectedResponse,
	})

	ctx := context.Background()
	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: utils.MustParseURL(httpServer.URL)})
	resp, err := c.SubscriberApi.Get(ctx, subscriberID)

	require.NoError(t, err)
	require.Equal(t, resp, expectedResponse)
}

func TestSubscriberService_GetPreferences_Success(t *testing.T) {
	var expectedResponse *lib.SubscriberPreferencesResponse
	fileToStruct(filepath.Join("../testdata", "subscriber_preferences_response.json"), &expectedResponse)

	httpServer := createTestServer(t, TestServerOptions[map[string]string, *lib.SubscriberPreferencesResponse]{
		expectedURLPath:    fmt.Sprintf("/v1/subscribers/%s/preferences", subscriberID),
		expectedSentMethod: http.MethodGet,
		expectedSentBody:   map[string]string{},
		responseStatusCode: http.StatusOK,
		responseBody:       expectedResponse,
	})

	ctx := context.Background()
	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: utils.MustParseURL(httpServer.URL)})
	resp, err := c.SubscriberApi.GetPreferences(ctx, subscriberID)

	require.NoError(t, err)
	require.Equal(t, resp, expectedResponse)
}

func TestSubscriberService_UpdatePreferences_Success(t *testing.T) {
	var topicID = "topicId"

	var expectedResponse *lib.SubscriberPreferencesResponse
	fileToStruct(filepath.Join("../testdata", "subscriber_preferences_response.json"), &expectedResponse)

	var opts *lib.UpdateSubscriberPreferencesOptions = &lib.UpdateSubscriberPreferencesOptions{
		Enabled: true,
		Channel: []lib.UpdateSubscriberPreferencesChannel{
			{
				Type:    "email",
				Enabled: true,
			},
		},
	}
	httpServer := createTestServer(t, TestServerOptions[*lib.UpdateSubscriberPreferencesOptions, *lib.SubscriberPreferencesResponse]{
		expectedURLPath:    fmt.Sprintf("/v1/subscribers/%s/preferences/%s", subscriberID, topicID),
		expectedSentMethod: http.MethodPatch,
		expectedSentBody:   opts,
		responseStatusCode: http.StatusOK,
		responseBody:       expectedResponse,
	})

	ctx := context.Background()
	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: utils.MustParseURL(httpServer.URL)})
	resp, err := c.SubscriberApi.UpdatePreferences(ctx, subscriberID, topicID, opts)

	require.NoError(t, err)
	require.Equal(t, resp, expectedResponse)
}
