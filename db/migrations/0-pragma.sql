-- +migrate Up notransaction
-- enable foreign key support
pragma foreign_keys = on;
-- set jornal mode to write ahead log
pragma  journal_mode = WAL;

-- +migrate Down notransaction
pragma foreign_keys = off;

pragma journal_mode = DELETE;