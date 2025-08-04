package utils

import "github.com/gofiber/fiber/v2/log"

func ExitOnErr(err error, errorMessage string, a... any) {
	if err != nil {
		log.Fatalf(errorMessage + "\n", a...)
	}
}
