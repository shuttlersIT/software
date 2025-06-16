-- Table: departments
CREATE TABLE departments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Table: teams
CREATE TABLE teams (
    id INT AUTO_INCREMENT PRIMARY KEY,
    department_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE
);

-- Table: staff
CREATE TABLE staff (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    status ENUM('Active', 'Inactive') NOT NULL,
    department_id INT NOT NULL,
    team_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY idx_unique_email (email),
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
);

-- Table: software
CREATE TABLE software (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    type ENUM('License', 'SaaS', 'Subscription', 'Other') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Table: software_assignments
CREATE TABLE software_assignments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    software_id INT NOT NULL,
    scope_type ENUM('Department', 'Team', 'Staff') NOT NULL,
    scope_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (software_id) REFERENCES software(id) ON DELETE CASCADE,
    INDEX idx_scope (scope_type, scope_id)
);

-- Table: assigned_software
CREATE TABLE assigned_software (
    id INT AUTO_INCREMENT PRIMARY KEY,
    staff_id INT NOT NULL,
    software_id INT NOT NULL,
    source ENUM('manual', 'department', 'team', 'organization') NOT NULL DEFAULT 'manual',
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (staff_id) REFERENCES staff(id) ON DELETE CASCADE,
    FOREIGN KEY (software_id) REFERENCES software(id) ON DELETE CASCADE,
    UNIQUE (staff_id, software_id)
);

-- Table: software_assignment_logs
CREATE TABLE software_assignment_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    staff_id INT NOT NULL,
    software_id INT NOT NULL,
    action ENUM('Assigned', 'Unassigned') NOT NULL,
    changed_by INT NOT NULL,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (staff_id) REFERENCES staff(id) ON DELETE CASCADE,
    FOREIGN KEY (software_id) REFERENCES software(id) ON DELETE CASCADE,
    FOREIGN KEY (changed_by) REFERENCES staff(id) ON DELETE CASCADE,
    INDEX (staff_id),
    INDEX (changed_at)
);

-- Table: software_department_matches
CREATE TABLE software_organization_matches (
    id INT AUTO_INCREMENT PRIMARY KEY,
    software_id INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (software_id) REFERENCES software(id) ON DELETE CASCADE
);

-- Table: software_department_matches
CREATE TABLE software_department_matches (
    id INT AUTO_INCREMENT PRIMARY KEY,
    software_id INT NOT NULL,
    department_id INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (software_id) REFERENCES software(id) ON DELETE CASCADE,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE
);

-- Table: software_team_matches
CREATE TABLE software_team_matches (
    id INT AUTO_INCREMENT PRIMARY KEY,
    software_id INT NOT NULL,
    team_id INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (software_id) REFERENCES software(id) ON DELETE CASCADE,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
);
