CREATE TABLE workflow_templates (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id     UUID REFERENCES organizations(id),
    name       VARCHAR(255) NOT NULL,
    module     VARCHAR(100) NOT NULL,
    steps      JSONB NOT NULL,
    is_active  BOOLEAN DEFAULT TRUE,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE approval_requests (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id       UUID REFERENCES organizations(id),
    workflow_id  UUID REFERENCES workflow_templates(id),
    module       VARCHAR(100) NOT NULL,
    resource_id  UUID NOT NULL,
    requested_by UUID REFERENCES users(id),
    current_step INT DEFAULT 1,
    status       VARCHAR(50) DEFAULT 'pending' CHECK (status IN
                   ('pending','approved','rejected','cancelled')),
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE approval_steps (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    approval_request_id UUID REFERENCES approval_requests(id) ON DELETE CASCADE,
    step_no             INT NOT NULL,
    approver_id         UUID REFERENCES users(id),
    action              VARCHAR(20) CHECK (action IN ('approved','rejected')),
    comment             TEXT,
    actioned_at         TIMESTAMPTZ
);
