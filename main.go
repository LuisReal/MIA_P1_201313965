package main

//"bufio"
//"encoding/binary"
//"flag"

//"time"

func main() {

	var disk MBR
	disk.Mbr_tamano = 1

	archivo_binario := "./binario.dsk" // verificar que sea extension dsk

	if err := crearArchivo(archivo_binario); err != nil {
		return
	}

	//Agregar informacion a un Archivo sin sobreescribir el anterior
	//file, err:= os.OpenFile("archivo.txt", os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)

	//ABRIENDO EL ARCHIVO BINARIO
	/*
		file_binario, err := abrirArchivoBinario(archivo_binario)

		if err != nil {
			return
		}*/

}
