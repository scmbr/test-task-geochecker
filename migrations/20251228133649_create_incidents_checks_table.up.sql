CREATE TABLE incidents_checks(
    incident_id UUID NOT NULL,
    check_id UUID NOT NULL,
    PRIMARY KEY (incident_id, check_id),
    FOREIGN KEY (incident_id) REFERENCES incidents(incident_id) ON DELETE CASCADE,
    FOREIGN KEY (check_id) REFERENCES checks(check_id) ON DELETE CASCADE
);