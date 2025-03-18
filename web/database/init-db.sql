-- Script d'initialisation pour la base de données MySQL du forum

-- Utiliser la base de données forum_db
USE forum_db;

-- Table des utilisateurs
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    role ENUM('user', 'moderator', 'admin') DEFAULT 'user'
);

-- Table des catégories
CREATE TABLE categories (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(36),
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Table des topics
CREATE TABLE topics (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    category_id VARCHAR(36) NOT NULL,
    created_by VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Table des posts (réponses dans les topics)
CREATE TABLE posts (
    id VARCHAR(36) PRIMARY KEY,
    content TEXT NOT NULL,
    topic_id VARCHAR(36) NOT NULL,
    created_by VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (topic_id) REFERENCES topics(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Table pour stocker les likes/dislikes des posts
CREATE TABLE reactions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    post_id VARCHAR(36) NOT NULL,
    reaction_type ENUM('like', 'dislike') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY user_post_unique (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

-- Table pour les sessions utilisateur
CREATE TABLE sessions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Insertion d'un utilisateur admin par défaut (mot de passe: admin123)
-- Note: En production, il faudrait utiliser un mot de passe plus sécurisé et le hasher correctement
INSERT INTO users (id, username, email, password, role)
VALUES (
    UUID(),
    'admin',
    'admin@forum.com',
    '$2a$10$XgU2GZP8.7WEQ9y5WEr4XeRjxqp.p16QA08ELnJw3EBqLwUHeXOPi', -- admin123 hashé avec bcrypt
    'admin'
);

-- Insertion de quelques catégories par défaut
INSERT INTO categories (id, name, description, created_by)
VALUES 
    (UUID(), 'Général', 'Discussions générales sur divers sujets', (SELECT id FROM users WHERE username = 'admin')),
    (UUID(), 'Technologie', 'Discussions sur les technologies, la programmation, et l\'informatique', (SELECT id FROM users WHERE username = 'admin')),
    (UUID(), 'Loisirs', 'Discussions sur les hobbies, les sports, et les activités de loisir', (SELECT id FROM users WHERE username = 'admin'));