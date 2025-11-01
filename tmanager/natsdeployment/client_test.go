//  Copyright (c) 2025 Metaform Systems, Inc
//
//  This program and the accompanying materials are made available under the
//  terms of the Apache License, Version 2.0 which is available at
//  https://www.apache.org/licenses/LICENSE-2.0
//
//  SPDX-License-Identifier: Apache-2.0
//
//  Contributors:
//       Metaform Systems, Inc. - initial API and implementation
//

package natsdeployment

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/natsclient"
	"github.com/metaform/connector-fabric-manager/common/natstestfixtures"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/common/types"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testTimeout  = 30 * time.Second
	streamName   = "cfm-deployment"
	cfmBucker    = "cfm-bucket"
	waitDuration = 300 * time.Millisecond
	tickDuration = 5 * time.Millisecond
)

func TestNatsDeploymentClient_Deploy(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	nt, err := natstestfixtures.SetupNatsContainer(ctx, cfmBucker)
	require.NoError(t, err)

	defer natstestfixtures.TeardownNatsContainer(ctx, nt)

	stream := natstestfixtures.SetupTestStream(t, ctx, nt.Client, streamName)
	natstestfixtures.SetupTestConsumer(t, ctx, stream, natsclient.CFMDeployment)

	msgClient := natsclient.NewMsgClient(nt.Client)
	dispatcher := &testDeploymentDispatcher{}

	client := newNatsDeploymentClient(msgClient, dispatcher, system.NoopMonitor{})

	manifest := model.DeploymentManifest{
		ID:             "test-deployment-123",
		DeploymentType: model.VPADeploymentType,
		Payload:        make(map[string]any),
	}

	// Send the manifest
	err = client.Send(ctx, manifest)
	require.NoError(t, err)

	// Verify the message was published by consuming it
	consumer, err := stream.Consumer(ctx, natsclient.CFMDeployment)
	require.NoError(t, err)

	messageBatch, err := consumer.Fetch(1, jetstream.FetchMaxWait(time.Second))
	require.NoError(t, err)

	messageFound := false
	for message := range messageBatch.Messages() {
		messageFound = true

		// Verify the message contains the deployment manifest
		var receivedManifest model.DeploymentManifest
		err = json.Unmarshal(message.Data(), &receivedManifest)
		require.NoError(t, err)

		assert.Equal(t, manifest.ID, receivedManifest.ID)
		break
	}
	assert.True(t, messageFound, "Should have received a deployment message")
}

func TestNatsDeploymentClient_ProcessMessage_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	nt, err := natstestfixtures.SetupNatsContainer(ctx, "cfm-bucket")
	require.NoError(t, err)

	defer natstestfixtures.TeardownNatsContainer(ctx, nt)

	stream := natstestfixtures.SetupTestStream(t, ctx, nt.Client, streamName)
	consumer := natstestfixtures.SetupTestConsumer(t, ctx, stream, natsclient.CFMDeployment)

	// Setup dispatcher with expectations
	dispatcher := &testDeploymentDispatcher{
		responses: make(chan model.DeploymentResponse, 1),
	}

	msgClient := natsclient.NewMsgClient(nt.Client)
	client := newNatsDeploymentClient(msgClient, dispatcher, system.NoopMonitor{})

	err = client.Init(ctx, consumer)
	require.NoError(t, err)

	// Create and publish the deployment response
	response := model.DeploymentResponse{
		ID:             "test-deployment-response-123",
		Success:        true,
		ManifestID:     "manifest-456",
		CorrelationID:  "test-correlation-id",
		DeploymentType: model.VPADeploymentType,
		Properties:     map[string]any{"test": "value"},
	}

	payload, err := json.Marshal(response)
	require.NoError(t, err)

	_, err = nt.Client.JetStream.Publish(ctx, natsclient.CFMDeploymentSubject, payload)
	require.NoError(t, err)

	// Verify the message was processed
	select {
	case receivedResponse := <-dispatcher.responses:
		assert.Equal(t, response.ID, receivedResponse.ID)
		assert.Equal(t, response.Success, receivedResponse.Success)
		assert.Equal(t, response.ManifestID, receivedResponse.ManifestID)
		assert.Equal(t, response.DeploymentType, receivedResponse.DeploymentType)
		assert.Equal(t, response.Properties, receivedResponse.Properties)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for deployment response")
	}
}

func TestNatsDeploymentClient_ProcessMessage_RecoverableError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	nt, err := natstestfixtures.SetupNatsContainer(ctx, "cfm-bucket")
	require.NoError(t, err)

	defer natstestfixtures.TeardownNatsContainer(ctx, nt)

	stream := natstestfixtures.SetupTestStream(t, ctx, nt.Client, streamName)
	consumer := natstestfixtures.SetupTestConsumer(t, ctx, stream, natsclient.CFMDeployment)

	// Setup dispatcher that returns recoverable error
	dispatcher := &testDeploymentDispatcher{
		responses:     make(chan model.DeploymentResponse, 1),
		shouldError:   true,
		errorToReturn: types.NewRecoverableError("test recoverable error"),
	}

	msgClient := natsclient.NewMsgClient(nt.Client)
	client := newNatsDeploymentClient(msgClient, dispatcher, system.NoopMonitor{})
	err = client.Init(ctx, consumer)
	require.NoError(t, err)

	// Create and publish the deployment response
	response := model.DeploymentResponse{
		ID:             "test-deployment-response-456",
		Success:        false,
		ErrorDetail:    "deployment failed",
		ManifestID:     "manifest-789",
		CorrelationID:  "test-correlation-id",
		DeploymentType: model.VPADeploymentType,
		Properties:     map[string]any{},
	}

	payload, err := json.Marshal(response)
	require.NoError(t, err)

	_, err = nt.Client.JetStream.Publish(ctx, natsclient.CFMDeploymentSubject, payload)
	require.NoError(t, err)

	// Verify the message was processed (should be NAKed due to recoverable error)
	select {
	case receivedResponse := <-dispatcher.responses:
		assert.Equal(t, response.ID, receivedResponse.ID)
		assert.Equal(t, response.Success, receivedResponse.Success)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for deployment response")
	}
}

func TestNatsDeploymentClient_ProcessMessage_FatalError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	nt, err := natstestfixtures.SetupNatsContainer(ctx, "cfm-bucket")
	require.NoError(t, err)

	defer natstestfixtures.TeardownNatsContainer(ctx, nt)

	stream := natstestfixtures.SetupTestStream(t, ctx, nt.Client, streamName)
	consumer := natstestfixtures.SetupTestConsumer(t, ctx, stream, natsclient.CFMDeployment)

	// Setup dispatcher that returns fatal error
	dispatcher := &testDeploymentDispatcher{
		responses:     make(chan model.DeploymentResponse, 1),
		shouldError:   true,
		errorToReturn: types.NewFatalError("test fatal error"),
	}

	msgClient := natsclient.NewMsgClient(nt.Client)
	client := newNatsDeploymentClient(msgClient, dispatcher, system.NoopMonitor{})
	err = client.Init(ctx, consumer)
	require.NoError(t, err)

	response := model.DeploymentResponse{
		ID:             "test-deployment-response-789",
		Success:        false,
		ErrorDetail:    "fatal deployment error",
		ManifestID:     "manifest-999",
		CorrelationID:  "test-correlation-id",
		DeploymentType: model.VPADeploymentType,
		Properties:     map[string]any{},
	}

	payload, err := json.Marshal(response)
	require.NoError(t, err)

	_, err = nt.Client.JetStream.Publish(ctx, natsclient.CFMDeploymentSubject, payload)
	require.NoError(t, err)

	// Verify the message was processed (should be ACKed despite fatal error)
	select {
	case receivedResponse := <-dispatcher.responses:
		assert.Equal(t, response.ID, receivedResponse.ID)
		assert.Equal(t, response.Success, receivedResponse.Success)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for deployment response")
	}
}

func TestNatsDeploymentClient_ProcessLoop_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	nt, err := natstestfixtures.SetupNatsContainer(ctx, "cfm-bucket")
	require.NoError(t, err)

	defer natstestfixtures.TeardownNatsContainer(ctx, nt)

	stream := natstestfixtures.SetupTestStream(t, ctx, nt.Client, streamName)
	consumer := natstestfixtures.SetupTestConsumer(t, ctx, stream, natsclient.CFMDeployment)

	dispatcher := &testDeploymentDispatcher{
		responses: make(chan model.DeploymentResponse, 1),
	}

	msgClient := natsclient.NewMsgClient(nt.Client)
	client := newNatsDeploymentClient(msgClient, dispatcher, system.NoopMonitor{})

	// Create a context that can be cancelled
	shortCtx, shortCancel := context.WithCancel(context.Background())

	err = client.Init(shortCtx, consumer)
	require.NoError(t, err)

	// Cancel the context
	shortCancel()

	// Check if Processing finished
	assert.Eventually(t, func() bool {
		return !client.Processing.Load()
	}, waitDuration, tickDuration, "Processing should have stopped after context cancellation")
}

func TestNatsDeploymentClient_MultipleMessages(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	nt, err := natstestfixtures.SetupNatsContainer(ctx, "cfm-bucket")
	require.NoError(t, err)

	defer natstestfixtures.TeardownNatsContainer(ctx, nt)

	stream := natstestfixtures.SetupTestStream(t, ctx, nt.Client, streamName)
	consumer := natstestfixtures.SetupTestConsumer(t, ctx, stream, natsclient.CFMDeployment)

	const messageCount = 5
	dispatcher := &testDeploymentDispatcher{
		responses: make(chan model.DeploymentResponse, messageCount),
	}

	msgClient := natsclient.NewMsgClient(nt.Client)
	client := newNatsDeploymentClient(msgClient, dispatcher, system.NoopMonitor{})

	err = client.Init(ctx, consumer)
	require.NoError(t, err)

	// Publish multiple messages
	var expectedResponses []model.DeploymentResponse
	for i := 0; i < messageCount; i++ {
		response := model.DeploymentResponse{
			ID:             fmt.Sprintf("test-deployment-response-%d", i),
			Success:        true,
			ErrorDetail:    "",
			ManifestID:     fmt.Sprintf("manifest-%d", i),
			CorrelationID:  "test-correlation-id",
			DeploymentType: model.VPADeploymentType,
			Properties:     map[string]any{"index": float64(i)},
		}
		expectedResponses = append(expectedResponses, response)

		payload, err := json.Marshal(response)
		require.NoError(t, err)

		_, err = nt.Client.JetStream.Publish(ctx, natsclient.CFMDeploymentSubject, payload)
		require.NoError(t, err)
	}

	// Collect all received responses
	var receivedResponses []model.DeploymentResponse
	for i := 0; i < messageCount; i++ {
		select {
		case response := <-dispatcher.responses:
			receivedResponses = append(receivedResponses, response)
		case <-time.After(5 * time.Second):
			t.Fatalf("Timeout waiting for response %d", i)
		}
	}

	// Verify all messages were processed
	assert.Len(t, receivedResponses, messageCount)

	// Verify each expected response was received (order may vary)
	for _, expected := range expectedResponses {
		found := false
		for _, received := range receivedResponses {
			if received.ID == expected.ID {
				assert.Equal(t, expected.ManifestID, received.ManifestID)
				assert.Equal(t, expected.Success, received.Success)
				assert.Equal(t, expected.Properties, received.Properties)
				found = true
				break
			}
		}
		assert.True(t, found, "Expected response %s not found", expected.ID)
	}
}

func TestNatsDeploymentClient_ProcessMessage_InvalidJSON(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	nt, err := natstestfixtures.SetupNatsContainer(ctx, "cfm-bucket")
	require.NoError(t, err)

	defer natstestfixtures.TeardownNatsContainer(ctx, nt)

	stream := natstestfixtures.SetupTestStream(t, ctx, nt.Client, streamName)
	consumer := natstestfixtures.SetupTestConsumer(t, ctx, stream, natsclient.CFMDeployment)

	// Setup dispatcher that should NOT be called as the test sends invalid JSON
	dispatcher := &testDeploymentDispatcher{
		onDispatch: func(ctx context.Context, response model.DeploymentResponse) error {
			t.Error("Dispatcher should not be called for invalid JSON")
			return nil
		},
	}

	msgClient := natsclient.NewMsgClient(nt.Client)
	client := newNatsDeploymentClient(msgClient, dispatcher, system.NoopMonitor{})
	err = client.Init(ctx, consumer)
	require.NoError(t, err)

	// Get initial NATS consumer info to track message Processing
	initialInfo, err := consumer.Info(ctx)
	require.NoError(t, err)
	initialAckCount := initialInfo.AckFloor.Consumer

	// Publish the invalid message
	invalidJSON := []byte(`{"invalid": json}`)
	_, err = nt.Client.JetStream.Publish(ctx, natsclient.CFMDeploymentSubject, invalidJSON)
	require.NoError(t, err)

	// Wait for message Processing and verify it was ACKed
	assert.Eventually(t, func() bool {
		info, err := consumer.Info(ctx)
		if err != nil {
			return false
		}
		// Check if the message was acknowledged (processed)
		return info.AckFloor.Consumer > initialAckCount
	}, waitDuration, tickDuration, "Invalid message should be acknowledged")

	// Verify that no more messages are pending
	finalInfo, err := consumer.Info(ctx)
	require.NoError(t, err)
	assert.Equal(t, finalInfo.NumPending, uint64(0), "No messages should be pending after Processing invalid message")
}

func TestNatsDeploymentClient_ProcessMessage_DispatcherSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Set up NATS container
	nt, err := natstestfixtures.SetupNatsContainer(ctx, "cfm-bucket")
	require.NoError(t, err)

	defer natstestfixtures.TeardownNatsContainer(ctx, nt)

	stream := natstestfixtures.SetupTestStream(t, ctx, nt.Client, streamName)
	consumer := natstestfixtures.SetupTestConsumer(t, ctx, stream, natsclient.CFMDeployment)

	// Track successful Processing
	var processedCount int
	var mu sync.Mutex

	dispatcher := &testDeploymentDispatcher{
		onDispatch: func(ctx context.Context, response model.DeploymentResponse) error {
			mu.Lock()
			processedCount++
			mu.Unlock()
			return nil // Success
		},
	}

	msgClient := natsclient.NewMsgClient(nt.Client)
	client := newNatsDeploymentClient(msgClient, dispatcher, system.NoopMonitor{})

	// Initialize client with consumer
	err = client.Init(ctx, consumer)
	require.NoError(t, err)

	// Create and publish deployment response message
	response := model.DeploymentResponse{
		ID:             "test-success-response",
		Success:        true,
		ErrorDetail:    "",
		ManifestID:     "success-manifest",
		CorrelationID:  "test-correlation-id",
		DeploymentType: model.VPADeploymentType,
		Properties:     map[string]any{"status": "success"},
	}

	payload, err := json.Marshal(response)
	require.NoError(t, err)

	_, err = nt.Client.JetStream.Publish(ctx, natsclient.CFMDeploymentSubject, payload)
	require.NoError(t, err)

	// Wait for Processing
	assert.Eventually(t, func() bool {
		mu.Lock()
		count := processedCount
		mu.Unlock()
		return count == 1
	}, waitDuration, tickDuration, "Message should be processed successfully")
}

// testDeploymentDispatcher implements api.deploymentCallbackDispatcher for testing
type testDeploymentDispatcher struct {
	responses     chan model.DeploymentResponse
	shouldError   bool
	errorToReturn error
	onDispatch    func(ctx context.Context, response model.DeploymentResponse) error
	mu            sync.Mutex
}

func (t *testDeploymentDispatcher) Dispatch(ctx context.Context, response model.DeploymentResponse) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.onDispatch != nil {
		return t.onDispatch(ctx, response)
	}

	if t.responses != nil {
		select {
		case t.responses <- response:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if t.shouldError {
		return t.errorToReturn
	}

	return nil
}
