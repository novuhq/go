package lib

import (
	"io"
	"time"
)

type ChannelType string
type GeneralError error
type ProviderIdType string

const Version = "v1"

const (
	HTTPStatusOk       = 200
	HTTPStatusCreated  = 201
	HTTPStatusNoChange = 204
	HTTPRedirectOk     = 300
)

const (
	INAPP ChannelType = "in_app"
	EMAIL ChannelType = "email"
	SMS   ChannelType = "sms"
	CHAT  ChannelType = "chat"
	PUSH  ChannelType = "push"
)

const (
	slack       ProviderIdType = "slack"
	discord     ProviderIdType = "discord"
	msteams     ProviderIdType = "msteams"
	mattermost  ProviderIdType = "mattermost"
	fcm         ProviderIdType = "fcm"
	apns        ProviderIdType = "apns"
	expo        ProviderIdType = "expo"
	oneSignal   ProviderIdType = "one-signal"
	pushWebhook ProviderIdType = "push-webhook"
)

type Data struct {
	Acknowledged bool   `json:"acknowledged"`
	Status       string `json:"status"`
}

type Response struct {
	Data Data `json:"data"`
}

type ITriggerPayloadOptions struct {
	To            interface{} `json:"to,omitempty"`
	Payload       interface{} `json:"payload,omitempty"`
	Overrides     interface{} `json:"overrides,omitempty"`
	TransactionId string      `json:"transactionId,omitempty"`
	Actor         interface{} `json:"actor,omitempty"`
}

type TriggerRecipientsTypeArray interface {
	[]string | []SubscriberPayload
}
type TriggerRecipientsTypeSingle interface {
	string | SubscriberPayload
}

type TriggerTopicRecipientsTypeSingle struct {
	TopicKey string `json:"topicKey,omitempty"`
	Type     string `json:"type,omitempty"`
}

type SubscriberPayload struct {
	FirstName    string                 `json:"firstName,omitempty"`
	LastName     string                 `json:"lastName,omitempty"`
	Email        string                 `json:"email,omitempty"`
	Phone        string                 `json:"phone,omitempty"`
	Avatar       string                 `json:"avatar,omitempty"`
	Locale       string                 `json:"locale,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
	SubscriberId string                 `json:"subscriberId"`
}

type SubscriberBulkPayload struct {
	Subscribers []SubscriberPayload `json:"subscribers"`
}

type TriggerRecipientsType interface {
	TriggerRecipientsTypeSingle | TriggerRecipientsTypeArray | TriggerTopicRecipientsTypeSingle | []TriggerTopicRecipientsTypeSingle
}

type ITriggerPayload interface {
	string | []string | bool | int64 | IAttachmentOptions | []IAttachmentOptions
}

type IAttachmentOptions struct {
	Mime     string        `json:"mime,omitempty"`
	File     io.Reader     `json:"file,omitempty"`
	Name     string        `json:"name,omitempty"`
	Channels []ChannelType `json:"channels,omitempty"`
}
type JsonResponse struct {
	Data interface{} `json:"data"`
}
type ExecutionsQueryParams struct {
	NotificationId string
	SubscriberId   string
}

type EventResponse struct {
	JsonResponse
}

type EventRequest struct {
	Name          string      `json:"name"`
	To            interface{} `json:"to"`
	Payload       interface{} `json:"payload"`
	Overrides     interface{} `json:"overrides,omitempty"`
	TransactionId string      `json:"transactionId,omitempty"`
	Actor         interface{} `json:"actor,omitempty"`
}

type MessagesQueryParams struct {
	Channel       string
	SubscriberId  string
	TransactionId []string
	Page          int
	Limit         int
}

// QueryBuilder gives us an interface to pass as arg to our API methods.
// See messages.go for an example of implementing this interface
type QueryBuilder interface {
	BuildQuery() string
}

type SubscriberResponse struct {
	JsonResponse
}

type SubscriberBulkCreateResponse struct {
	Data struct {
		Updated []struct {
			SubscriberId string `json:"subscriberId"`
		} `json:"updated"`
		Created []struct {
			SubscriberId string `json:"subscriberId"`
		} `json:"created"`
		Failed []interface{} `json:"failed"`
	} `json:"data"`
}

type Template struct {
	ID       string `json:"_id"`
	Name     string `json:"name"`
	Critical bool   `json:"critical"`
}

type Preference struct {
	Enabled  bool    `json:"enabled"`
	Channels Channel `json:"channels"`
}

type Channel struct {
	Email *bool `json:"email,omitempty"`
	Sms   *bool `json:"sms,omitempty"`
	Chat  *bool `json:"chat,omitempty"`
	InApp *bool `json:"in_app,omitempty"`
	Push  *bool `json:"push,omitempty"`
}

type SubscriberPreferencesResponse struct {
	Data []struct {
		Template   Template   `json:"template"`
		Preference Preference `json:"preference"`
	} `json:"data"`
}

type UpdateSubscriberPreferencesResponse struct {
	Data struct {
		Template   Template   `json:"template"`
		Preference Preference `json:"preference"`
	} `json:"data"`
}

type UpdateSubscriberPreferencesChannel struct {
	Type    ChannelType `json:"type"`
	Enabled bool        `json:"enabled"`
}

type UpdateSubscriberPreferencesOptions struct {
	Channel *UpdateSubscriberPreferencesChannel `json:"channel,omitempty"`
	Enabled *bool                               `json:"enabled,omitempty"`
}

type ListTopicsResponse struct {
	Page       int                `json:"name"`
	PageSize   int                `json:"pageSize"`
	TotalCount int                `json:"totalCount"`
	Data       []GetTopicResponse `json:"data"`
}

type GetTopicResponse struct {
	Id             string   `json:"_id"`
	OrganizationId string   `json:"_organizationId"`
	EnvironmentId  string   `json:"_environmentId"`
	Key            string   `json:"key"`
	Name           string   `json:"name"`
	Subscribers    []string `json:"subscribers"`
}

type CheckTopicSubscriberResponse struct {
	OrganizationId       string `json:"_organizationId"`
	EnvironmentId        string `json:"_environmentId"`
	SubsriberId          string `json:"_subscriberId"`
	Id                   string `json:"_topicId"`
	Key                  string `json:"topicKey"`
	ExternalSubscriberId string `json:"externalSubscriberId"`
}

type ListTopicsOptions struct {
	Page     *int    `json:"page,omitempty"`
	PageSize *int    `json:"pageSize,omitempty"`
	Key      *string `json:"key,omitempty"`
}

type CreateTopicRequest struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type RenameTopicRequest struct {
	Name string `json:"name"`
}

type SubscribersTopicRequest struct {
	Subscribers []string `json:"subscribers"`
}

type SubscriberNotificationFeedOptions struct {
	Page           int         `queryKey:"page"`
	FeedIdentifier string      `queryKey:"feedIdentifier"`
	Seen           bool        `queryKey:"seen"`
	Payload        interface{} `queryKey:"payload"`
}
type Base64Payload struct {
	Payload string `queryKey:"payload"`
}

type SubscriberUnseenCountOptions struct {
	Seen *bool `json:"seen"`
}

type SubscriberMarkMessageSeenOptions struct {
	MessageID string `json:"messageId"`
	Seen      bool   `json:"seen"`
	Read      bool   `json:"read"`
}

type NotificationFeedData struct {
	CTA              CTA       `json:"cta"`
	Channel          string    `json:"channel"`
	Content          string    `json:"content"`
	CreatedAt        time.Time `json:"createdAt"`
	Deleted          bool      `json:"deleted"`
	DeviceTokens     []string  `json:"deviceTokens"`
	DirectWebhookURL string    `json:"directWebhookUrl"`
	EnvironmentID    string    `json:"_environmentId"`
	ErrorID          string    `json:"errorId"`
	ErrorText        string    `json:"errorText"`
	FeedID           string    `json:"_feedId"`
	ID               string    `json:"_id"`
	JobID            string    `json:"_jobId"`
	LastReadDate     time.Time `json:"lastReadDate"`
	LastSeenDate     time.Time `json:"lastSeenDate"`
	MessageTemplate  string    `json:"_messageTemplateId"`
	NotificationID   string    `json:"_notificationId"`
	OrganizationID   string    `json:"_organizationId"`
	Payload          struct {
		UpdateMessage string `json:"updateMessage"`
	} `json:"payload"`
	ProviderID string `json:"providerId"`
	Read       bool   `json:"read"`
	ResponseID string `json:"id"`
	Seen       bool   `json:"seen"`
	Status     string `json:"status"`
	Subscriber struct {
		ID           string `json:"_id"`
		SubscriberID string `json:"subscriberId"`
	} `json:"subscriber"`
	SubscriberID       string    `json:"_subscriberId"`
	TemplateID         string    `json:"_templateId"`
	TemplateIdentifier string    `json:"templateIdentifier"`
	TransactionID      string    `json:"transactionId"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type SubscriberNotificationFeedResponse struct {
	TotalCount int                    `json:"totalCount"`
	Data       []NotificationFeedData `json:"data"`
	PageSize   int                    `json:"pageSize"`
	Page       int                    `json:"page"`
}

type SubscriberUnseenCountResponse struct {
	Data struct {
		Count int `json:"count"`
	} `json:"data"`
}

type Credentials struct {
	WebhookUrl   string   `json:"webhookUrl,omitempty"`
	Channel      string   `json:"channel,omitempty"`
	DeviceTokens []string `json:"deviceTokens,omitempty"`
}

type SubscriberCredentialPayload struct {
	Credentials           Credentials    `json:"credentials"`
	IntegrationIdentifier string         `json:"integrationIdentifier,omitempty"`
	ProviderId            ProviderIdType `json:"providerId"`
}

type CTA struct {
	Type   string `json:"type"`
	Action struct {
		Status  string `json:"status"`
		Buttons []struct {
			Type          string `json:"type"`
			Content       string `json:"content"`
			ResultContent string `json:"resultContent"`
		} `json:"buttons"`
		Result struct {
			Payload map[string]interface{} `json:"payload"`
			Type    string                 `json:"type"`
		} `json:"result"`
	}
}
type IntegrationCredentials struct {
	ApiKey           string                 `json:"apiKey,omitempty"`
	User             string                 `json:"user,omitempty"`
	SecretKey        string                 `json:"secretKey,omitempty"`
	Domain           string                 `json:"domain,omitempty"`
	Password         string                 `json:"password,omitempty"`
	Host             string                 `json:"host,omitempty"`
	Port             string                 `json:"port,omitempty"`
	Secure           bool                   `json:"secure,omitempty"`
	Region           string                 `json:"region,omitempty"`
	AccountSID       string                 `json:"accountSid,omitempty"`
	MessageProfileID string                 `json:"messageProfileId,omitempty"`
	Token            string                 `json:"token,omitempty"`
	From             string                 `json:"from,omitempty"`
	SenderName       string                 `json:"senderName,omitempty"`
	ProjectName      string                 `json:"projectName,omitempty"`
	ApplicationID    string                 `json:"applicationId,omitempty"`
	ClientID         string                 `json:"clientId,omitempty"`
	RequireTls       bool                   `json:"requireTls,omitempty"`
	IgnoreTls        bool                   `json:"ignoreTls,omitempty"`
	TlsOptions       map[string]interface{} `json:"tlsOptions,omitempty"`
}

type CreateIntegrationRequest struct {
	ProviderID  string                 `json:"providerId"`
	Channel     ChannelType            `json:"channel"`
	Credentials IntegrationCredentials `json:"credentials,omitempty"`
	Active      bool                   `json:"active"`
	Check       bool                   `json:"check"`
}

type UpdateIntegrationRequest struct {
	Credentials IntegrationCredentials `json:"credentials"`
	Active      bool                   `json:"active"`
	Check       bool                   `json:"check"`
}

type Integration struct {
	Id             string                 `json:"_id"`
	EnvironmentID  string                 `json:"_environmentId"`
	OrganizationID string                 `json:"_organizationId"`
	ProviderID     string                 `json:"providerId"`
	Channel        ChannelType            `json:"channel"`
	Credentials    IntegrationCredentials `json:"credentials"`
	Active         bool                   `json:"active"`
	Deleted        bool                   `json:"deleted"`
	UpdatedAt      string                 `json:"updatedAt"`
	DeletedAt      string                 `json:"deletedAt"`
	DeletedBy      string                 `json:"deletedBy"`
}

type IntegrationResponse struct {
	Data Integration `json:"data"`
}

type GetIntegrationsResponse struct {
	Data []Integration `json:"data"`
}

type IntegrationChannelLimitResponse struct {
	Data struct {
		Limit int `json:"limit"`
		Count int `json:"count"`
	} `json:"data"`
}

type SetIntegrationAsPrimaryResponse struct {
	Data struct {
		ID             string                 `json:"_id"`
		EnvironmentID  string                 `json:"_environmentId"`
		OrganizationID string                 `json:"_organizationId"`
		Name           string                 `json:"name"`
		Identifier     string                 `json:"identifier"`
		ProviderID     string                 `json:"providerId"`
		Channel        string                 `json:"channel"`
		Credentials    IntegrationCredentials `json:"credentials"`
		Active         bool                   `json:"active"`
		Deleted        bool                   `json:"deleted"`
		DeletedAt      string                 `json:"deletedAt"`
		DeletedBy      string                 `json:"deletedBy"`
		Primary        bool                   `json:"primary"`
	} `json:"data"`
}
type BulkTriggerOptions struct {
	Name          interface{} `json:"name,omitempty"`
	To            interface{} `json:"to,omitempty"`
	Payload       interface{} `json:"payload,omitempty"`
	Overrides     interface{} `json:"overrides,omitempty"`
	TransactionId string      `json:"transactionId,omitempty"`
	Actor         interface{} `json:"actor,omitempty"`
}

type BulkTriggerEvent struct {
	Events []BulkTriggerOptions `json:"events"`
}

type BroadcastEventToAll struct {
	Name          interface{} `json:"name,omitempty"`
	Payload       interface{} `json:"payload,omitempty"`
	Overrides     interface{} `json:"overrides,omitempty"`
	TransactionId string      `json:"transactionId,omitempty"`
	Actor         interface{} `json:"actor,omitempty"`
}
type MxRecordConfiguredStatus struct {
	MxRecordConfigured bool `json:"mxRecordConfigured"`
}
type InboundParserResponse struct {
	Data MxRecordConfiguredStatus `json:"data"`
}

type CreateLayoutRequest struct {
	Name        string        `json:"name"`
	Identifier  string        `json:"identifier"`
	Description string        `json:"description"`
	Content     string        `json:"content"`
	Variables   []interface{} `json:"variables,omitempty"`
	IsDefault   bool          `json:"isDefault,omitempty"`
}

type CreateLayoutResponse struct {
	Data struct {
		Id string `json:"_id"`
	} `json:"data"`
}
type LayoutRequestOptions struct {
	Page     *int    `json:"page,omitempty"`
	PageSize *int    `json:"pageSize,omitempty"`
	Key      *string `json:"key,omitempty"`
	OrderBy  *int    `json:"orderBy,omitempty"`
}
type LayoutResponse struct {
	Id             string        `json:"_id"`
	OrganizationId string        `json:"_organizationId"`
	EnvironmentId  string        `json:"_environmentId"`
	CreatorId      string        `json:"_creatorId"`
	Name           string        `json:"name"`
	Identifier     string        `json:"identifier"`
	Description    string        `json:"description"`
	Channel        string        `json:"channel"`
	Content        string        `json:"content"`
	ContentType    string        `json:"contentType"`
	Variables      []interface{} `json:"variables"`
	IsDefault      bool          `json:"isDefault"`
	IsDeleted      bool          `json:"isDeleted"`
	CreatedAt      string        `json:"createdAt"`
	UpdatedAt      string        `json:"updatedAt"`
	ParentId       string        `json:"_parentId"`
}
type LayoutsResponse struct {
	TotalCount int              `json:"totalCount"`
	Data       []LayoutResponse `json:"data"`
	PageSize   int              `json:"pageSize"`
	Page       int              `json:"page"`
}
type BlueprintByTemplateIdResponse struct {
	Id                  string        `json:"_id,omitempty"`
	Name                string        `json:"name,omitempty"`
	Description         string        `json:"description,omitempty"`
	Active              bool          `json:"active,omitempty"`
	Draft               bool          `json:"draft,omitempty"`
	PreferenceSettings  interface{}   `json:"preferenceSettings,omitempty"`
	Critical            bool          `json:"critical,omitempty"`
	Tags                []string      `json:"tags,omitempty"`
	Steps               []interface{} `json:"steps,omitempty"`
	OrganizationID      string        `json:"_organizationId,omitempty"`
	CreatorID           string        `json:"_creatorId,omitempty"`
	EnvironmentID       string        `json:"_environmentId,omitempty"`
	Triggers            []interface{} `json:"triggers,omitempty"`
	NotificationGroupID string        `json:"_notificationGroupId,omitempty"`
	ParentId            string        `json:"_parentId,omitempty"`
	Deleted             bool          `json:"deleted,omitempty"`
	DeletedAt           string        `json:"deletedAt,omitempty"`
	DeletedBy           string        `json:"deletedBy,omitempty"`
	CreatedAt           string        `json:"createdAt,omitempty"`
	UpdatedAt           string        `json:"updatedAt,omitempty"`
	NotificationGroup   interface{}   `json:"notificationGroup,omitempty"`
	IsBlueprint         bool          `json:"isBlueprint,omitempty"`
	BlueprintID         string        `json:"blueprintId,omitempty"`
}

type BlueprintGroupByCategoryResponse struct {
	General []interface{} `json:"general,omitempty"`
	Popular interface{}   `json:"popular,omitempty"`
}

type ChangesGetQuery struct {
	Page     int    `json:"page,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Promoted string `json:"promoted,omitempty"`
}

type ChangesGetResponseData struct {
	Id             string      `json:"_id,omitempty"`
	CreatorId      string      `json:"_creatorId,omitempty"`
	EnvironmentId  string      `json:"_environmentId,omitempty"`
	OrganizationId string      `json:"_organizationId,omitempty"`
	EntityId       string      `json:"_entityId,omitempty"`
	Enabled        bool        `json:"enabled,omitempty"`
	Type           string      `json:"type,omitempty"`
	Change         interface{} `json:"change,omitempty"`
	CreatedAt      string      `json:"createdAt,omitempty"`
	ParentId       string      `json:"_parentId,omitempty"`
}

type ChangesGetResponse struct {
	TotalCount int                      `json:"totalCount,omitempty"`
	Data       []ChangesGetResponseData `json:"data"`
	PageSize   int                      `json:"pageSize,omitempty"`
	Page       int                      `json:"page,omitempty"`
}

type ChangesCountResponse struct {
	Data int `json:"data"`
}

type ChangesBulkApplyPayload struct {
	ChangeIds []string `json:"changeIds"`
}

type ChangesApplyResponse struct {
	Data []ChangesGetResponseData `json:"data,omitempty"`
}


type UpdateTenantRequest struct {
	Name 	 string `json:"name"`
	Data 	 map[string]interface{} `json:"data"`
	Identifier string `json:"identifier"`
}
