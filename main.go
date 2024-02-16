package main

import (
	//"bufio"
	//"encoding/binary"
	//"flag"

	//"math/rand"
	"fmt"
	"os"
	//"time"
)

func main() {

	var disk MBR
	disk.Mbr_tamano = 1

	archivo_binario := "./binario.bin" // verificar que sea extension dsk

	if err := crearArchivo(archivo_binario); err != nil {
		return
	}

	//ABRIENDO EL ARCHIVO BINARIO
	/*
		file_binario, err := abrirArchivoBinario(archivo_binario)

		if err != nil {
			return
		}*/

}

func crearArchivo(name string) error {

	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)

		if err != nil {
			fmt.Println(err)
			return err
		}

		defer file.Close()
	}
	return nil
}

func abrirArchivoBinario(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)

	if err != nil {
		fmt.Println("el error es: ", err)

		return nil, err
	}

	return file, nil
}
