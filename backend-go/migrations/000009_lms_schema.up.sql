CREATE TABLE courses (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id       UUID REFERENCES organizations(id),
    title        VARCHAR(255) NOT NULL,
    description  TEXT,
    category     VARCHAR(100),
    type         VARCHAR(50) CHECK (type IN ('video','pdf','scorm','classroom')),
    content_url  TEXT,
    duration_hrs DECIMAL(5,2),
    is_mandatory BOOLEAN DEFAULT FALSE,
    created_by   UUID REFERENCES users(id),
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE course_assignments (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id    UUID REFERENCES courses(id) ON DELETE CASCADE,
    emp_id       UUID REFERENCES employees(id) ON DELETE CASCADE,
    assigned_by  UUID REFERENCES users(id),
    due_date     DATE,
    status       VARCHAR(50) DEFAULT 'assigned' CHECK (status IN
                   ('assigned','in_progress','completed','overdue')),
    progress_pct INT DEFAULT 0 CHECK (progress_pct BETWEEN 0 AND 100),
    score        DECIMAL(5,2),
    completed_at TIMESTAMPTZ,
    cert_url     TEXT,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (course_id, emp_id)
);

CREATE TABLE skills (
    id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id   UUID REFERENCES organizations(id),
    name     VARCHAR(100) NOT NULL,
    category VARCHAR(100),
    UNIQUE (org_id, name)
);

CREATE TABLE employee_skills (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id            UUID REFERENCES employees(id) ON DELETE CASCADE,
    skill_id          UUID REFERENCES skills(id),
    proficiency_level INT CHECK (proficiency_level BETWEEN 1 AND 5),
    assessed_by       UUID REFERENCES users(id),
    assessed_at       TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (emp_id, skill_id)
);
