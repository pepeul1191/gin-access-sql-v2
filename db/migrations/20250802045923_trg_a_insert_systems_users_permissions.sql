-- migrate:up

CREATE TRIGGER trg_a_insert_systems_users_permissions
AFTER INSERT ON systems_users_permissions
BEGIN
  INSERT INTO systems_users_roles (system_id, user_id, role_id, created)
  SELECT
    NEW.system_id,
    NEW.user_id,
    p.role_id,
    NEW.created
  FROM permissions p
  WHERE p.id = NEW.permission_id
    AND NOT EXISTS (
      SELECT 1
      FROM systems_users_roles sur
      WHERE sur.system_id = NEW.system_id
        AND sur.user_id = NEW.user_id
        AND sur.role_id = p.role_id
    );
END;

-- migrate:down

DROP TRIGGER IF EXISTS trg_a_insert_systems_users_permissions;
