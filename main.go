package main

//imports usados
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql" //traer repositorio de github
)

//conexion a base de datos
func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"   //tipo de db
	Usuario := "root"   //usuario
	Contrasena := ""    //contrasena
	Nombre := "sistema" //nombre db
	//se indica el driver usuario + contrasena concatenada +@tcp(direccion de db)+ nombre
	conexion, err := sql.Open(Driver, Usuario+":"+Contrasena+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil { //diferente a error
		panic(err.Error()) //imprimir error
	}
	return conexion //regresar conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id") //traer dato de formulario
	fmt.Println(idEmpleado)               //imprimir por consola

	conexionEstablecida := conexionDB() //establecer conexion
	//insertar datos en db
	borrarRegistro, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?")

	if err != nil {
		panic(err.Error())
	}

	borrarRegistro.Exec(idEmpleado)
	http.Redirect(w, r, "/", 301) //redireccionar
}

func main() {
	http.HandleFunc("/", Inicio)               //inicio
	http.HandleFunc("/crear", Crear)           //crear
	http.HandleFunc("/insertar", insertar)     //insertar
	http.HandleFunc("/borrar", Borrar)         //borrar
	http.HandleFunc("/editar", Editar)         //pagina editar
	http.HandleFunc("/actualizar", Actualizar) //actualizar

	log.Println("Servidor corriendo...")
	http.ListenAndServe(":3000", nil) //puerto

}

type Empleado struct { //tipo de structura de datos
	Id     int    //tipo int
	Nombre string //tipo string
	Correo string
}

func Inicio(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "hola develoteca")//imprimir  en pantalla
	conexionEstablecida := conexionDB()
	registros, err := conexionEstablecida.Query("SELECT * FROM empleados")

	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}
	arreglaEmpleado := []Empleado{}

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
		arreglaEmpleado = append(arreglaEmpleado, empleado)
	}
	//fmt.Println(arreglaEmpleado)//imprimir por consola

	plantillas.ExecuteTemplate(w, "inicio", arreglaEmpleado)
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id") //traer dato de formulario
	fmt.Println(idEmpleado)               //imprimir por consola

	conexionEstablecida := conexionDB()
	registro, err := conexionEstablecida.Query("SELECT * FROM empleados WHERE id=?", idEmpleado)

	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{}

	for registro.Next() {
		var id int
		var nombre, correo string

		err = registro.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

	}
	//fmt.Println(empleado)
	plantillas.ExecuteTemplate(w, "editar", empleado)
}

//funcion para formulario crear
func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}

//funcion para insertar en db
func insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" { //traer datos de post
		nombre := r.FormValue("nombre") //nombre en formulario
		correo := r.FormValue("correo") //correo en formulario

		conexionEstablecida := conexionDB() //establecer conexion
		//insertar datos en db
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre,correo)VALUEs(?,?)")

		if err != nil {
			panic(err.Error())
		}
		//sentencia usada para insertar en db y parametros como van en sentencia sin cambiar de posicion
		insertarRegistros.Exec(nombre, correo)
		http.Redirect(w, r, "/", 301) //redireccionar
	}

}

//funcion actualizar
func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" { //traer datos de post
		id := r.FormValue("id")         //nombre en formulario
		nombre := r.FormValue("nombre") //nombre en formulario
		correo := r.FormValue("correo") //correo en formulario

		conexionEstablecida := conexionDB() //establecer conexion
		//insertar datos en db
		modificarRegistro, err := conexionEstablecida.Prepare("UPDATE empleados SET nombre=?, correo=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}
		//sentencia usada para insertar en db y parametros como van en sentencia sin cambiar de posicion
		modificarRegistro.Exec(nombre, correo, id)
		http.Redirect(w, r, "/", 301) //redireccionar
	}

}
