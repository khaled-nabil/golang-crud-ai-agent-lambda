CREATE TABLE IF NOT EXISTS chat (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    message TEXT NOT NULL,
    response TEXT NOT NULL,
    embedding vector (1536) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_user_id ON chat (user_id);

CREATE INDEX IF NOT EXISTS idx_embedding_hnsw ON chat USING hnsw (embedding vector_cosine_ops);