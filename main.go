package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Vano2903/vano-otp/internal/pkg/email"
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
}

var (
	c        Config
	files    []File
	u        users.Users
	pendings users.Users
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
}

//handler of the login (post), check if the user sent is a valid user and if it is will return the correct user page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var post PostContent

	//read post body
	_ = json.NewDecoder(r.Body).Decode(&post)

	//check if user is correct
	user, err := u.GetUser(post.Email, post.Password)

	//return response
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf(`{"accepted":false, "code": 401, "msg": %q}`, err.Error())))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(fmt.Sprintf(`{"accepted":true, "code": 202, "pfpUrl": %s}`, user.PfpUrl)))
}

//handler that let user register to the database
func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	var post PostContent

	_ = json.NewDecoder(r.Body).Decode(&post)

	id := ksuid.New()
	e := fmt.Sprintf(`<head>
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
		<h2>Experia</h2>
		<h1>Ciao, abbiamo quasi fatto, conferma la tua registrazione cliccando qui sotto!</h1>
		<a href='https://vano-otp.herokuapp.com/auth/confirm?email=%s;id=%s' id='submit'>Conferma la registrazione</a>
	</div>`, post.Email, id.String())

	if !email.IsValid(post.Email) {
		PrintErr(w, "email is not valid")
		return
	}

	err := email.SendEmail(c.Email, c.EmailPassword, post.Email, "Conferma la registrazione", e)
	if err != nil {
		PrintInternalErr(w, err.Error())
		return
	}

	err = pendings.AddUser(post.Email, id.String(), c.PendingFilePath)
	if err != nil {
		PrintInternalErr(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status": 202, "msg": "confirmation email correctly sent"}`))
}

func ConfirmAccountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	id := vars["id"]

	_, err := pendings.GetUser(email, id)
	if err != nil {
		PrintErr(w, err.Error())
		return
	}

	pendings.DeleteUser(email, id, c.PendingFilePath)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status": 201, "msg": "user registered correctly"}`))
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
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

			fmt.Println(file)
			err = u.UpdateUser(postData.Email, postData.Password, postData.Password, file.Url, c.UserFilePath)
			fmt.Println(u)
			if err != nil {
				PrintInternalErr(w, err.Error())
				return
			}

			files = append(files[:i], files[i+1:]...)
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte(`{"status": 202, "msg": "document added successfully"}`))
			log.Println("file uploaded successfully")
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

	//root
	// r.HandleFunc(root.String(), LoginPageHandler).Methods("GET", "OPTIONS")

	//user area
	r.HandleFunc(usersLogin.String(), LoginHandler).Methods("POST", "OPTIONS")
	r.HandleFunc(addUser.String(), AddUserHandler).Methods("POST", "OPTIONS")

	//email
	r.HandleFunc(ConfirmAccount.String(), ConfirmAccountHandler).Methods("GET", "OPTIONS")

	//document section
	r.HandleFunc(fileupload.String(), UploadFileHandler).Methods("POST", "OPTIONS")
	r.HandleFunc(fileBind.String(), DocumentBindHandler).Methods("POST", "OPTIONS")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST"})

	log.Println("starting on", ":"+port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
