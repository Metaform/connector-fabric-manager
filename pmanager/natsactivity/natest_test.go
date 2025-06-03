package natsactivity

import (
	"context"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewJetStreamPubSub(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	nt, err := setupNatsContainer(ctx, "test-bucket")
	require.NoError(t, err)
	//goland:noinspection GoUnhandledErrorResult
	defer teardownNatsContainer(ctx, nt)

	cfg := jetstream.StreamConfig{
		Name:      "test-activity",
		Retention: jetstream.WorkQueuePolicy,
		Subjects:  []string{"event.*"},
	}

	stream, err := nt.natsConnection.jetStream.CreateOrUpdateStream(ctx, cfg)
	require.NoError(t, err)

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "test-activity",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	require.NoError(t, err)

	_, err = nt.natsConnection.jetStream.Publish(ctx, "event.foo", []byte("Test message"))
	require.NoError(t, err)

	msg, err := consumer.Next(jetstream.FetchMaxWait(100 * time.Millisecond))

	require.NoError(t, err)
	require.NotNil(t, msg)

	cancel()
}
