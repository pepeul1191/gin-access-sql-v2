-- migrate:up

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

-- migrate:down

DROP TABLE systems_users_roles;