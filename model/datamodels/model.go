package datamodels

type (
	Gemini interface {
		Chat(userInput string, h []HistoryContext) (string, []HistoryContext, error)
	}

	DynamoDB interface {
		StoreItem(item interface{}) error
	}

	HistoryContext struct {
		UserInput *string `json:"user_input"`
		Response  *string `json:"response"`
	}
)
