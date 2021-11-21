package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Content struct {
	Password string `json:"password"`
	PfpUrl   string `json:"pfpUrl"`
}

type Users map[string]Content

//create a new Users map
func NewUsers(fileName string) (Users, error) {
	//check if fileName exists
	//if not create a file with the name fileName
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		//create file
		file, err := os.Create(fileName)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		file.WriteString("{}")
	} else {
		var u Users
		//read file
		file, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, err
		}
		//create map
		err = json.Unmarshal(file, &u)
		if err != nil {
			return nil, err
		}
		return u, nil
	}
	return make(Users), nil
}

//add a new user to the map and save the map on file
func (u *Users) AddUser(email, password, path string) error {
	if u.ExistUser(email) {
		return fmt.Errorf("User already exists")
	}
	(*u)[email] = Content{password, ""}

	err := u.saveOnFile(path)
	if err != nil {
		return err
	}
	return nil
}

//update a user on the map and save the map on file
func (u *Users) UpdateUser(email, password, newPassword, PfpUrl, path string) error {
	if !u.ExistUser(email) {
		return fmt.Errorf("User does not exist")
	}
	if (*u)[email].Password != password {
		return fmt.Errorf("Incorrect password")
	}

	(*u)[email] = Content{newPassword, PfpUrl}
	err := u.saveOnFile(path)
	if err != nil {
		return err
	}
	return nil
}

//return the Content of the user given the email
func (u *Users) GetUser(email, password string) (Content, error) {
	if !u.ExistUser(email) {
		return Content{}, fmt.Errorf("User does not exist")
	}
	if (*u)[email].Password != password {
		return Content{}, fmt.Errorf("Incorrect password")
	}
	return (*u)[email], nil
}

func (u *Users) GetUserNoPassword(email string) (Content, error) {
	if !u.ExistUser(email) {
		return Content{}, fmt.Errorf("User does not exist")
	}
	return (*u)[email], nil
}

//delete a user on the map and the file
func (u *Users) DeleteUser(email, password, path string) error {
	if !u.ExistUser(email) {
		return fmt.Errorf("User does not exist")
	}
	if (*u)[email].Password != password {
		return fmt.Errorf("Incorrect password")
	}
	delete(*u, email)
	err := u.saveOnFile(path)
	if err != nil {
		return err
	}
	return nil
}

func (u *Users) DeleteUserNoPassword(email, path string) error {
	if !u.ExistUser(email) {
		return fmt.Errorf("User does not exist")
	}
	delete(*u, email)
	err := u.saveOnFile(path)
	if err != nil {
		return err
	}
	return nil
}

//check if a user with given email exists
func (u *Users) ExistUser(email string) bool {
	_, ok := (*u)[email]
	return ok
}

//save the map on file
func (u Users) saveOnFile(path string) error {
	jsonByte, err := json.Marshal(&u)
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

func (u Users) PrintAllUsers() {
	for k, v := range u {
		fmt.Println(k, v)
	}
}
