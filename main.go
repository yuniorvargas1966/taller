package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	//"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var status = 301

func conexionBD() (conexion *sql.DB) {
	godotenv.Load()

	Driver := os.Getenv("Driver")
	Usuario := os.Getenv("Usuario")
	Contrasena := os.Getenv("Contrasena")
	Nombre := os.Getenv("Nombre")

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasena+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", Inicio)

	http.HandleFunc("/crear", Crear)

	http.HandleFunc("/insertar", Insertar)

	http.HandleFunc("/borrar", Borrar)

	http.HandleFunc("/editar", Editar)

	http.HandleFunc("/actualizar", Actualizar)

	fmt.Println("Servidor corriendo en el puerto 8080, entre a http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idServicio := r.URL.Query().Get("id")
	fmt.Println(idServicio)

	conexionEstablecida := conexionBD()

	borrarRegistro, err := conexionEstablecida.Prepare("DELETE FROM servicios WHERE id=?")

	if err != nil {
		panic(err.Error())
	}

	borrarRegistro.Exec(idServicio)
	http.Redirect(w, r, "/", status)
}

type Servicio struct {
	Id       int
	Nombre   string
	Equipo   string
	Trabajo  string
	Telefono string
	Correo   string
	Fecha    string
}

func Inicio(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionBD()

	registros, err := conexionEstablecida.Query("SELECT * FROM servicios")

	if err != nil {
		panic(err.Error())
	}
	servicio := Servicio{}
	arregloServicio := []Servicio{}

	for registros.Next() {
		var id int
		var nombre, equipo, trabajo, telefono, correo, fecha string
		err = registros.Scan(&id, &nombre, &equipo, &trabajo, &telefono, &correo, &fecha)
		if err != nil {
			panic(err.Error())
		}
		servicio.Id = id
		servicio.Nombre = nombre
		servicio.Equipo = equipo
		servicio.Trabajo = trabajo
		servicio.Telefono = telefono
		servicio.Correo = correo
		servicio.Fecha = fecha

		arregloServicio = append(arregloServicio, servicio)

	}
	// fmt.Println(arregloAtencion)
	plantillas.ExecuteTemplate(w, "inicio", arregloServicio)
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idServicio := r.URL.Query().Get("id")
	fmt.Println(idServicio)

	conexionEstablecida := conexionBD()
	registro, err := conexionEstablecida.Query("SELECT * FROM servicios WHERE id=?", idServicio)

	servicio := Servicio{}
	for registro.Next() {
		var id int
		var nombre, equipo, trabajo, telefono, correo, fecha string
		err = registro.Scan(&id, &nombre, &equipo, &trabajo, &telefono, &correo, &fecha)
		if err != nil {
			panic(err.Error())
		}
		servicio.Id = id
		servicio.Nombre = nombre
		servicio.Equipo = equipo
		servicio.Trabajo = trabajo
		servicio.Telefono = telefono
		servicio.Correo = correo
		servicio.Fecha = fecha
	}

	fmt.Println(servicio)
	plantillas.ExecuteTemplate(w, "editar", servicio)
}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}
func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		nombre := r.FormValue("nombre")
		equipo := r.FormValue("equipo")
		trabajo := r.FormValue("trabajo")
		telefono := r.FormValue("telefono")
		correo := r.FormValue("correo")
		fecha := r.FormValue("fecha")

		conexionEstablecida := conexionBD()

		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO servicios(nombre,equipo,trabajo,telefono,correo,fecha) VALUES(?, ?, ?, ?, ?, ?) ")

		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, equipo, trabajo, telefono, correo, fecha)

		http.Redirect(w, r, "/", status)

	}

}
func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		equipo := r.FormValue("equipo")
		trabajo := r.FormValue("trabajo")
		telefono := r.FormValue("telefono")
		correo := r.FormValue("correo")
		fecha := r.FormValue("fecha")

		conexionEstablecida := conexionBD()

		modificarRegistros, err := conexionEstablecida.Prepare("UPDATE servicios SET nombre=?,equipo=?,trabajo=?,telefono=?,correo=?,fecha=? WHERE id=? ")

		if err != nil {
			panic(err.Error())
		}
		modificarRegistros.Exec(nombre, equipo, trabajo, telefono, correo, fecha, id)

		http.Redirect(w, r, "/", status)

	}

}
