package mysql

type GormModel struct {
	ID        string `json:"id" gorm:"primaryKey;column:id""`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime;column:createdAt"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime;column:updatedAt"`
	IsDeleted bool   `json:"isDeleted" gorm:"default:false; column:isDeleted"`
}

type GormUserDTO struct {
	GormModel GormModel `gorm:"embedded"`
	Name      string    `json:"name"  gorm:"column:name"`
	Email     string    `json:"email"  gorm:"column:email"`
}

/*
 user 테이블 생성 기존꺼 일단 그대로 두고 임시로 작업
create table gorm_user_dtos (
	id varchar(200),
	name varchar(200),
	email varchar(200),
	isDeleted TINYINT(1) not null,
	createdAt int,
	updatedAt int,
    PRIMARY KEY (id)
	);
*/

type GormUserAuthDTO struct {
	GormModel  GormModel `gorm:"embedded"`
	Provider   string    `json:"provider" gorm:"column:provider"`
	UserID     string    `json:"userID" gorm:"column:userID"`
	LastSignIn int64     `json:"lastSignIn" gorm:"column:lastSignIn"`
	Password   string    `json:"password" gorm:"column:password"`
	Name       string    `json:"name" gorm"column:name"`
	Email      string    `json:"email" gorm"column:email"`
}

/*
 userAuth 테이블 생성 기존꺼 일단 그대로 두고 임시로 작업

  create table gorm_user_auth_dtos (
	id varchar(200),
	provider varchar(200) not null,
	userID varchar(200) not null,
    name varchar(200) not null,
    email varchar(200) not null,
	lastSignIn int,
	isDeleted TINYINT(1) not null,
	createdAt int,
	updatedAt int,
    PRIMARY KEY (id),
	FOREIGN KEY (userID)
	REFERENCES gorm_user_dtos(id) on update cascade
	);
*/
