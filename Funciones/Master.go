package Funciones

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
)

//PARA PRIMARIA Y EXTENDIDA SOLO SE VA A USAR EL MBR

func Mkdisk(size int, fit string, unit string, letra string) {

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
	err := crearArchivo("./archivos/" + letra + ".dsk")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// Open bin file
	file, err := abrirArchivo("./archivos/" + letra + ".dsk")
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

func CrearMBR(size int, fit string, letra string) {

	fmt.Println("***********CREANDO MBR Y ESCRIBIENDO EN EL ARCHIVO BINARIO******************")

	//Abriendo el archivo para usarlo y escribir el MBR
	file, err := abrirArchivo("./archivos/" + letra + ".dsk")
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

//FUNCION QUE ADMINISTRA LAS PARTICIONES

func Fdisk(size int, driveletter string, name string, unit string, type_ string, fit string, delete string, add int) {

	// validando que el tamano sea mayor que cero
	if size <= 0 {
		fmt.Println("Error: El tamano(size) debe ser mayor a cero")
		return
	}

	//Abriendo el archivo para usarlo y escribir el MBR
	file, err := abrirArchivo("./archivos/" + driveletter + ".dsk")
	if err != nil {
		return
	}

	var TemporalMBR MBR
	// Read object from bin file
	if err := LeerObjeto(file, &TemporalMBR, 0); err != nil {
		return
	}

	// validar unit sea igual a b/k/m
	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Println("Error: Unit must be b, k or m")
		return
	}

	// Configurar el size en bytes

	if unit == "k" {
		size = size * 1024
	} else if unit == "m" {
		size = size * 1024 * 1024
	}

	// valida el type puede ser (p=primaria e=extendida l=logica)

	if type_ == " " { // si el usuario no indica el type este sera primaria por defecto
		type_ = "p"
		fmt.Println("type sera por defecto p=primaria")
	} else if type_ != "p" && type_ != "e" && type_ != "l" {
		fmt.Println("Error: Type must be p, e or l")
		return
	}

	var count = 0
	var gap = int32(0)
	// Iterate over the partitions
	for i := 0; i < 4; i++ {
		if TemporalMBR.Mbr_partitions[i].Part_size != 0 {
			count++
			gap = TemporalMBR.Mbr_partitions[i].Part_start + TemporalMBR.Mbr_partitions[i].Part_size
		}
	}

	for i := 0; i < 4; i++ {
		if TemporalMBR.Mbr_partitions[i].Part_size == 0 { // si la particion no esta creada(por defecto tiene el part_size tiene valor 0)

			TemporalMBR.Mbr_partitions[i].Part_size = int32(size)

			if count == 0 {
				TemporalMBR.Mbr_partitions[i].Part_start = int32(binary.Size(TemporalMBR))
			} else {
				TemporalMBR.Mbr_partitions[i].Part_start = gap
			}

			TemporalMBR.Mbr_partitions[i].Part_status = true

			byteString_name := make([]byte, 16)
			byteString_fit := make([]byte, 1)
			byteString_type := make([]byte, 1)
			copy(byteString_name, name)
			copy(byteString_fit, fit)
			copy(byteString_type, type_)

			TemporalMBR.Mbr_partitions[i].Part_name = [16]byte(byteString_name)
			TemporalMBR.Mbr_partitions[i].Part_fit = [1]byte(byteString_fit)
			TemporalMBR.Mbr_partitions[i].Part_type = [1]byte(byteString_type)

			//copy(TemporalMBR.Mbr_partitions[i].Part_name[:], name)
			//copy(TemporalMBR.Mbr_partitions[i].Part_fit[:], fit) // el fit por defecto tomara el primer ajuste
			//copy(TemporalMBR.Mbr_partitions[i].Part_type[:], type_)
			TemporalMBR.Mbr_partitions[i].Part_correlative = int32(count + 1)
			break

		}
	}

	// Sobreescribe el MBR
	if err := escribirObjeto(file, TemporalMBR, 0); err != nil {
		return
	}

	var TemporalMBR2 MBR

	// Read object from bin file
	if err := LeerObjeto(file, &TemporalMBR2, 0); err != nil { //Leera el objeto desde la posicion 0
		return
	}

	// Print object
	PrintMBR(TemporalMBR2)

	// Close bin file
	defer file.Close()

}
