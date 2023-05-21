package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func EnvMySQL() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("MYSQL")
}

func EnvPort() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
		return "3000"
	}

	return os.Getenv("PORT")
}

func StripeApiKey() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("STRIPE_API_KEY")
}

func StripePublishableKey() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("STRIPE_PUBLIC_KEY")
}

func GetAmount() int {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	amountStr := os.Getenv("AMOUNT")

	amount, err := strconv.Atoi(amountStr)

	if err != nil {
		log.Fatal("Error converting Env Amount to Int")
	}

	return amount
}
