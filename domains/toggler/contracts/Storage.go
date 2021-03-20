package contracts

import (
	"testing"

	"github.com/adamluzsi/testcase"
	"github.com/toggler-io/toggler/domains/deployment"
	"github.com/toggler-io/toggler/domains/release"
	"github.com/toggler-io/toggler/domains/security"

	deplspecs "github.com/toggler-io/toggler/domains/deployment/contracts"
	rollspecs "github.com/toggler-io/toggler/domains/release/contracts"
	secuspecs "github.com/toggler-io/toggler/domains/security/contracts"

	"github.com/toggler-io/toggler/domains/toggler"
	sh "github.com/toggler-io/toggler/spechelper"
)

type Storage struct {
	Subject        func(testing.TB) toggler.Storage
	FixtureFactory sh.FixtureFactory
}

func (spec Storage) Test(t *testing.T) {
	spec.Spec(t)
}

func (spec Storage) Benchmark(b *testing.B) {
	spec.Spec(b)
}

func (spec Storage) Spec(tb testing.TB) {
	s := testcase.NewSpec(tb)
	s.Describe(`toggler#Storage`, func(s *testcase.Spec) {
		testcase.RunContract(s,
			rollspecs.Storage{
				Subject: func(tb testing.TB) release.Storage {
					return spec.Subject(tb)
				},
				FixtureFactory: spec.FixtureFactory,
			},
			secuspecs.Storage{
				Subject: func(tb testing.TB) security.Storage {
					return spec.Subject(tb)
				},
				FixtureFactory: spec.FixtureFactory,
			},
			deplspecs.Storage{
				Subject: func(tb testing.TB) deployment.Storage {
					return spec.Subject(tb)
				},
				FixtureFactory: spec.FixtureFactory,
			},
		)
	})
}
