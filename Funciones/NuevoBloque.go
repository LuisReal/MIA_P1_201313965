package Funciones

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CreateNewBlock(file *os.File, tempSuperblock Superblock, crrInode Inode, user string, group string, pass string) error {

	var bloque int
	var index int

	for i := 0; i < len(crrInode.I_block); i++ { //iterando bloques de inodo1

		if crrInode.I_block[i] != -1 {
			bloque = int(crrInode.I_block[i]) //obtiene el numero del ultimo bloque de archivos creado
			index = i
		}
	}
	fmt.Printf("\nbloque: %d, index: %d", bloque, index)
	fmt.Println()

	var fileblock Fileblock

	//fileblock_start := tempSuperblock.S_block_start + crrInode.I_block[0]*int32(binary.Size(Fileblock{})) // bloque1
	fileblock_start := tempSuperblock.S_block_start + int32(bloque)*int32(binary.Size(Fileblock{})) // bloque1

	if err := LeerObjeto(file, &fileblock, int64(fileblock_start)); err != nil { //bloque1
		return err
	}

	fmt.Println("Fileblock------------")
	//data := "1,G,root\n1,U,root,root,123\n"

	var cadena string = " "

	cadena = string(fileblock.B_content[:])

	fmt.Println("\n Imprimiendo cadena: ", string(fileblock.B_content[:]))

	lines := strings.Split(cadena, "\n")

	if len(lines) > 0 {
		lines = lines[:len(lines)-1]
	}

	fmt.Println("\n\nContenido del arreglo lines: ", lines)
	fmt.Println("\nEl tamano del arreglo lines es: ", len(lines))

	fmt.Println("\nImprimiendo ultimo elemento de arreglo lines: ", lines[len(lines)-1])
	//2, G, usuarios, \n
	var num_group int = 0
	var exist int = 0
	var datos []string
	//var linea_ int

	for i := 0; i < len(lines); i++ {

		datos = strings.Split(lines[i], ",")

		if len(datos) != 0 {

			if len(datos) > 3 {
				if string(datos[3]) == user {
					fmt.Println("\n EL usuario a crear ya existe")
					return nil
				}
			}

		}
	}

	for i := 0; i < len(lines); i++ {

		datos = strings.Split(lines[i], ",")

		num_group_, _ := strconv.Atoi(datos[0]) // contiene el numero de grupo

		num_group = num_group_

		if len(datos) != 0 {

			if string(datos[2]) == group {

				//2,U,usuarios,user1,usuario\n
				if num_group == 0 {
					fmt.Println("\nEl grupo no existe porque ya fue eliminado anteriormente")
					return nil
				} else {
					fmt.Println("\n\n      ********** El grupo si existe ************")

					exist++ // verifica que el grupo exista

					break
				}
			}
		}

	}

	if exist != 0 { // si el grupo donde se creara el usuario existe

		newCadena := strconv.Itoa(num_group) + ",U," + group + "," + user + "," + pass + "\n"

		fmt.Println("\n ********datos de la variable newCadena: ", newCadena)

		//Agregando nuevo usuario a users.txt en fileblock.B_content
		var c int
		var no_space bool

		for i := 0; i < len(fileblock.B_content); i++ {
			//fmt.Println(fileblock[i])

			if fileblock.B_content[i] == 0 { // si hay todavia espacio

				if c < len(newCadena) {

					fileblock.B_content[i] = byte(newCadena[c])
					//fmt.Printf("agregando letra:  %s   ", string(newCadena[c]))
					c++

				} else {
					break
				}

			} else {

				no_space = true

			}
		}

		if no_space { // si ya no existe espacio en el slice de fileblock.B_content (se crea un nuevo bloque)

			fmt.Println("\n\n ********** Escribiendo objeto FILEBLOCK en el archivo ******************")
			if err := EscribirObjeto(file, fileblock, int64(fileblock_start)); err != nil { //aqui solo escribi el primer EBR
				return err
			}

			CrearBloque(newCadena, c, crrInode, tempSuperblock, file)

		}

		fmt.Println("\n El contenido nuevo de B_content es: ", string(fileblock.B_content[:]))

		fmt.Println("\n\n ********** Escribiendo objeto FILEBLOCK en el archivo ******************")
		if err := EscribirObjeto(file, fileblock, int64(fileblock_start)); err != nil { //aqui solo escribi el primer EBR
			return err
		}

	}

	return nil
}

func CrearBloque(newCadena string, contador int, crrInode Inode, tempSuperblock Superblock, file *os.File) {

	fmt.Println("\n\n............Creando nuevo bloque de archivos")
	fmt.Println("\nLa cadena faltante es: ", newCadena[contador:])

	resto_cadena := newCadena[contador:]

	var bloque int
	var index int

	for i := 0; i < len(crrInode.I_block); i++ {

		if crrInode.I_block[i] != -1 {
			bloque = int(crrInode.I_block[i]) //obtiene el numero del ultimo bloque creado
			index = i
		}
	}

	newBlock := bloque + 1
	crrInode.I_block[index+1] = int32(newBlock)

	err := EscribirObjeto(file, crrInode, int64(tempSuperblock.S_inode_start+int32(binary.Size(Inode{})))) //Inode 1

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("\nCreando Nuevo bloque numero: ", newBlock)
	//Creando nuevo bloque de archivos

	var newFileblock Fileblock

	//newFileblock.B_content

	fileblock_start := tempSuperblock.S_block_start + int32(newBlock)*int32(binary.Size(Fileblock{})) // bloque1

	var c int
	var no_space bool

	for i := 0; i < len(newFileblock.B_content); i++ {
		//fmt.Println(fileblock[i])

		if newFileblock.B_content[i] == 0 { // si hay espacio

			if c < len(resto_cadena) {
				//fileblock.B_content [2,U,usuarios,user2,    contra2sena]
				newFileblock.B_content[i] = byte(resto_cadena[c])
				//fmt.Printf("agregando letra:  %s   ", string(newCadena[c]))
				c++

			} else {
				break
			}

		} else {

			no_space = true

		}
	}

	if no_space {
		err := EscribirObjeto(file, newFileblock, int64(fileblock_start)) //Bloque 1

		if err != nil {
			fmt.Println("Error: ", err)
		}

		//Recursivo para crear nuevo bloque
		CrearBloque(resto_cadena, c, crrInode, tempSuperblock, file)
	}

	err = EscribirObjeto(file, newFileblock, int64(fileblock_start)) //Bloque 1

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("\nExistSpace: ", no_space)
	//fmt.Println("\nEl valor de (c) es: ", c)
}
