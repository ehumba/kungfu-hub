-- +gooseUp
INSERT INTO martial_arts (name) VALUES
('Karate'),
('Judo'),
('Taekwondo'),
('Brazilian Jiu-Jitsu'),
('Sanda');

-- +gooseDown
DELETE FROM martial_arts WHERE name IN ('Karate', 'Judo', 'Taekwondo', 'Brazilian Jiu-Jitsu', 'Sanda');