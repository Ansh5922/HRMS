CREATE TABLE review_cycles (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id     UUID REFERENCES organizations(id),
    name       VARCHAR(255) NOT NULL,
    type       VARCHAR(50) CHECK (type IN ('annual','semi_annual','quarterly')),
    year       INT NOT NULL,
    quarter    INT CHECK (quarter BETWEEN 1 AND 4),
    start_date DATE NOT NULL,
    end_date   DATE NOT NULL,
    status     VARCHAR(50) DEFAULT 'upcoming' CHECK (status IN ('upcoming','active','completed','archived')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE goals (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id        UUID REFERENCES employees(id) ON DELETE CASCADE,
    cycle_id      UUID REFERENCES review_cycles(id),
    title         VARCHAR(255) NOT NULL,
    description   TEXT,
    metric        VARCHAR(255),
    target_value  DECIMAL(15,2),
    current_value DECIMAL(15,2) DEFAULT 0,
    weight        DECIMAL(5,2) DEFAULT 100,
    type          VARCHAR(20) CHECK (type IN ('kpi','okr','learning')),
    status        VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active','completed','missed','cancelled')),
    due_date      DATE,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE competencies (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id         UUID REFERENCES organizations(id),
    name           VARCHAR(255) NOT NULL,
    description    TEXT,
    level_criteria JSONB
);

CREATE TABLE performance_reviews (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id         UUID REFERENCES employees(id) ON DELETE CASCADE,
    cycle_id       UUID REFERENCES review_cycles(id),
    reviewer_id    UUID REFERENCES users(id),
    type           VARCHAR(30) CHECK (type IN ('self','peer','manager','subordinate')),
    overall_rating DECIMAL(3,2),
    feedback       TEXT,
    goal_ratings   JSONB,
    comp_ratings   JSONB,
    submitted_at   TIMESTAMPTZ,
    status         VARCHAR(30) DEFAULT 'pending' CHECK (status IN ('pending','submitted','acknowledged')),
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (emp_id, cycle_id, reviewer_id, type)
);

CREATE TABLE pip_plans (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emp_id     UUID REFERENCES employees(id) ON DELETE CASCADE,
    manager_id UUID REFERENCES users(id),
    start_date DATE NOT NULL,
    end_date   DATE NOT NULL,
    reason     TEXT,
    goals      JSONB,
    status     VARCHAR(50) CHECK (status IN ('active','completed','extended')),
    outcome    VARCHAR(50) CHECK (outcome IN ('improved','terminated')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
