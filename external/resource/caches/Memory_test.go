package caches_test

import (
	csh "github.com/adamluzsi/frameless/contracts"
	"github.com/adamluzsi/testcase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/toggler-io/toggler/domains/release"
	"github.com/toggler-io/toggler/domains/security"
	"github.com/toggler-io/toggler/domains/toggler"
	"github.com/toggler-io/toggler/domains/toggler/contracts"
	"github.com/toggler-io/toggler/external/resource/caches"
	"github.com/toggler-io/toggler/external/resource/storages"
	sh "github.com/toggler-io/toggler/spechelper"
	"testing"
)

var _ release.Storage = &caches.Memory{}

func TestMemory_smoke(t *testing.T) {
	storage := storages.NewEventLogMemoryStorage()
	storage.EventLog.Options.DisableAsyncSubscriptionHandling = true
	m, err := caches.NewMemory(storage)
	require.Nil(t, err)
	t.Cleanup(func() { require.Nil(t, m.Close()) })
	ff := sh.DefaultFixtureFactory
	ctx := ff.Context()
	ts := m.SecurityToken(ctx)
	token := ff.Create(security.Token{}).(*security.Token)
	csh.CreateEntity(t, ts, ctx, token)
	token.OwnerUID = uuid.New().String()
	csh.UpdateEntity(t, ts, ctx, token)
}

func TestMemory(t *testing.T)      { SpecMemory(t) }
func BenchmarkMemory(b *testing.B) { SpecMemory(b) }

func SpecMemory(tb testing.TB) {
	testcase.RunContract(tb, contracts.Storage{
		Subject: func(tb testing.TB) toggler.Storage {
			storage := storages.NewEventLogMemoryStorage()
			storage.EventLog.Options.DisableAsyncSubscriptionHandling = true
			m, err := caches.NewMemory(storage)
			require.Nil(tb, err)
			tb.Cleanup(func() { require.Nil(tb, m.Close()) })
			return m
		},
		FixtureFactory: sh.DefaultFixtureFactory,
	})
}
