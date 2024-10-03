-- Create the database
CREATE DATABASE blogdb;

-- Switch to the new database
USE blogdb;

-- Create the table for blog posts
CREATE TABLE posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
