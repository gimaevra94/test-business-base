INSERT INTO users (name, role) VALUES
('Dispatcher 2', 'dispatcher'),
('Dispatcher 3', 'dispatcher'),
('Master 3', 'master'),
('Master 4', 'master'),
('Master 5', 'master'),
('Master 6', 'master');
INSERT INTO requests (client_name, phone, address, problem_text, status, assigned_to) VALUES
('Client 3', '333-4567', 'Street 5, Apt 12', 'Water leakage in kitchen', 'in_progress', 3),
('Client 4', '444-8910', 'Park Ave 25', 'Heating not working', 'canceled', 4),
('Client 5', '555-1122', 'Central St 45, Floor 3', 'Broken window', 'new', NULL),
('Client 6', '666-3344', 'River Road 8', 'Electrical issues', 'assigned', 5),
('Client 7', '777-5566', 'Hilltop Dr 15', 'Clogged drain', 'in_progress', 6),
('Client 8', '888-7788', 'Lakeview Blvd 30', 'Door lock broken', 'new', NULL),
('Client 9', '999-9900', 'Maple Street 7', 'Bathroom renovation needed', 'done', 2);