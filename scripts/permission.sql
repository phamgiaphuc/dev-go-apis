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
    ('dashboard:user:read', 'Read or access user list'),
    ('dashboard:user:create', 'Create a user'),
    ('dashboard:user:edit', 'Edit a user'),
    ('dashboard:user:delete', 'Delete a user'),
    ('dashboard:role:read', 'Read or access role list'),
    ('dashboard:role:create', 'Create a role'),
    ('dashboard:role:edit', 'Edit a role'),
    ('dashboard:role:delete', 'Delete a role')
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
  'dashboard:user:delete',
  'dashboard:role:read',
  'dashboard:role:create',
  'dashboard:role:edit',
  'dashboard:role:delete'
)
WHERE r.name = 'admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

# Delete queries
DELETE * FROM "role_permissions";
DELETE * FROM "roles";
DELETE * FROM "permissions";
DELETE * FROM "permission_groups";