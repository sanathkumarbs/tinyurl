CREATE TABLE IF NOT EXISTS tinyurls 
(
    original    varchar         NOT NULL,
    hash        varchar         PRIMARY KEY,
    expiry      date            DEFAULT now() + interval '365 days',
    -- CHECK constraint at the table level
    CONSTRAINT non_empty_original CHECK (original != '')
);
