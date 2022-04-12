package tests

import (
	"boilerplate/apps/authApp/services"
	"boilerplate/core"
	"boilerplate/core/infrastructure"
	"boilerplate/core/validators"
	"testing"

	"github.com/joho/godotenv"
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
			core.RoutesModule,
			core.ControllerModule,
			core.SeviceModule,
			core.RepositoryModule,
			core.MiddlewaresModule,
			validators.Module,
			fx.Provide(NewTestSuiteEnv),
			fx.Supply(t),
			fx.Invoke(Setup),
		),
	).Done()
}

func Setup(t *testing.T, tse TestSuiteEnv, lc fx.Lifecycle,
	routes core.Routes, middlewares core.Middlewares, validators validators.Validators) {
	validators.Setup()
	routes.Setup()
	middlewares.Setup()
	suite.Run(t, &tse)
}
