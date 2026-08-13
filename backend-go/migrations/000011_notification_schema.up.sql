CREATE TABLE notifications (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id     UUID REFERENCES organizations(id),
    user_id    UUID REFERENCES users(id) ON DELETE CASCADE,
    type       VARCHAR(100) NOT NULL,
    title      VARCHAR(255) NOT NULL,
    message    TEXT,
    action_url TEXT,
    is_read    BOOLEAN DEFAULT FALSE,
    read_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE announcements (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id       UUID REFERENCES organizations(id),
    title        VARCHAR(255) NOT NULL,
    content      TEXT NOT NULL,
    author_id    UUID REFERENCES users(id),
    target_dept  UUID REFERENCES departments(id),
    published_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at   TIMESTAMPTZ
);

CREATE TABLE notification_preferences (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID REFERENCES users(id) ON DELETE CASCADE,
    event_type VARCHAR(100) NOT NULL,
    email      BOOLEAN DEFAULT TRUE,
    push       BOOLEAN DEFAULT TRUE,
    sms        BOOLEAN DEFAULT FALSE,
    UNIQUE (user_id, event_type)
);

CREATE TABLE push_tokens (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID REFERENCES users(id) ON DELETE CASCADE,
    token      TEXT NOT NULL,
    platform   VARCHAR(20) CHECK (platform IN ('android','ios','web')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (user_id, token)
);

CREATE INDEX idx_notifications_user ON notifications(user_id, is_read, created_at DESC);
CREATE INDEX idx_notifications_org  ON notifications(org_id, created_at DESC);
