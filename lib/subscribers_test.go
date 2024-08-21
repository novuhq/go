package lib_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/novuhq/go-novu/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(subscriberService.URL)})

	resp, err := c.SubscriberApi.Identify(ctx, subscriberID, subscriberPayload)
	require.Nil(t, err)
	assert.NotNil(t, resp)

	t.Run("Response is as expected", func(t *testing.T) {
		fileToStruct(filepath.Join("../testdata", "subscriber_response.json"), &expectedResponse)
		assert.Equal(t, expectedResponse, resp)
	})
}

func TestSubscriberService_BulkCreate_Success(t *testing.T) {
	var (
		subscriberBulkPayload lib.SubscriberBulkPayload
		receivedBody          lib.SubscriberBulkPayload
		expectedRequest       lib.SubscriberBulkPayload
		expectedResponse      lib.SubscriberBulkCreateResponse
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
			expectedURL := "/v1/subscribers/bulk"
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Equal(t, expectedURL, req.RequestURI)
		})

		t.Run("Request is as expected", func(t *testing.T) {
			fileToStruct(filepath.Join("../testdata", "subscriber_bulk.json"), &expectedRequest)
			assert.Equal(t, expectedRequest, receivedBody)
		})

		var resp lib.SubscriberResponse
		fileToStruct(filepath.Join("../testdata", "subscriber_bulk_response.json"), &resp)

		w.WriteHeader(http.StatusOK)
		bb, _ := json.Marshal(resp)
		w.Write(bb)
	}))

	defer subscriberService.Close()

	ctx := context.Background()
	fileToStruct(filepath.Join("../testdata", "subscriber_bulk.json"), &subscriberBulkPayload)

	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(subscriberService.URL)})

	resp, err := c.SubscriberApi.BulkCreate(ctx, subscriberBulkPayload)
	require.Nil(t, err)
	assert.NotNil(t, resp)

	t.Run("Response is as expected", func(t *testing.T) {
		fileToStruct(filepath.Join("../testdata", "subscriber_bulk_response.json"), &expectedResponse)
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

	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(subscriberService.URL)})

	resp, err := c.SubscriberApi.Update(ctx, subscriberID, updateSubscriber)
	require.Nil(t, err)
	assert.NotNil(t, resp)

	t.Run("Response is as expected", func(t *testing.T) {
		fileToStruct(filepath.Join("../testdata", "subscriber_response.json"), &expectedResponse)
		assert.Equal(t, expectedResponse, resp)
	})
}

func TestSubscriberService_Update_Credentials_Success(t *testing.T) {
	var (
		updateSubscriberCreds lib.SubscriberCredentialPayload
		receivedBody          lib.SubscriberCredentialPayload
		expectedRequest       lib.SubscriberCredentialPayload
		expectedResponse      lib.SubscriberResponse
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
			expectedURL := "/v1/subscribers/" + subscriberID + "/credentials"
			assert.Equal(t, http.MethodPut, req.Method)
			assert.Equal(t, expectedURL, req.RequestURI)
		})

		t.Run("Request is as expected", func(t *testing.T) {
			fileToStruct(filepath.Join("../testdata", "update_subscriber_credentials.json"), &expectedRequest)
			assert.Equal(t, expectedRequest, receivedBody)
		})

		var resp lib.SubscriberResponse
		fileToStruct(filepath.Join("../testdata", "subscriber_response.json"), &resp)

		w.WriteHeader(http.StatusOK)
		bb, _ := json.Marshal(resp)
		w.Write(bb)
	}))

	ctx := context.Background()
	fileToStruct(filepath.Join("../testdata", "update_subscriber_credentials.json"), &updateSubscriberCreds)

	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(subscriberService.URL)})

	resp, err := c.SubscriberApi.UpdateCredentials(ctx, subscriberID, updateSubscriberCreds)
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

	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(subscriberService.URL)})

	resp, err := c.SubscriberApi.Delete(ctx, subscriberID)
	require.Nil(t, err)
	assert.NotNil(t, resp)

	t.Run("Response is as expected", func(t *testing.T) {
		fileToStruct(filepath.Join("../testdata", "subscriber_response.json"), &expectedResponse)
		assert.Equal(t, expectedResponse, resp)
	})
}

func TestSubscriberService_GetNotificationFeed_Success(t *testing.T) {
	var expectedResponse *lib.SubscriberNotificationFeedResponse
	fileToStruct(filepath.Join("../testdata", "subscriber_notification_feed_response.json"), &expectedResponse)

	page := 1
	seen := true
	feedIdentifier := "feed_identifier"
	payload := map[string]interface{}{
		"name": "test",
	}

	opts := lib.SubscriberNotificationFeedOptions{
		Page:           page,
		Seen:           seen,
		FeedIdentifier: feedIdentifier,
		Payload:        payload,
	}

	httpServer := createTestServer(t, TestServerOptions[io.Reader, *lib.SubscriberNotificationFeedResponse]{
		expectedURLPath:    fmt.Sprintf("/v1/subscribers/%s/notifications/feed?feedIdentifier=%s&page=%s&payload=eyJuYW1lIjoidGVzdCJ9&seen=%s", subscriberID, feedIdentifier, strconv.Itoa(page), strconv.FormatBool(seen)),
		expectedSentMethod: http.MethodGet,
		expectedSentBody:   http.NoBody,
		responseStatusCode: http.StatusOK,
		responseBody:       expectedResponse,
	})

	ctx := context.Background()
	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})
	resp, err := c.SubscriberApi.GetNotificationFeed(ctx, subscriberID, &opts)

	require.NoError(t, err)
	require.Equal(t, resp, expectedResponse)
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
	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})
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
	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})
	resp, err := c.SubscriberApi.GetPreferences(ctx, subscriberID)

	require.NoError(t, err)
	require.Equal(t, resp, expectedResponse)
}

func TestSubscriberService_GetUnseenCount_Success(t *testing.T) {
	var expectedResponse *lib.SubscriberUnseenCountResponse
	fileToStruct(filepath.Join("../testdata", "subscriber_notification_feed_unseen.json"), &expectedResponse)

	seen := false

	opts := lib.SubscriberUnseenCountOptions{
		Seen: &seen,
	}

	httpServer := createTestServer(t, TestServerOptions[io.Reader, *lib.SubscriberUnseenCountResponse]{
		expectedURLPath:    fmt.Sprintf("/v1/subscribers/%s/notifications/unseen?seen=false", subscriberID),
		expectedSentMethod: http.MethodGet,
		expectedSentBody:   http.NoBody,
		responseStatusCode: http.StatusOK,
		responseBody:       expectedResponse,
	})

	ctx := context.Background()
	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})
	resp, err := c.SubscriberApi.GetUnseenCount(ctx, subscriberID, &opts)

	require.NoError(t, err)
	require.Equal(t, resp, expectedResponse)
}

func TestSubscriberService_MarkMessageSeen(t *testing.T) {
	var expectedResponse *lib.SubscriberNotificationFeedResponse
	fileToStruct(filepath.Join("../testdata", "subscriber_notification_feed_response.json"), &expectedResponse)

	opts := lib.SubscriberMarkMessageSeenOptions{
		MessageID: "message_id",
		Seen:      true,
		Read:      true,
	}

	httpServer := createTestServer(t, TestServerOptions[lib.SubscriberMarkMessageSeenOptions, *lib.SubscriberNotificationFeedResponse]{
		expectedURLPath:    fmt.Sprintf("/v1/subscribers/%s/messages/markAs", subscriberID),
		expectedSentMethod: http.MethodPost,
		expectedSentBody:   opts,
		responseStatusCode: http.StatusOK,
		responseBody:       expectedResponse,
	})

	ctx := context.Background()
	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})
	resp, err := c.SubscriberApi.MarkMessageSeen(ctx, subscriberID, opts)
	require.NoError(t, err)
	require.Equal(t, resp, expectedResponse)
}

func TestSubscriberService_UpdatePreferences_Success(t *testing.T) {
	var topicID = "topicId"

	var expectedResponse *lib.SubscriberPreferencesResponse
	fileToStruct(filepath.Join("../testdata", "subscriber_preferences_response.json"), &expectedResponse)

	enabled := true
	var opts *lib.UpdateSubscriberPreferencesOptions = &lib.UpdateSubscriberPreferencesOptions{
		Enabled: &enabled,
		Channel: &lib.UpdateSubscriberPreferencesChannel{
			Type:    "email",
			Enabled: true,
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
	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})
	resp, err := c.SubscriberApi.UpdatePreferences(ctx, subscriberID, topicID, opts)

	require.NoError(t, err)
	require.Equal(t, resp, expectedResponse)
}
