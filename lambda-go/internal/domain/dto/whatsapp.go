package dto

type WhatsAppEvent struct {
	Object string          `json:"object"`
	Entry  []WhatsAppEntry `json:"entry"`
}

type WhatsAppEntry struct {
	ID      string           `json:"id"`
	Changes []WhatsAppChange `json:"changes"`
}

type WhatsAppChange struct {
	Value WhatsAppValue `json:"value"`
	Field string        `json:"field"`
}

type WhatsAppValue struct {
	MessagingProduct string            `json:"messaging_product"`
	Metadata         WhatsAppMetadata  `json:"metadata"`
	Contacts         []WhatsAppContact `json:"contacts"`
	Messages         []WhatsAppMessage `json:"messages"`
	Statuses         []WhatsAppStatus  `json:"statuses"`
}

type WhatsAppMetadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type WhatsAppContact struct {
	Profile WhatsAppProfile `json:"profile"`
	WaID    string          `json:"wa_id"`
}

type WhatsAppProfile struct {
	Name string `json:"name"`
}

type WhatsAppMessage struct {
	From      string              `json:"from"`
	ID        string              `json:"id"`
	Timestamp string              `json:"timestamp"`
	Text      WhatsAppMessageText `json:"text"`
	Button    WhatsAppButton      `json:"button"`
	Type      string              `json:"type"`
}

type WhatsAppMessageText struct {
	Body string `json:"body"`
}

type WhatsAppButton struct {
	Payload string `json:"payload"`
	Text    string `json:"text"`
}

type WhatsAppStatus struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	Timestamp    string `json:"timestamp"`
	RecipientID  string `json:"recipient_id"`
	Conversation struct {
		ID     string `json:"id"`
		Origin struct {
			Type string `json:"type"`
		} `json:"origin"`
	} `json:"conversation"`
	Pricing struct {
		Billable     bool   `json:"billable"`
		PricingModel string `json:"pricing_model"`
		Category     string `json:"category"`
	} `json:"pricing"`
}
