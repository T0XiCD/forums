package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"

	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

//je crée une struct pour garder mes info de mo utilisateur, que je récupérerais plus tard dans ma bdd, la struct aurait put etre remplacer par des variable ça marche aussi mais cette fois-ci que pour un utilisateur a la fois ce qui es pas ouf, si sur notre site on veut avoir plusieurs utilisateur connecter en meme temp, ont utilisera un struct comme le RPG en go, ou on crée un struct et plusieurs joueurs totalement différents ders un des autres pour la meme struct (c'était des objets --> es joueurs dans le RPG)
type User struct {
	Pseudo    string
	FirstName string
	LastName  string
	Password  string
	Email     string
}

const (
	Host = "localhost"
	Port = "4445"
)

var pseudo string
var password string
var db *sql.DB              // crée une variable qui pointe a ta base de donnée en l'occurence wampserver qui servira de pont entre les deux
var tmpl *template.Template // pareil mais pointe vers nos template
func registerDB(pseudo string, firstname string, name string, pass string, email string) error { //je passe des argument qui seront completer plus tard la il vaut tous un string vide = ""
	create_user, err := db.Query(fmt.Sprintf("INSERT INTO `register` (`pseudo`,`firstname`, `lastname`, `password`, `email`) VALUES('%s','%s','%s','%s','%s')", pseudo, firstname, name, pass, email)) // %s est un placeholder en sql il existe aussi le "?" qui marche aussi ça veut dire il attend un retour d'une variable de type str --> %s || %d --> int || %t --> bool // fmt.Sprintf --> f car il prend les modulo donc les placeholder, la c'est la partie requetes sql tu peut te réfèrer a ta bdd pour plus de détails
	if err != nil {
		err.Error() //revoie l'erreur s'il y'en a une, pour debug plus facile
	}
	defer create_user.Close()
	return nil
}
func register(w http.ResponseWriter, r *http.Request) {
	pseudo := r.FormValue("pseudo")       //je crée une variable qui va aller dans mon html et prendre ce qui a pour name="pseudo" et je veut récupérer ça value, comme c'est un button par défaut, j'ai délibérement pas mis de value car c'est a l'utilisateur de la rentrer c'est sont pseudo, j'ai juste mis un placeholder, pour dire ce qu'il doit mettre a l'intérieur
	firstname := r.FormValue("firstname") //je récupère sont prénom
	lastname := r.FormValue("lastname")
	password := r.FormValue("password")
	email := r.FormValue("email")
	if r.FormValue("submit") == "Register" { ///ici c'est mal dit mais ça veut dire chope moi le name avec "submit" si il a comme value "Submit" tu lance ma condition, vue que c'est un  type="submit" il ne récupère la value que si tu clique pour envoyer ton formulaire
		registerDB(pseudo, firstname, lastname, password, email)
		http.Redirect(w, r, "http://"+Host+":"+Port+"/login", http.StatusMovedPermanently) // redirection instantané vers un autre url si tout marche
	}
	tmpl.ExecuteTemplate(w, "register", nil) //on execute le template define "register" en l'occurence moi c'est index.html
}
func connectionDB(name string, pass string) (User, error) {
	person := User{}
	getUser := fmt.Sprintf("SELECT * FROM register WHERE pseudo='%s'", name)
	err := db.QueryRow(getUser).Scan(&person.Pseudo, &person.FirstName, &person.LastName, &person.Password, &person.Email) //je récupère tous ce qui ce trouve dans ma bdd qui a pour name celui rentrer, et j'enregistre c'est info dans ma struct
	password = person.Password
	pseudo = person.Pseudo
	return person, err
}
func connection(w http.ResponseWriter, r *http.Request) {
	pseudo := r.FormValue("pseudo")
	password := r.FormValue("password")
	if r.FormValue("submit") == "Submit" {
		log.Println(password)
		if password == password {
			connectionDB(pseudo, password)
			log.Println("GG") //j'envoie un gg dans mon terminal si tout marche
			http.Redirect(w, r, "http://"+Host+":"+Port+"/home", http.StatusMovedPermanently)
		}
	}
	tmpl.ExecuteTemplate(w, "connection", nil)
}
func main() {
	log.Println("lancement du forum")
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/forumv2")
	defer db.Close() // defer comme gorutine mes en mieu permet de mettre cette fonction en priorité, n'attend pas que la fonction en cours ce termie, si tu veut crée un nuvelle utilisateur par exemple, tu veut que ce soie instantané, meme si une function et déja en cours sur golang, eh bien defer permet de faire abstraction des autres func
	log.Println("connection réussie")
	var err error
	tmpl, err = template.New("").Funcs(template.FuncMap{}).ParseGlob("forums-main/*.html") //parse glob permet de récupérer tous les fichier html "*.html" précise que l'on veut tous ce qui finit par ".html" comme dans linux en faite
	if err != nil {
		err.Error()
	}
	http.HandleFunc("/register", register) // exemple : 127.0.0.1:4444/register || le register après la virgule ça la function qui est appeler quand tu va sur cette url
	http.HandleFunc("/login", connection)
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "home", nil)
	})
	http.ListenAndServe(Host+":"+Port, nil)
}
