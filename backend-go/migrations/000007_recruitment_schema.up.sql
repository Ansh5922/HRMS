CREATE TABLE job_postings (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id         UUID REFERENCES organizations(id),
    dept_id        UUID REFERENCES departments(id),
    designation_id UUID REFERENCES designations(id),
    title          VARCHAR(255) NOT NULL,
    description    TEXT,
    requirements   TEXT,
    type           VARCHAR(50) CHECK (type IN ('full_time','part_time','contract','internship')),
    openings       INT DEFAULT 1,
    min_experience DECIMAL(4,1),
    max_experience DECIMAL(4,1),
    min_salary     DECIMAL(15,2),
    max_salary     DECIMAL(15,2),
    location       VARCHAR(255),
    is_remote      BOOLEAN DEFAULT FALSE,
    status         VARCHAR(50) DEFAULT 'open' CHECK (status IN ('draft','open','paused','closed','cancelled')),
    posted_by      UUID REFERENCES users(id),
    posted_at      TIMESTAMPTZ,
    closed_at      TIMESTAMPTZ,
    created_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE candidates (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id       UUID REFERENCES organizations(id),
    first_name   VARCHAR(100) NOT NULL,
    last_name    VARCHAR(100),
    email        VARCHAR(255) NOT NULL,
    phone        VARCHAR(20),
    resume_url   TEXT,
    linkedin_url TEXT,
    source       VARCHAR(100),
    skills       JSONB,
    parsed_data  JSONB,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (org_id, email)
);

CREATE TABLE applications (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id            UUID REFERENCES job_postings(id),
    candidate_id      UUID REFERENCES candidates(id),
    status            VARCHAR(50) DEFAULT 'applied' CHECK (status IN
                        ('applied','screening','phone_screen','interview',
                         'technical','hr_round','offered','hired','rejected','withdrawn')),
    current_ctc       DECIMAL(15,2),
    expected_ctc      DECIMAL(15,2),
    notice_period_days INT,
    referral_emp_id   UUID REFERENCES employees(id),
    applied_at        TIMESTAMPTZ DEFAULT NOW(),
    updated_at        TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (job_id, candidate_id)
);

CREATE TABLE interviews (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    application_id UUID REFERENCES applications(id) ON DELETE CASCADE,
    round_no       INT NOT NULL,
    round_name     VARCHAR(100),
    interviewer_id UUID REFERENCES users(id),
    scheduled_at   TIMESTAMPTZ,
    duration_mins  INT DEFAULT 60,
    mode           VARCHAR(20) CHECK (mode IN ('online','offline','phone')),
    meet_link      TEXT,
    feedback       TEXT,
    rating         INT CHECK (rating BETWEEN 1 AND 5),
    result         VARCHAR(20) CHECK (result IN ('selected','rejected','on_hold','no_show')),
    completed_at   TIMESTAMPTZ,
    created_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE job_offers (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    application_id UUID REFERENCES applications(id),
    ctc_offered    DECIMAL(15,2) NOT NULL,
    joining_date   DATE,
    offer_url      TEXT,
    valid_till     DATE,
    status         VARCHAR(50) DEFAULT 'sent' CHECK (status IN
                     ('sent','accepted','declined','expired','revoked')),
    issued_by      UUID REFERENCES users(id),
    responded_at   TIMESTAMPTZ,
    created_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_applications_job    ON applications(job_id);
CREATE INDEX idx_applications_status ON applications(status);
CREATE INDEX idx_interviews_app      ON interviews(application_id);
