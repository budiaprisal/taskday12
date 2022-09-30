package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"personal-web/middleware"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	connection.DatabaseConnect()
	route := mux.NewRouter()
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))
	route.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")

	route.HandleFunc("/project", myProject).Methods("GET")
	route.HandleFunc("/project/{id}", myProjectDetail).Methods("GET")
	route.HandleFunc("/form-project", myProjectForm).Methods("GET")
	route.HandleFunc("/add-project", middleware.UploadFile(myProjectData)).Methods("POST")
	route.HandleFunc("/edit-project/{id}", myProjectEdited).Methods("POST")
	route.HandleFunc("/delete-project/{id}", myProjectDelete).Methods("GET")
	route.HandleFunc("/form-edit-project/{id}", myProjectFormEditProject).Methods("GET")
	route.HandleFunc("/contact", contact).Methods(("GET"))

	route.HandleFunc("/form-register", formRegister).Methods(("GET"))
	route.HandleFunc("/register", register).Methods("POST")


	route.HandleFunc("/form-login", formLogin).Methods("GET")
	route.HandleFunc("/login", login).Methods("POST")

	route.HandleFunc("/logout", logout).Methods("GET")

	fmt.Println("Server running at localhost port 8000")
	http.ListenAndServe("localhost:8000", route)
}
type SessionData struct {
	IsLogin   bool
	UserName  string
	FlashData string
}

var Data = SessionData{}

type StructInputDataForm struct {
	Id              int
	ProjectName     string
	StartDate       time.Time
	EndDate         time.Time
	StartDateFormat string
	EndDateFormat   string
	Image 			string
	Description     string
	Techno          []string
	Duration        string
	Author  		 string
	IsLogin  		 bool
	
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

func home(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("views/index.html")
	if err != nil {
		panic(err)
	}
	template.Execute(w, nil)
}

func myProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("views/myProject.html")
	
	
//

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}



	fm := session.Flashes("message")
//perlu loping karena nanti ketika refresh si alertny masih ada
	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, f1 := range fm {
			// meamasukan flash message
			flashes = append(flashes, f1.(string))
		}
	}


	Data.FlashData = strings.Join(flashes, "")






//
	var result []StructInputDataForm
	data, _ := connection.Conn.Query(context.Background(), "SELECT db_myprojects.id, projectname, startdate, enddate, description, technology,image, tb_user.name as author FROM db_myprojects LEFT JOIN tb_user ON db_myprojects.author_id = tb_user.id ORDER BY id Desc")

	for data.Next() {
		var each = StructInputDataForm{}
		err := data.Scan(&each.Id, &each.ProjectName, &each.StartDate, &each.EndDate, &each.Description, &each.Techno, &each.Author, &each.Image)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		each.Duration = ""

		hour := 1
		day := hour * 24
		week := hour * 24 * 7
		month := hour * 24 * 30
		year := hour * 24 * 365
		differHour := each.EndDate.Sub(each.StartDate).Hours()
		var differHours int = int(differHour)
		days := differHours / day
		weeks := differHours / week
		months := differHours / month
		years := differHours / year
		if differHours < week {
			each.Duration = strconv.Itoa(int(days)) + " Days"
		} else if differHours < month {
			each.Duration = strconv.Itoa(int(weeks)) + " Weeks"
		} else if differHours < year {
			each.Duration = strconv.Itoa(int(months)) + " Months"
		} else if differHours > year {
			each.Duration = strconv.Itoa(int(years)) + " Years"
		}

		result = append(result, each)
	}

	response := map[string]interface{}{
		"DataSession": Data,
		"Projects": result,
	}

	if err == nil {
		tmpl.Execute(w, response)
	} else {
		w.Write([]byte("Message: "))
		w.Write([]byte(err.Error()))
	}
}

func myProjectDetail(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/myProjectDetail.html")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	ProjectDetail := StructInputDataForm{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT db_myprojects.id, projectname, startdate, enddate, description, technology,  tb_user.name as author FROM db_myprojects LEFT JOIN tb_user ON db_myprojects.author_id = tb_user.id  WHERE db_myprojects.id=$1", id).Scan(
		&ProjectDetail.Id, &ProjectDetail.ProjectName, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Description, &ProjectDetail.Techno, &ProjectDetail.Author)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}
	ProjectDetail.StartDateFormat = ProjectDetail.StartDate.Format("2006-01-02")
	ProjectDetail.EndDateFormat = ProjectDetail.EndDate.Format("2006-01-02")
	ProjectDetail.Duration = ""

	hour := 1
	day := hour * 24
	week := hour * 24 * 7
	month := hour * 24 * 30
	year := hour * 24 * 365
	differHour := ProjectDetail.EndDate.Sub(ProjectDetail.StartDate).Hours()
	var differHours int = int(differHour)
	days := differHours / day
	weeks := differHours / week
	months := differHours / month
	years := differHours / year
	if differHours < week {
		ProjectDetail.Duration = strconv.Itoa(int(days)) + " Days"
	} else if differHours < month {
		ProjectDetail.Duration = strconv.Itoa(int(weeks)) + " Weeks"
	} else if differHours < year {
		ProjectDetail.Duration = strconv.Itoa(int(months)) + " Months"
	} else if differHours > year {
		ProjectDetail.Duration = strconv.Itoa(int(years)) + " Years"
	}

	response := map[string]interface{}{
		"Project": ProjectDetail,
	}

	if err == nil {
		tmpl.Execute(w, response)
	} else {
		panic(err)
	}
}

func myProjectForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/myProjectForm.html")
	if err == nil {
		tmpl.Execute(w, nil)
	} else {
		panic(err)
	}
}


func myProjectData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	
	author := session.Values["ID"].(int)
	fmt.Println(author)

	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

 

	var projectName string
	var startDate string
	var endDate string
	var description string
	var techno []string
	fmt.Println(r.Form)
	for i, values := range r.Form {
		fmt.Printf("type of values is %T\n", values)
		fmt.Println(values)
		fmt.Println(i)
		for _, value := range values {
			if i == "projectName" {
				projectName = value
			}
			if i == "startDate" {
				startDate = value
			}
			if i == "endDate" {
				endDate = value
			}
			if i == "description" {
				description = value
			}
			if i == "techno" {
				techno = append(techno, value)
				fmt.Printf("type of value is %T\n", value)
			}
		}
	}
	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO db_myprojects(projectname, startdate, enddate, description, technology,image, author_id ) VALUES ($1, $2, $3, $4, $5, $6, &7)", projectName, startDate, endDate, description, techno, image, author )
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	http.Redirect(w, r, "/project", http.StatusMovedPermanently)
}

func myProjectFormEditProject(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/myProjectFormEditProject.html")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	ProjectEdit := StructInputDataForm{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, projectname, startdate, enddate, description, technology FROM db_myprojects WHERE id=$1", id).Scan(
		&ProjectEdit.Id, &ProjectEdit.ProjectName, &ProjectEdit.StartDate, &ProjectEdit.EndDate, &ProjectEdit.Description, &ProjectEdit.Techno)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}
	ProjectEdit.StartDateFormat = ProjectEdit.StartDate.Format("2006-01-02")
	ProjectEdit.EndDateFormat = ProjectEdit.EndDate.Format("2006-01-02")

	response := map[string]interface{}{
		"Project": ProjectEdit,
	}

	if err == nil {
		tmpl.Execute(w, response)
	} else {
		panic(err)
	}
}

func myProjectEdited(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	var projectName string
	var startDate string
	var endDate string
	var description string
	var techno []string
	fmt.Println(r.Form)
	for i, values := range r.Form {
		for _, value := range values {
			if i == "projectName" {
				projectName = value
			}
			if i == "startDate" {
				startDate = value
			}
			if i == "endDate" {
				endDate = value
			}
			if i == "description" {
				description = value
			}
			if i == "techno" {
				techno = append(techno, value)
			}
		}
	}
	_, err = connection.Conn.Exec(context.Background(), "UPDATE db_myprojects SET projectname=$1, startdate=$2, enddate=$3, description=$4, technology=$5 WHERE id=$6", projectName, startDate, endDate, description, techno, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	http.Redirect(w, r, "/project", http.StatusMovedPermanently)
}

func myProjectDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM db_myprojects WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}
	http.Redirect(w, r, "/project", http.StatusMovedPermanently)
}

func contact(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/contact.html")
	if err == nil {
		tmpl.Execute(w, nil)
	} else {
		panic(err)
	}
}

func formRegister(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/form-register.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var name = r.PostForm.Get("inputName")
	var email = r.PostForm.Get("inputEmail")
	var password = r.PostForm.Get("inputPassword")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	// fmt.Println(passwordHash)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
}

func formLogin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/form-login.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, f1 := range fm {
			// meamasukan flash message
			flashes = append(flashes, f1.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")
	tmpl.Execute(w, Data)
}

func login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	var email = r.PostForm.Get("inputEmail")
	var password = r.PostForm.Get("inputPassword")

	user := User{}

	// mengambil data email, dan melakukan pengecekan email
	err = connection.Conn.QueryRow(context.Background(),
		"SELECT * FROM tb_user WHERE email=$1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {

		// fmt.Println("Email belum terdaftar")
		var store = sessions.NewCookieStore([]byte("SESSION_KEY"))// cookiedari browser
		session, _ := store.Get(r, "SESSION_KEY")

		//session = menyimpan data
		// _  = menampilkan data
		session.AddFlash("Email belum terdaftar!", "message")
		session.Save(r, w)

		http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
		// w.WriteHeader(http.StatusBadRequest)
		// w.Write([]byte("message : Email belum terdaftar " + err.Error()))
		return
	}

	// melakukan pengecekan password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// fmt.Println("Password salah")
		var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")

		session.AddFlash("Password Salah!", "message")
		session.Save(r, w)

		
		
		http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
		// w.WriteHeader(http.StatusBadRequest)
		// w.Write([]byte("message : Email belum terdaftar " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	// berfungsi untuk menyimpan data kedalam session browser
	session.Values["Name"] = user.Name
	session.Values["Email"] = user.Email
	session.Values["ID"] = user.ID
	session.Values["IsLogin"] = true
	session.Options.MaxAge = 10800 // 3 JAM expred

	session.AddFlash("LOGIN BERHASIL", "message")
	
	session.Save(r, w)

	http.Redirect(w, r, "/project", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/form-login", http.StatusSeeOther)
}