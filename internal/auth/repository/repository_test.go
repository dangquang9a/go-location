package repository

import (
	"context"
	"dangquang9a/go-location/internal/auth"
	"dangquang9a/go-location/internal/models"
	"database/sql"
	"regexp"
	"testing"

	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository auth.UserRepository
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	// s.DB, err = gorm.Open("postgres", db)
	s.DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	require.NoError(s.T(), err)

	// s.DB.LogMode(true)

	s.repository = NewUserRepository(s.DB)
}

func (s *Suite) TestGetUserByUsername() {
	var (
		id       = uuid.New().String()
		username = "test"
		password = "testpass"
	)
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE "users"."username" = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).
			AddRow(id, username, password))

	res, err := s.repository.GetUserByUsername(context.Background(), username)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(&models.User{Id: id, Username: username, Password: password}, res))
}

func (s *Suite) TestGetUserByUserId() {
	var (
		id       = uuid.New().String()
		username = "test"
		password = "testpass"
	)
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).
			AddRow(id, username, password))

	res, err := s.repository.GetUserById(context.Background(), id)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(&models.User{Id: id, Username: username, Password: password}, res))
}

func (s *Suite) TestCreateUser() {
	var (
		id       = uuid.New().String()
		username = "test"
		password = "testpass"
	)
	user := &models.User{
		Id:       id,
		Username: username,
		Password: password,
	}

	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.
		QuoteMeta(`INSERT INTO "users" ("id","username","password") VALUES ($1,$2,$3)`)).
		WithArgs(id, username, password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectCommit()

	err := s.repository.CreateUser(context.Background(), user)
	require.NoError(s.T(), err)
}
