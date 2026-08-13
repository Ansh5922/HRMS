CREATE TABLE shifts (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id       UUID REFERENCES organizations(id) ON DELETE CASCADE,
    name         VARCHAR(100) NOT NULL,
    start_time   TIME NOT NULL,
    end_time     TIME NOT NULL,
    grace_mins   INT DEFAULT 15,
    is_overnight BOOLEAN DEFAULT FALSE,
    created_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE shift_assignments (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id         UUID REFERENCES employees(id) ON DELETE CASCADE,
    shift_id       UUID REFERENCES shifts(id),
    effective_from DATE NOT NULL,
    effective_to   DATE,
    created_by     UUID REFERENCES users(id)
);

CREATE TABLE holidays (
    id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id  UUID REFERENCES organizations(id) ON DELETE CASCADE,
    name    VARCHAR(255) NOT NULL,
    date    DATE NOT NULL,
    type    VARCHAR(50) CHECK (type IN ('national','optional','restricted')),
    UNIQUE (org_id, date)
);

CREATE TABLE geofences (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id         UUID REFERENCES organizations(id),
    name           VARCHAR(255) NOT NULL,
    latitude       DECIMAL(10,8) NOT NULL,
    longitude      DECIMAL(11,8) NOT NULL,
    radius_meters  INT DEFAULT 100,
    is_active      BOOLEAN DEFAULT TRUE
);

CREATE TABLE attendance_records (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id          UUID REFERENCES organizations(id),
    emp_id          UUID REFERENCES employees(id) ON DELETE CASCADE,
    date            DATE NOT NULL,
    check_in        TIMESTAMPTZ,
    check_out       TIMESTAMPTZ,
    method          VARCHAR(50) CHECK (method IN ('face','manual','geo','biometric')),
    status          VARCHAR(50) CHECK (status IN ('present','absent','half_day','late','wfh','on_leave')),
    duration_mins   INT,
    overtime_mins   INT DEFAULT 0,
    location_lat    DECIMAL(10,8),
    location_lng    DECIMAL(11,8),
    geofence_id     UUID REFERENCES geofences(id),
    face_conf_score DECIMAL(5,4),
    is_regularized  BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (emp_id, date)
);

CREATE TABLE attendance_regularizations (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id      UUID REFERENCES employees(id) ON DELETE CASCADE,
    date        DATE NOT NULL,
    check_in    TIMESTAMPTZ,
    check_out   TIMESTAMPTZ,
    reason      TEXT NOT NULL,
    status      VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending','approved','rejected')),
    reviewer_id UUID REFERENCES users(id),
    review_note TEXT,
    reviewed_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE face_embeddings (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id        UUID REFERENCES employees(id) ON DELETE CASCADE UNIQUE,
    embedding     vector(512) NOT NULL,
    model_version VARCHAR(100),
    enrolled_by   UUID REFERENCES users(id),
    enrolled_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_face_embeddings_hnsw
    ON face_embeddings
    USING hnsw (embedding vector_cosine_ops)
    WITH (m = 16, ef_construction = 64);

CREATE INDEX idx_attendance_emp_date ON attendance_records(emp_id, date DESC);
CREATE INDEX idx_attendance_org_date ON attendance_records(org_id, date DESC);
