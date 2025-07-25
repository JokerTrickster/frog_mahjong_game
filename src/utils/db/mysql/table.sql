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
    description VARCHAR(255),
    image VARCHAR(500)
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
    timer INT,
    game_id INT,    -- 개굴작 :1 , 윙스팬 : 2, 틀린그림 찾기 : 3
    start_time TIMESTAMP
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
    coin INT,   -- 게임 플레이에 필요한 코인 
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
    type VARCHAR(100),
    project varchar(200),
    is_active BOOLEAN DEFAULT FALSE
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


CREATE TABLE items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    name VARCHAR(255), -- 아이템 이름
    description VARCHAR(500), -- 아이템 설명
    max_uses INT DEFAULT 0 -- 아이템 최대 사용 가능 횟수
);

INSERT INTO items (name, description, max_uses) 
VALUES 
('cards change', '카드덱 교체',  3),
('missions change', '미션 교체',  1),
('get discarded card', '버린 카드 가져오기', 3);


CREATE TABLE user_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    user_id INT, -- 유저 ID
    item_id INT, -- 아이템 ID
    room_id INT, -- 룸 ID
    remaining_uses INT DEFAULT 0, -- 남은 사용 횟수
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
	FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

CREATE TABLE game_users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    name VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255),
    coin INT,   -- 게임 플레이에 필요한 코인 
    state VARCHAR(50),
    profile_id INT default 1,
	room_id INT,
    alert_enabled TINYINT(1) DEFAULT 1, 
    provider VARCHAR(50) -- email, google, kakao
);
CREATE TABLE game_rooms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    current_count INT default 0,
    max_count INT default 2,
    min_count INT default 2,
    name VARCHAR(255),
    password VARCHAR(255),
    state VARCHAR(50),
    owner_id INT,
    game_id INT,    -- 틀린그림 찾기 : 1
    start_time TIMESTAMP
);

create table frog_room_users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    user_id INT,
    room_id INT,
    score INT,
    owned_card_count INT,
    player_state VARCHAR(50),       -- 유저 게임 상태
    turn_number INT default 0,
    FOREIGN KEY (user_id) REFERENCES game_users(id),
    FOREIGN KEY (room_id) REFERENCES game_rooms(id)
);

create table frog_cards(
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    count VARCHAR(255),
    color VARCHAR(50),
    image varchar(255)
);

INSERT INTO frog_cards (count, color, image)
VALUES
    (1,  'red',    'red1.png'),
    (2,  'red',    'red2.png'),
    (3,  'red',    'red3.png'),
    (4,  'red',    'red4.png'),
    (5,  'red',    'red5.png'),
    (6,  'red',    'red6.png'),
    (7,  'red',    'red7.png'),
    (8,  'red',    'red8.png'),
    (9,  'red',    'red9.png'),
    (10, 'red',    'red10.png'),
    (11, 'green',  'green11.png'),
    (1,  'normal', 'normal1.png'),
    (2,  'green',  'green2.png'),
    (3,  'green',  'green3.png'),
    (4,  'green',  'green4.png'),
    (5,  'normal', 'normal5.png'),
    (6,  'green',  'green6.png'),
    (7,  'normal', 'normal7.png'),
    (8,  'green',  'green8.png'),
    (9,  'normal', 'normal9.png'),
    (10, 'red',    'red10.png'),
    (11, 'green',  'green11.png'),
    (1,  'normal', 'normal1.png'),
    (2,  'green',  'green2.png'),
    (3,  'green',  'green3.png'),
    (4,  'green',  'green4.png'),
    (5,  'normal', 'normal5.png'),
    (6,  'green',  'green6.png'),
    (7,  'normal', 'normal7.png'),
    (8,  'green',  'green8.png'),
    (9,  'normal', 'normal9.png'),
    (10, 'red',    'red10.png'),
    (11, 'green',  'green11.png'),
    (1,  'normal', 'normal1.png'),
    (2,  'green',  'green2.png'),
    (3,  'green',  'green3.png'),
    (4,  'green',  'green4.png'),
    (5,  'normal', 'normal5.png'),
    (6,  'green',  'green6.png'),
    (7,  'normal', 'normal7.png'),
    (8,  'green',  'green8.png'),
    (9,  'normal', 'normal9.png'),
    (10, 'red',    'red10.png'),
    (11, 'green',  'green11.png');


create table frog_user_cards (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    user_id INT,
    room_id INT,
    card_id INT,
    state VARCHAR(50),
    FOREIGN KEY (card_id) REFERENCES frog_cards(id) ON DELETE CASCADE,
    FOREIGN KEY (room_id) REFERENCES game_rooms(id) ON DELETE CASCADE
);


create table frog_game_room_settings(
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    room_id INT, 
    timer INT, -- 게임 타이머
    remaining_card_count INT, -- 남은 카드 수 
    current_round INT,  -- 현재 라운드
    FOREIGN KEY (room_id) REFERENCES game_rooms(id) ON DELETE CASCADE
);


-- find it 

create table game_profiles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL, 
    image VARCHAR(255),  -- 프로필 이미지 경로
    description VARCHAR(255) -- 프로필 획득 설명
);
CREATE TABLE game_user_profiles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL, 
    user_id INT,
    profile_id INT,
    is_achieved BOOLEAN DEFAULT FALSE, -- 달성 여부
    FOREIGN KEY (user_id) REFERENCES game_users(id) ON DELETE CASCADE,
    FOREIGN KEY (profile_id) REFERENCES game_profiles(id) ON DELETE CASCADE
);

CREATE TABLE game_room_users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    user_id INT,
    room_id INT,
    player_state VARCHAR(50),
    FOREIGN KEY (user_id) REFERENCES game_users(id),
    FOREIGN KEY (room_id) REFERENCES game_rooms(id)
);

------------- find it 게임 테이블 -------------
CREATE TABLE find_it_room_settings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    room_id INT NOT NULL,
    timer INT,                  -- 게임 타이머
    lifes INT DEFAULT 3,          -- 초기 목숨
    item_hint_count INT DEFAULT 2, -- 힌트 아이템 개수
    item_timer_stop_count INT DEFAULT 2, -- 타이머 정지 아이템 개수
    round INT,                 -- 게임 라운드
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (room_id) REFERENCES game_rooms(id) ON DELETE CASCADE
);

CREATE TABLE find_it_images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    level INT NOT NULL DEFAULT 1,  -- 난이도
    normal_image_url VARCHAR(500) NOT NULL,  -- 정상 이미지 URL
    abnormal_image_url VARCHAR(500) NOT NULL -- 비정상 이미지 URL
);

CREATE TABLE find_it_image_correct_positions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    image_id INT NOT NULL,  -- 어떤 이미지의 정답인지 (find_it_images 테이블 참조)
    x_position DOUBLE NOT NULL, -- 정답 X 좌표
    y_position DOUBLE NOT NULL, -- 정답 Y 좌표
    FOREIGN KEY (image_id) REFERENCES find_it_images(id)
);
CREATE TABLE find_it_round_images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    round INT NOT NULL,          -- 라운드 번호
    image_set_id INT NOT NULL,   -- 사용된 이미지 세트 (find_it_images 테이블과 연결)
    room_id INT NOT NULL,        -- 어떤 게임방에서 사용되는지
    FOREIGN KEY (room_id) REFERENCES game_rooms(id),
    FOREIGN KEY (image_set_id) REFERENCES find_it_images(id)
);

CREATE TABLE find_it_user_correct_positions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    user_id INT NOT NULL,          -- 정답을 맞춘 유저 ID
    room_id INT NOT NULL,          -- 해당 정답이 속한 게임 방 ID
    round INT NOT NULL,            -- 라운드 정보
    image_id INT NOT NULL,         -- 정답을 맞춘 이미지 ID (find_it_images 테이블 참조)
    correct_position_id INT NOT NULL, -- 맞춘 정답의 ID (find_it_image_correct_positions 테이블 참조)
    
    FOREIGN KEY (user_id) REFERENCES game_users(id) ON DELETE CASCADE,
    FOREIGN KEY (room_id) REFERENCES game_rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (image_id) REFERENCES find_it_images(id) ON DELETE CASCADE,
    FOREIGN KEY (correct_position_id) REFERENCES find_it_image_correct_positions(id) ON DELETE CASCADE,

    UNIQUE (user_id, room_id, round, correct_position_id) -- 중복 방지
);


create table games (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    title VARCHAR(255),
    description VARCHAR(255),
    category varchar(255),
    hash_tag varchar(255),
    youtube_url varchar(255),
    image VARCHAR(500),
    is_enabled BOOLEAN DEFAULT FALSE
);
INSERT INTO games (title, description, category, hash_tag, youtube_url, is_enabled, image)
VALUES 
    ('틀린그림찾기', '협력해서 틀린그림을 찾아보세요', '협력', '보드게임', 'https://www.youtube.com/shorts/h6wIckelzpk', TRUE, 'find-it.png'),
    ('슬라임전쟁', '자신의 영역을 넓혀서 상대방을 이겨보세요.', '전략', '보드게임', 'https://www.youtube.com/shorts/h6wIckelzpk', TRUE, 'slime-war.png'),
    ('카후나', '자신의 영역을 넓혀서 상대방을 이겨보세요.', '전략', '보드게임', 'https://www.youtube.com/shorts/h6wIckelzpk', FALSE, 'default.png'),
    ('스플랜더 듀얼', '미션을 빨리 달성해서 상대방을 이겨보세요.', '전략', '보드게임', 'https://www.youtube.com/shorts/h6wIckelzpk', FALSE, 'default.png');



INSERT INTO find_it_images (level, normal_image_url, abnormal_image_url)
VALUES 
(1, 'normal1-level1.jpg', 'abnormal1-level1.jpg'),
(1, 'normal2-level1.jpg', 'abnormal2-level1.jpg'),
(1, 'normal3-level1.jpg', 'abnormal3-level1.jpg'),
(1, 'normal4-level1.jpg', 'abnormal4-level1.jpg'),
(1, 'normal5-level1.jpg', 'abnormal5-level1.jpg'),
(1, 'normal6-level1.jpg', 'abnormal6-level1.jpg'),
(1, 'normal7-level1.jpg', 'abnormal7-level1.jpg'),
(1, 'normal8-level1.jpg', 'abnormal8-level1.jpg'),
(1, 'normal9-level1.jpg', 'abnormal9-level1.jpg'),
(1, 'normal10-level1.jpg', 'abnormal10-level1.jpg');


INSERT INTO find_it_image_correct_positions (image_id, x_position, y_position)
VALUES 
(1, 60.2, 40.3), (1, 130.6, 85.7), (1, 190.4, 120.8), (1, 220.2, 160.1), (1, 310.7, 230.9),
(2, 45.3, 55.9), (2, 115.8, 95.6), (2, 175.2, 130.7), (2, 205.5, 145.8), (2, 290.3, 205.4),
(3, 60.2, 40.3), (3, 130.6, 85.7), (3, 190.4, 120.8), (3, 220.2, 160.1), (3, 310.7, 230.9),
(4, 45.3, 55.9), (4, 115.8, 95.6), (4, 175.2, 130.7), (4, 205.5, 145.8), (4, 290.3, 205.4),
(5, 55.7, 60.2), (5, 125.3, 90.4), (5, 185.5, 125.9), (5, 215.8, 155.2), (5, 295.4, 215.6),
(6, 48.2, 42.8), (6, 128.6, 92.3), (6, 176.9, 127.6), (6, 220.4, 165.7), (6, 305.3, 225.1),
(7, 52.3, 57.1), (7, 132.9, 98.4), (7, 180.5, 135.6), (7, 225.2, 172.3), (7, 308.7, 235.8),
(8, 42.7, 50.5), (8, 138.3, 99.7), (8, 188.2, 142.3), (8, 230.9, 178.6), (8, 312.5, 240.9),
(9, 58.5, 45.3), (9, 140.2, 102.5), (9, 190.8, 145.9), (9, 235.3, 185.2), (9, 315.6, 248.1),
(10, 60.5, 48.8), (10, 144.6, 105.3), (10, 195.7, 150.4), (10, 240.2, 190.7), (10, 320.9, 255.2);

-- 게임 카드 정보
create table slime_war_cards (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    direction INT,   -- 방향 0 : 왼쪽 위, 1: 위, 2: 오른쪽 위, 3: 왼쪽, 4: 오른쪽, 5: 왼쪽 아래, 6: 아래, 7: 오른쪽 아래
    image VARCHAR(255),
    move INT -- 1 ~ 3

);
INSERT INTO slime_war_cards (direction, image, move)
VALUES
(0, 'card_0_1.png', 1), (0, 'card_0_1.png', 1),
(0, 'card_0_2.png', 2), (0, 'card_0_2.png', 2),
(0, 'card_0_3.png', 3), (0, 'card_0_3.png', 3),

(1, 'card_1_1.png', 1), (1, 'card_1_1.png', 1),
(1, 'card_1_2.png', 2), (1, 'card_1_2.png', 2),
(1, 'card_1_3.png', 3), (1, 'card_1_3.png', 3),

(2, 'card_2_1.png', 1), (2, 'card_2_1.png', 1),
(2, 'card_2_2.png', 2), (2, 'card_2_2.png', 2),
(2, 'card_2_3.png', 3), (2, 'card_2_3.png', 3),

(3, 'card_3_1.png', 1), (3, 'card_3_1.png', 1),
(3, 'card_3_2.png', 2), (3, 'card_3_2.png', 2),
(3, 'card_3_3.png', 3), (3, 'card_3_3.png', 3),

(4, 'card_4_1.png', 1), (4, 'card_4_1.png', 1),
(4, 'card_4_2.png', 2), (4, 'card_4_2.png', 2),
(4, 'card_4_3.png', 3), (4, 'card_4_3.png', 3),

(5, 'card_5_1.png', 1), (5, 'card_5_1.png', 1),
(5, 'card_5_2.png', 2), (5, 'card_5_2.png', 2),
(5, 'card_5_3.png', 3), (5, 'card_5_3.png', 3),

(6, 'card_6_1.png', 1), (6, 'card_6_1.png', 1),
(6, 'card_6_2.png', 2), (6, 'card_6_2.png', 2),
(6, 'card_6_3.png', 3), (6, 'card_6_3.png', 3),

(7, 'card_7_1.png', 1), (7, 'card_7_1.png', 1),
(7, 'card_7_2.png', 2), (7, 'card_7_2.png', 2),
(7, 'card_7_3.png', 3), (7, 'card_7_3.png', 3);

-- 게임 맵 정보 
create table slime_war_maps(
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    idx INT
);

INSERT INTO slime_war_maps (idx)
VALUES
(0), (1), (2), (3), (4), (5), (6), (7), (8), (9),
(10), (11), (12), (13), (14), (15), (16), (17), (18), 
(19), (20), (21), (22), (23), (24), (25), (26), (27), 
(28), (29), (30), (31), (32), (33), (34), (35), (36), 
(37), (38), (39), (40), (41), (42), (43), (44), (45), 
(46), (47), (48), (49),(50), (51), (52), (53), (54), 
(55), (56), (57), (58), (59),(60), (61), (62), (63), 
(64), (65), (66), (67), (68), (69),(70), (71), (72), 
(73), (74), (75), (76), (77), (78), (79),(80), (81), 
(82), (83), (84), (85), (86), (87), (88), (89),(90); 

-- 게임 유저 정보
create table slime_war_users(
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    user_id INT,
    room_id INT,
    hero_count INT, -- 남은 히어로 수
    turn INT, -- 플레이 순서 
    color_type INT, -- 0: 빨강, 1: 파랑
    FOREIGN KEY (user_id) REFERENCES game_users(id) ON DELETE CASCADE,
    FOREIGN KEY (room_id) REFERENCES game_rooms(id) ON DELETE CASCADE
);

create table slime_war_game_room_settings(
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    room_id INT, 
    timer INT, -- 게임 타이머
    remaining_card_count INT, -- 남은 카드 수 
    king_index INT, -- 왕 인덱스
    current_round INT,  -- 현재 라운드
    remaining_slime_count INT, -- 남은 슬라임 수
    FOREIGN KEY (room_id) REFERENCES game_rooms(id) ON DELETE CASCADE
);

create table slime_war_room_cards(
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    room_id INT,
    card_id INT,
    state varchar(100), -- 남은 카드, 소윺 카드, 버려진 카드
    user_id INT,
    FOREIGN KEY (room_id) REFERENCES game_rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (card_id) REFERENCES slime_war_cards(id) ON DELETE CASCADE
);

-- 게임 맵 정보
create table slime_war_room_maps(
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    room_id INT,
    user_id INT,
    map_id INT, -- (10 ~ 90)
    FOREIGN KEY (room_id) REFERENCES game_rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (map_id) REFERENCES slime_war_maps(id) ON DELETE CASCADE
);


-- 게임 결과 테이블 
create table game_results(
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    room_id INT,
    user_id INT,
    score INT,
    result INT, -- 0: 패배, 1: 승리
    game_type INT, -- 0: 틀린그림찾기, 1: 슬라임전쟁, 2: 시퀀스
    FOREIGN KEY (room_id) REFERENCES game_rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES game_users(id) ON DELETE CASCADE
);



-- 게임 카드 정보
create table sequence_cards (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    image VARCHAR(255),
    type varchar(255),
    count varchar(255)
);

INSERT INTO sequence_cards (image, type, count) VALUES
('all.png', 'all', 'all'),
('spacea.png', 'space', 'a'),
('spacea.png', 'space', 'a'),
('space2.png', 'space', '2'),
('space2.png', 'space', '2'),
('space3.png', 'space', '3'),
('space3.png', 'space', '3'),
('space4.png', 'space', '4'),
('space4.png', 'space', '4'),
('all.png', 'all', 'all'),
('space5.png', 'space', '5'),
('space5.png', 'space', '5'),
('space6.png', 'space', '6'),
('space6.png', 'space', '6'),
('space7.png', 'space', '7'),
('space7.png', 'space', '7'),
('space8.png', 'space', '8'),
('space8.png', 'space', '8'),
('space9.png', 'space', '9'),
('space9.png', 'space', '9'),
('space10.png', 'space', '10'),
('space10.png', 'space', '10'),
('spacek.png', 'space', 'k'),
('spacek.png', 'space', 'k'),
('spaceq.png', 'space', 'q'),
('spaceq.png', 'space', 'q'),
('hearta.png', 'heart', 'a'),
('hearta.png', 'heart', 'a'),
('heart2.png', 'heart', '2'),
('heart2.png', 'heart', '2'),
('heart3.png', 'heart', '3'),
('heart3.png', 'heart', '3'),
('heart4.png', 'heart', '4'),
('heart4.png', 'heart', '4'),
('heart5.png', 'heart', '5'),
('heart5.png', 'heart', '5'),
('heart6.png', 'heart', '6'),
('heart6.png', 'heart', '6'),
('heart7.png', 'heart', '7'),
('heart7.png', 'heart', '7'),
('heart8.png', 'heart', '8'),
('heart8.png', 'heart', '8'),
('heart9.png', 'heart', '9'),
('heart9.png', 'heart', '9'),
('heart10.png', 'heart', '10'),
('heart10.png', 'heart', '10'),
('heartk.png', 'heart', 'k'),
('heartk.png', 'heart', 'k'),
('heartq.png', 'heart', 'q'),
('heartq.png', 'heart', 'q'),
('clovaa.png', 'clova', 'a'),
('clovaa.png', 'clova', 'a'),
('clova2.png', 'clova', '2'),
('clova2.png', 'clova', '2'),
('clova3.png', 'clova', '3'),
('clova3.png', 'clova', '3'),
('clova4.png', 'clova', '4'),
('clova4.png', 'clova', '4'),
('clova5.png', 'clova', '5'),
('clova5.png', 'clova', '5'),
('clova6.png', 'clova', '6'),
('clova6.png', 'clova', '6'),
('clova7.png', 'clova', '7'),
('clova7.png', 'clova', '7'),
('clova8.png', 'clova', '8'),
('clova8.png', 'clova', '8'),
('clova9.png', 'clova', '9'),
('clova9.png', 'clova', '9'),
('clova10.png', 'clova', '10'),
('clova10.png', 'clova', '10'),
('clovak.png', 'clova', 'k'),
('clovak.png', 'clova', 'k'),
('clovaq.png', 'clova', 'q'),
('clovaq.png', 'clova', 'q'),
('diamonda.png', 'diamond', 'a'),
('diamonda.png', 'diamond', 'a'),
('diamond2.png', 'diamond', '2'),
('diamond2.png', 'diamond', '2'),
('diamond3.png', 'diamond', '3'),
('diamond3.png', 'diamond', '3'),
('diamond4.png', 'diamond', '4'),
('diamond4.png', 'diamond', '4'),
('diamond5.png', 'diamond', '5'),
('diamond5.png', 'diamond', '5'),
('diamond6.png', 'diamond', '6'),
('diamond6.png', 'diamond', '6'),
('diamond7.png', 'diamond', '7'),
('diamond7.png', 'diamond', '7'),
('diamond8.png', 'diamond', '8'),
('diamond8.png', 'diamond', '8'),
('all.png', 'all', 'all'),
('diamond9.png', 'diamond', '9'),
('diamond9.png', 'diamond', '9'),
('diamond10.png', 'diamond', '10'),
('diamond10.png', 'diamond', '10'),
('diamondk.png', 'diamond', 'k'),
('diamondk.png', 'diamond', 'k'),
('diamondq.png', 'diamond', 'q'),
('diamondq.png', 'diamond', 'q'),
('all.png', 'all', 'all'),
('joker1.png', 'joker', 'j1'),
('joker1.png', 'joker', 'j1'),
('joker2.png', 'joker', 'j2'),
('joker2.png', 'joker', 'j2');


-- 게임 맵 정보 
create table sequence_maps(
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    idx INT,
    card_id INT
);

INSERT INTO sequence_maps (idx, card_id)
VALUES
 (1, 1), (2, 2), (3, 3), (4, 4), (5, 5), (6, 6), (7, 7), (8, 8), (9, 9),(10, 10), 
(11, 11), (12, 12), (13, 13), (14, 14), (15, 15), (16, 16), (17, 17), (18, 18), (19, 19), (20, 20), 
(21, 21), (22, 22), (23, 23), (24, 24), (25, 25), (26, 26), (27, 27), (28, 28), (29, 29), (30, 30), 
(31, 31), (32, 32), (33, 33), (34, 34), (35, 35), (36, 36), (37, 37), (38, 38), (39, 39), (40, 40), 
(41, 41), (42, 42), (43, 43), (44, 44), (45, 45), (46, 46), (47, 47), (48, 48), (49, 49), (50, 50), 
(51, 51), (52, 52), (53, 53), (54, 54), (55, 55), (56, 56), (57, 57), (58, 58), (59, 59), (60, 60), 
(61, 61), (62, 62), (63, 63), (64, 64), (65, 65), (66, 66), (67, 67), (68, 68), (69, 69), (70, 70),
(71, 71), (72, 72), (73, 73), (74, 74), (75, 75), (76, 76), (77, 77), (78, 78), (79, 79), (80, 80),
(81, 81), (82, 82), (83, 83), (84, 84), (85, 85), (86, 86), (87, 87), (88, 88), (89, 89), (90, 90),
(91, 91), (92, 92), (93, 93), (94, 94), (95, 95), (96, 96), (97, 97), (98, 98), (99, 99), (100, 100);


-- 게임 유저 정보
create table sequence_users(
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    user_id INT,
    room_id INT,
    turn INT, -- 플레이 순서 
    color_type INT, -- 0: 빨강, 1: 파랑
    FOREIGN KEY (user_id) REFERENCES game_users(id) ON DELETE CASCADE,
    FOREIGN KEY (room_id) REFERENCES game_rooms(id) ON DELETE CASCADE
);

create table sequence_game_room_settings(
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    room_id INT, 
    timer INT, -- 게임 타이머
    current_round INT,  -- 현재 라운드
    FOREIGN KEY (room_id) REFERENCES game_rooms(id) ON DELETE CASCADE
);

create table sequence_room_cards(
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    room_id INT,
    card_id INT,
    user_id INT,
    state varchar(100), -- 남은 카드, 소윺 카드, 버려진 카드
    FOREIGN KEY (room_id) REFERENCES game_rooms(id) ON DELETE CASCADE
);

-- 게임 맵 정보
create table sequence_room_maps(
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    room_id INT,
    user_id INT,
    map_id INT, -- (1 ~ 100)
    FOREIGN KEY (room_id) REFERENCES game_rooms(id) ON DELETE CASCADE
);