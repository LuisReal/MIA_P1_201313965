package Funciones

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

var contador int = 0
var abecedario = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func getCommandAndParams(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		command := strings.ToLower(parts[0])
		params := strings.Join(parts[1:], " ")
		return command, params
	}
	return "", input
}

func Analyze() {

	var archivo *os.File

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

		command, params := getCommandAndParams(linea)

		fmt.Println("\nCommand: ", command, "Params: ", params)

		AnalyzeCommnad(command, params, contador)

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
	} else {
		fmt.Println("Error: Command not found")
	}
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
		case "size", "fit", "unit", "driveletter", "name", "type":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	//Funciones.Fdisk(10, "A", "Particion1", "b", " ", "bf", "", 0)
	// Call the function
	Fdisk(*size, *driveletter, *name, *unit, *type_, *fit, "", 0)
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
	Login(*user, *pass, *id)

}
