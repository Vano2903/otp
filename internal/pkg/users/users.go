package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Content struct {
	Password string
	PfpPath  string
}

type Users map[string]Content

var filePath string

//create a new Users map
func NewUsers(fileName string) (Users, error) {
	//check if fileName exists
	//if not create a file with the name fileName
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		//create file
		file, err := os.Create("./" + fileName)
		if err != nil {
			return nil, err
		}
		defer file.Close()
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
	filePath = fileName
	return make(Users), nil
}

//add a new user to the map and save the map on file
func (u *Users) AddUser(username, password, pfpPath string) error {
	if u.ExistUser(username) {
		return fmt.Errorf("User already exists")
	}
	(*u)[username] = Content{password, pfpPath}
	fmt.Println(username, password)
	u.PrintAllUsers()
	err := u.saveOnFile()
	if err != nil {
		return err
	}
	return nil
}

//update a user on the map and save the map on file
func (u *Users) UpdateUser(username, password, pfpPath string) error {
	if !u.ExistUser(username) {
		return fmt.Errorf("User does not exist")
	}
	(*u)[username] = Content{password, pfpPath}
	err := u.saveOnFile()
	if err != nil {
		return err
	}
	return nil
}

//return the Content of the user given the username
func (u *Users) GetUser(username string) (Content, error) {
	if !u.ExistUser(username) {
		return Content{}, fmt.Errorf("User does not exist")
	}
	return (*u)[username], nil
}

//delete a user on the map and the file
func (u *Users) DeleteUser(username string) error {
	if !u.ExistUser(username) {
		return fmt.Errorf("User does not exist")
	}
	delete(*u, username)
	err := u.saveOnFile()
	if err != nil {
		return err
	}
	return nil
}

//check if a user with given username exists
func (u *Users) ExistUser(username string) bool {
	_, ok := (*u)[username]
	return ok
}

//save the map on file
func (u Users) saveOnFile() error {
	jsonByte, err := json.Marshal(&u)
	if err != nil {
		return err
	}
	file, err := os.Create(filePath)
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
