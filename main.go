package main

import (
	"MIA_P1_201313965/Funciones"
)

func main() {

	contador := 0

	size := 1
	fit := "b"
	unit := "k"

	abecedario := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	Funciones.Mkdisk(size, fit, unit, abecedario[contador])
	Funciones.CrearMBR(size, fit, abecedario[contador])

	//size int, driveletter string, name string, unit string, type_ string, fit string, delete string, add int
	Funciones.Fdisk(0, "A", "disco", "B", "P", "B", "FULL", 500)

	contador++

	//Agregar informacion a un Archivo sin sobreescribir el anterior
	//file, err:= os.OpenFile("archivo.txt", os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)

}
