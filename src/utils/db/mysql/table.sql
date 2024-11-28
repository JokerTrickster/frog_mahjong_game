
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

INSERT INTO missions (name, description) VALUES ('새 크기가 60cm 이상', '새 크기가 60cm 이상인 새를 완성해주세요');
INSERT INTO missions (name, description) VALUES ('새 크기가 60cm 미만', '새 크기가 60cm 미만인 새를 완성해주세요');
INSERT INTO missions (name, description) VALUES ('숲에서 사는 새', '서식지가 숲인 새를 모아주세요');
INSERT INTO missions (name, description) VALUES ('초원에서 사는 새', '서식지가 초원인 새를 모아주세요');
INSERT INTO missions (name, description) VALUES ('물에서 사는 새', '서식지가 물인 새를 모아주세요');
INSERT INTO missions (name, description) VALUES ('이름에 신체 부위가 들어간 새', '새 이름에 신체 부위가 들어간 새를 모아주세요 귀, 눈, 날개, 다리, 부리, 꼬리, 몸통');
INSERT INTO missions (name, description) VALUES ('이름에 색깔이 들어간 새', '새 이름에 색깔이 들어간 새를 모아주세요 빨간, 주황, 노란, 초록, 파란, 남색, 보라');
INSERT INTO missions (name, description) VALUES ('서식지가 2곳 이상인 새', '서식지가 2곳 이상인 새를 모아주세요');
INSERT INTO missions (name, description) VALUES ('부리 방향이 오른쪽인 새', '새 부리 방향이 오른쪽인 새를 모아주세요');
INSERT INTO missions (name, description) VALUES ('부리 방향이 왼쪽인 새', '새 부리 방향이 왼쪽인 새를 모아주세요');




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
    timer INT
);

INSERT INTO rooms (current_count, max_count, min_count, name, password, state, owner_id)
VALUES (0, 10, 1, 'Example Room', 'room_password', 'waiting', 1);	

-- 중간 테이블로 rooms와 missions의 다대다 관계를 정의
CREATE TABLE room_missions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    room_id INT,
    mission_id INT,
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (mission_id) REFERENCES missions(id) ON DELETE CASCADE
);

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
    alert_enabled TINYINT(1) DEFAULT 1, 
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

CREATE TABLE user_tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    user_id INT,
    token VARCHAR(1000),
    FOREIGN KEY (user_id) REFERENCES users(id)
);


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

-- 이름, 이미지, 설명, 크기, 서식지, 부리 방향, 둥지 형태
create table bird_cards (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL, 
    name VARCHAR(255),
    image VARCHAR(255),
    description VARCHAR(255),
    size INT,
    habitat VARCHAR(255),
    beak_direction VARCHAR(255),
    nest varchar(255) -- 그릇형 bowl, 구멍둥지 cavity, 자유형 wild, 땅둥지 ground, 평평형 platform
);


CREATE TABLE user_missions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    room_id INT,
    mission_id INT,
    user_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (mission_id) REFERENCES missions(id) ON DELETE CASCADE
);

CREATE TABLE user_mission_cards (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    user_mission_id INT,
    card_id INT,
    FOREIGN KEY (user_mission_id) REFERENCES user_missions(id) ON DELETE CASCADE,
    FOREIGN KEY (card_id) REFERENCES bird_cards(id) ON DELETE CASCADE
);

CREATE TABLE user_bird_cards (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    user_id INT,
    room_id INT,
    card_id INT,
    state VARCHAR(50),
    FOREIGN KEY (card_id) REFERENCES bird_cards(id) ON DELETE CASCADE,
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);


CREATE TABLE foods (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    name VARCHAR(255) NOT NULL,
    food_image_id INT,
    FOREIGN KEY (food_image_id) REFERENCES food_images(id)
);

CREATE TABLE categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    name VARCHAR(255) NOT NULL UNIQUE,
    type ENUM('time', 'type', 'scenario', 'theme') NOT NULL COMMENT '카테고리 유형'
);

CREATE TABLE food_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    food_id INT NOT NULL,
    category_id INT NOT NULL,
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);
