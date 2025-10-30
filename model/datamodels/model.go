package datamodels

type (
	Gemini interface {
		Chat(userInput string, h []HistoryContext) (*HistoryContext, error)
	}

	DynamoDB interface {
		StoreItem(item interface{}) error
		RetrieveItems(id string, limit int32) ([]HistoryContext, error)
	}

	HistoryContext struct {
		UserInput string `json:"user"`
		Response  string `json:"response"`
	}
)
