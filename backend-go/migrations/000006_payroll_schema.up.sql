CREATE TABLE salary_structures (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id      UUID REFERENCES organizations(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    is_active   BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE salary_components (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    structure_id  UUID REFERENCES salary_structures(id) ON DELETE CASCADE,
    name          VARCHAR(100) NOT NULL,
    code          VARCHAR(50) NOT NULL,
    type          VARCHAR(20) CHECK (type IN ('earning','deduction','stat_deduction')),
    calc_type     VARCHAR(50) CHECK (calc_type IN ('fixed','pct_of_basic','pct_of_gross','formula')),
    value         DECIMAL(10,4),
    formula       TEXT,
    is_taxable    BOOLEAN DEFAULT TRUE,
    display_order INT DEFAULT 0
);

CREATE TABLE employee_salaries (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id        UUID REFERENCES employees(id) ON DELETE CASCADE,
    structure_id  UUID REFERENCES salary_structures(id),
    ctc           DECIMAL(15,2) NOT NULL,
    basic         DECIMAL(15,2) NOT NULL,
    effective_from DATE NOT NULL,
    effective_to  DATE,
    approved_by   UUID REFERENCES users(id),
    approved_at   TIMESTAMPTZ,
    created_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE payroll_runs (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id           UUID REFERENCES organizations(id),
    month            INT NOT NULL CHECK (month BETWEEN 1 AND 12),
    year             INT NOT NULL,
    status           VARCHAR(50) DEFAULT 'draft' CHECK (status IN
                       ('draft','processing','pending_approval','approved','disbursed')),
    total_gross      DECIMAL(15,2),
    total_deductions DECIMAL(15,2),
    total_net        DECIMAL(15,2),
    run_by           UUID REFERENCES users(id),
    approved_by      UUID REFERENCES users(id),
    approved_at      TIMESTAMPTZ,
    disbursed_at     TIMESTAMPTZ,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (org_id, month, year)
);

CREATE TABLE payslips (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id           UUID REFERENCES employees(id) ON DELETE CASCADE,
    payroll_run_id   UUID REFERENCES payroll_runs(id),
    working_days     INT,
    lop_days         DECIMAL(5,2) DEFAULT 0,
    paid_days        DECIMAL(5,2),
    gross            DECIMAL(15,2) NOT NULL,
    total_deductions DECIMAL(15,2) NOT NULL,
    net_pay          DECIMAL(15,2) NOT NULL,
    breakdown        JSONB,
    payslip_url      TEXT,
    generated_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE reimbursements (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id         UUID REFERENCES employees(id) ON DELETE CASCADE,
    category       VARCHAR(100) NOT NULL,
    amount         DECIMAL(10,2) NOT NULL,
    description    TEXT,
    receipt_url    TEXT,
    status         VARCHAR(50) DEFAULT 'pending' CHECK (status IN
                     ('pending','approved','rejected','paid')),
    approved_by    UUID REFERENCES users(id),
    payroll_run_id UUID REFERENCES payroll_runs(id),
    created_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE tax_declarations (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id       UUID REFERENCES employees(id),
    year         INT NOT NULL,
    regime       VARCHAR(20) CHECK (regime IN ('old','new')),
    sections     JSONB,
    submitted_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (emp_id, year)
);

CREATE INDEX idx_payslips_emp          ON payslips(emp_id);
CREATE INDEX idx_payslips_run          ON payslips(payroll_run_id);
CREATE INDEX idx_employee_salaries_emp ON employee_salaries(emp_id);
