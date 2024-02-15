package subscriber

import "github.com/gofiber/fiber/v2/log"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
