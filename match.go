package main

import (
	"errors"
	"github.com/urfave/cli/v2"
	"gome/engine"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "match",
		Usage: "消费下单队列里的订单并撮合",
		Action: func(c *cli.Context) error {
			symbol := c.Args().Get(0)
			if symbol == "" {
				return errors.New("请输入需要消费的队列名称")
			}
			mq := engine.NewSimpleRabbitMQ(symbol)
			mq.ConsumeNewOrder()

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
