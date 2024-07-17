package repository

import (
	"context"
	"testing"

	sql "main/utils/db/mysql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestFindOneUser(t *testing.T) {
	// Mock the database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Open a GORM DB connection using the mock DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Create a new repository instance
	repo := NewGetUsersRepository(gormDB)

	// Define the expected behavior of the SQL mock
	expectedUser := sql.Users{Model: gorm.Model{
		ID: 1,
	}, Name: "John Doe"}
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(expectedUser.ID, expectedUser.Name)

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE id = ?").
		WithArgs(expectedUser.ID).
		WillReturnRows(rows)

	// Execute the method under test
	ctx := context.Background()
	user, err := repo.FindOneUser(ctx, int(expectedUser.ID))

	// Validate the results
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}
