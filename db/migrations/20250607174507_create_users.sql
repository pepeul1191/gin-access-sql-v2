-- migrate:up

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

-- migrate:down

DROP TABLE users;