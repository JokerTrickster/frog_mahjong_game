
CREATE TABLE tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    refresh_expired_at INT,
    user_id INT,
    access_token VARCHAR(255),
    refresh_token VARCHAR(255)
);

create table missions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL, 
    name VARCHAR(255),
    description VARCHAR(255)
);

INSERT INTO missions (name, description) VALUES ('연속된 숫자 2쌍', '연속된 숫자 2쌍을 완성해주세요 ex) 123 567');
INSERT INTO missions (name, description) VALUES ('동일한 숫자 2쌍', '동일한 숫자 2쌍을 완성해주세요 ex) 111 222');


CREATE TABLE rooms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    current_count INT default 0,
    play_turn int default 0,
    max_count INT default 2,
    min_count INT default 2,
    name VARCHAR(255),
    password VARCHAR(255),
    state VARCHAR(50),
    owner_id INT,
    timer INT,
    mission_id INT default 1
);

INSERT INTO rooms (current_count, max_count, min_count, name, password, state, owner_id)
VALUES (0, 10, 1, 'Example Room', 'room_password', 'waiting', 1);	

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    name VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255),
    coin INT,
    state VARCHAR(50),
    profile_id INT default 1,
	room_id INT,
    provider VARCHAR(50),
    FOREIGN KEY (room_id) REFERENCES rooms(id)
);

INSERT INTO users (name, email, password, coin, state, profile_id, room_id, provider) 
VALUES ('test', 'test@test.com', 'asd123', 100, 'logout', 1, 1, 'email');
INSERT INTO users (name, email, password, coin, state, profile_id, room_id, provider) 
VALUES ('test2', 'test2@test.com', 'asd123', 100, 'logout', 1, 1, 'email');
INSERT INTO users (name, email, password, coin, state, profile_id, room_id, provider) 
VALUES ('test3', 'test3@test.com', 'asd123', 100, 'logout', 1, 1, 'email');
INSERT INTO users (name, email, password, coin, state, profile_id, room_id, provider) 
VALUES ('test4', 'test4@test.com', 'asd123', 100, 'logout', 1, 1, 'email');



CREATE TABLE room_users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    user_id INT,
    room_id INT,
    score INT,
    owned_card_count INT,
    player_state VARCHAR(50),
    turn_number INT default 0,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (room_id) REFERENCES rooms(id)
);

CREATE TABLE cards (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    room_id INT,
	user_id INT,
    card_id INT,
    name VARCHAR(255),
    color VARCHAR(50),
    state VARCHAR(50),
    FOREIGN KEY (room_id) REFERENCES rooms(id)
);


CREATE TABLE chats (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    user_id INT,
    room_id INT,
    name varchar(255),
    message varchar(255),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (room_id) REFERENCES rooms(id)
);

CREATE TABLE user_auths (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    auth_code VARCHAR(255),
    email VARCHAR(255),
    type VARCHAR(100)
);

-- 메타 데이터 테이블
  CREATE TABLE meta_tables (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    table_name VARCHAR(255) NOT NULL UNIQUE,
    table_description VARCHAR(255)
);
INSERT INTO meta_tables (table_name, table_description) VALUES ('times', '게임 타이머');

-- 게임 시간 타이머 관리 메타 테이블
CREATE TABLE times (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    timer INT,
    description VARCHAR(255)
);
-- times 테이블에 15초, 30초, 60초 데이터를 넣어줘
INSERT INTO times (timer, description) VALUES (15, '15초'),(30, '30초'),(60, '60초');

-- 신고 테이블
CREATE TABLE reports (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    target_user_id int,
    reporter_user_id int,
    category_id int,
    reason varchar(1000)
);
-- 신고하기 메타 테이블
CREATE TABLE categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    reason VARCHAR(1000)
);

INSERT INTO categories (type, reason) VALUES ('report', '도배 및 불건전한 언어 사용'),('report', '불법 프로그램 사용'),('report', '비매너 행위'),('report', '기타');



create table profiles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL, 
    total_count INT,        -- 프로필 달성 조건 횟수
    image VARCHAR(255),  -- 프로필 이미지 경로
    description VARCHAR(255) -- 프로필 획득 설명
);
INSERT INTO profiles (name, total_count, image, description) VALUES ('소라게 개굴', 0, '1.png', '기본 이미지');
INSERT INTO profiles (name, total_count, image, description) VALUES ('증명 사진 개굴', 0, '2.png', '기본 이미지');
INSERT INTO profiles (name, total_count, image, description) VALUES ('뽀또', 0, '3.png', '기본 이미지');
INSERT INTO profiles (name, total_count, image, description) VALUES ('분홍 개굴', 0, '4.png', '기본 이미지');

CREATE TABLE user_profiles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL, 
    user_id INT,
    profile_id INT,
    current_count INT DEFAULT 0,  -- 해당 프로필 달성한 횟수
    is_achieved BOOLEAN DEFAULT FALSE, -- 달성 여부
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);