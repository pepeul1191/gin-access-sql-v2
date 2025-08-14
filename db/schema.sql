CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username VARCHAR(20) NOT NULL,
  password VARCHAR(100) NOT NULL,
  activation_key VARCHAR(30),
  reset_key VARCHAR(30),
  email VARCHAR(50) UNIQUE NOT NULL,
  activated BOOLEAN NOT NULL DEFAULT 0,
  created DATETIME NOT NULL,
  updated DATETIME NOT NULL
);
CREATE TABLE systems (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(40) NOT NULL,
  description TEXT,
  repository VARCHAR(100),
  created DATETIME NOT NULL,
  updated DATETIME NOT NULL
);
CREATE TABLE roles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(40) NOT NULL,
  created DATETIME NOT NULL,
  updated DATETIME NOT NULL,
  system_id INTEGER NOT NULL,
  FOREIGN KEY(system_id) REFERENCES systems(id)
);
CREATE TABLE permissions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(40) NOT NULL,
  created DATETIME NOT NULL,
  updated DATETIME NOT NULL,
  role_id INTEGER NOT NULL,
  FOREIGN KEY(role_id) REFERENCES roles(id)
);
CREATE TABLE systems_users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  system_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  created DATETIME NOT NULL,
  FOREIGN KEY(system_id) REFERENCES systems(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE TABLE systems_users_permissions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  system_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  permission_id INTEGER NOT NULL,
  created DATETIME NOT NULL,
  FOREIGN KEY(system_id) REFERENCES systems(id),
  FOREIGN KEY(user_id) REFERENCES users(id),
  FOREIGN KEY(permission_id) REFERENCES permissions(id)
);
CREATE VIEW vw_system_users AS
SELECT
    SU.user_id AS id,
    SU.system_id,
    U.username,
    U.password,
    U.email,
    U.activated
FROM systems_users SU
INNER JOIN users U ON SU.user_id = U.id
/* vw_system_users(id,system_id,username,password,email,activated) */;
CREATE TABLE systems_users_roles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  system_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  role_id INTEGER NOT NULL,
  created DATETIME NOT NULL,
  FOREIGN KEY(system_id) REFERENCES systems(id),
  FOREIGN KEY(user_id) REFERENCES users(id),
  FOREIGN KEY(role_id) REFERENCES roles(id)
);
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
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20250607174507'),
  ('20250607174516'),
  ('20250607174522'),
  ('20250607174528'),
  ('20250607174542'),
  ('20250607174552'),
  ('20250607175027'),
  ('20250607175057'),
  ('20250624051821'),
  ('20250802045248'),
  ('20250802045923'),
  ('20250802050229');
