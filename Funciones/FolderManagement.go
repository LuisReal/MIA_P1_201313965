package Funciones

import (
	"fmt"
	"strings"
)

func Mkdir(path string, r string) {
	fmt.Println("\n\n=========================Creando Carpeta (MKDIR)===========================")
	fmt.Println("\n**********************El path ingresado es: ", path)

	carpetas := strings.Split(path, "/")

	for i := 0; i < len(carpetas); i++ {
		fmt.Println("\n  Carpeta: ", carpetas[i]) // el primer valor no contiene nada en la posicion 0
	}

	fmt.Println("\n\n=======================Finalizando Creacion Carpeta (MKDIR)===========================")
}
