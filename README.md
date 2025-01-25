[![Coverage Status](https://coveralls.io/repos/github/JokerTrickster/frog_mahjong_game/badge.svg)](https://coveralls.io/github/JokerTrickster/frog_mahjong_game)

# 도메인
[https://www.frog-mahjong.xyz/](https://www.frog-mahjong.xyz/)

---

## 회의 일지
[회의 일지 링크](https://github.com/JokerTrickster/frog_mahjong_game/wiki/%ED%9A%8C%EC%9D%98-%EC%9D%BC%EC%A7%80)

---

# 토리덱 보드게임

## 프로젝트 설명
미션에 맞는 카드를 상대방보다 빠르게 수집하여 미션을 달성하는 카드 수집 게임

### 게임 화면

| 매칭 화면 | 카드 수집 화면 | 미션 달성 화면 | 결과 화면 |
|-----------|----------------|----------------|-----------|
| <img src="https://github.com/user-attachments/assets/03b12b4f-a216-4aec-8ce0-ae855208f72f" alt="매칭 화면" width="150"/> | <img src="https://github.com/user-attachments/assets/fed6603f-7d46-428a-ac84-4ef7ae115f66" alt="카드 수집 화면" width="150"/> | <img src="https://github.com/user-attachments/assets/a9073b44-5987-4fa0-9376-7e84ff6d2af5" alt="미션 달성 화면" width="150"/> | <img src="https://github.com/user-attachments/assets/92e773f2-988d-4796-8816-91805d6b615e" alt="결과 화면" width="150"/> |

### MVP 기능
1. 빠른 매칭, 함께하기  
2. 카드 가져오기  
3. 카드 버리기  
4. 미션 달성  
5. 아이템 사용 (오픈 카드 교체)  
6. 튜토리얼 화면  

---

# 개굴작 보드게임

## 프로젝트 설명
참새작 보드게임을 모티브로 턴제 보드게임

### 게임 화면

| 매칭 화면 | 게임 화면 | 결과 화면 |
|-----------|-----------|-----------|
| <img src="https://github.com/user-attachments/assets/03b12b4f-a216-4aec-8ce0-ae855208f72f" alt="매칭 화면" width="150"/> | <img src="https://github.com/user-attachments/assets/6539470b-12b9-4857-ac4c-6a9afe8d45ab" alt="게임 화면" width="150"/> | <img src="https://github.com/user-attachments/assets/2a8238e2-61a4-42d1-a23d-3b6a9d0b3583" alt="결과 화면" width="150"/> |

### MVP 기능
1. 인증 기능  
2. 매칭 기능  
3. 함께하기 기능 & 방 생성(인증코드 발급)  
4. 게임 플레이  
    - 패 가져오기  
    - 패 버리기  
    - 도라 선택하기  
    - 쯔모 외치기  
    - 론 외치기  
5. 신고하기 기능 (이메일 및 대시보드 연동)  

---

## 아키텍처

### DB 스키마
<img src="https://github.com/user-attachments/assets/9da7f847-e047-4d50-886d-cb52f06bc54c" alt="DB 스키마" width="700"/>

### 인프라 구성도
<img src="https://github.com/JokerTrickster/frog_mahjong_game/assets/140731661/74245fc3-d3cb-4d06-a9c3-022ec4514c8f" alt="인프라 구성도" width="700"/>

---

## 역할 및 기술
- **백엔드**: 조현준 / Golang, AWS, MySQL, Redis  
- **프론트**: 이다익 / JavaScript  

---

## 프로젝트 문서들
[문서 링크](https://drive.google.com/drive/folders/1km1pTM_KVxDrc0HCSJ-pdL-DT5wblw1D?usp=drive_link)
