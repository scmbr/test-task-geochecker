CREATE TABLE incidents (
    incident_id UUID PRIMARY KEY,
    operator_id UUID NOT NULL,
    location geometry(POINT,4326) NOT NULL,
    radius INTEGER NOT NULL DEFAULT 5,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE INDEX idx_incidents_location ON incidents USING GIST(location);