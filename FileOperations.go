package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

func crearArchivo(name string) error {

	//Verifica si el directorio existe
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Err CreateFile dir==", err)
		return err
	}

	//Crea el archivo

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

func abrirArchivo(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)

	if err != nil {
		fmt.Println("el error es: ", err)

		return nil, err
	}

	return file, nil
}

func escribirObjeto(file *os.File, disk interface{}, position int64) error {

	file.Seek(position, 0)

	err := binary.Write(file, binary.LittleEndian, disk)

	fmt.Println("**********ESCRIBIENDO EL ARCHIVO****************")
	if err != nil {
		return err
	}

	return nil
}

func LeerObjeto(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)

	err := binary.Read(file, binary.LittleEndian, data)

	if err != nil {
		fmt.Println("Err ReadObject==", err)
		//fmt.Println("ME ENCUENTRO POR ACA")
		return err
	}
	return nil
}
