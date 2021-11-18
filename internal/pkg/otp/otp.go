package otp

import (
	"fmt"
	"math/rand"
	"time"
)

type Otp struct {
	Secret   string    `json:"secret"`   //rappresent the secret key
	Duration int       `json:"duration"` //rappresent the duration of the otp
	Creation time.Time `json:"creation"` //rappresent the creation time of the otp
	User     string    `json:"user"`     //rappresent the user which the otp code belongs to
}

type OtpHandler []Otp

func NewOtpHandler() OtpHandler {
	return make(OtpHandler, 0)
}

func generateSecret() string {
	//generate a random number from 100000 to 999999
	return fmt.Sprintf("%06d", rand.Intn(899999)+100000)
}

func (h OtpHandler) CreateNew(user string) string {
	secret := generateSecret()
	otp := Otp{
		Secret:   secret,
		Duration: 60,
		Creation: time.Now(),
		User:     user,
	}
	h = append(h, otp)
	return secret
}

func (h OtpHandler) CheckOtp(user, secret string) (bool, error) {
	otp, err := h.get(user)
	if err != nil {
		return false, err
	}

	if !otp.validate(secret) {
		return false, fmt.Errorf("wrong secret")
	}

	if !otp.isStillValid() {
		return false, fmt.Errorf("otp expired")
	}

	h.remove(otp)

	return true, nil
}

func (h *OtpHandler) remove(otp Otp) {
	for i, o := range *h {
		if o.Secret == otp.Secret {
			*h = append((*h)[:i], (*h)[i+1:]...)
			return
		}
	}
}

func (h OtpHandler) get(user string) (Otp, error) {
	for _, otp := range h {
		if otp.User == user {
			return otp, nil
		}
	}
	return Otp{}, fmt.Errorf("user doesnt have and otp active")
}

func (o Otp) validate(code string) bool {
	return o.Secret == code
}

func (o Otp) isStillValid() bool {
	return o.Duration > 0
}
