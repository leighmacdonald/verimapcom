INSERT INTO agency(agency_name, created_on, slots, invite_key)
VALUES ('Verimap Plus Inc.', now(), 10, 'asdf');

INSERT INTO person (agency_id, email, password_hash, first_name, last_name, created_on, deleted)
VALUES (1, 'leigh.macdonald@gmail.com', '$2a$14$RLM7L5eJFxUTP1Vpby0JEesCWzd2LGRegUfJaCYiuETIhr1058dW6',
        'Leigh', 'MacDonald', now(), false);