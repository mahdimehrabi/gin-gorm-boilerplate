package tests

import (
	"boilerplate/api/controllers"
	"boilerplate/api/middlewares"
	"boilerplate/api/repositories"
	"boilerplate/api/routes"
	"boilerplate/api/services"
	"boilerplate/api/validators"
	"boilerplate/infrastructure"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type TestSuiteEnv struct {
	suite.Suite
	router      infrastructure.Router
	database    infrastructure.Database
	encryption  infrastructure.Encryption
	logger      infrastructure.Logger
	env         infrastructure.Env
	migrations  infrastructure.Migrations
	authService services.AuthService
}

func NewTestSuiteEnv(router infrastructure.Router, database infrastructure.Database,
	encryption infrastructure.Encryption, logger infrastructure.Logger,
	migrations infrastructure.Migrations, env infrastructure.Env,
	authService services.AuthService) TestSuiteEnv {
	suite := new(suite.Suite)
	migrations.Migrate()
	return TestSuiteEnv{
		*suite,
		router,
		database,
		encryption,
		logger,
		env,
		migrations,
		authService,
	}
}

// Tests are run before they start
func (suite *TestSuiteEnv) SetupSuite() {

}

// Running after each test
func (suite *TestSuiteEnv) TearDownTest() {
	err := suite.database.DB.Exec(`
	DO $$ DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP
        EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || '';
    END LOOP;
END $$;
	`).Error
	if err != nil {
		suite.logger.Zap.Error("Failed to truncate tables", err)
	}
}

// Running after all tests are completed
func (suite *TestSuiteEnv) TearDownSuite() {
	// os.Exit(0)
	// suite.db.Close()
}

// This gets run automatically by `go test` so we call `suite.Run` inside it
func TestSuite(t *testing.T) {
	err := godotenv.Overload("../.env.test")
	if err != nil {
		panic(err)
	}
	fx.New(
		fx.Options(
			infrastructure.Module,
			routes.Module,
			controllers.Module,
			services.Module,
			repositories.Module,
			middlewares.Module,
			validators.Module,
			fx.Provide(NewTestSuiteEnv),
			fx.Supply(t),
			fx.Invoke(Setup),
		),
	).Done()
}

func Setup(t *testing.T, tse TestSuiteEnv, lc fx.Lifecycle,
	routes routes.Routes, middlewares middlewares.Middlewares, validators validators.Validators) {
	validators.Setup()
	routes.Setup()
	middlewares.Setup()
	suite.Run(t, &tse)
}

type SendEmailMock struct {
	mock.Mock
}

func (m *SendEmailMock) SendEmail(from string, to string, subject string) error {
	var ch chan bool
	htmlFilePath := "/file1.txt"
	args := m.Called(ch, from, to, subject, htmlFilePath)
	return args.Error(0)
}
