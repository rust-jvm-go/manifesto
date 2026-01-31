package otp

import (
	"crypto/rand"
	"fmt"
	"time"
)

type OTPPurpose string

const (
	OTPPurposeJobApplication OTPPurpose = "JOB_APPLICATION"
	OTPPurposeVerification   OTPPurpose = "VERIFICATION"
)

type OTP struct {
	ID          string
	Contact     string // Email or phone
	Code        string
	Purpose     OTPPurpose
	ExpiresAt   time.Time
	VerifiedAt  *time.Time
	Attempts    int
	CreatedAt   time.Time
	MaxAttempts int
}

func (o *OTP) IsValid() bool {
	return time.Now().Before(o.ExpiresAt) && o.VerifiedAt == nil && o.Attempts < o.MaxAttempts
}

func (o *OTP) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

func (o *OTP) Verify() error {
	if o.Attempts >= o.MaxAttempts {
		return ErrTooManyAttempts()
	}
	if o.IsExpired() {
		return ErrOTPExpired()
	}
	if o.Attempts >= 5 {
		return ErrTooManyAttempts()
	}
	now := time.Now()
	o.VerifiedAt = &now
	return nil
}

func (o *OTP) IncrementAttempts() {
	o.Attempts++
}

func GenerateOTPCode(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	code := fmt.Sprintf("%0*d", length*2, int(bytes[0])<<16|int(bytes[1])<<8|int(bytes[2]))
	return code[:length*2][:6], nil // Return first 6 digits
}
