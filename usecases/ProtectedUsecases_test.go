package usecases_test

import (
	"testing"

	"github.com/adamluzsi/FeatureFlags/services/rollouts"
	"github.com/adamluzsi/FeatureFlags/services/security"
	"github.com/adamluzsi/FeatureFlags/usecases"
	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"

	. "github.com/adamluzsi/FeatureFlags/testing"
)

func TestUseCases_ProtectedUsecases(t *testing.T) {

	s := testcase.NewSpec(t)
	SetupSpecCommonVariables(s)
	SetupSpec(s)
	s.Parallel()

	subject := func(t *testcase.T) (*usecases.ProtectedUsecases, error) {
		return GetUseCases(t).ProtectedUsecases(t.I(`TokenString`).(string))
	}

	s.When(`token doesn't exist`, func(s *testcase.Spec) {
		s.Let(`TokenString`, func(t *testcase.T) interface{} {
			return `The answer is 42`
		})

		s.Before(func(t *testcase.T) {
			require.Nil(t, GetStorage(t).Truncate(security.Token{}))
		})

		s.Then(`we receive back error`, func(t *testcase.T) {
			protectedUsecases, err := subject(t)
			require.Nil(t, protectedUsecases)
			require.Equal(t, usecases.ErrInvalidToken, err)
		})
	})

	s.When(`token exist`, func(s *testcase.Spec) {
		s.Let(`TokenString`, func(t *testcase.T) interface{} { return CreateToken(t, `manager`).Token })

		s.Then(`protected use-cases returned`, func(t *testcase.T) {
			pu, err := subject(t)
			require.Nil(t, err)
			require.NotNil(t, pu)

			// var explicit type check creates direct reference
			// which in testing equal to say, the object behaves the same as the
			// type requirement.
			var (
				_ *rollouts.RolloutManager = pu.RolloutManager
				_ *security.Doorkeeper     = pu.Doorkeeper
				_ *security.Issuer         = pu.Issuer
			)
		})

	})
}