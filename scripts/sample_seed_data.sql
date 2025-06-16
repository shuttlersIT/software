-- Clear tables in proper order to prevent foreign key constraint errors
SET FOREIGN_KEY_CHECKS = 0;

TRUNCATE TABLE software_assignment_logs;
TRUNCATE TABLE assigned_software;
TRUNCATE TABLE software_assignments;
TRUNCATE TABLE software_team_matches;
TRUNCATE TABLE software_department_matches;
TRUNCATE TABLE software_organization_matches;
TRUNCATE TABLE software;
TRUNCATE TABLE staff;
TRUNCATE TABLE teams;
TRUNCATE TABLE departments;

SET FOREIGN_KEY_CHECKS = 1;

-- Sample Seed Data for departments
INSERT INTO departments (name) VALUES 
('Engineering'),     -- id = 1
('Operations'),      -- id = 2
('Human Resources'); -- id = 3

-- Sample Seed Data for teams
INSERT INTO teams (department_id, name) VALUES 
(1, 'Backend'),       -- id = 1
(1, 'Frontend'),      -- id = 2
(2, 'Logistics'),     -- id = 3
(3, 'Recruitment');   -- id = 4

-- Sample Seed Data for staff
INSERT INTO staff (first_name, last_name, email, status, department_id, team_id) VALUES 
('Alice', 'Ngugi', 'alice.ngugi@company.com', 'Active', 1, 1),    -- id = 1
('Brian', 'Okoro', 'brian.okoro@company.com', 'Active', 1, 2),    -- id = 2
('Chinedu', 'Umeh', 'chinedu.umeh@company.com', 'Active', 2, 3),  -- id = 3
('Diana', 'Adeyemi', 'diana.adeyemi@company.com', 'Active', 3, 4);-- id = 4

-- Sample Seed Data for software
INSERT INTO software (name, description, type) VALUES 
('Slack', 'Team communication tool', 'SaaS'),             -- id = 1
('Jira', 'Project and issue tracking', 'SaaS'),           -- id = 2
('Figma', 'Collaborative design tool', 'Subscription'),   -- id = 3
('Microsoft Office', 'Productivity suite', 'License');    -- id = 4

-- Sample Seed Data for software_assignments (high-level logical mapping)
INSERT INTO software_assignments (software_id, scope_type, scope_id) VALUES 
(1, 'Department', 1), -- Slack to Engineering
(2, 'Department', 1), -- Jira to Engineering
(3, 'Team', 2),       -- Figma to Frontend
(4, 'Staff', 3);      -- MS Office to Chinedu

-- Sample Seed Data for software_department_matches
INSERT INTO software_department_matches (software_id, department_id) VALUES 
(1, 1), -- Slack → Engineering
(2, 1); -- Jira  → Engineering

-- Sample Seed Data for software_team_matches
INSERT INTO software_team_matches (software_id, team_id) VALUES 
(3, 2); -- Figma → Frontend

-- Sample Seed Data for software_organization_matches
INSERT INTO software_organization_matches (software_id) VALUES 
(4); -- MS Office assigned to all staff org-wide

-- Sample Seed Data for assigned_software
INSERT INTO assigned_software (staff_id, software_id, source) VALUES 
(1, 1, 'department'), -- Alice gets Slack
(1, 2, 'department'), -- Alice gets Jira
(2, 1, 'department'), -- Brian gets Slack
(2, 3, 'team'),       -- Brian gets Figma
(3, 4, 'manual'),     -- Chinedu manually assigned MS Office
(4, 4, 'organization'); -- Diana receives org-wide Office

-- Sample Seed Data for software_assignment_logs
INSERT INTO software_assignment_logs (staff_id, software_id, action, changed_by) VALUES 
(1, 1, 'Assigned', 3),
(1, 2, 'Assigned', 3),
(2, 1, 'Assigned', 3),
(2, 3, 'Assigned', 4),
(3, 4, 'Assigned', 1),
(4, 4, 'Assigned', 1);
