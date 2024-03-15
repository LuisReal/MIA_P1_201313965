package Funciones

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Login(user string, pass string, id string) (string, error) {
	fmt.Println("\n\n========================= LOGIN ===========================")

	fmt.Printf("\nUser: %s, pass: %s, id: %s\n", user, pass, id)

	driveletter := string(id[0])

	// Open bin file
	filepath := "./archivos/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := AbrirArchivo(filepath)
	if err != nil {
		return "", err
	}

	var TempMBR MBR
	// Read object from bin file
	if err := LeerObjeto(file, &TempMBR, 0); err != nil {
		return "", err
	}

	// Print object
	fmt.Println("\n***********Imprimiendo MBR")
	fmt.Println()
	PrintMBR(TempMBR)
	fmt.Println("\n********Finalizando Impresion de MBR")

	var index int = -1
	// Iterate over the partitions
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_partitions[i].Part_size != 0 {
			if strings.Contains(string(TempMBR.Mbr_partitions[i].Part_id[:]), id) {
				fmt.Println("\n****Particion Encontrada*****")
				if TempMBR.Mbr_partitions[i].Part_status { // si la particion esta montada = true
					fmt.Println("\n*******La particion esta montada*****")
					index = i
				} else {
					fmt.Println("\n*******La particion NO esta montada*****")
					return "", err
				}
				break
			}
		}
	}

	if index != -1 {
		ImprimirParticion(TempMBR.Mbr_partitions[index])
		fmt.Println()
	} else {
		fmt.Println("\n*****Particion NO encontrada******")
		return "", err
	}

	var tempSuperblock Superblock

	if err := LeerObjeto(file, &tempSuperblock, int64(TempMBR.Mbr_partitions[index].Part_start)); err != nil {
		return "", err
	}

	// initSearch /users.txt -> regresa no Inodo
	// initSearch -> 1

	indexInode := InitSearch("/users.txt", file, tempSuperblock) // devuelve el valor = 1 que se usara para encontrar el inodo 1

	fmt.Println("\nindexInode el valor que devuelve InitSearch: ", indexInode)

	var crrInode Inode //Inodo 1

	if err := LeerObjeto(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Inode{})))); err != nil {
		return "", err
	}

	// getInodeFileData -> Iterate the I_Block n concat the data

	var fileblock Fileblock

	fileblock_start := tempSuperblock.S_block_start + crrInode.I_block[0]*int32(binary.Size(Fileblock{}))

	if err := LeerObjeto(file, &fileblock, int64(fileblock_start)); err != nil {
		return "", err
	}

	fmt.Println("Fileblock------------")
	//data := "1,G,root\n1,U,root,root,123\n"
	data := string(fileblock.B_content[:])

	lines := strings.Split(data, "\n")

	datos := strings.Split(lines[1], ",")

	usuario := datos[3]
	contrasena := datos[4]

	fmt.Println("Inode", crrInode.I_block)

	// Close bin file
	defer file.Close()

	if usuario == user && pass == contrasena {
		fmt.Println("\n **********Usuario encontrado***********")

		fmt.Println("\n\n========================= FIN LOGIN ===========================")

		return usuario, err
	} else {
		fmt.Println("\n*********Usuario NO encontrado**********")

		fmt.Println("\n\n========================= FIN LOGIN ===========================")

		return "failed", err
	}

}

func Mkgrp(name string, id string) error {
	fmt.Println("\n\n========================= Inicio MKGRP ===========================")

	fmt.Printf("El grupo a crear sera: %s, El id es: %s", name, id)
	fmt.Println()

	//return file, fileblock, fileblock_start, nil
	file, tempSuperblock, err := getUsersTXT(id)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	indexInode := int32(1) // para poder buscar el inodo1

	inode_start := tempSuperblock.S_inode_start + indexInode*int32(binary.Size(Inode{})) //Inode 1

	var crrInode Inode //inodo que contiene los bloques de archivos de users.txt (Inodo 1)

	if err := LeerObjeto(file, &crrInode, int64(inode_start)); err != nil { //Inode1
		return err
	}

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

	/*
		var fileblock Fileblock

		fileblock_start := tempSuperblock.S_block_start + crrInode.I_block[0]*int32(binary.Size(Fileblock{})) // bloque1

		if err := LeerObjeto(file, &fileblock, int64(fileblock_start)); err != nil { //bloque1
			return err
		}

		if err != nil {
			fmt.Println("Error: ", err)
		}

		fmt.Println("Fileblock------------")
		//data := "1,G,root\n1,U,root,root,123\n"

		var cadena string = " "
		cadena = string(fileblock.B_content[:])

		fmt.Println("\n Imprimiendo cadena: ", string(fileblock.B_content[:]))
	*/

	lines := strings.Split(cadena, "\n")

	if len(lines) > 0 {
		lines = lines[:len(lines)-1]
	}

	fmt.Println("\n\nContenido del arreglo lines: ", lines)
	fmt.Println("\nEl tamano del arreglo lines es: ", len(lines))

	fmt.Println("\nImprimiendo ultimo elemento de arreglo lines: ", lines[len(lines)-1])
	//2, G, usuarios, \n
	var contador int = 0
	var exist int = 0
	var datos []string

	for i := 0; i < len(lines); i++ {

		datos = strings.Split(lines[i], ",")

		contador_, _ := strconv.Atoi(datos[0])

		contador = contador_
		contador++

		if len(datos) != 0 {

			if string(datos[2]) == name {

				fmt.Println("\n\n      ********** El Grupo ya existe ************")

				fmt.Println("\n\n========================= Fin MKGRP ===========================")
				exist++
				return nil
			}
		}

	}

	if exist == 0 { // si el grupo a crear no existe

		newCadena := strconv.Itoa(contador) + ",G," + name + "\n"

		fmt.Println("\n ********datos de la variable newCadena: ", newCadena)

		var contador int

		for i := 0; i < len(fileblock.B_content); i++ {
			if fileblock.B_content[i] == 0 { //verifica si todavia hay espacio
				contador++
			}

		}

		if contador < len(newCadena) {
			//fmt.Println("\nEl contador es: ", contador)
			fmt.Println("\nYa no hay suficiente espacio en users.txt que esta en fileblock.B_content")
			return nil
		}
		//Agregando nuevo grupo a users.txt en fileblock.B_content
		var c int

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

		fmt.Println("\n El contenido nuevo de B_content es: ", string(fileblock.B_content[:]))

		fmt.Println("\n\n ********** Escribiendo objeto FILEBLOCK en el archivo ******************")
		if err := EscribirObjeto(file, fileblock, int64(fileblock_start)); err != nil { //aqui solo escribi el primer EBR
			return err
		}

	}

	var tempfileblock Fileblock

	fmt.Println("\n\n ********** Recuperando y Leyendo objeto FILEBLOCK del archivo binario ******************")
	if err := LeerObjeto(file, &tempfileblock, int64(fileblock_start)); err != nil {
		return err
	}

	printFileblock(tempfileblock)

	//fmt.Println("\n\nLo que se guardo en fileblock.B_content es: ", string(fileblock.B_content[:]))

	fmt.Println("\n\n========================= Fin MKGRP ===========================")

	return nil
}

//
//
//
//
//
//
//

func Mkusr(user string, pass string, group string, id string) error {

	fmt.Println("\n\n========================= Inicio MKUSR ===========================")

	fmt.Printf("El usuario a crear sera: %s, El password es: %s, el grupo al que pertenecera es: %s, El id es: %s", user, pass, group, id)
	fmt.Println()

	//return file, fileblock, fileblock_start, nil
	file, tempSuperblock, err := getUsersTXT(id)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	indexInode := int32(1)

	inode_start := tempSuperblock.S_inode_start + indexInode*int32(binary.Size(Inode{})) //Inode 1

	var crrInode Inode //inodo que contiene los bloques de archivos de users.txt

	if err := LeerObjeto(file, &crrInode, int64(inode_start)); err != nil { //Inode1
		return err
	}

	//data := "1,G,root\n1,U,root,root,123\n"

	CreateNewBlock(file, tempSuperblock, crrInode, user, group, pass)

	/*
		var tempfileblock Fileblock

		fmt.Println("\n\n ********** Recuperando y Leyendo objeto FILEBLOCK del archivo binario ******************")
		if err := LeerObjeto(file, &tempfileblock, int64(fileblock_start)); err != nil {
			return err
		}

		printFileblock(tempfileblock)*/

	fmt.Println("\n\n========================= Fin MKUSR ===========================")

	return nil
}

//
//
//
//
//
//
//
//

func Rmgrp(name string, id string) error {
	fmt.Println("\n\n========================= Inicio MKGRP ===========================")

	fmt.Printf("El grupo a crear sera: %s, El id es: %s", name, id)
	fmt.Println()

	//return file, fileblock, fileblock_start, nil
	file, tempSuperblock, err := getUsersTXT(id)

	indexInode := int32(1)

	inode_start := tempSuperblock.S_inode_start + indexInode*int32(binary.Size(Inode{})) //Inode 1

	var crrInode Inode //inodo que contiene los bloques de archivos de users.txt

	if err := LeerObjeto(file, &crrInode, int64(inode_start)); err != nil { //Inode1
		return err
	}

	var fileblock Fileblock

	fileblock_start := tempSuperblock.S_block_start + crrInode.I_block[0]*int32(binary.Size(Fileblock{})) // bloque1

	if err := LeerObjeto(file, &fileblock, int64(fileblock_start)); err != nil { //bloque1
		return err
	}

	if err != nil {
		fmt.Println("Error: ", err)
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

		contador_, _ := strconv.Atoi(datos[0]) // contiene el numero de grupo

		num_group = contador_

		if len(datos) != 0 {

			if string(datos[2]) == name {

				if num_group == 0 {
					fmt.Println("\nEl grupo no existe porque ya fue eliminado anteriormente")
					return nil
				} else {
					fmt.Println("\n\n      ********** Eliminando grupo " + name + " ************")

					datos[0] = "0"
					lines[i] = strings.Join(datos, ",")

					fmt.Println("\nImprimiendo la linea: ", lines)

					exist++
				}
			}
		}

	}

	newCadena := strings.Join(lines, "\n") // convirtiendo slice lines a cadena de texto
	newCadena += "\n"

	var cadena_bytes [64]byte
	copy(cadena_bytes[:], []byte(newCadena))

	fileblock.B_content = cadena_bytes

	fmt.Println("\nImprimiendo fileblock.B_content con nuevo users.txt: ", string(fileblock.B_content[:]))

	//ESCRIBIENDO FILEBLOCK
	fmt.Println("\n\n ********** Escribiendo objeto FILEBLOCK en el archivo ******************")

	if err := EscribirObjeto(file, fileblock, int64(fileblock_start)); err != nil { //aqui solo escribi el primer EBR
		return err

	}

	var tempfileblock Fileblock

	fmt.Println("\n\n ********** Recuperando y Leyendo objeto FILEBLOCK del archivo binario ******************")
	if err := LeerObjeto(file, &tempfileblock, int64(fileblock_start)); err != nil {
		return err
	}

	printFileblock(tempfileblock)

	fmt.Println("\n\n========================= Fin RMGRP ===========================")

	return nil
}

func Rmusr(user string, id string) error {
	fmt.Println("\n\n========================= Inicio RMUSR ===========================")

	fmt.Printf("El usuario a eliminar sera: %s, El id es: %s", user, id)
	fmt.Println()

	//return file, fileblock, fileblock_start, nil
	file, tempSuperblock, err := getUsersTXT(id)

	indexInode := int32(1)

	inode_start := tempSuperblock.S_inode_start + indexInode*int32(binary.Size(Inode{})) //Inode 1

	var crrInode Inode //inodo que contiene los bloques de archivos de users.txt

	if err := LeerObjeto(file, &crrInode, int64(inode_start)); err != nil { //Inode1
		return err
	}

	var fileblock Fileblock
	fmt.Println("\n\n Buscando el bloque: ", crrInode.I_block[0])
	fileblock_start := tempSuperblock.S_block_start + crrInode.I_block[0]*int32(binary.Size(Fileblock{})) // bloque1

	if err := LeerObjeto(file, &fileblock, int64(fileblock_start)); err != nil { //bloque1
		return err
	}

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("Fileblock------------")
	//data := "1,G,root\n1,U,root,root,123\n"

	var cadena string = " "

	cadena = string(fileblock.B_content[:])

	fmt.Println("\n Imprimiendo cadena: ", string(fileblock.B_content[:]))

	lines := strings.Split(cadena, "\n")

	if len(lines) > 0 { // NO lee el ultimo salto de linea \n
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

		contador_, _ := strconv.Atoi(datos[0]) // contiene el numero de grupo

		num_group = contador_

		if len(datos) != 0 {

			if len(datos) > 3 {

				if string(datos[3]) == user {
					fmt.Println("\nEL usuario a eliminar si existe")

					if num_group == 0 {
						fmt.Println("\nEl usuario no existe porque ya fue eliminado anteriormente")
						return nil
					} else {
						fmt.Println("\n\n      ********** Eliminando usuario " + user + " ************")

						datos[0] = "0"
						lines[i] = strings.Join(datos, ",")

						fmt.Println("\nImprimiendo la linea: ", lines)

						exist++

						break
					}
				}
			}

		}

	}

	if exist != 0 {

		newCadena := strings.Join(lines, "\n") // convirtiendo slice lines a cadena de texto
		newCadena += "\n"

		var cadena_bytes [64]byte
		copy(cadena_bytes[:], []byte(newCadena))

		fileblock.B_content = cadena_bytes

		fmt.Println("\nImprimiendo fileblock.B_content con nuevo users.txt: ", string(fileblock.B_content[:]))

		//ESCRIBIENDO FILEBLOCK
		fmt.Println("\n\n ********** Escribiendo objeto FILEBLOCK en el archivo ******************")

		if err := EscribirObjeto(file, fileblock, int64(fileblock_start)); err != nil { //aqui solo escribi el primer EBR
			return err

		}
	} else {
		//SI NO EXISTE EL USUARIO HAY QUE SEGUIR BUSCANDO EN EL SIGUIENTE BLOQUE DE ARCHIVOS
		fmt.Println("\n°°°°°°°°°°°°°°°°°El usuario no existe°°°°°°°°°°°°°°")
		return nil
	}

	var tempfileblock Fileblock

	fmt.Println("\n\n ********** Recuperando y Leyendo objeto FILEBLOCK del archivo binario ******************")
	if err := LeerObjeto(file, &tempfileblock, int64(fileblock_start)); err != nil {
		return err
	}

	printFileblock(tempfileblock)

	fmt.Println("\n\n========================= Fin RMUSR ===========================")

	return nil
}

func getUsersTXT(id string) (*os.File, Superblock, error) {

	driveletter := string(id[0])

	// Open bin file
	filepath := "./archivos/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := AbrirArchivo(filepath)
	if err != nil {
		return nil, Superblock{}, err
	}

	var TempMBR MBR
	// Read object from bin file
	if err := LeerObjeto(file, &TempMBR, 0); err != nil {
		return nil, Superblock{}, err
	}

	// Print object
	fmt.Println("\n***********Imprimiendo MBR")
	fmt.Println()
	PrintMBR(TempMBR)
	fmt.Println("\n********Finalizando Impresion de MBR")

	var index int = -1
	// Iterate over the partitions
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_partitions[i].Part_size != 0 {
			if strings.Contains(string(TempMBR.Mbr_partitions[i].Part_id[:]), id) {
				fmt.Println("\n****Particion Encontrada*****")
				if TempMBR.Mbr_partitions[i].Part_status { // si la particion esta montada = true
					fmt.Println("\n*******La particion esta montada*****")
					index = i
				} else {
					fmt.Println("\n*******La particion NO esta montada*****")
					return nil, Superblock{}, err
				}
				break
			}
		}
	}

	if index != -1 {
		ImprimirParticion(TempMBR.Mbr_partitions[index])
		fmt.Println()
	} else {
		fmt.Println("\n*****Particion NO encontrada******")
		return nil, Superblock{}, err
	}

	var tempSuperblock Superblock

	if err := LeerObjeto(file, &tempSuperblock, int64(TempMBR.Mbr_partitions[index].Part_start)); err != nil {
		return nil, Superblock{}, err
	}

	// initSearch /users.txt -> regresa no Inodo
	// initSearch -> 1

	// getInodeFileData -> Iterate the I_Block n concat the data

	return file, tempSuperblock, nil

}

/*	 CODIGO PARA MANEJAR LOS SLICES DE BYTES DE TIPO [SIZE]BYTE

cadena := "1,U,root,123\n"
  //usuario := "2,U,user,dracker"
  var fileblock [32]byte
  copy(fileblock[:], []byte(cadena))

  //data := string(fileblock[:])
  //fmt.Println("\nLa data es: ",data)
  data := "2,U,usuario,562\n"
  //cadena += "3,U,user,002"
  fmt.Println("\nLa NUEVA data es: ",data)
  fmt.Println("la longitud de la cadena es: ", len(data))
  //Data := make([]byte,3)

  //fmt.Println(Data) //output is [0,0,0]

  var c int
  for i := 0; i < len(fileblock); i++ {
      //fmt.Println(fileblock[i])

      if fileblock[i] ==0 {

          if c < len(data){
              fileblock[i] = byte(data[c])
              fmt.Printf("letra:  %s   ", string(data[c]))
              c++

          }else{
              break
          }


      }
  }

  var contador int
  for i := 0; i < len(fileblock); i++ {
      if fileblock[i] ==0 {
        contador++
      }

  }


  fmt.Println("\nfileblock: ", string(fileblock[:]))
  fmt.Println("\nEl nuevo tamano de fileblock es: ", len(fileblock))

  if contador < len(data){
      fmt.Println("\nEl contador es: ", contador)
      fmt.Println("\nYa no hay suficiente espacio")
  }
*/
