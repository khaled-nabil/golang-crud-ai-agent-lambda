package datamodels

type (
	Chat struct {
		ID       *string `json:"id,omitempty"`
		UserID   string  `json:"user_id"`
		Message  string  `json:"message"`
		Response string  `json:"response"`
		CreateAt *int64  `json:"created_at,omitempty"`
	}
)
