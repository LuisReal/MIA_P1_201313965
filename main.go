package main

import (
	"MIA_P1_201313965/Funciones"
)

func main() {

	contador := 0

	size := 1
	fit := "b"
	unit := "m"

	abecedario := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	Funciones.Mkdisk(size, fit, unit, abecedario[contador])
	contador++
	Funciones.Mkdisk(size, fit, unit, abecedario[contador])
	contador++
	Funciones.Mkdisk(size, fit, unit, abecedario[contador])
	//Funciones.CrearMBR(size, fit, abecedario[contador])

	//size int, driveletter string, name string, unit string, type_ string, fit string, delete string, add int
	Funciones.Fdisk(10, "A", "Particion1", "b", " ", "bf", "", 0) //(el tamano size sera en unit b= bytes, tipo particion= primaria)

	Funciones.Fdisk(20, "A", "Particion2", "b", " ", "bf", "", 0) //(el tamano size sera en unit b= bytes,  tipo particion= extendida)

	Funciones.Fdisk(5, "A", "Particion3", "k", "e", "bf", "", 0) //(el tamano size sera en unit b= bytes,  tipo particion= extendida)

	//fdisk -size=1 -type=L -unit=M -fit=bf -driveletter=A -name="Particion3" ejemplo al crear una particion logica

	/*
		Funciones.Fdisk(1, "A", "ParticionLogica1", "k", "l", "bf", "", 0)
		Funciones.Fdisk(2, "A", "ParticionLogica2", "k", "l", "bf", "", 0)
		Funciones.Fdisk(1, "A", "ParticionLogica3", "k", "l", "bf", "", 0)*/

	//fdisk -delete=full -name="Particion2" -driveletter=A

	//Funciones.Fdisk(1, "A", "Particion3", "", "", "", "FULL", 0) // esto eliminara una particion

	//fdisk -add=500 -size=10 -unit=K -driveletter=D -name=”Particion4”
	Funciones.Fdisk(1, "A", "Particion3", "k", "", "", "", -5)

	Funciones.Mount("A", "Particion2")

	//Funciones.Fdisk(0, "A", "Particion3", "", " ", "", "full", 0) elimina una particion ya se rapida o completa

	//fdisk -delete=full -name="Particion1" -driveletter=A

	//Funciones.Fdisk(40, "A", "Particion2", "b", " ", "bf", "full", 500)

	//Agregar informacion a un Archivo sin sobreescribir el anterior
	//file, err:= os.OpenFile("archivo.txt", os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)

}
