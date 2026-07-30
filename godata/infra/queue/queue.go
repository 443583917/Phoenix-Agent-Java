package queue

import (
    "github.com/phoenix-agent-go/internal/config"
    amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ(cfg *config.RabbitMQConfig) (*amqp.Connection, *amqp.Channel, error) {
    conn, err := amqp.Dial(cfg.Addr)
    if err != nil {
        return nil, nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, nil, err
    }

    if err := ch.Qos(cfg.PrefetchCount, 0, false); err != nil {
        ch.Close()
        conn.Close()
        return nil, nil, err
    }

    err = ch.ExchangeDeclare(
        cfg.Exchange,
        "topic",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        ch.Close()
        conn.Close()
        return nil, nil, err
    }

    return conn, ch, nil
}
