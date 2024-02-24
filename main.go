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
	Funciones.Fdisk(10, "A", "Particion1", "b", " ", "bf", "FULL", 500) //(el tamano size sera en unit b= bytes, tipo particion= primaria)

	Funciones.Fdisk(20, "A", "Particion2", "b", "e", "bf", "FULL", 500) //(el tamano size sera en unit b= bytes,  tipo particion= extendida)

	Funciones.Fdisk(30, "A", "Particion3", "b", "e", "bf", "FULL", 500) //(el tamano size sera en unit b= bytes,  tipo particion= extendida)

	contador++

	//Agregar informacion a un Archivo sin sobreescribir el anterior
	//file, err:= os.OpenFile("archivo.txt", os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)

}
