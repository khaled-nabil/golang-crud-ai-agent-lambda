CREATE TABLE documents_gemini (
    chat_id VARCHAR(255) PRIMARY KEY REFERENCES chat (id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL,
    embedding vector (1536) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT documents_gemini_user_chat_uniq UNIQUE (user_id, chat_id)
);

CREATE INDEX IF NOT EXISTS idx_documents_gemini_user_created ON documents_gemini (user_id, chat_id);

CREATE INDEX IF NOT EXISTS idx_documents_gemini_embedding_hnsw ON documents_gemini USING hnsw (embedding vector_cosine_ops);