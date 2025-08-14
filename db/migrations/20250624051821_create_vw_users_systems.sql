-- migrate:up

CREATE VIEW vw_system_users AS
SELECT 
    SU.user_id AS id, 
    SU.system_id, 
    U.username, 
    U.password, 
    U.email, 
    U.activated
FROM systems_users SU 
INNER JOIN users U ON SU.user_id = U.id;

-- migrate:down

DROP VIEW vw_system_users;