package geminimodel

type (
	Gemini interface {
		Chat(userInput string, h []HistoryContext) (string, []HistoryContext, error)
	}

	HistoryContext struct {
		UserInput *string `json:"user_input"`
		Response  *string `json:"response"`
	}
)
