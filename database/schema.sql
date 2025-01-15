CREATE TABLE IF NOT EXISTS agents
(
    id                  INTEGER PRIMARY KEY,
    uid                 text      NOT NULL,
    name                text      NOT NULL,
    status              text      NOT NULL,
    category            text      NOT NULL,
    mcap      text      NOT NULL,

    notified            bool      NOT NULL
);
