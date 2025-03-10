CREATE TABLE `user` (
  `user_id` integer PRIMARY KEY,
  `username` varchar(255),
  `email` varchar(255),
  `password` varchar(255),
  `profil_picture` varchar(255)
);

CREATE TABLE `categories` (
  `category_id` integer PRIMARY KEY,
  `title` varchar(255),
  `description` varchar(255)
);

CREATE TABLE `topics` (
  `topic_id` integer PRIMARY KEY,
  `category_id` integer,
  `user_id` integer,
  `title` varchar(255),
  FOREIGN KEY (`category_id`) REFERENCES `categories` (`category_id`),
  FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`)
  
);

CREATE TABLE `posts` (
  `post_id` integer PRIMARY KEY,
  `topic_id` integer,
  `user_id` integer,
  `title` varchar(255),
  `content` varchar(255),
  `image` image,
  FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`),
  FOREIGN KEY (`topic_id`) REFERENCES `topics` (`topic_id`)
);

CREATE TABLE `answers` (
  `answer_id` integer PRIMARY KEY,
  `post_id` integer,
  `user_id` integer,
  `content` varchar(255),
  FOREIGN KEY (`post_id`) REFERENCES `posts` (`post_id`)
  FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`)
);


