package lib

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type ISubscribers interface {
	Identify(ctx context.Context, subscriberID string, data interface{}) (SubscriberResponse, error)
	BulkCreate(ctx context.Context, subscribers SubscriberBulkPayload) (SubscriberBulkCreateResponse, error)
	Get(ctx context.Context, subscriberID string) (SubscriberResponse, error)
	Update(ctx context.Context, subscriberID string, data interface{}) (SubscriberResponse, error)
	UpdateCredentials(ctx context.Context, subscriberID string, payload SubscriberCredentialPayload) (SubscriberResponse, error)
	Delete(ctx context.Context, subscriberID string) (SubscriberResponse, error)
	GetNotificationFeed(ctx context.Context, subscriberID string, opts *SubscriberNotificationFeedOptions) (*SubscriberNotificationFeedResponse, error)
	GetUnseenCount(ctx context.Context, subscriberID string, opts *SubscriberUnseenCountOptions) (*SubscriberUnseenCountResponse, error)
	MarkMessageSeen(ctx context.Context, subscriberID string, opts SubscriberMarkMessageSeenOptions) (*SubscriberNotificationFeedResponse, error)
	GetPreferences(ctx context.Context, subscriberID string) (*SubscriberPreferencesResponse, error)
	UpdatePreferences(ctx context.Context, subscriberID string, templateId string, opts *UpdateSubscriberPreferencesOptions) (*UpdateSubscriberPreferencesResponse, error)
}

type SubscriberService service

func (s *SubscriberService) Identify(ctx context.Context, subscriberID string, data interface{}) (SubscriberResponse, error) {
	var resp SubscriberResponse
	URL := s.client.config.BackendURL.JoinPath("subscribers")

	reqBody, err := s.client.mergeStruct(data, map[string]interface{}{"subscriberId": subscriberID})
	if err != nil {
		return resp, errors.Wrap(err, "unable to merge struct")
	}

	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return resp, err
	}

	_, err = s.client.sendRequest(req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *SubscriberService) BulkCreate(ctx context.Context, subscribers SubscriberBulkPayload) (SubscriberBulkCreateResponse, error) {
	var resp SubscriberBulkCreateResponse
	URL := s.client.config.BackendURL.JoinPath("subscribers", "bulk")
	jsonBody, err := json.Marshal(subscribers)
	if err != nil {
		return resp, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return resp, err
	}

	_, err = s.client.sendRequest(req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *SubscriberService) Get(ctx context.Context, subscriberID string) (SubscriberResponse, error) {
	var resp SubscriberResponse
	URL := s.client.config.BackendURL.JoinPath("subscribers", subscriberID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL.String(), http.NoBody)
	if err != nil {
		return resp, err
	}

	_, err = s.client.sendRequest(req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *SubscriberService) Update(ctx context.Context, subscriberID string, data interface{}) (SubscriberResponse, error) {
	var resp SubscriberResponse
	URL := s.client.config.BackendURL.JoinPath("subscribers", subscriberID)

	jsonBody, _ := json.Marshal(data)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, URL.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return resp, err
	}

	_, err = s.client.sendRequest(req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *SubscriberService) UpdateCredentials(ctx context.Context, subscriberID string, data SubscriberCredentialPayload) (SubscriberResponse, error) {
	var resp SubscriberResponse
	URL := s.client.config.BackendURL.JoinPath("subscribers", subscriberID, "credentials")

	jsonBody, _ := json.Marshal(data)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, URL.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return resp, err
	}

	_, err = s.client.sendRequest(req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *SubscriberService) Delete(ctx context.Context, subscriberID string) (SubscriberResponse, error) {
	var resp SubscriberResponse
	URL := s.client.config.BackendURL.JoinPath("subscribers", subscriberID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, URL.String(), http.NoBody)
	if err != nil {
		return resp, err
	}

	_, err = s.client.sendRequest(req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *SubscriberService) GetNotificationFeed(ctx context.Context, subscriberID string, opts *SubscriberNotificationFeedOptions) (*SubscriberNotificationFeedResponse, error) {
	var resp SubscriberNotificationFeedResponse
	URL := s.client.config.BackendURL.JoinPath("subscribers", subscriberID, "notifications", "feed")

	if opts != nil {
		queryValues := URL.Query()
		if opts.Payload != nil {
			var payloadOpts Base64Payload
			payloadString, err := json.Marshal(opts.Payload)
			if err != nil {
				return nil, err
			}
			opts.Payload = nil

			payloadOpts.Payload = base64.StdEncoding.EncodeToString(payloadString)
			params, err := GenerateQueryParamsFromStruct(payloadOpts)
			if err != nil {
				return nil, err
			}
			for _, param := range params {
				queryValues.Add(param.Key, param.Value)
			}
		}

		params, err := GenerateQueryParamsFromStruct(*opts)
		if err != nil {
			return nil, err
		}

		for _, param := range params {
			queryValues.Add(param.Key, param.Value)
		}

		URL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL.String(), http.NoBody)
	if err != nil {
		return nil, err
	}
	_, err = s.client.sendRequest(req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SubscriberService) GetPreferences(ctx context.Context, subscriberID string) (*SubscriberPreferencesResponse, error) {
	var resp SubscriberPreferencesResponse
	URL := s.client.config.BackendURL.JoinPath("subscribers", subscriberID, "preferences")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL.String(), http.NoBody)
	if err != nil {
		return &resp, err
	}

	_, err = s.client.sendRequest(req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SubscriberService) GetUnseenCount(ctx context.Context, subscriberID string, opts *SubscriberUnseenCountOptions) (*SubscriberUnseenCountResponse, error) {
	var resp SubscriberUnseenCountResponse
	URL := s.client.config.BackendURL.JoinPath("subscribers", subscriberID, "notifications", "unseen")

	if opts != nil {
		if opts.Seen != nil {
			queryValues := URL.Query()
			queryValues.Add("seen", strconv.FormatBool(*opts.Seen))
			URL.RawQuery = queryValues.Encode()
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL.String(), http.NoBody)
	if err != nil {
		return nil, err
	}

	_, err = s.client.sendRequest(req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SubscriberService) UpdatePreferences(ctx context.Context, subscriberID string, templateId string, opts *UpdateSubscriberPreferencesOptions) (*UpdateSubscriberPreferencesResponse, error) {
	var resp UpdateSubscriberPreferencesResponse
	URL := s.client.config.BackendURL.JoinPath("subscribers", subscriberID, "preferences", templateId)

	var reqBody io.Reader = http.NoBody

	if opts != nil {
		jsonBody, err := json.Marshal(opts)
		if err != nil {
			return nil, err
		}

		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, URL.String(), reqBody)
	if err != nil {
		return &resp, err
	}

	_, err = s.client.sendRequest(req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SubscriberService) MarkMessageSeen(ctx context.Context, subscriberID string, opts SubscriberMarkMessageSeenOptions) (*SubscriberNotificationFeedResponse, error) {
	var resp SubscriberNotificationFeedResponse
	URL := s.client.config.BackendURL.JoinPath("subscribers", subscriberID, "messages", "markAs")

	jsonBody, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	_, err = s.client.sendRequest(req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

var _ ISubscribers = &SubscriberService{}
