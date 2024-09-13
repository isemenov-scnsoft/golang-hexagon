package rmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"golang-hexagon/internal/adapter/config"
	"golang-hexagon/internal/core/domain"
	"golang-hexagon/internal/core/port"
	"log/slog"
)

// message types
const (
	msgTypeLogin  = "login"
	msgTypeSignup = "signup"
	msgTypeUpdate = "update"
	msgTypeDelete = "delete"
	msgTypeList   = "list"
)

const connFormat = "amqp://%s:%s@%s:%s/%s"

// MessageHandler is a RabbitMQ message service
type (
	MessageHandler struct {
		authSvc port.AuthService
		userSvc port.UserService
		conf    *config.Container
		conn    *amqp.Connection
		ch      *amqp.Channel
	}

	msg struct {
		Type     string           `json:"type"`
		Name     *string          `json:"name"`
		Email    *string          `json:"email"`
		Password *[]byte          `json:"password"`
		Role     *domain.UserRole `json:"role"`
		UID      *uint64          `json:"uid"`
		Token    *string          `json:"token"`
		Offset   *uint64          `json:"offset"`
		Limit    *uint64          `json:"limit"`
	}
)

// New creates a new RabbitMQ message service
func New(conf *config.Container, authSvc port.AuthService, userSvc port.UserService) *MessageHandler {
	connection, err := amqp.Dial(fmt.Sprintf(connFormat, conf.RMQ.User, conf.RMQ.Password, conf.RMQ.Host, conf.RMQ.Port, conf.RMQ.Vhost))
	if err != nil {
		slog.Error("Error connecting to RabbitMQ instance", err)
		panic(err)
	}
	slog.Info("Successfully connected to RabbitMQ instance!")

	channel, err := connection.Channel()
	if err != nil {
		slog.Error("Error opening channel", err)
		panic(err)
	}

	return &MessageHandler{
		authSvc: authSvc,
		userSvc: userSvc,
		conf:    conf,
		conn:    connection,
		ch:      channel,
	}
}

// Consume consumes messages from queue
func (r *MessageHandler) Consume(ctx context.Context) {
	defer func() {
		if err := r.conn.Close(); err != nil {
			slog.Error("Error closing connection", err)
		}
		slog.Info("Config Connection closed...")
	}()

	// declaring consumer with its properties over channel opened
	msgs, err := r.ch.Consume(
		r.conf.RMQ.InQueue,     // queue
		r.conf.RMQ.ConsumerTag, // consumer
		false,                  // auto ack
		false,                  // exclusive
		false,                  // no local
		false,                  // no wait
		nil,                    //args
	)
	if err != nil {
		slog.Error("Error consuming messages", err)
		panic(err)
	}
	gorutineCtx, cancel := context.WithCancel(ctx)

	go func(ctx context.Context) {
		select {
		case msg := <-msgs:
			r.processMessage(msg)

		case <-ctx.Done():
			if err := r.ch.Cancel(r.conf.RMQ.ConsumerTag, true); err != nil {
				slog.Error("Error canceling consumer", err)
			}
			slog.Info("Consumer cancelled. Stopping...")
			return
		}
	}(gorutineCtx)

	select {
	case <-ctx.Done():
		cancel()
		if err := r.ch.Close(); err != nil {
			slog.Error("Error closing channel", err)
		}
		slog.Info("Channel closed. Stopping...")
	default: // do nothing
	}
}

// processMessage processes message
func (r *MessageHandler) processMessage(delivery amqp.Delivery) {
	var (
		m       msg
		message []byte
		err     error
		u       *domain.User
		us      []*domain.User
	)
	if err = json.Unmarshal(delivery.Body, &m); err != nil {
		slog.Error("Error unmarshalling delivery", err)
		return
	}
	ctx := context.Background()
	switch m.Type {
	case msgTypeLogin:
		message, err = r.authSvc.Login(ctx, asVal(m.Email), string(asVal(m.Password)))
	case msgTypeSignup:
		user := toUser(&m)
		u, err = r.userSvc.Register(ctx, user)
		if u != nil {
			message, _ = json.Marshal(u)
		}
	case msgTypeUpdate:
		user := toUser(&m)
		u, err = r.userSvc.UpdateUser(ctx, user)
		if u != nil {
			message, _ = json.Marshal(u)
		}
	case msgTypeDelete:
		err = r.userSvc.DeleteUser(ctx, asVal(m.UID))
	case msgTypeList:
		us, err = r.userSvc.ListUsers(ctx, asVal(m.Offset), asVal(m.Limit))
		if us != nil {
			message, _ = json.Marshal(us)
		}
	}

	if err = r.sendMessage(newResponseMessage(string(message), err)); err != nil {
		slog.Error("Error sending message", err)
	}
	if err = delivery.Ack(false); err != nil {
		slog.Error("Error acknowledging message", err)
	}
}

// sendMessage sends message back to responses queue
func (r *MessageHandler) sendMessage(msg *amqp.Publishing) error {
	err := r.ch.Publish(
		r.conf.RMQ.OutExchange,   // exchange
		r.conf.RMQ.OutRoutingKey, // routing key
		false,                    // mandatory
		false,                    // immediate
		asVal(msg),               // message
	)
	if err != nil {
		return err
	} else {
		return nil
	}
}
