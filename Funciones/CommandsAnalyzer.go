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

	for {
		var input string
		fmt.Println("Enter command: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input = scanner.Text()

		command, params := getCommandAndParams(input)

		fmt.Println("Command: ", command, "Params: ", params)

		AnalyzeCommnad(command, params, contador)

	}
}

func AnalyzeCommnad(command string, params string, contador int) {

	if strings.Contains(command, "mkdisk") {

		fmt.Println("\n       El valor del contador es: ", contador)
		fmt.Println()

		bn_mkdisk(params, abecedario[contador])
		contador++

	} else if strings.Contains(command, "fdisk") {
		bn_fdisk(params)
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
