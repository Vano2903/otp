package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Vano2903/vano-otp/internal/pkg/email"
	"github.com/Vano2903/vano-otp/internal/pkg/otp"
	"github.com/Vano2903/vano-otp/internal/pkg/users"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

type File struct {
	Url string `json:"url, omitempty"`
	ID  string `json:"id, omitempty"`
}

type PostContent struct {
	Email    string `json:"email, omitempty"`
	Password string `json:"password, omitempty"`
	ID       string `json:"id, omitempty"`
}

var (
	c          Config
	files      []File
	u          users.Users
	pendings   users.Users
	otpHandler otp.OtpHandler
)

func init() {
	var err error
	//load config
	if err = c.Load(); err != nil {
		log.Fatal(err)
	}

	//load pending users
	pendings, err = users.NewUsers(c.PendingFilePath)
	if err != nil {
		log.Fatal(err)
	}

	//load users
	u, err = users.NewUsers(c.UserFilePath)
	if err != nil {
		log.Fatal(err)
	}

	otpHandler, err = otp.NewOtpHandler(c.OtpFilePath)
	if err != nil {
		log.Fatal(err)
	}
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	home, err := os.ReadFile("pages/login.html")
	if err != nil {
		UnavailablePage(w)
		return
	}
	w.Write(home)
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	home, err := os.ReadFile("pages/register.html")
	if err != nil {
		UnavailablePage(w)
		return
	}
	w.Write(home)
}

func OtpPageHandler(w http.ResponseWriter, r *http.Request) {
	home, err := os.ReadFile("pages/otp.html")
	if err != nil {
		UnavailablePage(w)
		return
	}
	w.Write(home)
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	home, err := os.ReadFile("pages/home.html")
	if err != nil {
		UnavailablePage(w)
		return
	}
	w.Write(home)
}

//handler of the login (post), check if the user sent is a valid user and if it is will return the correct user page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var post PostContent

	//read post body
	_ = json.NewDecoder(r.Body).Decode(&post)

	//check if user is correct
	_, err := u.GetUser(post.Email, post.Password)
	if err != nil {
		PrintErr(w, err.Error())
		return
	}

	otpSecret := otpHandler.CreateNew(post.Email, c.OtpFilePath)

	err = email.SendEmail(c.Email, c.EmailPassword, post.Email, "Confirmation Code", fmt.Sprintf("Confirmation code is: <br><br> <b>%s</b>", otpSecret))
	if err != nil {
		PrintInternalErr(w, err.Error())
		return
	}

	//return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"sent":true, "code": 200}`))
}

//handle the otp request, get as post the users credentials and the otp code as endpoint
func OtpHandler(w http.ResponseWriter, r *http.Request) {
	otpSecret := mux.Vars(r)["otp"]
	var post PostContent

	//read post body
	_ = json.NewDecoder(r.Body).Decode(&post)

	//check if user is correct
	_, err := u.GetUser(post.Email, post.Password)
	if err != nil {
		PrintErr(w, err.Error())
		return
	}

	//check if the otp is correct
	err = otpHandler.CheckOtp(post.Email, otpSecret, c.OtpFilePath)

	//if the otp is incorrect return 400
	if err != nil {
		PrintErr(w, err.Error())
		return
	}

	//TODO respond the user page
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusAccepted)
	HomePageHandler(w, r)
}

func GetUserPfP(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]

	//check if user is correct
	u, err := u.GetUserNoPassword(email)
	if err != nil {
		PrintErr(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"code":200, "url":"%s"}`, u.PfpUrl)))
}

//handler that let user register to the database
func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post PostContent

	_ = json.NewDecoder(r.Body).Decode(&post)

	if !u.ExistUser(post.Email) {

		id := ksuid.New()
		emailHead := `<head>
	<style>
		div {
			background-color: #1e1e1e;
			display: grid;
			padding: 0 1rem 1rem 1rem;
			justify-content: center;
			align-items: center;
			border-radius: .2rem;
		}
		#submit, #submit:visited, #submit:active {
			margin: 1rem auto;
			cursor: pointer;
			font-family: inherit;
			font-size: 1rem;
			border-radius: .2rem;
			padding: 1rem 3rem;
			transition: .2s;
			outline: none;
			height: fit-content;
			background-color: #ffcc80;
			border: none;
			color: #000000;
			text-decoration: none;
		}
	
		#submit:hover {
			background-color: #ca9b52;
		}
	
		h1 {
			margin: 0 auto;
			color: #ffffff;
		}
		p {
			margin-top: 2rem;
			width: 100%;
			color: white;
		}
		#delete, #delete:hover, #delete:visited, #delete:active {
			color: #9c64a6;
			text-decoration: none;
		}
		h2 {
			width: 100%;
			color: #ffffff;
			margin: 0 0 1rem 0;
		}
	</style>
	</head>
	<div>
		<h1>Hi, we are almost done, confirma your registration clicking the button below</h1>`
		emailHead += fmt.Sprintf(`
		<a href='https://vano-otp.herokuapp.com/auth/confirm?email=%s&id=%s' id='submit'>Confirm your registration</a>
	</div>`, post.Email, id.String())

		if !email.IsValid(post.Email) {
			PrintErr(w, "email is not valid")
			return
		}

		err := email.SendEmail(c.Email, c.EmailPassword, post.Email, "Confirm your registration", emailHead)
		if err != nil {
			PrintInternalErr(w, err.Error())
			return
		}

		err = pendings.AddUser(post.Email, id.String()+";"+post.Password, c.PendingFilePath)
		if err != nil {
			PrintErr(w, err.Error())
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"code": 202, "msg": "confirmation email correctly sent"}`))
		return
	}
	PrintErr(w, "user already exist")
}

func ConfirmAccountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()

	var email, id string
	for k, v := range r.Form {
		switch k {
		case "id":
			id = v[0]
		case "email":
			email = v[0]
		default:
			PrintErr(w, "missing id or email parameters")
			return
		}
	}
	fmt.Println(id)

	user, err := pendings.GetUserNoPassword(email)
	if err != nil {
		PrintErr(w, err.Error())
		return
	}

	err = pendings.DeleteUserNoPassword(email, c.PendingFilePath)
	fmt.Println("deleting pending user:", err)
	u.AddUser(email, strings.Split(user.Password, ";")[1], c.UserFilePath)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status": 201, "msg": "user registered correctly"}`))
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("document")

	if err != nil {
		PrintErr(w, err.Error())
		return
	}

	defer file.Close()
	fmt.Printf("Uploading File: %+v\n", handler.Filename)
	// fmt.Printf("File Size: %+v\n", handler.Size)
	// fmt.Printf("MIME Header: %+v\n", handler.Header)
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(buf, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url, err := UploadFile(buf.Bytes(), handler.Filename)
	if err != nil {
		PrintInternalErr(w, err.Error())
		return
	}
	id := ksuid.New()
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(fmt.Sprintf(`{"code": 202, "fileID": "%s"}`, id.String())))

	files = append(files, File{url, id.String()})
	fmt.Println(url)
}

func DocumentBindHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(files)
	kuid := mux.Vars(r)["kuid"]
	for i, file := range files {
		if file.ID == kuid {
			var postData PostContent
			_ = json.NewDecoder(r.Body).Decode(&postData)

			_, err := u.GetUser(postData.Email, postData.Password)
			if err != nil {
				PrintErr(w, err.Error())
				return
			}

			err = u.UpdateUser(postData.Email, postData.Password, postData.Password, file.Url, c.UserFilePath)
			if err != nil {
				PrintInternalErr(w, err.Error())
				return
			}

			files = append(files[:i], files[i+1:]...)
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte(fmt.Sprintf(`{"code": 202, "msg": "document added successfully", "pfpUrl":"%s"}`, file.Url)))
			return
		}
	}
	PrintErr(w, "invalid KUID")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	//router
	r := mux.NewRouter()

	//statics
	r.PathPrefix(statics.String()).Handler(http.StripPrefix(statics.String(), http.FileServer(http.Dir("static/"))))

	r.HandleFunc("/", HomePageHandler).Methods("GET")

	//pages handlers
	r.HandleFunc(loginPage.String(), LoginPageHandler).Methods("GET", "OPTIONS")
	r.HandleFunc(registerPage.String(), RegisterPageHandler).Methods("GET", "OPTIONS")
	r.HandleFunc(otpPage.String(), OtpPageHandler).Methods("GET", "OPTIONS")

	//user area
	r.HandleFunc(usersLogin.String(), LoginHandler).Methods("POST", "OPTIONS")
	r.HandleFunc(addUser.String(), AddUserHandler).Methods("POST", "OPTIONS")
	r.HandleFunc(getUserPfP.String(), GetUserPfP).Methods("GET", "OPTIONS")

	//email
	r.HandleFunc(ConfirmAccount.String(), ConfirmAccountHandler).Methods("GET", "OPTIONS")

	//otp
	r.HandleFunc(otpConfirmation.String(), OtpHandler).Methods("POST", "OPTIONS")

	//document section
	r.HandleFunc(fileupload.String(), UploadFileHandler).Methods("POST", "OPTIONS")
	r.HandleFunc(fileBind.String(), DocumentBindHandler).Methods("POST", "OPTIONS")

	//api mode
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST"})

	log.Println("starting on", ":"+port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
