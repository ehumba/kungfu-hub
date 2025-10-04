-- +gooseUp
INSERT INTO martial_arts (name) VALUES
('Judo'),
('Taekwondo'),
('Wushu');

-- +gooseDown
DELETE FROM martial_arts WHERE name IN ('Judo', 'Taekwondo', 'Wushu');