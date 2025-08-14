-- migrate:up

CREATE TRIGGER trg_a_delete_systems_users_permissions
AFTER DELETE ON systems_users_permissions
BEGIN
  DELETE FROM systems_users_roles
  WHERE system_id = OLD.system_id
    AND user_id = OLD.user_id
    AND role_id = (
      SELECT role_id FROM permissions WHERE id = OLD.permission_id
    )
    AND NOT EXISTS (
      SELECT 1
      FROM systems_users_permissions sup
      JOIN permissions p ON sup.permission_id = p.id
      WHERE sup.system_id = OLD.system_id
        AND sup.user_id = OLD.user_id
        AND p.role_id = (
          SELECT role_id FROM permissions WHERE id = OLD.permission_id
        )
    );
END;


-- migrate:down

DROP TRIGGER IF EXISTS trg_a_delete_systems_users_permissions;