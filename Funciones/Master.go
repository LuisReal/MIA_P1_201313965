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

	CrearMBR(size, fit, letra)

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
		//fmt.Println("type sera por defecto p=primaria")
	} else if type_ != "p" && type_ != "e" && type_ != "l" {
		fmt.Println("Error: Type must be p, e or l")
		return
	}

	var TemporalMBR MBR
	// Read object from bin file
	if err := LeerObjeto(file, &TemporalMBR, 0); err != nil {
		return
	}

	//ELimina una particion si recibe un FULL en parametro delete
	if delete == "full" {

		var name_bytes [16]byte

		copy(name_bytes[:], []byte(name)) // esto convierte de string a []byte y luego a [16]byte

		cont_existPartition := 0

		for i := 0; i < 4; i++ {

			if name_bytes == TemporalMBR.Mbr_partitions[i].Part_name {

				if TemporalMBR.Mbr_partitions[i].Part_size == 0 {
					fmt.Println("\n************LA PARTICION NO EXISTE***************")
					fmt.Println()

				} else {
					fmt.Println("Eliminando la particion: ", string(TemporalMBR.Mbr_partitions[i].Part_name[:]))
					fmt.Println("\n         ***********Formateo rapido************")

					TemporalMBR.Mbr_partitions[i].Part_size = int32(0)
					TemporalMBR.Mbr_partitions[i].Part_name = [16]byte([]byte(""))
					TemporalMBR.Mbr_partitions[i].Part_type = [1]byte([]byte(""))
					TemporalMBR.Mbr_partitions[i].Part_correlative = int32(0)
					TemporalMBR.Mbr_partitions[i].Part_fit = [1]byte([]byte(""))
					TemporalMBR.Mbr_partitions[i].Part_id = int32(0)
					TemporalMBR.Mbr_partitions[i].Part_status = false
					TemporalMBR.Mbr_partitions[i].Part_type = [1]byte([]byte(""))

					//En este espacio va el codigo que llena el espacio con ceros en el archivo
				}
			}

		}

		for i := 0; i < 4; i++ {
			//fmt.Printf("\n************Part_name es: %s ***********\n", string(TemporalMBR.Mbr_partitions[i].Part_name[:]))
			//fmt.Printf("\n************El name ingresado es: %s\n", name)

			//name_bytes = []byte(name)

			if name_bytes == TemporalMBR.Mbr_partitions[i].Part_name {

				fmt.Printf("\n\n*******************************Eliminando la particion: %s...........\n", name)

				start := TemporalMBR.Mbr_partitions[i].Part_start

				//gap := TemporalMBR.Mbr_partitions[i].Part_start + TemporalMBR.Mbr_partitions[i].Part_size

				if err := escribirObjeto(file, byte(0), int64(start)); err != nil { //esto elimina donde empieza la particion y hasta el final del archivo.
					return
				}
				/*
					if err := escribirObjeto(file, TemporalMBR, int64(gap)); err != nil {
						return
					}*/

				/*
					TemporalMBR.Mbr_partitions[i].Part_correlative = int32(0)
					TemporalMBR.Mbr_partitions[i].Part_fit = [1]byte(byteString_fit)
					TemporalMBR.Mbr_partitions[i].Part_id = int32(0)
					TemporalMBR.Mbr_partitions[i].Part_name = [16]byte(byteString_name)
					TemporalMBR.Mbr_partitions[i].Part_size = int32(0)
					TemporalMBR.Mbr_partitions[i].Part_start = int32(0)
					TemporalMBR.Mbr_partitions[i].Part_status = false
					TemporalMBR.Mbr_partitions[i].Part_type = [1]byte(byteString_type)*/

			}

			if string(TemporalMBR.Mbr_partitions[i].Part_type[:]) == "e" {
				fmt.Println("\n**************aumentando el contador cont_existPartition++")
				cont_existPartition++
			}
		}

		if cont_existPartition == 0 {
			fmt.Println("*************\nLa particion NO existe************")
			fmt.Println()
			return
		}

		var TemporalMBR3 MBR
		if err := LeerObjeto(file, &TemporalMBR3, 0); err != nil {
			return
		}
		// Print object
		//PrintMBR(TemporalMBR3)
		//fmt.Println("\n\nLa particion se ha eliminado con exito")
		return
	}

	// Read object from bin file

	cont_extendida := 0

	//verifica que solo haya una particion extendida (ya que no pueden haber 2 o mas extendidas)
	for i := 0; i < 4; i++ {
		if string(TemporalMBR.Mbr_partitions[i].Part_type[:]) == "e" {

			cont_extendida++
		}

	}

	if cont_extendida == 1 {

		if type_ == "e" {
			fmt.Println("\n**************************No puede haber mas de una particion extendida en el disco************************")
			fmt.Println()
			return
		}

	}

	var count = 0
	var gap = int32(0)
	// Iterate over the partitions

	if type_ != "l" {

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

				var TempEBR EBR
				if type_ == "e" { // se crea el primer EBR al crear una extendida
					fmt.Println("\n\n            -----------------Creando el primer EBR---------------------------")
					if err := escribirObjeto(file, TempEBR, int64(TemporalMBR.Mbr_partitions[i].Part_start)); err != nil {
						return
					}
					fmt.Println("\n           -----------------Finalizo Creacion del primer EBR---------------------------")
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
	} else if type_ == "l" { // AQUI SE CREAN PARTICIONES LOGICAS
		// si el tipo de particion ingresada es LOGICA
		var name_bytes [16]byte
		copy(name_bytes[:], []byte(name))

		var fit_bytes [16]byte
		copy(fit_bytes[:], []byte(fit))

		/*
			if count == 0 {
				TemporalMBR.Mbr_partitions[i].Part_start = int32(binary.Size(TemporalMBR))
			} else {
				TemporalMBR.Mbr_partitions[i].Part_start = gap
			}*/

		var TempEBR2 EBR

		for i := 0; i < 4; i++ {

			if string(TemporalMBR.Mbr_partitions[i].Part_type[:]) == "e" { // primero busca una extendida

				start := TemporalMBR.Mbr_partitions[i].Part_start // donde inicia la particion extendida

				//Recupera toda la informacion escrita en el EBR(solo el EBR recupera) en la particion extendida en el archivo binario
				if err := LeerObjeto(file, &TempEBR2, int64(start)); err != nil {
					return
				}

				var gap int32

				if TempEBR2.Part_size != int32(0) { // Valida si existe una particion logica (en los datos del EBR)

					fmt.Println("\n==========================*****RECUPERANDO Y LEYENDO EL EBR SIGUIENTE*****=================================")

					gap = TempEBR2.Part_next + int32(binary.Size(TempEBR2))

					if err := LeerObjeto(file, &TempEBR2, int64(TempEBR2.Part_next)); err != nil {
						return
					}

					TempEBR2.Part_start = gap
					TempEBR2.Part_size = int32(size)
					TempEBR2.Part_name = name_bytes
					TempEBR2.Part_fit = fit_bytes
					TempEBR2.Part_mount = false
					TempEBR2.Part_next = TempEBR2.Part_start + int32(size) // donde empieza el siguiente EBR

					if err := escribirObjeto(file, TempEBR2, int64(TempEBR2.Part_start)); err != nil { //aqui solo escribi el EBR
						return
					}

					PrintEBR(TempEBR2)
					fmt.Println("\n==========================*****FINALIZANDO RECUPERACION y LECTURA DEl EBR SIGUIENTE*****=================================")

				} else {
					TempEBR2.Part_size = int32(size)
					TempEBR2.Part_start = int32(binary.Size(TempEBR2)) // la primera particion logica se coloca donde termina el primer EBR
					TempEBR2.Part_name = name_bytes
					TempEBR2.Part_fit = fit_bytes
					TempEBR2.Part_mount = false
					TempEBR2.Part_next = TempEBR2.Part_start + int32(size) // donde empieza el siguiente EBR

					if err := escribirObjeto(file, TempEBR2, int64(start)); err != nil { //aqui solo escribi el EBR
						return
					}
				}

				//ESCRIBIENDO EL SIGUIENTE EBR
				var TempEBRnext EBR

				if err := escribirObjeto(file, TempEBRnext, int64(TempEBR2.Part_next)); err != nil { //aqui solo escribi el siguiente EBR (con info vacia)
					return
				}

				/*
					fmt.Println("\n==========================*****LEYENDO EL EBR*****=================================")
					var TemporalEBR2 EBR
					if err := LeerObjeto(file, &TemporalEBR2, int64(start)); err != nil {
						return
					}

					PrintEBR(TemporalEBR2)
					fmt.Println("\n==========================*****FINALIZANDO LECTURA DEL EBR*****=================================")*/

				//------------------------------AQUI SE LEE EL SIGUENTE EBR---------------------------------------------

				fmt.Println("\n==========================*****LEYENDO EL SIGUIENTE EBR*****=================================")

				var TemporalEBR3 EBR
				if err := LeerObjeto(file, &TemporalEBR3, int64(TempEBR2.Part_next)); err != nil {
					return
				}

				PrintEBR(TemporalEBR3)
				fmt.Println("\n==========================*****FINALIZANDO LECTURA DEL SIGUIENTE EBR*****=================================")
			}
		}
	}

	// Sobreescribe el MBR
	if err := escribirObjeto(file, TemporalMBR, 0); err != nil {
		return
	}

	var TemporalMBR2 MBR
	if err := LeerObjeto(file, &TemporalMBR2, 0); err != nil {
		return
	}

	fmt.Println()
	fmt.Println()
	// Print object
	PrintMBR(TemporalMBR2)

	// Close bin file
	defer file.Close()

}
