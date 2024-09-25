BEGIN;

ALTER TABLE restaurants ADD COLUMN created_by UUID REFERENCES users(id);

COMMIT;
