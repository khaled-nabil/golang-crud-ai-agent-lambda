package datamodels

type (
	Gemini interface {
		Chat(userInput string, h []HistoryContext) (*HistoryContext, error)
		EmbedMessage(message, response string) ([]float32, error)
	}

	HistoryContext struct {
		UserInput string `json:"user"`
		Response  string `json:"response"`
	}

	AppConfig struct {
		GeminiAPIKey string `json:"GEMINI_API_KEY"`
		ModelID      string `json:"MODEL_ID"`
		DBHost       string `json:"DB_HOST"`
		DBPort       string `json:"DB_PORT"`
		DBName       string `json:"DB_NAME"`
		DBUser       string `json:"DB_USER"`
		DBPassword   string `json:"DB_PASSWORD"`
	}
)

func (c *Chat) ChatToHistoryContext() *HistoryContext {
	return &HistoryContext{
		UserInput: c.Message,
		Response:  c.Response,
	}
}

func ChatListToHistoryContextList(c []Chat) []HistoryContext {
	var history []HistoryContext

	for _, chat := range c {
		history = append(history, *chat.ChatToHistoryContext())
	}

	return history
}
