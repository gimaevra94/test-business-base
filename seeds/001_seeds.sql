INSERT INTO users (name, role) VALUES ('Dispatcher 1', 'dispatcher'), ('Master 1', 'master'), ('Master 2', 'master');
INSERT INTO requests (client_name, phone, address, problem_text, status, assigned_to) VALUES 
('Client 1', '111', 'Addr 1', 'Problem 1', 'new', NULL),
('Client 2', '222', 'Addr 2', 'Problem 2', 'assigned', 2);