package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/yourusername/twitter-services-app/shared/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Producer interface {
	Publish(ctx context.Context, subject string, event Event) error
}

type natsProducer struct {
	js jetstream.JetStream
}

func NewProducer(js jetstream.JetStream) Producer {
	return &natsProducer{
		js: js,
	}
}

func (p *natsProducer) Publish(ctx context.Context, subject string, event Event) error {
	// Start Span
	tracer := otel.Tracer("nats-producer")
	ctx, span := tracer.Start(ctx, fmt.Sprintf("publish %s", subject), trace.WithAttributes(
		attribute.String("messaging.system", "nats"),
		attribute.String("messaging.destination", subject),
		attribute.String("messaging.destination_kind", "subject"),
	))
	defer span.End()

	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Create NATS message with headers
	msg := &natsMsg{
		subject: subject,
		data:    value,
		headers: make(map[string][]string),
	}

	msg.headers[HeaderCorrelationID] = []string{event.CorrelationID}
	msg.headers[EventTypeHeader] = []string{event.Type}

	// Propagate Trace Context using OpenTelemetry observability wrapper
	observability.InjectNATSHeaders(ctx, msg.headers)

	// Publish to Jetstream
	nMsg := &nats.Msg{
		Subject: subject,
		Data:    value,
		Header:  nats.Header(msg.headers),
	}

	_, err = p.js.PublishMsg(ctx, nMsg, jetstream.WithMsgID(event.AggregateID))

	return err
}

// Helper struct for bridging OTEL injector interface, similar to kafka.Message but for NATS Msg
type natsMsg struct {
	subject string
	data    []byte
	headers map[string][]string
}
