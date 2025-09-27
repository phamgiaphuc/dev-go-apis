# Insert queries
INSERT INTO "permission_groups" (name, description)
VALUES 
  ('dashboard:group', 'Permissions related to admin dashboard')
ON CONFLICT (name) DO NOTHING;

INSERT INTO "permissions" (group_id, name, description)
SELECT pg.id, p.name, p.description
FROM "permission_groups" pg
JOIN (
  VALUES
    ('dashboard:admin', 'Access admin dashboard'),
    ('dashboard:user:read', 'Read user dashboard'),
    ('dashboard:user:create', 'Create user dashboard'),
    ('dashboard:user:edit', 'Edit user dashboard'),
    ('dashboard:user:delete', 'Delete user dashboard')
) AS p(name, description)
ON (pg.name = 'dashboard:group')
ON CONFLICT (name) DO NOTHING;

INSERT INTO "roles" (name, description)
VALUES 
  ('admin', 'Administrator role with full permissions'),
  ('user', 'User role with specific permissions')
ON CONFLICT (name) DO NOTHING;

INSERT INTO "role_permissions" (role_id, permission_id)
SELECT r.id, p.id
FROM "roles" r
JOIN "permissions" p ON p.name IN (
  'dashboard:admin',
  'dashboard:user:read',
  'dashboard:user:create',
  'dashboard:user:edit',
  'dashboard:user:delete'
)
WHERE r.name = 'admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

# Delete queries
DELETE * FROM "role_permissions";
DELETE * FROM "roles";
DELETE * FROM "permissions";
DELETE * FROM "permission_groups";