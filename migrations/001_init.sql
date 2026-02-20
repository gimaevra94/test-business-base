CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100),
    role ENUM('dispatcher', 'master') NOT NULL
);

CREATE TABLE requests (
    id INT PRIMARY KEY AUTO_INCREMENT,
    client_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    address VARCHAR(255) NOT NULL,
    problem_text TEXT NOT NULL,
    status ENUM('new', 'assigned', 'in_progress', 'done', 'canceled') DEFAULT 'new',
    assigned_to INT NULL,
    version INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (assigned_to) REFERENCES users(id)
);