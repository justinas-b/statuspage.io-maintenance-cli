package client

import (
	"net/http"
	"time"
)

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
	Pages      []Page
}

type Page struct {
	ID                       string    `json:"id,omitempty"`
	CreatedAt                time.Time `json:"created_at,omitempty"`
	UpdatedAt                time.Time `json:"updated_at,omitempty"`
	Name                     string    `json:"name,omitempty"`
	PageDescription          string    `json:"page_description,omitempty"`
	Headline                 string    `json:"headline,omitempty"`
	Branding                 string    `json:"branding,omitempty"`
	Subdomain                string    `json:"subdomain,omitempty"`
	Domain                   string    `json:"domain,omitempty"`
	URL                      string    `json:"url,omitempty"`
	SupportURL               string    `json:"support_url,omitempty"`
	HiddenFromSearch         bool      `json:"hidden_from_search,omitempty"`
	AllowPageSubscribers     bool      `json:"allow_page_subscribers,omitempty"`
	AllowIncidentSubscribers bool      `json:"allow_incident_subscribers,omitempty"`
	AllowEmailSubscribers    bool      `json:"allow_email_subscribers,omitempty"`
	AllowSMSSubscribers      bool      `json:"allow_sms_subscribers,omitempty"`
	AllowRSSAtomFeeds        bool      `json:"allow_rss_atom_feeds,omitempty"`
	AllowWebhookSubscribers  bool      `json:"allow_webhook_subscribers,omitempty"`
	NotificationsFromEmail   string    `json:"notifications_from_email,omitempty"`
	NotificationsEmailFooter string    `json:"notifications_email_footer,omitempty"`
	ActivityScore            int       `json:"activity_score,omitempty"`
	TwitterUsername          string    `json:"twitter_username,omitempty"`
	ViewersMustBeTeamMembers bool      `json:"viewers_must_be_team_members,omitempty"`
	IPRestrictions           string    `json:"ip_restrictions,omitempty"`
	City                     string    `json:"city,omitempty"`
	State                    string    `json:"state,omitempty"`
	Country                  string    `json:"country,omitempty"`
	TimeZone                 string    `json:"time_zone,omitempty"`
	CSSBodyBackgroundColor   string    `json:"css_body_background_color,omitempty"`
	CSSFontColor             string    `json:"css_font_color,omitempty"`
	CSSLightFontColor        string    `json:"css_light_font_color,omitempty"`
	CSSGreens                string    `json:"css_greens,omitempty"`
	CSSYellows               string    `json:"css_yellows,omitempty"`
	CSSOranges               string    `json:"css_oranges,omitempty"`
	CSSBlues                 string    `json:"css_blues,omitempty"`
	CSSReds                  string    `json:"css_reds,omitempty"`
	CSSBorderColor           string    `json:"css_border_color,omitempty"`
	CSSGraphColor            string    `json:"css_graph_color,omitempty"`
	CSSLinkColor             string    `json:"css_link_color,omitempty"`
	CSSNoData                string    `json:"css_no_data,omitempty"`
	FaviconLogo              Logo      `json:"favicon_logo,omitempty"`
	TransactionalLogo        Logo      `json:"transactional_logo,omitempty"`
	HeroCover                Logo      `json:"hero_cover,omitempty"`
	EmailLogo                Logo      `json:"email_logo,omitempty"`
	TwitterLogo              Logo      `json:"twitter_logo,omitempty"`
	client                   *Client
}

type Logo struct {
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Size        int       `json:"size,omitempty"`
	URL         string    `json:"url,omitempty"`
	OriginalURL string    `json:"original_url,omitempty"`
	RetinaURL   string    `json:"retina_url,omitempty"`
	NormalURL   string    `json:"normal_url,omitempty"`
}

type Incident struct {
	Name                                      string   `json:"name"`
	Status                                    string   `json:"status,omitempty"`
	ImpactOverride                            string   `json:"impact_override,omitempty"`
	ScheduledFor                              string   `json:"scheduled_for,omitempty"`
	ScheduledUntil                            string   `json:"scheduled_until,omitempty"`
	ScheduledRemindPrior                      bool     `json:"scheduled_remind_prior,omitempty"`
	AutoTransitionToMaintenanceState          bool     `json:"auto_transition_to_maintenance_state,omitempty"`
	AutoTransitionToOperationalState          bool     `json:"auto_transition_to_operational_state,omitempty"`
	ScheduledAutoInProgress                   bool     `json:"scheduled_auto_in_progress,omitempty"`
	ScheduledAutoCompleted                    bool     `json:"scheduled_auto_completed,omitempty"`
	AutoTransitionDeliverNotificationsAtStart bool     `json:"auto_transition_deliver_notifications_at_start,omitempty"`
	AutoTransitionDeliverNotificationsAtEnd   bool     `json:"auto_transition_deliver_notifications_at_end,omitempty"`
	ReminderIntervals                         string   `json:"reminder_intervals,omitempty"`
	DeliverNotifications                      bool     `json:"deliver_notifications,omitempty"`
	AutoTweetAtBeginning                      bool     `json:"auto_tweet_at_beginning,omitempty"`
	AutoTweetOnCompletion                     bool     `json:"auto_tweet_on_completion,omitempty"`
	AutoTweetOnCreation                       bool     `json:"auto_tweet_on_creation,omitempty"`
	AutoTweetOneHourBefore                    bool     `json:"auto_tweet_one_hour_before,omitempty"`
	BackfillDate                              string   `json:"backfill_date,omitempty"`
	Backfilled                                bool     `json:"backfilled,omitempty"`
	Body                                      string   `json:"body,omitempty"`
	Metadata                                  struct{} `json:"metadata,omitempty"`
	Components                                struct{} `json:"components,omitempty"`
	ComponentIds                              []string `json:"component_ids,omitempty"`
	ScheduledAutoTransition                   bool     `json:"scheduled_auto_transition,omitempty"`
}
