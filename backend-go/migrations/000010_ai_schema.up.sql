CREATE TABLE kb_documents (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id      UUID REFERENCES organizations(id),
    title       VARCHAR(255) NOT NULL,
    file_url    TEXT NOT NULL,
    file_type   VARCHAR(50),
    category    VARCHAR(100),
    chunk_count INT DEFAULT 0,
    ingested_at TIMESTAMPTZ,
    ingested_by UUID REFERENCES users(id),
    is_active   BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE kb_chunks (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    doc_id      UUID REFERENCES kb_documents(id) ON DELETE CASCADE,
    org_id      UUID REFERENCES organizations(id),
    content     TEXT NOT NULL,
    embedding   vector(1536) NOT NULL,
    chunk_index INT,
    token_count INT,
    metadata    JSONB,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_kb_chunks_hnsw
    ON kb_chunks
    USING hnsw (embedding vector_cosine_ops)
    WITH (m = 16, ef_construction = 64);

CREATE INDEX idx_kb_chunks_org ON kb_chunks(org_id);
CREATE INDEX idx_kb_chunks_doc ON kb_chunks(doc_id);

CREATE TABLE chat_sessions (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID REFERENCES users(id) ON DELETE CASCADE,
    org_id     UUID REFERENCES organizations(id),
    title      VARCHAR(255),
    started_at TIMESTAMPTZ DEFAULT NOW(),
    ended_at   TIMESTAMPTZ
);

CREATE TABLE chat_messages (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id       UUID REFERENCES chat_sessions(id) ON DELETE CASCADE,
    role             VARCHAR(20) CHECK (role IN ('user','assistant','system')),
    content          TEXT NOT NULL,
    retrieved_chunks JSONB,
    model_used       VARCHAR(100),
    latency_ms       INT,
    created_at       TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_chat_messages_session ON chat_messages(session_id, created_at);
