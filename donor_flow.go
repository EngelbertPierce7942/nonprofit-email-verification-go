package main

import "fmt"

type Donor struct {
	ID            string
	Name          string
	Email         string
	EmailVerified bool
}

type Receipt struct {
	DonorID string
	Amount  int
}

type VolunteerReminder struct {
	DonorID string
	Shift   string
}

type CampaignReport struct {
	Campaign string
	Receipts int
	Total    int
}

func ShouldSendVerification(donor Donor) bool {
	return donor.Email != "" && !donor.EmailVerified
}

func verificationEmail(donor Donor, verificationURL string) map[string]string {
	return map[string]string{
		"to":      donor.Email,
		"subject": "Verify your nonprofit account",
		"html":    fmt.Sprintf("<p>Hello %s,</p><p><a href=\"%s\">Verify your email</a> to finish signing up.</p>", donor.Name, verificationURL),
	}
}

func SendVerification(client *emailClient, donor Donor, verificationURL string) (emailResult, error) {
	if !ShouldSendVerification(donor) {
		return emailResult{}, nil
	}
	return client.send(verificationEmail(donor, verificationURL), "donor-verification-"+donor.ID)
}
