package spechelper

import (
	"fmt"

	"github.com/adamluzsi/frameless/fixtures"
	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"

	"github.com/toggler-io/toggler/domains/release"
)

const (
	LetVarExampleReleaseManualPilotEnrollment = `ExampleReleaseManualPilotEnrollment`

	LetVarExampleReleaseRollout  = `example release rollout`
	LetVarExampleReleaseFlag     = `example release flag`
	LetVarExamplePilotExternalID = `PilotExternalID`
	LetVarExamplePilot           = `ManualPilot`
	LetVarExamplePilotEnrollment = `PilotEnrollment`
)

func init() {
	setups = append(setups, func(s *testcase.Spec) {
		s.Let(LetVarExampleReleaseManualPilotEnrollment, func(t *testcase.T) interface{} {
			mpe := Create(release.Pilot{}).(*release.Pilot)
			mpe.FlagID = ExampleReleaseFlag(t).ID
			mpe.EnvironmentID = ExampleDeploymentEnvironment(t).ID
			mpe.PublicID = ExampleExternalPilotID(t)
			storae := StorageGet(t).ReleasePilot(ContextGet(t))
			require.Nil(t, storae.Create(ContextGet(t), mpe))
			t.Defer(storae.DeleteByID, ContextGet(t), mpe.ID)
			return mpe
		})

		s.Let(LetVarExamplePilotExternalID, func(t *testcase.T) interface{} {
			return fixtures.Random.StringN(100)
		})

		s.Let(LetVarExamplePilotEnrollment, func(t *testcase.T) interface{} {
			return fixtures.Random.Bool()
		})

		s.Let(LetVarExamplePilot, func(t *testcase.T) interface{} {
			// domains/release/specs/FlagFinder.go:53:1: DEPRECATED, clean it up
			return &release.Pilot{
				FlagID:          ExampleReleaseFlag(t).ID,
				EnvironmentID:   ExampleDeploymentEnvironment(t).ID,
				PublicID:        t.I(LetVarExamplePilotExternalID).(string),
				IsParticipating: t.I(LetVarExamplePilotEnrollment).(bool),
			}
		})

		GivenWeHaveReleaseFlag(s, LetVarExampleReleaseFlag)

		GivenWeHaveReleaseRollout(s,
			LetVarExampleReleaseRollout,
			LetVarExampleReleaseFlag,
			LetVarExampleDeploymentEnvironment,
		)
	})
}

func ExampleReleaseManualPilotEnrollment(t *testcase.T) *release.Pilot {
	return t.I(LetVarExampleReleaseManualPilotEnrollment).(*release.Pilot)
}

func ExampleExternalPilotID(t *testcase.T) string {
	return t.I(LetVarExamplePilotExternalID).(string)
}

func FindStoredReleaseFlagByName(t *testcase.T, name string) *release.Flag {
	f, err := StorageGet(t).ReleaseFlag(ContextGet(t)).FindReleaseFlagByName(ContextGet(t), name)
	require.Nil(t, err)
	require.NotNil(t, f)
	return f
}

func ExampleReleaseRollout(t *testcase.T) *release.Rollout {
	return GetReleaseRollout(t, LetVarExampleReleaseRollout)
}

func getReleaseRolloutPlanLetVar(vn string) string {
	return fmt.Sprintf(`%s.plan`, vn)
}

func GetReleaseRolloutPlan(t *testcase.T, rolloutLVN string) release.RolloutDefinition {
	return t.I(getReleaseRolloutPlanLetVar(rolloutLVN)).(release.RolloutDefinition)
}

func GetReleaseRollout(t *testcase.T, vn string) *release.Rollout {
	return t.I(vn).(*release.Rollout)
}

func GivenWeHaveReleaseRollout(s *testcase.Spec, vn, flagLVN, envLVN string) {
	s.Let(getReleaseRolloutPlanLetVar(vn), func(t *testcase.T) interface{} {
		return *Create(release.RolloutDecisionByPercentage{}).(*release.RolloutDecisionByPercentage)
	})

	s.Let(vn, func(t *testcase.T) interface{} {
		rf := GetReleaseFlag(t, flagLVN)
		de := GetDeploymentEnvironment(t, envLVN)

		rollout := FixtureFactory{}.Create(release.Rollout{}).(*release.Rollout)
		rollout.FlagID = rf.ID
		rollout.DeploymentEnvironmentID = de.ID
		rollout.Plan = GetReleaseRolloutPlan(t, vn)
		require.Nil(t, rollout.Plan.Validate())

		// TODO: replace when rollout manager has function for this
		storage := StorageGet(t).ReleaseRollout(ContextGet(t))
		require.Nil(t, storage.Create(ContextGet(t), rollout))
		t.Defer(storage.DeleteByID, ContextGet(t), rollout.ID)
		t.Logf(`%#v`, rollout)
		return rollout
	})
}

func GivenWeHaveReleaseFlag(s *testcase.Spec, vn string) {
	s.Let(vn, func(t *testcase.T) interface{} {
		rf := FixtureFactory{}.Create(release.Flag{}).(*release.Flag)
		rf.Name = fmt.Sprintf(`%s - %s`, vn, rf.Name)
		storage := StorageGet(t).ReleaseFlag(ContextGet(t))
		require.Nil(t, storage.Create(ContextGet(t), rf))
		t.Defer(ExampleRolloutManager(t).DeleteFeatureFlag, ContextGet(t), rf.ID)
		t.Defer(storage.DeleteByID, ContextGet(t), rf.ID)
		t.Logf(`%#v`, rf)
		return rf
	})
}

func AndReleaseFlagRolloutPercentageIs(s *testcase.Spec, rolloutLVN string, percentage int) {
	s.Before(func(t *testcase.T) {
		rollout := GetReleaseRollout(t, LetVarExampleReleaseRollout)
		byPercentage, ok := rollout.Plan.(release.RolloutDecisionByPercentage)
		require.True(t, ok, `unexpected release rollout plan definition for AndReleaseFlagRolloutPercentageIs helper`)
		byPercentage.Percentage = percentage
		t.Logf(`and the release rollout percentage is set to %d`, percentage)
		//
		// please note that this will eager load the rollout value in the testing framework
		// as it makes no sense to have a percentage set to something that doesn't even exists.
		//
		// And in case if we already initialized such context where rollout entry exists,
		// we need to update its rollout plan as well.
		t.Set(getReleaseRolloutPlanLetVar(rolloutLVN), byPercentage)
		rollout.Plan = GetReleaseRolloutPlan(t, LetVarExampleReleaseRollout)
		require.Nil(t, StorageGet(t).ReleaseRollout(ContextGet(t)).Update(ContextGet(t), rollout))
	})
}

func GetReleaseFlag(t *testcase.T, lvn string) *release.Flag {
	ff := t.I(lvn)
	if ff == nil {
		return nil
	}
	return ff.(*release.Flag)
}

func ExampleReleaseFlag(t *testcase.T) *release.Flag {
	return GetReleaseFlag(t, LetVarExampleReleaseFlag)
}

func ExampleRolloutManager(t *testcase.T) *release.RolloutManager {
	return release.NewRolloutManager(StorageGet(t))
}

func SpecPilotEnrolmentIs(t *testcase.T, enrollment bool) {
	if ExampleReleaseFlag(t).ID == `` {
		require.Nil(t, StorageGet(t).ReleaseFlag(ContextGet(t)).Create(ContextGet(t), ExampleReleaseFlag(t)))
	}

	rm := release.NewRolloutManager(StorageGet(t))
	require.Nil(t, rm.SetPilotEnrollmentForFeature(ContextGet(t),
		ExampleReleaseFlag(t).ID,
		ExampleDeploymentEnvironment(t).ID,
		ExampleExternalPilotID(t),
		enrollment))
}

func NoReleaseFlagPresentInTheStorage(s *testcase.Spec) {
	s.Before(func(t *testcase.T) {
		// TODO: replace with flag manager list+delete
		require.Nil(t, StorageGet(t).ReleaseFlag(ContextGet(t)).DeleteAll(ContextGet(t)))
	})
}

func NoReleaseRolloutPresentInTheStorage(s *testcase.Spec) {
	s.Before(func(t *testcase.T) {
		// TODO: replace with rollout manager list+delete
		require.Nil(t, StorageGet(t).ReleaseRollout(ContextGet(t)).DeleteAll(ContextGet(t)))
	})
}

func AndExamplePilotManualParticipatingIsSetTo(s *testcase.Spec, isParticipating bool) {
	s.Before(func(t *testcase.T) {
		require.Nil(t, ExampleRolloutManager(t).SetPilotEnrollmentForFeature(
			ContextGet(t),
			ExampleReleaseFlag(t).ID,
			ExampleDeploymentEnvironment(t).ID,
			ExampleExternalPilotID(t),
			isParticipating,
		))

		t.Defer(ExampleRolloutManager(t).UnsetPilotEnrollmentForFeature,
			ContextGet(t),
			ExampleReleaseFlag(t).ID,
			ExampleDeploymentEnvironment(t).ID,
			ExampleExternalPilotID(t),
		)
	})
}
