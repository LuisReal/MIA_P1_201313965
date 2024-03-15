package Funciones

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CreateNewBlock(file *os.File, tempSuperblock Superblock, crrInode Inode, user string, group string, pass string) error {

	fmt.Println("\n\n========================= Inicio CreateNewBlock ===========================")

	var bloque int
	var index int
	var fileblock Fileblock
	var cadena string = " "
	var fileblock_start int32

	for i := 0; i < len(crrInode.I_block); i++ { //iterando bloques de inodo1

		if crrInode.I_block[i] != -1 {

			bloque = int(crrInode.I_block[i]) //obtiene el numero del ultimo bloque de archivos creado
			index = i

			fileblock_start = tempSuperblock.S_block_start + int32(bloque)*int32(binary.Size(Fileblock{}))

			if err := LeerObjeto(file, &fileblock, int64(fileblock_start)); err != nil { //bloque1
				return err
			}

			cadena += string(fileblock.B_content[:])
		}
	}
	fmt.Printf("\nel ultimo bloque creado es: %d, index: %d", bloque, index)
	fmt.Println()

	fmt.Println("Fileblock------------")
	//data := "1,G,root\n1,U,root,root,123\n"

	fmt.Println("\n Imprimiendo cadena: ", cadena)

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

		fmt.Println("\ndatos: ", datos)
		fmt.Println("\nLongitud de datos es : ", len(datos))

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

			}
		}

		var espacios int

		for i := 0; i < len(fileblock.B_content); i++ {

			if fileblock.B_content[i] == 0 {
				espacios++
			}
		}

		//data := "1,G,root\n1,U,root,root,123\n"
		if espacios > 0 {

			fmt.Println("\n Todavia sobra espacio despues de escribir la cadena en el slice")
			no_space = false

		} else { // si ya no hay espacios en el slice para ingresar la cadena

			cadena_restante := newCadena[c:]
			fmt.Println("\n cadena restante es: ", cadena_restante)

			no_space = true
		}

		if no_space { // si ya no existe espacio en el slice de fileblock.B_content (se crea un nuevo bloque)

			fmt.Println("\n\n ****Escribiendo objeto FILEBLOCK en el archivo *****")
			if err := EscribirObjeto(file, fileblock, int64(fileblock_start)); err != nil { //aqui solo escribi el primer EBR
				return err
			}

			fmt.Println("\n La longitud de la cadena newCadena[c] es: ", len(newCadena[c:]))

			if len(newCadena[c:]) != 0 { //si todavia hay caracteres en newCadena para seguir ingresando en slice de fileblock.Bcontent
				fmt.Println("\n      LLamando funcion CrearBloque .......")
				CrearBloque(newCadena, c, crrInode, tempSuperblock, file)
			}

			return nil
		}

		fmt.Println("\n El contenido nuevo de B_content es: ", string(fileblock.B_content[:]))

		fmt.Println("\n\n ********** Escribiendo objeto FILEBLOCK en el archivo ******************")
		if err := EscribirObjeto(file, fileblock, int64(fileblock_start)); err != nil { //aqui solo escribi el primer EBR
			return err
		}

	}

	fmt.Println("\n\n========================= Fin CreateNewBlock ===========================")

	return nil
}

func CrearBloque(newCadena string, contador int, crrInode Inode, tempSuperblock Superblock, file *os.File) {
	fmt.Println("\n\n========================= Inicio CrearBloque ===========================")

	fmt.Println("\n............Creando nuevo bloque de archivos................")
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
	//Escribiendo Inode1
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
	var last_index int

	//data := "1,G,root\n1,U,root,root,123\n"
	for i := 0; i < len(newFileblock.B_content); i++ {

		last_index = i
		if newFileblock.B_content[i] == 0 { // si hay espacio

			if c < len(resto_cadena) {
				//fileblock.B_content [2,U,usuarios,user2,    contra2sena]
				newFileblock.B_content[i] = byte(resto_cadena[c])

				c++

			}

		}
	}

	//obtiene cantidad de espacios restantes en el slice
	var espacios int

	for i := 0; i < len(newFileblock.B_content); i++ {

		if newFileblock.B_content[i] == 0 {
			espacios++
		}
	}

	fmt.Println("\n La cantidad de espacios restantes es: ", espacios+1)

	fmt.Println("\nLast_index: ", last_index+1)

	fmt.Println("\n Escribiendo newFileblock en el archivo..........")

	err = EscribirObjeto(file, newFileblock, int64(fileblock_start)) //Bloque 1

	if err != nil {
		fmt.Println("Error: ", err)
	}

	//validando si todavia hay espacio en el slice
	if newFileblock.B_content[last_index] == 0 {
		fmt.Println("\n todavia hay espacio ")

		var fileblock Fileblock

		//fileblock_start := tempSuperblock.S_block_start + crrInode.I_block[0]*int32(binary.Size(Fileblock{})) // bloque1
		fileblock_start := tempSuperblock.S_block_start + int32(newBlock)*int32(binary.Size(Fileblock{})) // bloque1

		if err := LeerObjeto(file, &fileblock, int64(fileblock_start)); err != nil { //bloque1
			return
		}

		fmt.Println("\n Imprimiendo fileblock.B_content: ", string(fileblock.B_content[:]))

		return
	} else {
		no_space = true
	}

	if no_space {

		//Recursivo para crear nuevo bloque
		CrearBloque(resto_cadena, c, crrInode, tempSuperblock, file)
	}

	fmt.Println("\n\n========================= Fin CrearBloque ===========================")

}
