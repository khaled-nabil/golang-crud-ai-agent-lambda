CREATE TABLE domain_book (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    subtitle TEXT,
    description TEXT,
    embedding vector(768),
    thumbnail TEXT,
    published_year SMALLINT NOT NULL,
    average_rating FLOAT NOT NULL,
    rating_count INTEGER NOT NULL,
    num_pages SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE domain_book_category (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL
);

CREATE TABLE domain_book_author (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL
);

CREATE TABLE domain_book_author_map (
    book_id UUID REFERENCES domain_book(id),
    author_id UUID REFERENCES domain_book_author(id),
    PRIMARY KEY (book_id, author_id)
);

CREATE TABLE domain_book_category_map (
    book_id UUID REFERENCES domain_book(id),
    category_id UUID REFERENCES domain_book_category(id),
    PRIMARY KEY (book_id, category_id)
);

CREATE INDEX IF NOT EXISTS idx_domain_id ON domain (id);
CREATE INDEX idx_book_vector ON domain_book USING hnsw (embedding vector_cosine_ops);