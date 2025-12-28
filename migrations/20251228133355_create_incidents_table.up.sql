CREATE TABLE incidents (
    incident_id UUID PRIMARY KEY,
    operator_id UUID NOT NULL,
    latitude NUMERIC(8,6) NOT NULL,
    longitude NUMERIC(9,6) NOT NULL,
    radius INTEGER NOT NULL DEFAULT 5,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);