package tests

import (
	"boilerplate/api/controllers"
	"boilerplate/api/middlewares"
	"boilerplate/api/repositories"
	"boilerplate/api/routes"
	"boilerplate/api/services"
	"boilerplate/infrastructure"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type TestSuiteEnv struct {
	suite.Suite
	router infrastructure.Router
}

// Tests are run before they start
func (suite *TestSuiteEnv) SetupSuite(routes routes.Routes, middlewares middlewares.Middlewares) {
	routes.Setup()
	middlewares.Setup()
}

// Running after each test
func (suite *TestSuiteEnv) TearDownTest() {
	fmt.Println("teart down test")
	// database.ClearTable()
}

// Running after all tests are completed
func (suite *TestSuiteEnv) TearDownSuite() {
	fmt.Println("test suite tear down----")
	// os.Exit(0)
	// suite.db.Close()
}

// This gets run automatically by `go test` so we call `suite.Run` inside it
func TestSuite(t *testing.T) {
	// This is what actually runs our suite
	fx.New(
		fx.Options(
			infrastructure.Module,
			routes.Module,
			controllers.Module,
			services.Module,
			repositories.Module,
			middlewares.Module,
			fx.Supply(t),
			fx.Invoke(Setup),
		),
	).Done()
}

func Setup(t *testing.T, lc fx.Lifecycle) {
	suite.Run(t, new(TestSuiteEnv))
}
