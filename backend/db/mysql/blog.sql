-- Create the database
CREATE DATABASE go_blog_db;

-- Switch to the new database
USE go_blog_db;

-- Create the table for blog posts
CREATE TABLE posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
