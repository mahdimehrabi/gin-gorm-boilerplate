package tests

import (
	"boilerplate/api/controllers"
	"boilerplate/api/middlewares"
	"boilerplate/api/repositories"
	"boilerplate/api/routes"
	"boilerplate/api/services"
	"boilerplate/infrastructure"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type TestSuiteEnv struct {
	suite.Suite
	router   infrastructure.Router
	database infrastructure.Database
}

func NewTestSuiteEnv(router infrastructure.Router, database infrastructure.Database,
	logger infrastructure.Logger, migrations infrastructure.Migrations) TestSuiteEnv {
	suite := new(suite.Suite)
	migrations.Migrate()
	return TestSuiteEnv{
		*suite,
		router,
		database,
	}
}

// Tests are run before they start
func (suite *TestSuiteEnv) SetupSuite() {

}

// Running after each test
func (suite *TestSuiteEnv) TearDownTest() {
	// database.ClearTable()
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
			fx.Provide(NewTestSuiteEnv),
			fx.Supply(t),
			fx.Invoke(Setup),
		),
	).Done()
}

func Setup(t *testing.T, tse TestSuiteEnv, lc fx.Lifecycle,
	routes routes.Routes, middlewares middlewares.Middlewares) {
	routes.Setup()
	middlewares.Setup()
	suite.Run(t, &tse)
}
