package testing

import (
	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"

	"github.com/toggler-io/toggler/domains/deployment"
)

const LetVarExampleDeploymentEnvironment = `example deployment environment`

func init() {
	setups = append(setups, func(s *testcase.Spec) {
		GivenWeHaveDeploymentEnvironment(s, LetVarExampleDeploymentEnvironment)
	})
}

func GetDeploymentEnvironment(t *testcase.T, vn string) *deployment.Environment {
	return t.I(vn).(*deployment.Environment)
}

func ExampleDeploymentEnvironment(t *testcase.T) *deployment.Environment {
	return GetDeploymentEnvironment(t, LetVarExampleDeploymentEnvironment)
}

func GivenWeHaveDeploymentEnvironment(s *testcase.Spec, vn string) {
	s.Let(vn, func(t *testcase.T) interface{} {
		de := FixtureFactory{}.Create(deployment.Environment{}).(*deployment.Environment)
		require.Nil(t, ExampleStorage(t).Create(GetContext(t), de))
		t.Defer(ExampleStorage(t).DeleteByID, GetContext(t), *de, de.ID)
		t.Logf(`%#v`, de)
		return de
	})
}

func NoDeploymentEnvironmentPresentInTheStorage(s *testcase.Spec) {
	s.Before(func(t *testcase.T) {
		require.Nil(t, ExampleStorage(t).DeleteAll(GetContext(t), deployment.Environment{}))
	})
}