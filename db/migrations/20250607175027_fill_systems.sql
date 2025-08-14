-- migrate:up

INSERT INTO systems (id, name, description, repository, created, updated) VALUES (1, 'Allen, Cervantes and Brown', 'Picture lay executive politics off rise suddenly. Central voice appear kind many.
Who play already say foreign purpose expect current. Growth pick far.
Finish care ground raise easy.', 'https://murray.org/', '2024-05-10 12:41:54', '2024-05-30 12:41:54');
INSERT INTO systems (id, name, description, repository, created, updated) VALUES (2, 'Charles and Sons', 'Side enjoy provide however central condition culture. Teach least outside include present above financial.', 'http://www.kelley.net/', '2023-05-23 15:37:03', '2023-06-17 15:37:03');

-- migrate:down

DELETE FROM systems;