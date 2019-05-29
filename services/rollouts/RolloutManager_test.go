package rollouts_test

import (
	"math/rand"
	"net/url"
	"testing"
	"time"

	"github.com/adamluzsi/FeatureFlags/services/rollouts"
	. "github.com/adamluzsi/FeatureFlags/testing"
	"github.com/adamluzsi/testcase"

	"github.com/adamluzsi/frameless/iterators"
	"github.com/stretchr/testify/require"
)

func TestRolloutManager(t *testing.T) {
	s := testcase.NewSpec(t)
	s.Parallel()
	SetupSpecCommonVariables(s)

	s.Let(`GeneratedRandomSeed`, func(t *testcase.T) interface{} {
		return time.Now().Unix()
	})

	s.Let(`RolloutManager`, func(t *testcase.T) interface{} {
		return &rollouts.RolloutManager{
			Storage: GetStorage(t),

			RandSeedGenerator: func() int64 {
				return GetGeneratedRandomSeed(t)
			},
		}
	})

	s.Before(func(t *testcase.T) {
		require.Nil(t, GetStorage(t).Truncate(rollouts.FeatureFlag{}))
		require.Nil(t, GetStorage(t).Truncate(rollouts.Pilot{}))
	})

	SpecRolloutManagerSetFeatureFlag(s)
	SpecRolloutManagerListFeatureFlags(s)
	SpecRolloutManagerSetPilotEnrollmentForFeature(s)
}

func SpecRolloutManagerSetFeatureFlag(s *testcase.Spec) {
	s.Describe(`SetFeatureFlagJSON`, func(s *testcase.Spec) {
		subjectWithArgs := func(t *testcase.T, f *rollouts.FeatureFlag) error {
			return manager(t).SetFeatureFlag(f)
		}

		subject := func(t *testcase.T) error {
			return subjectWithArgs(t, GetFeatureFlag(t))
		}

		s.Let(`FeatureFlagName`, func(t *testcase.T) interface{} { return ExampleFeatureName() })
		s.Let(`RolloutApiURL`, func(t *testcase.T) interface{} { return nil })
		s.Let(`RolloutPercentage`, func(t *testcase.T) interface{} { return rand.Intn(101) })
		s.Let(`RolloutSeedSalt`, func(t *testcase.T) interface{} { return int64(42) })

		s.Let(`FeatureFlag`, func(t *testcase.T) interface{} {
			ff := &rollouts.FeatureFlag{Name: t.I(`FeatureName`).(string)}
			ff.Rollout.RandSeedSalt = t.I(`RolloutSeedSalt`).(int64)
			ff.Rollout.Strategy.Percentage = t.I(`RolloutPercentage`).(int)
			ff.Rollout.Strategy.DecisionLogicAPI = GetRolloutApiURL(t)
			return ff
		})

		s.Then(`on valid input the values persisted`, func(t *testcase.T) {
			require.Nil(t, subject(t))
			require.NotNil(t, FindStoredFeatureFlag(t))
			require.Equal(t, GetFeatureFlag(t), FindStoredFeatureFlag(t))
		})

		s.When(`name is empty`, func(s *testcase.Spec) {
			s.Let(`FeatureName`, func(t *testcase.T) interface{} { return "" })

			s.Then(`it will fail with invalid feature name`, func(t *testcase.T) {
				require.Equal(t, rollouts.ErrInvalidFeatureName, subject(t))
			})
		})

		s.When(`url`, func(s *testcase.Spec) {
			s.Context(`is not a valid request url`, func(s *testcase.Spec) {
				s.Let(`RolloutApiURL`, func(t *testcase.T) interface{} { return `http//example.com` })

				s.Then(`it will fail with invalid url`, func(t *testcase.T) {
					require.Equal(t, rollouts.ErrInvalidURL, subject(t))
				})
			})

			s.Context(`is not defined or nil`, func(s *testcase.Spec) {
				s.Let(`RolloutApiURL`, func(t *testcase.T) interface{} { return nil })

				s.Then(`it will be saved and will represent that no custom domain decision url used`, func(t *testcase.T) {
					require.Nil(t, subject(t))

					require.Nil(t, FindStoredFeatureFlag(t).Rollout.Strategy.DecisionLogicAPI)
				})
			})

			s.Context(`is a valid request URL`, func(s *testcase.Spec) {
				s.Let(`RolloutApiURL`, func(t *testcase.T) interface{} { return `https://example.com` })

				s.Then(`it will persist for future usage`, func(t *testcase.T) {
					require.Nil(t, subject(t))

					require.Equal(t, `https://example.com`, FindStoredFeatureFlag(t).Rollout.Strategy.DecisionLogicAPI.String())
				})
			})
		})

		s.When(`percentage`, func(s *testcase.Spec) {
			s.Context(`less than 0`, func(s *testcase.Spec) {
				s.Let(`RolloutPercentage`, func(t *testcase.T) interface{} { return -1 + (rand.Intn(1024) * -1) })

				s.Then(`it will report error regarding the percentage`, func(t *testcase.T) {
					require.Equal(t, rollouts.ErrInvalidPercentage, subject(t))
				})
			})

			s.Context(`greater than 100`, func(s *testcase.Spec) {
				s.Let(`RolloutPercentage`, func(t *testcase.T) interface{} { return 101 + rand.Intn(1024) })

				s.Then(`it will report error regarding the percentage`, func(t *testcase.T) {
					require.Equal(t, rollouts.ErrInvalidPercentage, subject(t))
				})
			})

			s.Context(`is a number between 0 and 100`, func(s *testcase.Spec) {
				s.Let(`RolloutPercentage`, func(t *testcase.T) interface{} { return rand.Intn(101) })

				s.Then(`it will persist the received percentage`, func(t *testcase.T) {
					require.Nil(t, subject(t))

					require.Equal(t, t.I(`RolloutPercentage`).(int), FindStoredFeatureFlag(t).Rollout.Strategy.Percentage)
				})
			})
		})

		s.When(`pseudo random seed salt`, func(s *testcase.Spec) {
			s.Context(`is 0`, func(s *testcase.Spec) {
				s.Let(`RolloutSeedSalt`, func(t *testcase.T) interface{} { return int64(0) })

				s.And(`feature is not stored before`, func(s *testcase.Spec) {
					s.Before(func(t *testcase.T) {
						require.Nil(t, GetStorage(t).Truncate(rollouts.FeatureFlag{}))
					})

					s.Then(`random seed generator used for setting seed value`, func(t *testcase.T) {
						require.Nil(t, subject(t))

						require.Equal(t, GetGeneratedRandomSeed(t), FindStoredFeatureFlag(t).Rollout.RandSeedSalt)
					})
				})

				s.And(`feature is already stored before`, func(s *testcase.Spec) {
					s.Before(func(t *testcase.T) {
						require.Nil(t, subject(t))
					})

					s.Then(`random seed used from the persisted object`, func(t *testcase.T) {
						f := *GetFeatureFlag(t) // pass by value copy
						f.ID = ``
						f.Rollout.RandSeedSalt = 0
						require.Nil(t, subjectWithArgs(t, &f))
						require.Equal(t, GetGeneratedRandomSeed(t), FindStoredFeatureFlag(t).Rollout.RandSeedSalt)
					})
				})
			})

			s.Context(`something else`, func(s *testcase.Spec) {
				s.Let(`RolloutSeedSalt`, func(t *testcase.T) interface{} { return int64(12) })

				s.Then(`it will persist the value`, func(t *testcase.T) {
					require.Nil(t, subject(t))

					require.Equal(t, int64(12), FindStoredFeatureFlag(t).Rollout.RandSeedSalt)
				})
			})
		})

		s.When(`feature flag`, func(s *testcase.Spec) {
			s.Context(`is nil`, func(s *testcase.Spec) {
				s.Let(`FeatureFlag`, func(t *testcase.T) interface{} { return nil })

				s.Then(`it will return error about it`, func(t *testcase.T) {
					require.Error(t, subject(t))
				})
			})

			s.Context(`was not stored until now`, func(s *testcase.Spec) {
				s.Before(func(t *testcase.T) {
					require.Nil(t, GetStorage(t).Truncate(rollouts.FeatureFlag{}))
				})

				s.Then(`it will be persisted`, func(t *testcase.T) {
					require.Nil(t, subject(t))
					require.NotNil(t, FindStoredFeatureFlag(t))
					require.Equal(t, GetFeatureFlag(t), FindStoredFeatureFlag(t))
				})
			})

			s.Context(`had been persisted previously`, func(s *testcase.Spec) {
				s.Before(func(t *testcase.T) {
					require.Nil(t, subject(t))
				})

				s.Then(`latest values are persisted`, func(t *testcase.T) {
					flag := *GetFeatureFlag(t) // pass by value copy
					flag.Rollout.Strategy.Percentage = 42
					u, err := url.Parse(`https://example.com`)
					require.Nil(t, err)
					flag.Rollout.Strategy.DecisionLogicAPI = u
					require.Nil(t, subjectWithArgs(t, &flag))

					storedFlag := FindStoredFeatureFlag(t)
					require.Equal(t, u, storedFlag.Rollout.Strategy.DecisionLogicAPI)
					require.Equal(t, 42, storedFlag.Rollout.Strategy.Percentage)
				})
			})
		})

	})
}

func SpecRolloutManagerListFeatureFlags(s *testcase.Spec) {
	s.Describe(`ListFeatureFlags`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) ([]*rollouts.FeatureFlag, error) {
			return manager(t).ListFeatureFlags()
		}

		onSuccess := func(t *testcase.T) []*rollouts.FeatureFlag {
			ffs, err := subject(t)
			require.Nil(t, err)
			return ffs
		}

		s.When(`features are in the system`, func(s *testcase.Spec) {

			s.Before(func(t *testcase.T) {
				require.Nil(t, manager(t).SetFeatureFlag(&rollouts.FeatureFlag{Name:`a`}))
				require.Nil(t, manager(t).SetFeatureFlag(&rollouts.FeatureFlag{Name:`b`}))
				require.Nil(t, manager(t).SetFeatureFlag(&rollouts.FeatureFlag{Name:`c`}))
			})

			s.Then(`feature flags are returned`, func(t *testcase.T) {
				flags := onSuccess(t)

				expectedFlagNames := []string{`a`, `b`, `c`}

				for _, ff := range flags {
					require.Contains(t, expectedFlagNames, ff.Name)
				}
			})

		})

		s.When(`no feature present in the system`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				require.Nil(t, GetStorage(t).Truncate(rollouts.FeatureFlag{}))
			})

			s.Then(`feature flags are returned`, func(t *testcase.T) {
				flags := onSuccess(t)

				require.Equal(t, []*rollouts.FeatureFlag{}, flags)
			})
		})

	})
}

func SpecRolloutManagerSetPilotEnrollmentForFeature(s *testcase.Spec) {
	s.Describe(`SetPilotEnrollmentForFeature`, func(s *testcase.Spec) {

		GetNewEnrollment := func(t *testcase.T) bool {
			return t.I(`NewEnrollment`).(bool)
		}

		subject := func(t *testcase.T) error {
			return manager(t).SetPilotEnrollmentForFeature(
				GetFeatureFlagName(t),
				GetExternalPilotID(t),
				GetNewEnrollment(t),
			)
		}

		s.Let(`NewEnrollment`, func(t *testcase.T) interface{} {
			return rand.Intn(2) == 0
		})

		findFlag := func(t *testcase.T) *rollouts.FeatureFlag {
			iter := GetStorage(t).FindAll(&rollouts.FeatureFlag{})
			require.NotNil(t, iter)
			require.True(t, iter.Next())
			var ff rollouts.FeatureFlag
			require.Nil(t, iter.Decode(&ff))
			require.False(t, iter.Next())
			require.Nil(t, iter.Err())
			return &ff
		}

		s.When(`no feature flag is seen ever before`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				require.Nil(t, GetStorage(t).Truncate(rollouts.FeatureFlag{}))
			})

			s.Then(`feature flag created`, func(t *testcase.T) {
				require.Nil(t, subject(t))

				flag := findFlag(t)
				require.Equal(t, GetFeatureFlagName(t), flag.Name)
				require.Nil(t, flag.Rollout.Strategy.DecisionLogicAPI)
				require.Equal(t, 0, flag.Rollout.Strategy.Percentage)
				require.Equal(t, GetGeneratedRandomSeed(t), flag.Rollout.RandSeedSalt)
			})

			s.Then(`pilot is enrollment is set for the feature is set`, func(t *testcase.T) {
				require.Nil(t, subject(t))

				flag := findFlag(t)
				pilot, err := GetStorage(t).FindFlagPilotByExternalPilotID(flag.ID, GetExternalPilotID(t))
				require.Nil(t, err)
				require.NotNil(t, pilot)
				require.Equal(t, GetNewEnrollment(t), pilot.Enrolled)
				require.Equal(t, GetExternalPilotID(t), pilot.ExternalID)
			})
		})

		s.When(`feature flag already configured`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				require.Nil(t, GetStorage(t).Save(GetFeatureFlag(t)))
			})

			s.Then(`flag is will not be recreated`, func(t *testcase.T) {
				require.Nil(t, subject(t))

				count, err := iterators.Count(GetStorage(t).FindAll(rollouts.FeatureFlag{}))
				require.Nil(t, err)
				require.Equal(t, 1, count)

				flag := findFlag(t)
				require.Equal(t, GetFeatureFlag(t), flag)
			})

			s.And(`pilot already exists`, func(s *testcase.Spec) {
				s.Before(func(t *testcase.T) {
					require.Nil(t, GetStorage(t).Save(GetPilot(t)))
				})

				s.And(`and pilot is has the opposite enrollment status`, func(s *testcase.Spec) {
					s.Let(`PilotEnrollment`, func(t *testcase.T) interface{} {
						return !GetNewEnrollment(t)
					})

					s.Then(`the original pilot is updated to the new enrollment status`, func(t *testcase.T) {
						require.Nil(t, subject(t))
						flag := findFlag(t)

						pilot, err := GetStorage(t).FindFlagPilotByExternalPilotID(flag.ID, GetExternalPilotID(t))
						require.Nil(t, err)

						require.NotNil(t, pilot)
						require.Equal(t, GetNewEnrollment(t), pilot.Enrolled)
						require.Equal(t, GetExternalPilotID(t), pilot.ExternalID)
						require.Equal(t, GetPilot(t), pilot)

						count, err := iterators.Count(GetStorage(t).FindAll(rollouts.Pilot{}))
						require.Nil(t, err)
						require.Equal(t, 1, count)
					})
				})

				s.And(`pilot already has the same enrollment status`, func(s *testcase.Spec) {
					s.Let(`PilotEnrollment`, func(t *testcase.T) interface{} {
						return GetNewEnrollment(t)
					})

					s.Then(`pilot remain the same`, func(t *testcase.T) {

						require.Nil(t, subject(t))
						ff := findFlag(t)

						pilot, err := GetStorage(t).FindFlagPilotByExternalPilotID(ff.ID, GetExternalPilotID(t))
						require.Nil(t, err)

						require.NotNil(t, pilot)
						require.Equal(t, GetNewEnrollment(t), pilot.Enrolled)
						require.Equal(t, GetExternalPilotID(t), pilot.ExternalID)

						count, err := iterators.Count(GetStorage(t).FindAll(rollouts.Pilot{}))
						require.Nil(t, err)
						require.Equal(t, 1, count)

					})
				})
			})

		})
	})
}

func GetGeneratedRandomSeed(t *testcase.T) int64 {
	return t.I(`GeneratedRandomSeed`).(int64)
}

func manager(t *testcase.T) *rollouts.RolloutManager {
	return t.I(`RolloutManager`).(*rollouts.RolloutManager)
}
