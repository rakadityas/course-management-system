-- Insert initial students
INSERT INTO students (email, create_time, update_time) VALUES 
('student1@example.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), 
('student2@example.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), 
('student3@example.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert initial courses
INSERT INTO courses (name, create_time, update_time) VALUES 
('Mathematics 101', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), 
('Introduction to Programming', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), 
('History of Art', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert initial course enrollments
INSERT INTO course_enrollments (student_id, course_id, status, create_time, update_time) VALUES 
(1, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),  -- Status 1 could represent "Active"
(2, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), 
(2, 2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), 
(3, 3, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);  -- Status 2 could represent "Cancelled"
