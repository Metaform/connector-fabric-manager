package natsorchestration

import (
	"context"
	"testing"

	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/stretchr/testify/assert"
)

func TestActivityContext_Delete(t *testing.T) {
	activity := api.Activity{ID: "test-activity"}
	activityContext := newActivityContext(context.TODO(), "test-oid", activity, map[string]any{}, map[string]any{})

	// Set a value
	activityContext.SetValue("key", "value")

	// Verify it exists
	_, exists := activityContext.Value("key")
	assert.True(t, exists)

	// Delete the key
	activityContext.Delete("key")

	// Verify it's deleted
	_, exists = activityContext.Value("key")
	assert.False(t, exists)
}

