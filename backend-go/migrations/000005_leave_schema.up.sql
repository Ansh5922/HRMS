CREATE TABLE leave_types (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id            UUID REFERENCES organizations(id) ON DELETE CASCADE,
    name              VARCHAR(100) NOT NULL,
    code              VARCHAR(20) NOT NULL,
    max_days_per_year INT,
    carry_forward     BOOLEAN DEFAULT FALSE,
    carry_forward_max INT,
    is_paid           BOOLEAN DEFAULT TRUE,
    requires_doc      BOOLEAN DEFAULT FALSE,
    min_notice_days   INT DEFAULT 0,
    gender_applicable VARCHAR(20) CHECK (gender_applicable IN ('all','male','female','other')),
    requires_approval BOOLEAN DEFAULT TRUE,
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (org_id, code)
);

CREATE TABLE leave_policies (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id         UUID REFERENCES organizations(id),
    leave_type_id  UUID REFERENCES leave_types(id),
    designation_id UUID REFERENCES designations(id),
    department_id  UUID REFERENCES departments(id),
    annual_quota   DECIMAL(5,2) NOT NULL,
    accrual_type   VARCHAR(50) CHECK (accrual_type IN ('upfront','monthly','quarterly')),
    created_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE leave_balances (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id        UUID REFERENCES employees(id) ON DELETE CASCADE,
    leave_type_id UUID REFERENCES leave_types(id),
    year          INT NOT NULL,
    total         DECIMAL(5,2) NOT NULL,
    used          DECIMAL(5,2) DEFAULT 0,
    pending       DECIMAL(5,2) DEFAULT 0,
    available     DECIMAL(5,2) GENERATED ALWAYS AS (total - used - pending) STORED,
    updated_at    TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (emp_id, leave_type_id, year)
);

CREATE TABLE leave_applications (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id        UUID REFERENCES employees(id) ON DELETE CASCADE,
    leave_type_id UUID REFERENCES leave_types(id),
    from_date     DATE NOT NULL,
    to_date       DATE NOT NULL,
    days          DECIMAL(5,2) NOT NULL,
    session       VARCHAR(20) CHECK (session IN ('full','first_half','second_half')),
    reason        TEXT,
    doc_url       TEXT,
    status        VARCHAR(50) DEFAULT 'pending' CHECK (status IN
                    ('pending','approved','rejected','cancelled','withdrawn')),
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE leave_approvals (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    application_id UUID REFERENCES leave_applications(id) ON DELETE CASCADE,
    approver_id    UUID REFERENCES users(id),
    level          INT NOT NULL,
    action         VARCHAR(20) CHECK (action IN ('approved','rejected')),
    comment        TEXT,
    actioned_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_leave_apps_emp    ON leave_applications(emp_id);
CREATE INDEX idx_leave_apps_status ON leave_applications(status);
CREATE INDEX idx_leave_apps_dates  ON leave_applications(from_date, to_date);
