package datamodels

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type (
	Chat struct {
		ID       *pgtype.UUID `json:"id,omitempty"`
		UserID   pgtype.UUID  `json:"user_id"`
		Message  string       `json:"message"`
		Response string       `json:"response"`
		CreateAt *time.Time   `json:"created_at,omitempty"`
	}
)
