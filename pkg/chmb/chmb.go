package chmb

import (
	"context"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type ChMessageBusConfig struct {
	BufferSize int `env:"BUFFER_SIZE" env-default:"128"`
}

func NewChMessageBusConfig() (ChMessageBusConfig, error) {
	var cfg ChMessageBusConfig
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return ChMessageBusConfig{}, err
	}

	return cfg, nil
}

type ChMessageBus[T any] struct {
	ch   chan T
	once sync.Once
}

func NewChMessageBus[T any]() (*ChMessageBus[T], error) {
	var cfg ChMessageBusConfig
	cfg, err := NewChMessageBusConfig()
	if err != nil {
		return nil, err
	}

	return &ChMessageBus[T]{
		ch: make(chan T, cfg.BufferSize),
	}, nil
}

func (c *ChMessageBus[T]) Send(
	ctx context.Context,
	msg T,
) error {
	return c.send(ctx, []T{msg})
}

func (c *ChMessageBus[T]) SendMany(
	ctx context.Context,
	msgs []T,
) error {
	return c.send(ctx, msgs)
}

func (c *ChMessageBus[T]) Recv(
	_ context.Context,
) (<-chan T, error) {
	return c.ch, nil
}

func (c *ChMessageBus[T]) Close(
	_ context.Context,
) error {
	c.once.Do(func() {
		close(c.ch)
	})

	return nil
}

func (c *ChMessageBus[T]) send(
	ctx context.Context,
	msgs []T,
) error {
	for _, msg := range msgs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case c.ch <- msg:
		}
	}

	return nil
}
