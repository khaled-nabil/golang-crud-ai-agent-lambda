package entity

type (
	ChatHistoryEntity struct {
		UserInput string `json:"user"`
		Response  string `json:"response"`
	}
)

