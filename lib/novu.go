package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
)

const (
	NovuURL     = "https://api.novu.co"
	NovuVersion = "v1"
)

type RetryConfigType struct {
	InitialDelay time.Duration // inital delay
	WaitMin      time.Duration // Minimum time to wait
	WaitMax      time.Duration // Maximum time to wait
	RetryMax     int           // Maximum number of retries
}

type Config struct {
	BackendURL        *url.URL
	HttpClient        *http.Client
	HttpClientHeaders map[string]string
	RetryConfig       *RetryConfigType
}

type APIClient struct {
	apiKey string
	config *Config
	common service

	// Api Service
	BlueprintApi     *BlueprintService
	ChangesApi       *ChangesService
	SubscriberApi    *SubscriberService
	EventApi         *EventService
	ExecutionsApi    *ExecutionsService
	MessagesApi      *MessagesService
	FeedsApi         *FeedsService
	TopicsApi        *TopicService
	IntegrationsApi  *IntegrationService
	InboundParserApi *InboundParserService
	LayoutApi        *LayoutService
	TenantApi        *TenantService
}

type service struct {
	client *APIClient
}

func NewAPIClient(apiKey string, cfg *Config) *APIClient {
	cfg.BackendURL = buildBackendURL(cfg)

	if cfg.HttpClient == nil {
		retyableClient := retryablehttp.NewClient()
		if cfg.RetryConfig != nil {
			retyableClient.RetryWaitMin = cfg.RetryConfig.WaitMin
			retyableClient.RetryWaitMax = cfg.RetryConfig.WaitMax
			retyableClient.RetryMax = cfg.RetryConfig.RetryMax
			retyableClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
				if resp != nil {
					if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
						if s, ok := resp.Header["Retry-After"]; ok {
							if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
								return time.Second * time.Duration(sleep)
							}
						}
					}
				}
				if attemptNum == 0 {
					return cfg.RetryConfig.InitialDelay //wait for InitialDelay on 1st retry
				}
				mult := math.Pow(2, float64(attemptNum)) * float64(min)
				sleep := time.Duration(mult)
				//float64(sleep) != mult is to make sure there is no conversion error
				//if there is a conversion error, number is huge and we set the sleep to max
				if float64(sleep) != mult || sleep > max {
					sleep = max
				}
				return sleep
			}
		} else {
			retyableClient.RetryMax = 0 //by default no retry
		}
		cfg.HttpClient = retyableClient.StandardClient()
	}

	c := &APIClient{apiKey: apiKey}
	c.config = cfg
	c.common.client = c

	// API Services
	c.ChangesApi = (*ChangesService)(&c.common)
	c.EventApi = (*EventService)(&c.common)
	c.ExecutionsApi = (*ExecutionsService)(&c.common)
	c.FeedsApi = (*FeedsService)(&c.common)
	c.SubscriberApi = (*SubscriberService)(&c.common)
	c.MessagesApi = (*MessagesService)(&c.common)
	c.TopicsApi = (*TopicService)(&c.common)
	c.IntegrationsApi = (*IntegrationService)(&c.common)
	c.InboundParserApi = (*InboundParserService)(&c.common)
	c.LayoutApi = (*LayoutService)(&c.common)
	c.BlueprintApi = (*BlueprintService)(&c.common)
	c.TenantApi = (*TenantService)(&c.common)
	return c
}

func (c APIClient) sendRequest(req *http.Request, resp interface{}) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("ApiKey %s", c.apiKey))
	req.Header.Set("Idempotency-Key", uuid.New().String())

	headers := c.config.HttpClientHeaders
	if headers != nil {
		for header, value := range headers {
			req.Header.Set(header, value)
		}
	}

	res, err := c.config.HttpClient.Do(req)
	if err != nil {
		return res, errors.Wrap(err, "failed to execute request")
	}

	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	if res.StatusCode >= http.StatusMultipleChoices {
		return res, errors.Errorf(
			`request was not successful, status code %d, %s`, res.StatusCode,
			string(body),
		)
	}

	if string(body) == "" {
		resp = map[string]string{}
		return res, nil
	}

	err = c.decode(&resp, body)
	if err != nil {
		return res, errors.Wrap(err, "unable to unmarshal response body")
	}

	return res, nil
}

func (c APIClient) mergeStruct(target, patch interface{}) (interface{}, error) {
	var m map[string]interface{}

	targetPayload, _ := json.Marshal(target)
	patchPayload, _ := json.Marshal(patch)

	_ = json.Unmarshal(targetPayload, &m)
	_ = json.Unmarshal(patchPayload, &m)

	return m, nil
}

func (c APIClient) decode(v interface{}, b []byte) (err error) {
	if err = json.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}

func buildBackendURL(cfg *Config) *url.URL {

	if cfg.BackendURL == nil {
		rawURL := fmt.Sprintf("%s/%s", NovuURL, NovuVersion)
		return MustParseURL(rawURL)
	}

	if strings.Contains(cfg.BackendURL.String(), "novu.co/v") {
		return cfg.BackendURL
	}

	return cfg.BackendURL.JoinPath(NovuVersion)
}
