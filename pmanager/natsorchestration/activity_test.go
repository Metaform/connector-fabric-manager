package natsorchestration

import (
	"context"
	"testing"

	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/stretchr/testify/assert"
)

func TestActivityContext_Delete(t *testing.T) {
	activity := api.Activity{ID: "test-activity"}
	activityContext := newActivityContext(context.TODO(), "test-oid", activity, map[string]any{}, map[string]any{}, map[string]any{})

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

func TestImmutableMap(t *testing.T) {
	data := createData()
	im := NewImmutableMap(data)

	value, exists := im.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)

	value, exists = im.Get("nonexistent")
	assert.False(t, exists)
	assert.Nil(t, value)
}

func TestImmutableMap_Keys(t *testing.T) {
	data := createData()
	im := NewImmutableMap(data)

	keys := im.Keys()
	assert.Len(t, keys, 2)
	assert.Contains(t, keys, "key1")
	assert.Contains(t, keys, "key2")
}

func TestImmutableMap_Size(t *testing.T) {
	data := createData()
	im := NewImmutableMap(data)

	assert.Equal(t, 2, im.Size())
}

func TestImmutableMap_Immutability(t *testing.T) {
	originalData := createData()
	im := NewImmutableMap(originalData)

	// Modify the original data
	originalData["key1"] = "modified"
	originalData["key3"] = "new"

	// ImmutableMap should not be affected
	value, exists := im.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)

	_, exists = im.Get("key3")
	assert.False(t, exists)

	assert.Equal(t, 2, im.Size())
}

func createData() map[string]any {
	originalData := map[string]any{
		"key1": "value1",
		"key2": "value2",
	}
	return originalData
}
