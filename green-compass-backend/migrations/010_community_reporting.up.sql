ALTER TABLE observations
    ADD COLUMN organization_id UUID REFERENCES organizations(id),
    ADD COLUMN verification_notes TEXT;
