package otp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type Otp struct {
	Secret   string    `json:"secret"`   //rappresent the secret key
	Duration int       `json:"duration"` //rappresent the duration of the otp
	Creation time.Time `json:"creation"` //rappresent the creation time of the otp
	User     string    `json:"user"`     //rappresent the user which the otp code belongs to
}

type OtpHandler []Otp

func NewOtpHandler(fileName string) (OtpHandler, error) {
	handler := make(OtpHandler, 0)

	//check if fileName exists
	//if not create a file with the name fileName
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		//create file
		file, err := os.Create(fileName)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		file.WriteString("[]")
	}

	//read file
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	//create map
	err = json.Unmarshal(file, &handler)
	if err != nil {
		return nil, err
	}
	return handler, nil
}

func (h OtpHandler) PrintAllOtps() {
	for _, otp := range h {
		fmt.Println(otp)
	}
}

func generateSecret() string {
	rand.Seed(time.Now().UnixNano())
	//generate a random number from 100000 to 999999
	return fmt.Sprintf("%06d", rand.Intn(899999)+100000)
}

func (h OtpHandler) CreateNew(user, filename string) string {
	secret := generateSecret()
	otp := Otp{
		Secret:   secret,
		Duration: 60,
		Creation: time.Now(),
		User:     user,
	}
	h = append(h, otp)
	h.saveOnFile(filename)
	return secret
}

func (h OtpHandler) CheckOtp(user, secret, filename string) error {
	otp, err := h.get(user)
	if err != nil {
		return err
	}

	if !otp.validate(secret) {
		return fmt.Errorf("wrong secret")
	}

	if !otp.isStillValid() {
		return fmt.Errorf("otp expired")
	}

	h.remove(otp, filename)

	return nil
}

func (h *OtpHandler) remove(otp Otp, filename string) {
	for i, o := range *h {
		if o.Secret == otp.Secret {
			*h = append((*h)[:i], (*h)[i+1:]...)
			h.saveOnFile(filename)
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
	//difference between now time and creation time
	diff := time.Since(o.Creation)
	//if difference is less than the duration of the otp
	//the otp is still valid
	return diff.Seconds() < float64(o.Duration)
}

//save the map on file
func (h OtpHandler) saveOnFile(path string) error {
	jsonByte, err := json.Marshal(&h)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(jsonByte)
	if err != nil {
		return err
	}
	return nil
}
