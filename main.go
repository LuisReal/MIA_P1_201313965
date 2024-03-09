package main

import (
	Funciones "MIA_P1_201313965/Funciones"
	"fmt"
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

	Funciones.Fdisk(10, "A", "Particion1", "b", " ", "bf", "", 0) //(el tamano size sera en unit b= bytes, tipo particion= primaria)

	Funciones.Fdisk(100, "A", "Particion2", "k", " ", "bf", "", 0) //(el tamano size sera en unit b= bytes,  tipo particion= extendida)

	Funciones.Fdisk(5, "A", "Particion3", "k", "e", "bf", "", 0) //(el tamano size sera en unit b= bytes,  tipo particion= extendida)

	Funciones.Fdisk(1, "A", "ParticionLogica1", "k", "l", "bf", "", 0)
	Funciones.Fdisk(2, "A", "ParticionLogica2", "k", "l", "bf", "", 0)
	Funciones.Fdisk(1, "A", "ParticionLogica3", "k", "l", "bf", "", 0)

	Funciones.Mount("A", "Particion2")

	Funciones.Mkfs("A265", "FULL", "2fs")

	//mkfile -size=15 -path=/home/user/docs/a.txt -r

	//Funciones.Mkfile("/home/user/docs/a.txt", "r", 15, "")

	//#Si no existen las carpetas home user o docs se crean
	//mkdir -r -path=/home/user/docs/usac
	//Funciones.Mkdir("/home/user/docs/usac", "r")

	//login -user=root -pass=123 -id=A118

	usuario, _ := Funciones.Login("root", "123", "A265")

	if usuario == "root" {
		fmt.Println("\n!!!!!!!!!!!!!!Usuario root encontrado!!!!!!!!!!!!!!!!!!")
	} else {
		fmt.Println("\n        El usuario es: ", usuario)
	}

}
