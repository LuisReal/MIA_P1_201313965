package Funciones

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var user_ User

var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

var contador int = 0
var abecedario = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func getCommandAndParams(input string) (string, string) {

	if input != " " {

		parts := strings.Fields(input)

		fmt.Println("\nImprimiendo parts: ", parts)
		if len(parts) > 0 {

			command := strings.ToLower(parts[0])
			params := strings.Join(parts[1:], " ")

			return command, params
		}
	}

	return "", input
}

func Analyze() {

	var archivo *os.File

	//se valida ejecucion de comando execute
	if len(os.Args) == 1 { // si no se pasa un argumento despues de go run main.go se ejecuta este if

		var input string
		fmt.Print("-> ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input = scanner.Text()

		command, params := getCommandAndParams(input)

		fmt.Println("Command: ", command, "Params: ", params)

		//input := bufio.NewScanner(archivo)

		execute := flag.NewFlagSet("execute", flag.ExitOnError)
		s := execute.String("path", "", "Ruta script")

		// Parse the flags
		execute.Parse(os.Args[1:])

		// find the flags in the input
		matches := re.FindAllStringSubmatch(params, -1)

		// Process the input
		for _, match := range matches {

			flagName := match[1]
			flagValue := match[2]

			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "path":
				execute.Set(flagName, flagValue)
			default:
				fmt.Println("Error: Flag not found")
			}
		}

		ruta := *s
		fmt.Println("\nLa ruta ingresada es: ", ruta)
		archivo_, err := os.Open(ruta)

		archivo = archivo_
		if err != nil {
			fmt.Println("Error al abrir el archivo:", err)
			return
		}
	}

	scanner := bufio.NewScanner(archivo)

	for scanner.Scan() {

		linea := scanner.Text()

		if linea != "" {
			command, params := getCommandAndParams(linea)
			comentario := command[0] // obtiene solo la primera letra de command

			if string(comentario) != "#" {
				fmt.Println("\nCommand: ", command, "Params: ", params)
				AnalyzeCommnad(command, params, contador)
			}
		} else {
			fmt.Println("\nLinea vacia")
		}

	}

	//execute -path=/home/darkun/Escritorio/scripts.sdaa

	//mkdisk -size=3000 -unit=K
	//fdisk -size=300 -driveletter=A -name=Particion1
	//mount -driveletter=A -name=Part1
	//mkfile -size=15 -path=/home/user/docs/a.txt -r

}

func AnalyzeCommnad(command string, params string, contador int) {

	if strings.Contains(command, "mkdisk") {

		fmt.Println("\n       El valor del contador es: ", contador)
		fmt.Println()

		bn_mkdisk(params, abecedario[contador])
		contador++

	} else if strings.Contains(command, "fdisk") {
		bn_fdisk(params)
	} else if strings.Contains(command, "mount") {
		bn_mount(params)
	} else if strings.Contains(command, "mkfs") {
		bn_mkfs(params)
	} else if strings.Contains(command, "login") {
		bn_login(params)
	} else if strings.Contains(command, "mkgrp") {
		bn_mkgrp(params)
	} else if strings.Contains(command, "rmgrp") {
		bn_rmgrp(params)
	} else if strings.Contains(command, "mkusr") {
		bn_mkusr(params)
	} else if strings.Contains(command, "rmusr") {
		bn_rmusr(params)
	} else if strings.Contains(command, "logout") {
		bn_logout()
	} else if strings.Contains(command, "rep") {
		bn_reportes(params)
	} else {
		fmt.Println("Error: Command not found")
	}
}

func bn_reportes(params string) {

	fs := flag.NewFlagSet("rep", flag.ExitOnError)
	name := fs.String("name", "", "Name")
	path := fs.String("path", "f", "Path")
	id := fs.String("id", "m", "Id")
	ruta := fs.String("ruta", "m", "ruta")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(params, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		//flagValue := strings.ToLower(match[2])
		flagValue := match[2]
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "name", "path", "id", "ruta":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	// Call the function
	Reportes(*name, *path, *id, *ruta)

}
func bn_mkdisk(params string, letra string) {
	// Define flags
	fs := flag.NewFlagSet("mkdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño")
	fit := fs.String("fit", "f", "Ajuste")
	unit := fs.String("unit", "m", "Unidad")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(params, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	// Call the function
	Mkdisk(*size, *fit, *unit, letra)

}

func bn_fdisk(input string) {
	// Define flags
	fs := flag.NewFlagSet("fdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño")
	driveletter := fs.String("driveletter", "", "Letra")
	name := fs.String("name", "", "Nombre")
	unit := fs.String("unit", "m", "Unidad")
	type_ := fs.String("type", "p", "Tipo")
	fit := fs.String("fit", "f", "Ajuste")
	delete := fs.String("delete", "", "Elimina particion")

	input_ := strings.Split(input, " ")
	fmt.Println("\nImprimiendo SLICE input: ", input_)
	var formateo string

	fmt.Println("\nImprimendo input_[1]: ", input_[1])
	if input_[0] == "-delete=full" {
		formateo = "rapido"
	} else {

		for i := 1; i < len(input_); i++ {

			if input_[i] == "-delete=full" {
				formateo = "completo"
				break
			}
		}

	}

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)
	fmt.Println("\nImprimiendo matches: ", matches)
	// Process the input

	for _, match := range matches {
		flagName := match[1]

		//fmt.Println("\nmatch[1]: ", match[1])
		//fmt.Println("\nmatch[2]: ", match[2])
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit", "driveletter", "name", "type", "delete":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	//Funciones.Fdisk(10, "A", "Particion1", "b", " ", "bf", "", 0)
	// Call the function

	fmt.Println("\nImprimiendo valor de formateo: ", formateo)
	Fdisk(*size, *driveletter, *name, *unit, *type_, *fit, *delete, 0, formateo)
}

func bn_mount(input string) {
	// Define flags
	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	driveletter := fs.String("driveletter", "", "Letra")
	name := fs.String("name", "", "Nombre")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "driveletter", "name":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	// Call the function
	Mount(*driveletter, *name)
}

func bn_mkfs(input string) {
	// Define flags
	fs := flag.NewFlagSet("mkfs", flag.ExitOnError)
	id := fs.String("id", "", "Id")
	type_ := fs.String("type", "", "Tipo")
	fs_ := fs.String("fs", "2fs", "Fs")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id", "type", "fs":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	// Call the function
	Mkfs(*id, *type_, *fs_)

}

func bn_login(input string) {
	// Define flags
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	user := fs.String("user", "", "Usuario")
	pass := fs.String("pass", "", "Contraseña")
	id := fs.String("id", "", "Id")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "user", "pass", "id":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	// Call the function
	usuario, err := Login(*user, *pass, *id)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	/*EL usuario root puede ejecutar los siguientes comandos:
	MKGRP
	RMGRP
	MKUSR
	RMUSR
	*/

	user_.Nombre = usuario
	user_.Status = true
	user_.Id = *id
}

func bn_logout() {
	fmt.Println("\n\n========================= Iniciando Logout =========================")

	fmt.Println("\n Cerrando sesion de usuario: ", user_.Nombre)
	user_.Nombre = ""
	user_.Status = false
	user_.Id = ""

	fmt.Println("\n\n========================= Finalizando Logout =========================")
}

func bn_mkgrp(input string) {

	if user_.Nombre == "root" && user_.Status { //si el usuario es root y esta logueado(true)
		// Define flags
		fs := flag.NewFlagSet("mkgrp", flag.ExitOnError)
		name := fs.String("name", "", "nombre de grupo")

		// Parse the flags
		fs.Parse(os.Args[1:])

		// find the flags in the input
		matches := re.FindAllStringSubmatch(input, -1)

		// Process the input
		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "name":
				fs.Set(flagName, flagValue)
			default:
				fmt.Println("Error: Flag not found")
			}
		}

		Mkgrp(*name, user_.Id)

	} else {
		fmt.Println("\n\n******************Necesita iniciar sesion como ususario ROOT***********************")
		return
	}
}

func bn_rmgrp(input string) {

	if user_.Nombre == "root" && user_.Status { //si el usuario es root y esta logueado(true)
		// Define flags
		fs := flag.NewFlagSet("rmgrp", flag.ExitOnError)
		name := fs.String("name", "", "nombre de grupo")

		// Parse the flags
		fs.Parse(os.Args[1:])

		// find the flags in the input
		matches := re.FindAllStringSubmatch(input, -1)

		// Process the input
		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "name":
				fs.Set(flagName, flagValue)
			default:
				fmt.Println("Error: Flag not found")
			}
		}

		Rmgrp(*name, user_.Id)

	} else {
		fmt.Println("\n\n******************Necesita iniciar sesion como ususario ROOT para poder REMOVER un grupo***********************")
		return
	}
}

func bn_mkusr(input string) {

	if user_.Nombre == "root" && user_.Status { //si el usuario es root y esta logueado(true)
		// Define flags
		fs := flag.NewFlagSet("mkusr", flag.ExitOnError)
		user := fs.String("user", "", "nombre de usuario")
		pass := fs.String("pass", "", "contrasena de usuario")
		group := fs.String("grp", "", "grupo al que pertenecera el usuario")

		// Parse the flags
		fs.Parse(os.Args[1:])

		// find the flags in the input
		matches := re.FindAllStringSubmatch(input, -1)

		// Process the input
		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "user", "pass", "grp":
				fs.Set(flagName, flagValue)
			default:
				fmt.Println("Error: Flag not found")
			}
		}

		Mkusr(*user, *pass, *group, user_.Id)

	} else {
		fmt.Println("\n\n******************Necesita iniciar sesion como ususario ROOT para poder crear un usuario ***********************")
		return
	}
}

func bn_rmusr(input string) {

	if user_.Nombre == "root" && user_.Status { //si el usuario es root y esta logueado(true)
		// Define flags
		fs := flag.NewFlagSet("rmusr", flag.ExitOnError)
		user := fs.String("user", "", "nombre de usuario")

		// Parse the flags
		fs.Parse(os.Args[1:])

		// find the flags in the input
		matches := re.FindAllStringSubmatch(input, -1)

		// Process the input
		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "user":
				fs.Set(flagName, flagValue)
			default:
				fmt.Println("Error: Flag not found")
			}
		}

		Rmusr(*user, user_.Id)

	} else {
		fmt.Println("\n\n******************Necesita iniciar sesion como ususario ROOT para poder REMOVER un grupo***********************")
		return
	}
}
