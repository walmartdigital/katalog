package kafka

import (
	"github.com/sirupsen/logrus"
	"github.com/walmartdigital/katalog/src/utils"
)

var log = logrus.New()

func init() {
	err := utils.LogInit(log)
	if err != nil {
		log.Fatal(err)
	}
}

// Consumer ...
type Consumer struct {
}

// CreateConsumer ...
func CreateConsumer() Consumer {
	return Consumer{}
}
