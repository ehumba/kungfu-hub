-- +gooseUp
INSERT INTO martial_arts (name) VALUES
('Judo'),
('Taekwondo'),
('Jiu-Jitsu');

-- +gooseDown
DELETE FROM martial_arts WHERE name IN ('Judo', 'Taekwondo', 'Jiu-Jitsu');