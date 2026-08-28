package main

import "testing"

func TestShouldSendVerificationOnlyForUnverifiedDonors(t *testing.T) {
	cases := []struct {
		name  string
		donor Donor
		want  bool
	}{
		{"new donor", Donor{Email: "donor@example.org"}, true},
		{"verified donor", Donor{Email: "donor@example.org", EmailVerified: true}, false},
		{"missing email", Donor{EmailVerified: false}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ShouldSendVerification(tc.donor); got != tc.want {
				t.Fatalf("ShouldSendVerification() = %v, want %v", got, tc.want)
			}
		})
	}
}
