-- Users seeds
INSERT INTO users (username, email, password, deposit_amount) VALUES
('rizky', 'rizky@mail.com', '$2a$10$examplehashedpassword001', 50000);

-- Memberships seeds
INSERT INTO memberships (name, costs, duration_days, availability, category) VALUES
('Basic Membership', 150000, 30, TRUE, 'basic'),
('Premium Membership', 200000, 30, TRUE, 'premium'),
('VIP Membership', 250000, 30, TRUE, 'vip');

-- Trainer seeds
INSERT INTO trainers (name, specialization) VALUES
('Andi Pratama', 'Strength Training'),
('Dewi Lestari', 'Yoga'),
('Fajar Ramadhan', 'Cardio');

-- Gym class seeds
INSERT INTO gym_classes ( trainer_id, class_name, schedule, quota, availability) VALUES
(1, 'Strength Beginner', '2026-12-01 09:00:00+07', 10, TRUE),
(2, 'Morning Yoga', '2026-12-01 07:00:00+07', 15, TRUE),
(3, 'Cardio Blast', '2026-12-01 17:00:00+07', 12, TRUE);