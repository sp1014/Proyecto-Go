package main

import (
	"database/sql"
	"fmt"

	//"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionDB() (conexion *sql.DB) {

	Driver := "mysql"
	Usuario := "root"
	Contrasena := ""
	Nombre := "go"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasena+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion

}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)

	fmt.Println("Servidor Corriendo...")
	http.ListenAndServe(":8080", nil)
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")

	conexonEstablecida := conexionDB()

	borrarRegistros, err := conexonEstablecida.Prepare("DELETE FROM empleados WHERE id=?")

	if err != nil {
		panic(err.Error())
	}
	borrarRegistros.Exec(idEmpleado)

	http.Redirect(w, r, "/", 301)

}

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func Inicio(w http.ResponseWriter, r *http.Request) {

	//fmt.Fprintf(w, "Hola develoteca")

	conexonEstablecida := conexionDB()

	registros, err := conexonEstablecida.Query("SELECT * FROM empleados")

	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for registros.Next() {
		var id int
		var nombre, correo string
		err = registros.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)
	}

	//	fmt.Println(arregloEmpleado)

	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado)
}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexonEstablecida := conexionDB()

		insertarRegistros, err := conexonEstablecida.Prepare("INSERT INTO empleados(nombre,correo) VALUES(?,?)")

		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, correo)

		http.Redirect(w, r, "/", 301)
	}
}
