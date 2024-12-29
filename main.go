package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
	"unicode"

	_ "github.com/lib/pq"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct {
	staffid   int
	Username  string
	password  []byte
	firstname string
	lastname  string
	hash      string
}

type Staffinfo struct {
	Id             int
	Firstname      string
	Lastname       string
	Position       string
	Age            float32
	Salary         float32
	Dateofbirth    time.Time
	Yearsofservice float32
	Officialcar    bool
}

var (
	templates = template.Must(template.New("").Funcs(fm).ParseGlob("templates/*"))
	// = template.Must(template.ParseGlob("templates/*"))
)

var db *sql.DB
var store = sessions.NewCookieStore([]byte("super-secret"))

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "charles"
	dbname   = "staffdb"
)

func main() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	fmt.Println("POSTGRESQL DATABASE CONNECTED!")

	http.HandleFunc("/", index)
	http.HandleFunc("/datavault", (datavault))
	http.HandleFunc("/register", register)
	http.HandleFunc("/registerhandler", registerhandler)
	http.HandleFunc("/login", login)
	http.HandleFunc("/loginhandler", loginhandler)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/createtable", createtable)
	http.HandleFunc("/createstafftable", createstafftable)
	http.HandleFunc("/dropstafftable", dropstafftable)
	http.HandleFunc("/alterstafftable", alterstafftable)
	http.HandleFunc("/droptable", droptable)
	http.HandleFunc("/alterfts", alterfts)
	http.HandleFunc("/insertdata", (insertdata))
	http.HandleFunc("/insertdone", (insertdone))
	http.HandleFunc("/updatestaffdata", Auth(updatestaffdata))
	http.HandleFunc("/updatestaffdone", Auth(updatestaffdone))
	http.HandleFunc("/deleteinfo", Auth(deleteinfo))
	http.HandleFunc("/searchdata", searchdata)
	http.HandleFunc("/search", search)

	err = http.ListenAndServe("0.0.0.0:8080", context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatalln(err)
	}
}

func createtable(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare(`create table IF NOT EXISTS bcrypt (
        staffid SERIAL PRIMARY KEY,
        username VARCHAR(50) NOT NULL,
        password VARCHAR(50) NOT NULL,
        firstname VARCHAR(50) NOT NULL,
        lastname VARCHAR(50) NOT NULL,
        hash VARCHAR(80) NOT NULL

  );`)
	if err != nil {
		fmt.Println("err: ", err)
	}

	defer stmt.Close()

	res, err := stmt.Exec()
	check(err)

	n, err := res.RowsAffected()
	check(err)

	fmt.Fprintln(w, "TABLE SUCCESSFULLY CREATED", n, " new table created")

}

func droptable(w http.ResponseWriter, r *http.Request) {
	drp, err := db.Prepare(`DROP TABLE IF EXISTS bcrypt;`)

	defer drp.Close()

	res, err := drp.Exec()
	check(err)

	n, err := res.RowsAffected()
	check(err)

	fmt.Fprintln(w, "TABLE SUCCESSFULLY DROPPED", n)

}

// main staff info
func createstafftable(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare(`create table IF NOT EXISTS staffinfo (
        id BIGSERIAL PRIMARY KEY,
        firstname VARCHAR(50) NOT NULL,
        lastname VARCHAR(50) NOT NULL,
        position VARCHAR(50) NOT NULL,
        age NUMERIC(4, 2) NOT NULL,
        salary NUMERIC(10, 2) NOT NULL,
        yearsofservice NUMERIC(10, 2) NOT NULL,
        dateofbirth DATE NOT NULL,
		officialcar BOOLEAN NOT NULL

		
  );`)
	if err != nil {
		fmt.Println("err: ", err)
	}

	defer stmt.Close()

	res, err := stmt.Exec()
	check(err)

	n, err := res.RowsAffected()
	check(err)

	fmt.Fprintln(w, "STAFF TABLE SUCCESSFULLY CREATED", n, " new table created")

}

// ALTER
func alterstafftable(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare(`
	ALTER TABLE staffinfo
	ALTER COLUMN officialcar TYPE BOOL
	SET DEFAULT FALSE
	;`)
	if err != nil {
		fmt.Println("err: ", err)
	}

	defer stmt.Close()

	res, err := stmt.Exec()
	check(err)

	n, err := res.RowsAffected()
	check(err)

	fmt.Fprintln(w, "STAFF TABLE IS SUCCESSFULLY ALTERED", n)

}

func dropstafftable(w http.ResponseWriter, r *http.Request) {
	drp, err := db.Prepare(`DROP TABLE IF EXISTS staffinfo;`)

	defer drp.Close()

	res, err := drp.Exec()
	check(err)

	n, err := res.RowsAffected()
	check(err)

	fmt.Fprintln(w, "TABLE SUCCESSFULLY DROPPED", n)

}

func index(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
}

// initializations
var usernamelen, alphanumeric bool
var passwordlen, passwordupper, passwordlower, passwordnum, passwordspecial, passwordspace bool

func register(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register.html", nil)
}

func registerhandler(w http.ResponseWriter, r *http.Request) {
	// username criteria
	// password criteria
	// Does username exist in database?
	// Create hash password
	// Inserting username and hash into database

	r.ParseForm()
	s := Bcrypt{}
	s.firstname = r.FormValue("firstname")
	s.lastname = r.FormValue("lastname")
	s.Username = r.FormValue("username")
	password := r.FormValue("password")

	// username length
	if 8 < len(s.Username) && len(s.Username) < 20 {
		usernamelen = true
	}

	// username to contain letters or numbers
	for _, char := range s.Username {
		if unicode.IsLetter(char) == false && unicode.IsNumber(char) == false {
			alphanumeric = false
		}
	}

	// password length
	if 8 < len(s.password) && len(s.password) < 20 {
		passwordlen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			passwordupper = true
		case unicode.IsLower(char):
			passwordlower = true
		case unicode.IsNumber(char):
			passwordnum = true
		case unicode.IsSymbol(char):
			passwordspecial = true
		case unicode.IsSpace(char):
			passwordspace = false

		}
	}

	//if !usernamelen || !alphanumeric || !passwordlen || !passwordupper || !passwordlower || !passwordnum || !passwordspecial || !passwordspace {
	//	templates.ExecuteTemplate(w, "register.html", "Follow the Username and Password criteria")
	//	return
	//}

	// Does username exist in database?
	// GET
	ext := `SELECT staffid FROM bcrypt WHERE username = $1;`
	row := db.QueryRow(ext, s.Username)
	err := row.Scan(&s.staffid)
	if err != sql.ErrNoRows {
		templates.ExecuteTemplate(w, "register.html", "Username has been taken!")
		return

	}

	// Create hash password
	hash, err := bcrypt.GenerateFromPassword(([]byte(password)), bcrypt.MinCost)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return

	}

	// Inserting username and hash into database
	if r.Method == http.MethodPost {
		if s.firstname == "" || s.lastname == "" || s.Username == "" || password == "" {
			templates.ExecuteTemplate(w, "register.html", "Ensure that you fill the fields")
			return
		}
		smt, err := db.Prepare(`INSERT into bcrypt (username, password, firstname, lastname, hash) VALUES ($1, $2, $3, $4, $5);`)
		if err != nil {
			panic(err)
		}
		defer smt.Close()

		res, err := smt.Exec(s.Username, s.password, s.firstname, s.lastname, hash)
		if err != nil {
			panic(err)
		}

		n, err := res.RowsAffected()
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		fmt.Fprintln(w, s.firstname+""+s.lastname, "You have successfully created a username and a password!", n)
		return

	}

}

func login(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func loginhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// form values
		//u := Bcrypt{}
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			templates.ExecuteTemplate(w, "login.html", "Ensure that you fill the fields")
			return

		}

		// check to see if username exist in the database
		var staffid, hash string

		smt := `SELECT staffid, hash FROM bcrypt WHERE username = $1;`
		row := db.QueryRow(smt, username)
		err := row.Scan(&staffid, &hash)
		fmt.Println("pass thru")
		if err != nil {
			templates.ExecuteTemplate(w, "login.html", "Username not in database. Do signup!")
			fmt.Println("cant login")
			return
		}

		// check to see if hash/password exist in the database
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		if err == nil {
			session, _ := store.Get(r, "session")
			session.Values["staffid"] = staffid
			session.Save(r, w)
			http.Redirect(w, r, "/datavault", http.StatusSeeOther)
			// templates.ExecuteTemplate(w, "datavault.html", "Logged In") templates works perfectly
			return
		}

		//templates.ExecuteTemplate(w, "login.html", "username or password not correct")
	}
	fmt.Fprintln(w, "Username or Password not correct")
	//templates.ExecuteTemplate(w, "login.html", "Check Username and Pass")

}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	delete(session.Values, "staffid")
	session.Save(r, w)
	templates.ExecuteTemplate(w, "login.html", "Logged Out")

}

func Auth(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		_, ok := session.Values["staffid"]
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		// ServeHTTP calls f(w, r)
		// func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)
		HandlerFunc.ServeHTTP(w, r)
	}
}

// type FuncMap map[string]interface{}
var fm = template.FuncMap{
	"Au": Auth,
}

func alterfts(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare(`
	ALTER table staffinfo
	ADD COLUMN document tsvector;
	UPDATE staffinfo
	SET document = to_tsvector(firstname || ' ' || lastname || ' ' || position)
	;`)
	if err != nil {
		fmt.Println("err: ", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec()
	check(err)

	n, err := res.RowsAffected()
	check(err)

	fmt.Fprintln(w, "ALTERED A TABLE FOR FULL TEXT SEARCH", n)
}

func insertdata(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "insertdata.html", nil)
}

func insertdone(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		age := r.FormValue("age")
		dateofbirth := r.FormValue("dateofbirth")
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lastname")
		officialcar := r.FormValue("officialcar")
		yearsofservice := r.FormValue("yearsofservice")
		position := r.FormValue("position")
		salary := r.FormValue("salary")

		if age == "" || dateofbirth == "" || firstname == "" || lastname == "" || officialcar == "" || yearsofservice == "" || position == "" || salary == "" {
			templates.ExecuteTemplate(w, "insertdata.html", "Error! Check to see that all fields have been completed")
			return
		}

		smt := `INSERT into staffinfo (firstname, lastname, position, age, salary, yearsofservice, dateofbirth, officialcar) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`
		ins, err := db.Prepare(smt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer ins.Close()

		res, err := ins.Exec(firstname, lastname, position, age, salary, yearsofservice, dateofbirth, officialcar)
		check(err)

		n, err := res.RowsAffected()
		check(err)

		//fmt.Fprintln(w, "Inserted successfully", n)
		//templates.ExecuteTemplate(w, "insertdone.html", "INSERTED SUCCESFULLY")
		fmt.Println(n)
		http.Redirect(w, r, "/datavault", 307)
		return

	}
	// templates.ExecuteTemplate(w, "insertdone.html", "Item inserted successfully")
	// had to mute that or it displays d content of inserdone.html after execution.
}

func datavault(w http.ResponseWriter, r *http.Request) {

	smt := `SELECT * FROM staffinfo;`
	rows, err := db.Query(smt)
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	var staffs []Staffinfo

	for rows.Next() {
		u := Staffinfo{}
		err = rows.Scan(&u.Id, &u.Firstname, &u.Lastname, &u.Position, &u.Age, &u.Salary, &u.Yearsofservice, &u.Dateofbirth, &u.Officialcar)
		if err != nil {
			panic(err.Error())
		}
		staffs = append(staffs, u)

	}
	err = templates.ExecuteTemplate(w, "datavault.html", staffs)
	if err != nil {
		panic(err.Error())
	}

}

// Get what you want to update then insert
func updatestaffdata(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	smt := `SELECT * FROM staffinfo WHERE id = $1;`

	d := Staffinfo{}
	row := db.QueryRow(smt, id)
	err := row.Scan(&d.Id, &d.Firstname, &d.Lastname, &d.Position, &d.Age, &d.Salary, &d.Yearsofservice, &d.Dateofbirth, &d.Officialcar)
	if err != nil {
		//http.Redirect(w, r, "/", 307)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//return
		panic(err.Error())
	}
	err = templates.ExecuteTemplate(w, "update.html", d)
	if err != nil {
		panic(err.Error())
	}
}

func updatestaffdone(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		id := r.FormValue("id")
		age := r.FormValue("age")
		dateofbirth := r.FormValue("dateofbirth")
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lastname")
		officialcar := r.FormValue("officialcar")
		yearsofservice := r.FormValue("yearsofservice")
		position := r.FormValue("position")
		salary := r.FormValue("salary")

		if age == "" || dateofbirth == "" || firstname == "" || lastname == "" || officialcar == "" || yearsofservice == "" || position == "" || salary == "" {
			templates.ExecuteTemplate(w, "update.html", "Error inserting data! Check all fields!")
			return
		}

		upt := `UPDATE staffinfo SET firstname = $1, lastname = $2, position = $3, age = $4, salary = $5, yearsofservice = $6, dateofbirth = $7, officialcar = $8 WHERE id = $9;`
		inst, err := db.Prepare(upt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer inst.Close()

		res, err := inst.Exec(firstname, lastname, position, age, salary, yearsofservice, dateofbirth, officialcar, id)
		check(err)

		n, err := res.RowsAffected()
		check(err)

		fmt.Println(n)

	}
	err := templates.ExecuteTemplate(w, "updatedatadone.html", "UPDATED SUCCESSFULLY")
	if err != nil {
		panic(err.Error())
	}
}

func deleteinfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	del := `DELETE FROM staffinfo WHERE id = $1;`
	smt, err := db.Prepare(del)
	if err != nil {
		panic(err.Error())
	}

	defer smt.Close()

	res, err := smt.Exec(id)
	if err != nil {
		panic(err.Error())
	}

	n, err := res.RowsAffected()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(n, "deleted")
	//fmt.Fprintln(w, "deleted")
	//templates.ExecuteTemplate(w, "datavault.html", "DATA DELETED")
	http.Redirect(w, r, "/datavault", 307)
	return

}

func searchdata(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "form", nil)
	if err != nil {
		panic(err)
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	firstname := r.FormValue("firstname")
	smt := `SELECT * FROM staffinfo WHERE firstname = $1;`
	row := db.QueryRow(smt, firstname)
	p := Staffinfo{}
	err := row.Scan(&p.Id, &p.Firstname, &p.Lastname, &p.Position, &p.Age, &p.Salary, &p.Dateofbirth, &p.Yearsofservice, &p.Officialcar)
	switch err {
	case nil:
		fmt.Println(err.Error())
	case sql.ErrNoRows:
		fmt.Fprintln(w, "No Data relating to you search found in the database!")
		return
	default:
		fmt.Println("Yes!!!")
	}
	err = templates.ExecuteTemplate(w, "result", nil)
	if err != nil {
		panic(err.Error())
	}

}

// password1 and password2 must be same
