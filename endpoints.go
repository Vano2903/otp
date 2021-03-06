package main

type endpoints string

const (
	root endpoints = "/"

	//*statics
	statics endpoints = "/static/" //endpoint for the statics document (js, css, html, images)

	//*pages
	loginPage    endpoints = "/login"
	registerPage endpoints = "/register"
	otpPage      endpoints = "/otp"

	//*users endpoints
	getUserPfP endpoints = "/users/pfp/{email}"
	usersLogin endpoints = "/users/login"  //endpoint to check the login informations
	checkEmail endpoints = "/users/check"  //check if a email is already registered
	addUser    endpoints = "/users/singup" //as get will return the signup page, as post will read the user's info and if valid add the user to the db
	updateUser endpoints = "/users/update"

	//*email
	ConfirmAccount endpoints = "/auth/confirm"

	//*document and upload endpoints
	fileupload endpoints = "/upload/file"
	fileBind   endpoints = "/upload/bind/{kuid}"

	//*otp endpoint
	otpConfirmation endpoints = "/otp/check/{otp}"
)

//convert endpoint to string
func (e endpoints) String() string {
	return string(e)
}
