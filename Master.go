//Here we have all the functions we need to manipulete and create what we need

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func mkdisk(size int, fit string, unit string) {

	fmt.Println("***********CREANDO ARCHIVO DSK (FUNCION MKDISK)*************")
	fmt.Println("Size:", size, " Fit: ", fit, " Unit: ", unit)

	// validando que el tamano sea mayor que cero
	if size <= 0 {
		fmt.Println("Error: El tamano(size) debe ser mayor a cero")
		return
	}

	// validando que el ajuste ingresado por el usuario sea el correcto
	if fit != "b" && fit != "w" && fit != "f" {
		fmt.Println("Error: Ingrese el ajuste correcto")
		return
	}

	// Validando que las unidades ingresadas por el usuario esten correctas
	if unit != "k" && unit != "m" {
		fmt.Println("Error: La unidad(unit) debe ser k o m")
		return
	}

	// Configurando el tamano en bytes
	if unit == "k" {
		size = size * 1024
	} else if unit == "m" {
		size = size * 1024 * 1024
	} else {
		fmt.Println("La unidad ingresada no es correcta")
	}

	// Creando el archivo
	err := crearArchivo("./archivos/A.dsk")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// Open bin file
	file, err := abrirArchivo("./archivos/A.dsk")
	if err != nil {
		return
	}

	//Creando el archivo binario con ceros

	for i := 0; i < size; i++ {
		err := escribirObjeto(file, byte(0), int64(i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	defer file.Close()

	fmt.Println("**********FINALIZO CREACION ARCHIVO EN FUNCION MKDISK****************")

}

func crearMBR(size int, fit string) {

	fmt.Println("***********CREANDO MBR Y ESCRIBIENDO EN EL ARCHIVO BINARIO******************")

	//Abriendo el archivo para usarlo y escribir el MBR
	file, err := abrirArchivo("./archivos/A.dsk")
	if err != nil {
		return
	}

	// Creando un nuevo objeto MBR
	var mbr MBR
	mbr.Mbr_tamano = int32(size)                  //number :=
	mbr.Mbr_dsk_signature = int32(rand.Intn(100)) // random

	copy(mbr.Dsk_fit[:], fit) //convierte de string a byte

	date := time.Now()
	//fmt.Println("La Fecha y Hora Actual es: ", date.Format("2006-01-02 15:04:05"))

	byteString := make([]byte, 16)
	copy(byteString, date.Format("2006-01-02 15:04:05")) //convierte de string(la fecha) a bytestring
	mbr.Mbr_fecha_creacion = [16]byte(byteString)

	// Escribiendo el objeto en el archivo binario
	if err := escribirObjeto(file, mbr, 0); err != nil {
		return
	}

	var TempMBR MBR
	// Leyendo el objeto del archivo binario
	if err := LeerObjeto(file, &TempMBR, 0); err != nil {
		return
	}

	// Imprimiendo el objeto (esta funcion se encuentra en MBR.go)
	PrintMBR(TempMBR)

	// cerrando el archivo binario
	defer file.Close()

	fmt.Println("**************FINALIZANDO CREACION DE MBR*****************")

}
