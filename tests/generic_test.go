package tests

import "fmt"

func (suite *TestSuiteEnv) TestPing() {
	fmt.Println("Test Ping")
	a := suite.Assert()
	a.True(true)
}
