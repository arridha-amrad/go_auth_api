package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"my-go-api/internal/models"
	"my-go-api/pkg/database"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	repo IUserRepository
	ids  []uuid.UUID
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	db, err := database.Connect("postgres://user_go_api_test:pg_pwd_go_api_test@localhost:5432/pg_db_go_api_test?sslmode=disable", "5m", 50, 25)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.db = db
	suite.repo = NewUserRepository(db)
	suite.db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	suite.db.Exec(`CREATE TYPE providers AS ENUM ('credentials', 'google')`)
	suite.db.Exec(`CREATE TYPE user_roles AS ENUM ('user', 'admin')`)
	// Create the users table
	_, err = suite.db.Exec(`
		CREATE TABLE
			users (
				id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
				username VARCHAR(50) UNIQUE NOT NULL,
				name VARCHAR(100) NOT NULL,
				email VARCHAR(100) UNIQUE NOT NULL,
				password TEXT,
				provider providers DEFAULT 'credentials',
				role user_roles DEFAULT 'user',
				jwt_version VARCHAR(20) NOT NULL,
				is_verified BOOLEAN NOT NULL DEFAULT false,
				created_at TIMESTAMP(0)
				WITH
					TIME ZONE NOT NULL DEFAULT NOW (),
					updated_at TIMESTAMP(0)
				WITH
					TIME ZONE NOT NULL DEFAULT NOW ()
			);
	`)
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
	// Drop the tokens table
	_, err := suite.db.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		suite.T().Fatal(err)
	}
	_, err = suite.db.Exec(`DROP EXTENSION IF EXISTS "uuid-ossp"`)
	if err != nil {
		suite.T().Fatal(err)
	}
	_, err = suite.db.Exec("DROP TYPE IF EXISTS providers")
	if err != nil {
		suite.T().Fatal(err)
	}
	_, err = suite.db.Exec("DROP TYPE IF EXISTS user_roles")
	if err != nil {
		suite.T().Fatal(err)
	}

	// Close the database connection
	suite.db.Close()
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	suite.ids = []uuid.UUID{}
	if _, err := suite.db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE"); err != nil {
		suite.T().Fatal(err)
	}
	users := []CreateOneParams{
		{Name: "dummy", Username: "dummy00", Email: "dummy@mail.com", Password: "pwd123", JWTVersion: "jwt_123_version"},
		{Name: "john", Username: "john00", Email: "john@mail.com", Password: "pwd123", JWTVersion: "jwt_123_version"},
		{Name: "jane", Username: "jane00", Email: "jane@mail.com", Password: "pwd123", JWTVersion: "jwt_123_version"},
	}
	query := fmt.Sprintf(`
		INSERT INTO users (
			name, 
			username, 
			email, 
			password, 
			jwt_version
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING %s
	`, userSelectedFields)
	for _, user := range users {
		model := &models.User{}
		if err := suite.db.QueryRowContext(context.Background(), query,
			user.Name,
			user.Username,
			user.Email,
			user.Password,
			user.JWTVersion,
		).Scan(scanUser(model)...); err != nil {
			suite.T().Fatal(err)
		}
		suite.ids = append(suite.ids, model.ID)
	}
}

func (suite *UserRepositoryTestSuite) TestGetAll() {
	suite.Run("It should get all users", func() {
		users, err := suite.repo.GetAll(context.Background())
		assert.NoError(suite.T(), err)
		assert.NotEmpty(suite.T(), users)
	})
	suite.Run("It should return empty slice if no users exist", func() {
		if _, err := suite.db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE"); err != nil {
			suite.T().Fatal(err)
		}
		users, err := suite.repo.GetAll(context.Background())
		assert.NoError(suite.T(), err)
		assert.Empty(suite.T(), users)
	})
}

func (suite *UserRepositoryTestSuite) TestCreateOne() {
	// insert params
	name := "ari"
	username := "ari08"
	email := "ari@mail.com"
	password := "12345"
	jwtVersion := "jwt-123-version"
	suite.Run("it should create new user", func() {
		// insert action
		newUser, err := suite.repo.CreateOne(context.Background(), CreateOneParams{
			Name:       name,
			Username:   username,
			Email:      email,
			Password:   password,
			JWTVersion: jwtVersion,
		})
		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), newUser)
		// verify returned data
		assert.Equal(suite.T(), name, newUser.Name)
		assert.Equal(suite.T(), username, newUser.Username)
		assert.Equal(suite.T(), email, newUser.Email)
		assert.Equal(suite.T(), password, newUser.Password)
		assert.Equal(suite.T(), "user", newUser.Role)
		assert.Equal(suite.T(), "credentials", newUser.Provider)
		// verify the new user in inserted into database
		var dbUser models.User
		err = suite.db.QueryRow(`SELECT id, username, name, email, password, provider, role, created_at, updated_at
														 FROM users
														 WHERE email = $1`, newUser.Email).
			Scan(&dbUser.ID, &dbUser.Username, &dbUser.Name, &dbUser.Email, &dbUser.Password, &dbUser.Provider, &dbUser.Role, &dbUser.CreatedAt, &dbUser.UpdatedAt)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), newUser.Username, dbUser.Username)
		assert.Equal(suite.T(), password, dbUser.Password)
	})
	suite.Run("it should fail, because duplicate email", func() {
		_, err := suite.repo.CreateOne(context.Background(), CreateOneParams{
			Name:       "vxcvx",
			Username:   "zxpopo",
			Email:      "ari@mail.com",
			Password:   "xcvxvx",
			JWTVersion: "xcvx",
		})
		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "duplicate")
	})
	suite.Run("it should fail, because duplicate username", func() {
		_, err := suite.repo.CreateOne(context.Background(), CreateOneParams{
			Name:       "vxcvx",
			Username:   "ari08",
			Email:      "asdasdxc@xcvxc.x",
			Password:   "xcvxvx",
			JWTVersion: "xcvx",
		})
		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "duplicate")
	})
}

func (suite *UserRepositoryTestSuite) TestGetOne() {
	suite.Run("It should find a user by username", func() {
		username := "john00"
		result, err := suite.repo.GetOne(context.Background(), GetOneParams{
			Username: &username,
		})
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), username, result.Username)
	})
	suite.Run("It should not find any user by username", func() {
		username := "mulyono"
		result, err := suite.repo.GetOne(context.Background(), GetOneParams{
			Username: &username,
		})
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), result)
	})
	suite.Run("It should find a user by email", func() {
		email := "john@mail.com"
		result, err := suite.repo.GetOne(context.Background(), GetOneParams{
			Email: &email,
		})
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), email, result.Email)
	})
	suite.Run("It should not find any user by email", func() {
		email := "nonexistent@example.com"
		result, err := suite.repo.GetOne(context.Background(), GetOneParams{
			Email: &email,
		})
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), result)
	})
	suite.Run("It should return error when no parameters are given", func() {
		result, err := suite.repo.GetOne(context.Background(), GetOneParams{})
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), result)
	})
	suite.Run("It should find a user by id", func() {
		for _, id := range suite.ids {
			result, err := suite.repo.GetOne(context.Background(), GetOneParams{
				Id: &id,
			})
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), result.ID, id)
		}
	})
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
