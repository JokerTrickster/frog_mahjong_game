package main

import (
	dto "main/utils/db/mysql"

	"gorm.io/gen"
)

type Querier interface {
	// SELECT * FROM @@table WHERE id=@id
	GetByID(id int) (gen.T, error) // GetByID query data by id and return it as *struct*

	// GetByRoles query data by roles and return it as *slice of pointer*
	//   (The below blank line is required to comment for the generated method)
	//
	// SELECT * FROM @@table WHERE role IN @rolesName
	GetByRoles(rolesName ...string) ([]*gen.T, error)

	// InsertValue insert value
	//
	// INSERT INTO @@table (name, age) VALUES (@name, @age)
	InsertValue(name string, age int) error
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(dto.GormMysqlDB) // reuse your gorm db

	// Generate the code
	g.ApplyBasic(dto.Users{})
	// Generate the code
	g.ApplyInterface(func(Querier) {}, dto.Users{}, dto.Rooms{}, dto.RoomUsers{}, dto.Cards{}, dto.Chats{}, dto.Tokens{})

	g.Execute()
}
