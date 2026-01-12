CREATE TABLE checks (
    check_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    location geometry(POINT,4326) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE INDEX idx_incidents_location ON incidents USING GIST(location);