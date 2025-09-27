-- +goose Up
-- +goose StatementBegin
CREATE TABLE "permission_groups" (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT
);

CREATE TABLE "permissions" (
  id SERIAL PRIMARY KEY,
  group_id INT NOT NULL REFERENCES "permission_groups" (id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL UNIQUE,
  description TEXT
);

CREATE TABLE "roles" (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  description TEXT
);

CREATE TABLE "role_permissions" (
  role_id INT NOT NULL REFERENCES "roles" (id) ON DELETE CASCADE,
  permission_id INT NOT NULL REFERENCES "permissions" (id) ON DELETE CASCADE,
  PRIMARY KEY (role_id, permission_id)
);

INSERT INTO "roles" (name, description)
VALUES 
  ('admin', 'Administrator role with full permissions'),
  ('user', 'User role with specific permissions')
ON CONFLICT (name) DO NOTHING;

CREATE OR REPLACE FUNCTION prevent_protected_roles_delete()
RETURNS trigger AS $$
BEGIN
  IF OLD.name IN ('admin', 'user') THEN
    RAISE EXCEPTION 'Cannot delete protected role: %', OLD.name;
  END IF;
  RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION prevent_protected_roles_update()
RETURNS trigger AS $$
BEGIN
  IF OLD.name IN ('admin', 'user') THEN
    RAISE EXCEPTION 'Cannot update protected role: %', OLD.name;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers
CREATE TRIGGER trg_prevent_roles_delete
BEFORE DELETE ON roles
FOR EACH ROW
EXECUTE FUNCTION prevent_protected_roles_delete();

CREATE TRIGGER trg_prevent_roles_update
BEFORE UPDATE ON roles
FOR EACH ROW
EXECUTE FUNCTION prevent_protected_roles_update();

ALTER TABLE "users" 
ADD COLUMN IF NOT EXISTS role_id INT REFERENCES "roles" (id);

UPDATE "users"
SET role_id = (SELECT id FROM roles WHERE name = 'user')
WHERE role_id IS NULL;

DO $$
DECLARE
  user_role_id INT;
BEGIN
  SELECT id INTO user_role_id FROM roles WHERE name = 'user';
  EXECUTE format('ALTER TABLE "users" ALTER COLUMN role_id SET DEFAULT %s', user_role_id);
END$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "users" DROP COLUMN IF EXISTS role_id;
DROP TRIGGER IF EXISTS trg_prevent_roles_delete ON roles;
DROP TRIGGER IF EXISTS trg_prevent_roles_update ON roles;
DROP FUNCTION IF EXISTS prevent_protected_roles_delete();
DROP FUNCTION IF EXISTS prevent_protected_roles_update();
DROP TABLE IF EXISTS "role_permissions";
DROP TABLE IF EXISTS "permissions";
DROP TABLE IF EXISTS "roles";
DROP TABLE IF EXISTS "permission_groups";

-- Cleanup sequences to fully reset SERIALs
DROP SEQUENCE IF EXISTS "permission_groups_id_seq";
DROP SEQUENCE IF EXISTS "permissions_id_seq";
DROP SEQUENCE IF EXISTS "roles_id_seq";
-- +goose StatementEnd
