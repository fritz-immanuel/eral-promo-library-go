package models

type WhatsappResponseBitu struct {
	Code    int    `json:"Code" db:"-"`
	Status  bool   `json:"Status" db:"-"`
	Message string `json:"Message" db:"-"`
	Data    struct {
		MessageID string `json:"MessageID" db:"-"`
	}
}

type WARequest struct {
	token  string
	number string
	msg    string
	file   string
}

type WhatsappResponse struct {
	Result  string `json:"result" db:"-"`
	ID      string `json:"id" db:"-"`
	Number  string `json:"number" db:"-"`
	Message string `json:"message" db:"-"`
	Status  string `json:"status" db:"-"`
}
