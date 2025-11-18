package vault

import (
	"context"
	"testing"

	"github.com/metaform/connector-fabric-manager/assembly/serviceapi"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestVaultServiceAssembly_Init(t *testing.T) {
	ctx := context.Background()
	result := setupTest(ctx, t)
	defer result.cleanup()

	assembly := &VaultServiceAssembly{}

	vConfig := viper.New()
	vConfig.Set(urlKey, result.url)
	vConfig.Set(roleIDKey, result.roleID)
	vConfig.Set(secretIDKey, result.secretID)

	ictx := &system.InitContext{
		StartContext: system.StartContext{
			Registry:   system.NewServiceRegistry(),
			LogMonitor: system.NoopMonitor{},
			Config:     vConfig,
			Mode:       system.DebugMode,
		},
	}
	err := assembly.Init(ictx)
	require.NoError(t, err)

	client := ictx.Registry.Resolve(serviceapi.VaultKey).(serviceapi.VaultClient)
	require.NotNil(t, client)

	err = client.StoreSecret(ctx, "test-secret", "test-value")
	require.NoError(t, err)

	val, err := client.ResolveSecret(ctx, "test-secret")
	require.NoError(t, err)
	require.Equal(t, "test-value", val, "Expected secret value to match")
}

type TestSetupResult struct {
	url      string
	roleID   string
	secretID string
	cleanup  func()
}

func setupTest(ctx context.Context, t *testing.T) TestSetupResult {
	containerResult, err := StartVaultContainer(ctx)
	require.NoError(t, err, "Failed to start Vault container")

	setupResult, err := SetupVault(containerResult.URL, containerResult.Token)
	if err != nil {
		containerResult.Cleanup()
		t.Fatalf("Failed to setup Vault: %v", err)
	}

	return TestSetupResult{
		url:      containerResult.URL,
		roleID:   setupResult.RoleID,
		secretID: setupResult.SecretID,
		cleanup:  containerResult.Cleanup,
	}
}
