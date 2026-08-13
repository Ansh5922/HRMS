CREATE TABLE departments (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id         UUID REFERENCES organizations(id) ON DELETE CASCADE,
    name           VARCHAR(255) NOT NULL,
    code           VARCHAR(50),
    parent_dept_id UUID REFERENCES departments(id),
    head_emp_id    UUID,
    description    TEXT,
    is_active      BOOLEAN DEFAULT TRUE,
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE designations (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id        UUID REFERENCES organizations(id) ON DELETE CASCADE,
    department_id UUID REFERENCES departments(id),
    title         VARCHAR(255) NOT NULL,
    grade         VARCHAR(50),
    level         INT,
    created_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE employees (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id           UUID REFERENCES organizations(id) ON DELETE CASCADE,
    user_id          UUID REFERENCES users(id) UNIQUE,
    emp_code         VARCHAR(50) NOT NULL,
    first_name       VARCHAR(100) NOT NULL,
    last_name        VARCHAR(100) NOT NULL,
    middle_name      VARCHAR(100),
    dob              DATE,
    gender           VARCHAR(20),
    blood_group      VARCHAR(10),
    phone            VARCHAR(20),
    personal_email   VARCHAR(255),
    address          TEXT,
    city             VARCHAR(100),
    state            VARCHAR(100),
    country          VARCHAR(100),
    pincode          VARCHAR(20),
    nationality      VARCHAR(100),
    aadhaar_no_enc   TEXT,
    pan_no_enc       TEXT,
    passport_no      VARCHAR(50),
    dept_id          UUID REFERENCES departments(id),
    designation_id   UUID REFERENCES designations(id),
    manager_id       UUID REFERENCES employees(id),
    employment_type  VARCHAR(50) CHECK (employment_type IN ('full_time','part_time','contract','intern')),
    work_location    VARCHAR(100),
    joining_date     DATE NOT NULL,
    confirmation_date DATE,
    status           VARCHAR(50) DEFAULT 'active' CHECK (status IN
                       ('active','probation','notice_period','resigned','terminated','absconding')),
    exit_date        DATE,
    exit_reason      TEXT,
    photo_url        TEXT,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (org_id, emp_code)
);

ALTER TABLE departments
    ADD CONSTRAINT fk_dept_head FOREIGN KEY (head_emp_id) REFERENCES employees(id) ON DELETE SET NULL;

CREATE TABLE employee_bank_details (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id          UUID REFERENCES employees(id) ON DELETE CASCADE,
    bank_name       VARCHAR(255) NOT NULL,
    account_no_enc  TEXT NOT NULL,
    ifsc_code       VARCHAR(20) NOT NULL,
    account_type    VARCHAR(50) CHECK (account_type IN ('savings','current')),
    is_primary      BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE emergency_contacts (
    id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id   UUID REFERENCES employees(id) ON DELETE CASCADE,
    name     VARCHAR(255) NOT NULL,
    relation VARCHAR(100),
    phone    VARCHAR(20) NOT NULL,
    email    VARCHAR(255),
    address  TEXT
);

CREATE TABLE employee_documents (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id      UUID REFERENCES employees(id) ON DELETE CASCADE,
    doc_type    VARCHAR(100) NOT NULL,
    file_url    TEXT NOT NULL,
    file_name   VARCHAR(255),
    is_verified BOOLEAN DEFAULT FALSE,
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMPTZ,
    uploaded_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE onboarding_tasks (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id      UUID REFERENCES employees(id) ON DELETE CASCADE,
    task_name   VARCHAR(255) NOT NULL,
    description TEXT,
    assigned_to UUID REFERENCES users(id),
    due_date    DATE,
    completed_at TIMESTAMPTZ,
    status      VARCHAR(50) DEFAULT 'pending'
);

CREATE TABLE custom_fields (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id      UUID REFERENCES organizations(id),
    module      VARCHAR(100) NOT NULL,
    field_name  VARCHAR(100) NOT NULL,
    field_type  VARCHAR(50) NOT NULL,
    options     JSONB,
    is_required BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE employee_custom_field_values (
    id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id   UUID REFERENCES employees(id) ON DELETE CASCADE,
    field_id UUID REFERENCES custom_fields(id) ON DELETE CASCADE,
    value    TEXT,
    UNIQUE (emp_id, field_id)
);

CREATE INDEX idx_employees_org     ON employees(org_id);
CREATE INDEX idx_employees_dept    ON employees(dept_id);
CREATE INDEX idx_employees_manager ON employees(manager_id);
CREATE INDEX idx_employees_status  ON employees(status);
