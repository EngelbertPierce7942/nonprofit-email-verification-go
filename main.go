package main

import (
	"fmt"
	"os"
)

func main() {
	to := os.Getenv("DEMO_EMAIL_TO")
	if to == "" {
		fmt.Fprintln(os.Stderr, "set DEMO_EMAIL_TO to a recipient address")
		os.Exit(1)
	}
	client, err := newEmailClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	donor := Donor{ID: "donor-42", Name: "Ari", Email: to}
	result, err := SendVerification(client, donor, "https://nonprofit.example/verify/donor-42")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("verification email sent:", result.MessageID)
}
