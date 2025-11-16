-- tables
CREATE TABLE teams (
    team_name TEXT PRIMARY KEY
);

CREATE TABLE users (
    user_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    team_name TEXT REFERENCES teams(team_name) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE pull_requests (
    pull_request_id TEXT PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    need_more_reviewers BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NULL,
    merged_at TIMESTAMPTZ NULL
);

CREATE TABLE pull_request_reviewers (
    pull_request_id TEXT NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    reviewer_id TEXT NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    PRIMARY KEY (pull_request_id, reviewer_id)
);

-- triggers
CREATE OR REPLACE FUNCTION limit_reviewers() RETURNS trigger AS $$
BEGIN
    IF (SELECT COUNT(*) FROM pull_request_reviewers
        WHERE pull_request_id = NEW.pull_request_id) > 2 THEN
    RAISE EXCEPTION 'Cannot assign more than 2 reviewers for a pull request';
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_limit_reviewers
BEFORE INSERT OR UPDATE ON pull_request_reviewers
FOR EACH ROW
EXECUTE FUNCTION limit_reviewers();

-- indexes
CREATE INDEX idx_pr_author_id ON pull_requests (author_id);
CREATE INDEX idx_pr_status ON pull_requests (status);
CREATE INDEX idx_rr_pr ON pull_request_reviewers (pull_request_id);
CREATE INDEX idx_rr_reviewer ON pull_request_reviewers (reviewer_id);
