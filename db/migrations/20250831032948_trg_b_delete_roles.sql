-- migrate:up

CREATE TRIGGER delete_role_permissions
BEFORE DELETE ON roles
FOR EACH ROW
BEGIN
    DELETE FROM permissions WHERE role_id = OLD.id;
END;


-- migrate:down

DROP TRIGGER IF EXISTS delete_role_permissions;