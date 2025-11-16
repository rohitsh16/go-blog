-- Create the database
CREATE DATABASE go_blog_db;

-- Switch to the new database
USE go_blog_db;

-- Create the table for blog posts
CREATE TABLE posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    slug VARCHAR(255) DEFAULT NULL,
    author_id INT DEFAULT NULL,
    author_name VARCHAR(100) DEFAULT NULL,
    published TINYINT(1) DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_posts_slug (slug)
);
