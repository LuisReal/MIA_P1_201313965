package Funciones

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var size_ int

var exp = regexp.MustCompile(`(\w+)\.(\w+)`) // para analizar archivos con extension

func Mkfile(path string, r string, size int, cont string) error {

	fmt.Println("\n\n=========================Creando Archivo (MKFILE)===========================")

	fmt.Println("\n**********************El path ingresado es: ", path)

	//fmt.Println("\nUsuario con sesion actual: ", user_.Nombre, " ID: ", user_.Id, " r: ", r)

	if size < 0 {
		fmt.Println("\n\n            ERROR: El tamano size no puede ser negativo")

		fmt.Println("\n\n=======================Finalizando Creacion Archivo (MKFILE)===========================")
		return nil
	}

	size_ = size

	id := user_.Id

	driveletter := string(id[0])

	//fmt.Println("\nEl disco es: ", driveletter)

	// Open bin file
	filepath := "./archivos/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := AbrirArchivo(filepath)
	if err != nil {
		return err
	}

	var TempMBR MBR
	// Read object from bin file
	if err := LeerObjeto(file, &TempMBR, 0); err != nil {
		return err
	}

	// Print object
	/*fmt.Println("\n***********Imprimiendo MBR")
	fmt.Println()
	PrintMBR(TempMBR)
	fmt.Println("\n********Finalizando Impresion de MBR")*/

	var index int = -1
	// Iterate over the partitions
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_partitions[i].Part_size != 0 {
			if strings.Contains(string(TempMBR.Mbr_partitions[i].Part_id[:]), strings.ToUpper(id)) {
				//fmt.Println("\n****Particion Encontrada*****")

				index = i
			}
		}
	}

	var tempSuperblock Superblock

	superblock_start := TempMBR.Mbr_partitions[index].Part_start

	if err := LeerObjeto(file, &tempSuperblock, int64(superblock_start)); err != nil {
		return err
	}

	carpetas := strings.Split(path, "/") // solo hay una carpeta que se creara en la raiz

	arreglo := carpetas

	if tempSuperblock.S_filesystem_type == 3 {

		fmt.Println("Registrando operacion al JOURNALING")

		var journaling Journaling

		if err := LeerObjeto(file, &journaling, int64(superblock_start+int32(binary.Size(Superblock{})))); err != nil {
			return err
		}

		for k := 0; k < len(journaling.Contenido); k++ {

			var operacion [10]byte

			//execute -path=/home/darkun/Escritorio/prueba.mia

			if journaling.Contenido[k].Operation == operacion { // verifica si la variable Operation esta vacia para ingresar nuevo valor

				copy(journaling.Contenido[k].Operation[:], "mkfile")
				copy(journaling.Contenido[k].Path[:], []byte(arreglo[len(arreglo)-1]))
				copy(journaling.Contenido[k].Content[:], "-")

				date := time.Now()
				//fmt.Println("La Fecha y Hora Actual es: ", date.Format("2006-01-02 15:04:05"))

				byteString := make([]byte, 17)
				copy(byteString, date.Format("2006-01-02 15:04:05"))

				copy(journaling.Contenido[k].Date[:], byteString)

				err := EscribirObjeto(file, journaling, int64(superblock_start+int32(binary.Size(Superblock{}))))

				if err != nil {
					fmt.Println("Error: ", err)
				}

				break
			}

		}

	}

	Inodo_start := tempSuperblock.S_inode_start

	var inodo0 Inode

	if err := LeerObjeto(file, &inodo0, int64(Inodo_start)); err != nil {
		return err
	}

	//mkfile -path=/home/archivos/user/docs/Tarea.txt -size=75

	carpetas = carpetas[1:]

	duplicado := DuplicateElement(carpetas)

	if duplicado != "-1" {
		fmt.Println("\n      ********** ERROR: recursivo *********")
		return nil
	}

	//execute -path=/home/darkun/Escritorio/prueba.mia

	fmt.Println("\nEl arreglo carpetas es: \n", carpetas)
	fmt.Println("\nEl tamano del arreglo carpetas es: \n", len(carpetas))

	fmt.Println("\nEl parametro -r es: ", r)

	var cont_folder int

	AddingNewFile(carpetas, inodo0, file, tempSuperblock, -1, superblock_start, cont_folder, cont, r)

	fmt.Println("\n\n=======================Finalizando Creacion Archivo (MKFILE)===========================")

	file.Close()

	return nil
}

func DuplicateElement(arr []string) string {
	visited := make(map[string]bool, 0)
	for i := 0; i < len(arr); i++ {
		if visited[arr[i]] {
			return arr[i]
		} else {
			visited[arr[i]] = true
		}
	}
	return "-1"
}

func AddingNewFile(carpeta []string, Inodo Inode, file *os.File, tempSuperblock Superblock, Inodo_actual int32, superblock_start int32, cont_folder int, cont string, r string) error {

	//mkfile -path=/home/archivos/user/docs/Tarea.txt -size=75

	fmt.Println("\n\n========================= Iniciando AddingNewFile ===========================")

	fmt.Println("\nEl valor de cont_folder es: ", cont_folder)

	if cont_folder == len(carpeta) {
		fmt.Println("\n return : debido a que cont_folder es igual que la longitud de la carpeta arrreglos")
		return nil
	}

	//execute -path=/home/darkun/Escritorio/prueba.mia

	fmt.Println("\nLeyendo Carpeta: ", carpeta[cont_folder])

	//fmt.Println("El tipo de inodo es: ", string(Inodo.I_type[:]))

	var folder_bytes [12]byte
	copy(folder_bytes[:], []byte(carpeta[cont_folder]))

	// Iterate over i_blocks from Inode
	for j, block := range Inodo.I_block {
		if block != -1 {

			if j < 12 {

				//CASO DIRECTO

				//// carpeta -> 0   archivo ->1

				if string(Inodo.I_type[:]) == "0" { // si es carpeta

					var crrFolderBlock Folderblock

					//grafo += `Bloque` + strconv.Itoa(int(block)) + `:0;` + "\n"
					// Read object from bin file                                       // un Inodo y un bloque las estructuras miden lo mismo
					if err := LeerObjeto(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))); err != nil {
						return err
					}

					var exist_folder int

					for k, folder := range crrFolderBlock.B_content {

						//fmt.Println("\nEl valor de k es: ", k)

						if string(folder.B_name[:]) == string(folder_bytes[:]) {

							fmt.Println("\nLa carpeta "+carpeta[cont_folder]+" si existe en el bloque: ", block)

							fmt.Println("\n ======= NextInode ======")
							var NextInode Inode
							// Read object from bin file
							if err := LeerObjeto(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(Inode{})))); err != nil {
								return err
							}

							fmt.Println("Enviando el INODO: " + strconv.Itoa(int(folder.B_inodo)))
							//grafo, _ = EnlazandoNodos(path, NextInode, file, tempSuperblock, disco, grafo, folder.B_inodo)
							exist_folder++

							AddingNewFile(carpeta, NextInode, file, tempSuperblock, folder.B_inodo, superblock_start, cont_folder+1, cont, r)

							return nil

						}

						params := carpeta[cont_folder]

						matches := exp.FindAllStringSubmatch(params, -1)

						//var Nombre string
						var Extension string

						for _, match := range matches {

							//Nombre = match[1]
							Extension = match[2]
						}

						//execute -path=/home/darkun/Escritorio/prueba.mia

						if r == "0" {

							if folder.B_inodo == -1 && Extension != "" {
								fmt.Println("\nEl parametro r NO esta incluido, llamando a funcion VerificaTipoArchivo")
								fmt.Println("\nEl bloque actual es: ", block)

								VerificaTipoArchivo(carpeta, cont_folder, &crrFolderBlock, tempSuperblock, file, block, cont, superblock_start, k)

								err := EscribirObjeto(file, crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))) //Inode 0 carpeta raiz

								if err != nil {
									fmt.Println("Error: ", err)
								}

								return nil

							}

						} else {

							if folder.B_inodo == -1 && Extension != "" {

								fmt.Println("\nEl parametro r SI esta incluido, llamando a funcion VerificaTipoArchivo")

								VerificaTipoArchivo(carpeta, cont_folder, &crrFolderBlock, tempSuperblock, file, block, cont, superblock_start, k)

								err := EscribirObjeto(file, crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))) //Inode 0 carpeta raiz

								if err != nil {
									fmt.Println("Error: ", err)
								}

								return nil

							} else { //si es una carpeta

								fmt.Println("\nIngresando la carpeta " + carpeta[cont_folder] + " en el bloque " + strconv.Itoa(int(block)))

								inodo := tempSuperblock.S_inodes_count - tempSuperblock.S_free_inodes_count

								fmt.Println("\nEl valor de inodo es: ", inodo)

								crrFolderBlock.B_content[k].B_inodo = inodo

								copy(crrFolderBlock.B_content[k].B_name[:], []byte(carpeta[cont_folder]))

								err := EscribirObjeto(file, crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))) //Inode 0 carpeta raiz

								if err != nil {
									fmt.Println("Error: ", err)
								}

								numBlocks := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count

								//Creando nuevo Inodo
								var newInode Inode

								user_id, _ := strconv.Atoi(user_.Uid)
								group_id, _ := strconv.Atoi(user_.Gid)

								newInode.I_uid = int32(user_id)
								newInode.I_gid = int32(group_id)
								copy(newInode.I_type[:], "0")

								copy(newInode.I_perm[:], "664")

								tempSuperblock.S_free_inodes_count -= 1

								for i := int32(0); i < 15; i++ {
									newInode.I_block[i] = -1 //-1 no han sido utilizados
								}

								fmt.Println("\nEl valor de numBlocks es: ", numBlocks)

								newInode.I_block[0] = numBlocks

								//fmt.Println("\nnewInode.Iblock[0]: ", newInode.I_block[0])

								err = EscribirObjeto(file, newInode, int64(tempSuperblock.S_inode_start+inodo*int32(binary.Size(Inode{})))) //Inode 0 carpeta raiz

								if err != nil {
									fmt.Println("Error: ", err)
								}

								//Escribiendo en bitmap inodos

								err = EscribirObjeto(file, byte(1), int64(tempSuperblock.S_bm_inode_start+inodo))
								if err != nil {
									fmt.Println("Error: ", err)
								}

								//creando NuevoBLoque
								var newfolder Folderblock

								for i := int32(0); i < 4; i++ {
									newfolder.B_content[i].B_inodo = -1
								}

								newfolder.B_content[0].B_inodo = inodo
								copy(newfolder.B_content[0].B_name[:], ".")

								newfolder.B_content[1].B_inodo = 0
								copy(newfolder.B_content[1].B_name[:], "..")

								tempSuperblock.S_free_blocks_count -= 1

								err = EscribirObjeto(file, newfolder, int64(tempSuperblock.S_block_start+numBlocks*int32(binary.Size(Folderblock{})))) //Inode 0 carpeta raiz

								if err != nil {
									fmt.Println("Error: ", err)
								}

								// Escribiendo en bitmap de bloques
								err = EscribirObjeto(file, byte(1), int64(tempSuperblock.S_bm_block_start+numBlocks))
								if err != nil {
									fmt.Println("Error: ", err)
								}

								// escribiendo y actualizando el superblock
								err = EscribirObjeto(file, tempSuperblock, int64(superblock_start))
								if err != nil {
									fmt.Println("Error: ", err)
								}

								fmt.Println("\n ======= NextInode ======")
								var NextInode Inode
								// Read object from bin file
								if err := LeerObjeto(file, &NextInode, int64(tempSuperblock.S_inode_start+inodo*int32(binary.Size(Inode{})))); err != nil {
									return err
								}

								AddingNewFile(carpeta, NextInode, file, tempSuperblock, inodo, superblock_start, cont_folder+1, cont, r)

								return nil
							}
						}

						/*
							err := EscribirObjeto(file, crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))) //Inode 0 carpeta raiz

							if err != nil {
								fmt.Println("Error: ", err)
							}*/

					}

				}

			} else { //aqui van los bloques de apuntadores, simples, dobles y triples

				//execute -path=/home/darkun/Escritorio/prueba.mia

				fmt.Println("\n ********** CASO INDIRECTO ************ ")

			}

		} else { //execute -path=/home/darkun/Escritorio/prueba.mia

			///home/archivos/user/docs/Tarea.txt

			fmt.Println("\nLeyendo el siguiente bloque " + strconv.Itoa(int(block)) + " vacio del inodo")
			fmt.Println("El indice J del bloque en el inodo es: ", j)
			fmt.Println("El cont_folder es: ", cont_folder)

			fmt.Println("*******La carpeta o ARCHIVO a ingresar es: ", carpeta[cont_folder])

			params := carpeta[cont_folder]

			matches := exp.FindAllStringSubmatch(params, -1)

			//var Nombre string
			var Extension string

			for _, match := range matches {

				//Nombre = match[1]
				Extension = match[2]
			}

			if Extension != "" { // solamente si es un archivo lo crea

				newBlock := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count
				Inodo.I_block[j] = int32(newBlock)
				//fmt.Println("\nEl nuevo bloque es: ", Inodo.I_block[j])

				fmt.Println("\nEl valor de Inodo actual es: ", Inodo_actual)
				fmt.Println("\nEl Inodo "+strconv.Itoa(int(Inodo_actual))+" en su posicion "+strconv.Itoa(j)+" apunta al bloque : ", Inodo.I_block[j])

				// escribiendo blocks
				err := EscribirObjeto(file, &Inodo, int64(tempSuperblock.S_inode_start+Inodo_actual*int32(binary.Size(Inode{})))) //Bloque 0
				if err != nil {
					fmt.Println("Error: ", err)
				}

				inodos_ocupados := tempSuperblock.S_inodes_count - tempSuperblock.S_free_inodes_count

				var newFolder1 Folderblock //Bloque 0 -> carpetas

				for i := int32(0); i < 4; i++ {
					newFolder1.B_content[i].B_inodo = -1
				}

				newFolder1.B_content[0].B_inodo = Inodo_actual
				copy(newFolder1.B_content[0].B_name[:], ".")

				newFolder1.B_content[1].B_inodo = 0
				copy(newFolder1.B_content[1].B_name[:], "..")

				newFolder1.B_content[2].B_inodo = inodos_ocupados
				copy(newFolder1.B_content[2].B_name[:], []byte(carpeta[cont_folder]))

				tempSuperblock.S_free_blocks_count -= 1

				// escribiendo blocks
				err = EscribirObjeto(file, newFolder1, int64(tempSuperblock.S_block_start+newBlock*int32(binary.Size(Folderblock{})))) //Bloque 0
				if err != nil {
					fmt.Println("Error: ", err)
				}

				// Escribiendo en bitmap de bloques
				err = EscribirObjeto(file, byte(1), int64(tempSuperblock.S_bm_block_start+newBlock))
				if err != nil {
					fmt.Println("Error: ", err)
				}

				// llenando con -1 los primeros 15 bloques del inodo

				var newInode Inode

				for i := int32(0); i < 15; i++ {
					newInode.I_block[i] = -1 //-1 no han sido utilizados
				}

				tempSuperblock.S_free_inodes_count -= 1

				userid, _ := strconv.Atoi(user_.Uid)
				newInode.I_uid = int32(userid)

				groupid, _ := strconv.Atoi(user_.Gid)
				newInode.I_gid = int32(groupid)

				copy(newInode.I_type[:], "1")
				copy(newInode.I_perm[:], "664")

				numBlocks := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count
				//fmt.Println("\nEl numero de bloques ocupados es: ", numBlocks)

				newInode.I_block[0] = numBlocks
				fmt.Println("\nEl nuevo numero de bloque es: ", newInode.I_block[0])

				err = EscribirObjeto(file, newInode, int64(tempSuperblock.S_inode_start+inodos_ocupados*int32(binary.Size(Inode{}))))
				if err != nil {
					fmt.Println("Error: ", err)
				}

				//Escribiendo en bitmap inodos

				err = EscribirObjeto(file, byte(1), int64(tempSuperblock.S_bm_inode_start+inodos_ocupados))
				if err != nil {
					fmt.Println("Error: ", err)
				}

				// escribiendo y actualizando el superblock
				err = EscribirObjeto(file, tempSuperblock, int64(superblock_start))
				if err != nil {
					fmt.Println("Error: ", err)
				}

				var newCadena string

				if cont != "" { // si hay parametro cont

					archivo_, err := os.Open(cont)

					if err != nil {
						fmt.Println("Error al abrir el archivo:", err)
						return err
					}

					scanner := bufio.NewScanner(archivo_)

					for scanner.Scan() {

						linea := scanner.Text()

						if linea != "" {

							for i := 0; i < len(linea); i++ {

								if string(linea[i]) == "#" {

									break
								}

								newCadena += string(linea[i])
							}

						}

					}
				} else { // si no hay parametro cont
					var cont int

					var cont_num int

					for cont < size_ {

						newCadena += strconv.Itoa(cont_num)
						cont_num++

						if cont_num == 10 {
							cont_num = 0
						}

						cont++
					}

				}

				var newFile Fileblock

				tempSuperblock.S_free_blocks_count -= 1

				var c int

				for j := 0; j < len(newFile.B_content[:]); j++ {

					if newFile.B_content[j] == 0 { // si hay todavia espacio

						if c < len(newCadena) {

							newFile.B_content[j] = byte(newCadena[c])
							//fmt.Printf("agregando letra:  %s   ", string(newCadena[c]))
							c++

						} else {
							break
						}

					}
				}

				var espacios int

				for i := 0; i < len(newFile.B_content); i++ {

					if newFile.B_content[i] == 0 {
						espacios++
					}
				}

				fmt.Println("\n\n ****Escribiendo objeto FILEBLOCK en el archivo *****")
				if err := EscribirObjeto(file, newFile, int64(tempSuperblock.S_block_start+numBlocks*int32(binary.Size(Fileblock{})))); err != nil { //aqui solo escribi el primer EBR
					return err
				}

				// Escribiendo en bitmap de bloques
				err = EscribirObjeto(file, byte(1), int64(tempSuperblock.S_bm_block_start+numBlocks))
				if err != nil {
					fmt.Println("Error: ", err)
				}

				// escribiendo y actualizando el superblock
				err = EscribirObjeto(file, tempSuperblock, int64(superblock_start))
				if err != nil {
					fmt.Println("Error: ", err)
				}

				//data := "1,G,root\n1,U,root,root,123\n"
				if espacios > 0 {

					fmt.Println("\n Todavia sobra espacio despues de escribir la cadena en el slice")

					return nil

				} else { // si ya no hay espacios en el slice para ingresar la cadena

					//fmt.Println("\n La longitud de la cadena newCadena[c] es: ", len(newCadena[c:]))

					if len(newCadena[c:]) != 0 { //si todavia hay caracteres en newCadena para seguir ingresando en slice de fileblock.Bcontent

						fmt.Println("\n Enviando cadena:\n", newCadena)
						fmt.Println("\nEnviando c: ", c, " Enviando inodos: ", inodos_ocupados)
						CreateNewBlock(newCadena, c, newInode, inodos_ocupados, &tempSuperblock, file, superblock_start)

						// escribiendo y actualizando el superblock
						err = EscribirObjeto(file, tempSuperblock, int64(superblock_start))
						if err != nil {
							fmt.Println("Error: ", err)
						}

						numBlocks := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count

						fmt.Println("\nEl nuevo numero de bloques ocupados es(despues de createnewblock): ", numBlocks)

						return nil
					}

				}
			} else {

				if r != "0" { // el parametro r esta incluido

					fmt.Println("Creando la carpeta")
					fmt.Println("*******La carpeta a ingresar es: ", carpeta[cont_folder])

					newBlock := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count
					Inodo.I_block[j] = int32(newBlock)
					//fmt.Println("\nEl nuevo bloque es: ", Inodo.I_block[j])

					fmt.Println("\nEl valor de Inodo actual es: ", Inodo_actual)
					fmt.Println("\nEl Inodo "+strconv.Itoa(int(Inodo_actual))+" en su posicion "+strconv.Itoa(j)+" apunta al bloque : ", Inodo.I_block[j])

					// escribiendo blocks
					err := EscribirObjeto(file, &Inodo, int64(tempSuperblock.S_inode_start+Inodo_actual*int32(binary.Size(Inode{})))) //Bloque 0
					if err != nil {
						fmt.Println("Error: ", err)
					}

					inodos_ocupados := tempSuperblock.S_inodes_count - tempSuperblock.S_free_inodes_count

					var newFolder Folderblock //Bloque 0 -> carpetas

					for i := int32(0); i < 4; i++ {
						newFolder.B_content[i].B_inodo = -1
					}

					newFolder.B_content[0].B_inodo = inodos_ocupados
					copy(newFolder.B_content[0].B_name[:], []byte(carpeta[cont_folder]))

					tempSuperblock.S_free_blocks_count -= 1

					// escribiendo blocks
					err = EscribirObjeto(file, newFolder, int64(tempSuperblock.S_block_start+newBlock*int32(binary.Size(Folderblock{})))) //Bloque 0
					if err != nil {
						fmt.Println("Error: ", err)
					}

					// Escribiendo en bitmap de bloques
					err = EscribirObjeto(file, byte(1), int64(tempSuperblock.S_bm_block_start+newBlock))
					if err != nil {
						fmt.Println("Error: ", err)
					}

					// llenando con -1 los primeros 15 bloques del inodo

					var newInode Inode

					for i := int32(0); i < 15; i++ {
						newInode.I_block[i] = -1 //-1 no han sido utilizados
					}

					tempSuperblock.S_free_inodes_count -= 1

					userid, _ := strconv.Atoi(user_.Uid)
					newInode.I_uid = int32(userid)

					groupid, _ := strconv.Atoi(user_.Gid)
					newInode.I_gid = int32(groupid)

					numBlocks := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count
					//fmt.Println("\nEl numero de bloques ocupados es: ", numBlocks)

					newInode.I_block[0] = numBlocks
					fmt.Println("\nEl nuevo numero de bloque es: ", newInode.I_block[0])

					copy(newInode.I_type[:], "0")
					copy(newInode.I_perm[:], "664")

					numInodo := newFolder.B_content[0].B_inodo

					err = EscribirObjeto(file, newInode, int64(tempSuperblock.S_inode_start+numInodo*int32(binary.Size(Inode{}))))
					if err != nil {
						fmt.Println("Error: ", err)
					}

					//Escribiendo en bitmap inodos

					err = EscribirObjeto(file, byte(1), int64(tempSuperblock.S_bm_inode_start+numInodo))
					if err != nil {
						fmt.Println("Error: ", err)
					}

					//creando NuevoBLoque
					var newfolder Folderblock

					//fmt.Println("\nEl numero de bloques ocupados es: ", numBlocks)

					for i := int32(0); i < 4; i++ {
						newfolder.B_content[i].B_inodo = -1
					}

					newfolder.B_content[0].B_inodo = inodos_ocupados
					copy(newfolder.B_content[0].B_name[:], ".")

					newfolder.B_content[1].B_inodo = 0
					copy(newfolder.B_content[1].B_name[:], "..")

					tempSuperblock.S_free_blocks_count -= 1

					err = EscribirObjeto(file, newfolder, int64(tempSuperblock.S_block_start+numBlocks*int32(binary.Size(Folderblock{})))) //Inode 0 carpeta raiz

					if err != nil {
						fmt.Println("Error: ", err)
					}

					// Escribiendo en bitmap de bloques
					err = EscribirObjeto(file, byte(1), int64(tempSuperblock.S_bm_block_start+numBlocks))
					if err != nil {
						fmt.Println("Error: ", err)
					}

					// escribiendo y actualizando el superblock
					err = EscribirObjeto(file, tempSuperblock, int64(superblock_start))
					if err != nil {
						fmt.Println("Error: ", err)
					}

					fmt.Println("\n ======= NextInode ======")
					var NextInode Inode
					// Read object from bin file
					if err := LeerObjeto(file, &NextInode, int64(tempSuperblock.S_inode_start+inodos_ocupados*int32(binary.Size(Inode{})))); err != nil {
						return err
					}

					AddingNewFile(carpeta, NextInode, file, tempSuperblock, inodos_ocupados, superblock_start, cont_folder+1, cont, r)

					return nil

				} else {
					fmt.Println("\n        ERROR: la carpeta \"" + carpeta[cont_folder] + "\" no se ha encontrado")

					break
				}

			}

			//AddingNewFile(carpeta, Inodo, file, tempSuperblock, inodos_ocupados, superblock_start, cont_folder+1)

			//execute -path=/home/darkun/Escritorio/prueba.mia

			///home/archivos/carpeta1/carpeta2/carpeta3/carpeta4/carpeta5

		}

	}

	fmt.Println("\n\n========================= Finalizando AddingNewFile ===========================")

	return nil
}

func CreateNewBlock(newCadena string, contador int, crrInode Inode, num_inodo int32, tempSuperblock *Superblock, file *os.File, superblock_start int32) {

	//fmt.Println("\n\n========================= Inicio CreateNewBlock ===========================")

	//execute -path=/home/darkun/Escritorio/prueba.mia

	resto_cadena := newCadena[contador:]

	//fmt.Println("\nEl resto_cadena es:\n", resto_cadena)
	//fmt.Println("\nLongitud de resto_cadena es: ", len(resto_cadena))

	var bloque int
	var index int

	for i := 0; i < 12; i++ {

		if crrInode.I_block[i] != -1 {
			bloque = int(crrInode.I_block[i]) //obtiene el numero del ultimo bloque creado
			index = i
		}
	}

	if index == 11 { //significa que los primeros 12 estan llenos (0-11)

		//fmt.Println("\nLlamando a la funcion CreateNewIndirectBlock")

		superbloque_ := CreateNewIndirectBlock(resto_cadena, contador, crrInode, num_inodo, *tempSuperblock, file, superblock_start, 12, 1)

		*tempSuperblock = superbloque_

		return
	}

	newBlock := bloque + 1
	crrInode.I_block[index+1] = int32(newBlock)

	tempSuperblock.S_free_blocks_count -= 1

	err := EscribirObjeto(file, crrInode, int64(tempSuperblock.S_inode_start+num_inodo*int32((binary.Size(Inode{}))))) //Inode 1

	if err != nil {
		fmt.Println("Error: ", err)
	}

	var newFileblock Fileblock

	//newFileblock.B_content

	var c int

	//data := "1,G,root\n1,U,root,root,123\n"
	for i := 0; i < len(newFileblock.B_content); i++ {

		if newFileblock.B_content[i] == 0 { // si hay espacio

			if c < len(resto_cadena) {
				//fileblock.B_content [2,U,usuarios,user2,    contra2sena]
				newFileblock.B_content[i] = byte(resto_cadena[c])

				c++

			}

		}
	}

	fileblock_start := tempSuperblock.S_block_start + int32(newBlock)*int32(binary.Size(Fileblock{})) // bloque1

	err = EscribirObjeto(file, newFileblock, int64(fileblock_start)) //Bloque 1

	if err != nil {
		fmt.Println("Error: ", err)
	}

	// Escribiendo en bitmap de bloques
	err = EscribirObjeto(file, byte(1), int64(tempSuperblock.S_bm_block_start+int32(newBlock)))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	//obtiene cantidad de espacios restantes en el slice
	var espacios int

	for i := 0; i < len(newFileblock.B_content); i++ {

		if newFileblock.B_content[i] == 0 {
			espacios++
		}
	}

	if espacios > 0 {

		fmt.Println("\n Todavia sobra espacio despues de escribir la cadena en el slice")

		//fmt.Println("\n\n========================= Fin CreateNewBlock ===========================")

		return

	} else { // si ya no hay espacios en el slice para ingresar la cadena

		if len(resto_cadena[c:]) != 0 { //si todavia hay caracteres en newCadena para seguir ingresando en slice de fileblock.Bcontent

			//fmt.Println("\n      LLamando funcion CrearBloque .......")

			CreateNewBlock(resto_cadena, c, crrInode, num_inodo, tempSuperblock, file, superblock_start)
		}

		return
	}

	//execute -path=/home/darkun/Escritorio/prueba.mia

}

func CreateNewIndirectBlock(resto_cadena string, contador int, crrInode Inode, num_inodo int32, tempSuperblock Superblock, file *os.File, superblock_start int32, block int, apuntador int) Superblock {

	fmt.Println("\n\n========================= Inicio CreateNewIndirectBlock ===========================")

	var indice int
	var x int

	if crrInode.I_block[block] == -1 {

		//fmt.Println("\n............Creando nuevo bloque apuntadores................")

		newBlock := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count

		//fmt.Println("\nEl numero de bloques ocupados es: ", newBlock)

		crrInode.I_block[block] = int32(newBlock)
		//fmt.Println("\nEl nuevo bloque es: ", Inodo.I_block[j])

		// escribiendo blocks
		err := EscribirObjeto(file, &crrInode, int64(tempSuperblock.S_inode_start+num_inodo*int32(binary.Size(Inode{})))) //Bloque 0
		if err != nil {
			fmt.Println("Error: ", err)
		}

		//execute -path=/home/darkun/Escritorio/prueba.mia

		if apuntador == 1 { // apuntador simple

			pointerblock := newBlock

			var newPointer Pointerblock

			for m := 0; m < len(newPointer.B_pointers); m++ {
				newPointer.B_pointers[m] = -1
			}

			tempSuperblock.S_free_blocks_count -= 1

			for j, bloque := range newPointer.B_pointers {

				indice = j

				if bloque == -1 {

					//Creando nuevo bloque de archivos

					newBlock = tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count

					//fmt.Println("\nEl nuevo bloque(newBlock) a crear es: ", newBlock)

					newPointer.B_pointers[j] = newBlock

					var newFileblock Fileblock

					tempSuperblock.S_free_blocks_count -= 1

					var c int

					//data := "1,G,root\n1,U,root,root,123\n"
					for i := 0; i < len(newFileblock.B_content); i++ {

						if newFileblock.B_content[i] == 0 { // si hay espacio

							if c < len(resto_cadena) {
								//fileblock.B_content [2,U,usuarios,user2,    contra2sena]
								newFileblock.B_content[i] = byte(resto_cadena[c])

								c++

							}

						}
					}

					x = c

					fileblock_start := tempSuperblock.S_block_start + int32(newBlock)*int32(binary.Size(Fileblock{})) // bloque1

					//fmt.Println("\n Escribiendo newFileblock en el archivo..........")

					err = EscribirObjeto(file, newFileblock, int64(fileblock_start)) //Bloque 1

					if err != nil {
						fmt.Println("Error: ", err)
					}

					// Escribiendo en bitmap de bloques
					err = EscribirObjeto(file, byte(1), int64(tempSuperblock.S_bm_block_start+int32(newBlock)))
					if err != nil {
						fmt.Println("Error: ", err)
					}

					var espacios int

					for i := 0; i < len(newFileblock.B_content); i++ {

						if newFileblock.B_content[i] == 0 {
							espacios++
						}
					}

					if espacios > 0 {

						//fmt.Println("\n Todavia sobra espacio despues de escribir la cadena en el slice")

						break

					} else { // si ya no hay espacios en el slice para ingresar la cadena

						//fmt.Println("\nLa longitud de resto_cadena es: ", len(resto_cadena[c:]))

						if len(resto_cadena[c:]) != 0 { // si todavia hay caracteres en la cadena para seguir ingresando en el slice

							resto_cadena = resto_cadena[c:] //si todavia hay caracteres en newCadena para seguir ingresando en slice de fileblock.Bcontent
						}
					}
				}
			}

			// escribiendo blocks
			err = EscribirObjeto(file, &newPointer, int64(tempSuperblock.S_block_start+pointerblock*int32(binary.Size(Pointerblock{})))) //Bloque 0
			if err != nil {
				fmt.Println("Error: ", err)
			}

			// escribiendo y actualizando el superblock
			err = EscribirObjeto(file, tempSuperblock, int64(superblock_start))
			if err != nil {
				fmt.Println("Error: ", err)
			}

		} else if apuntador == 2 { // apuntador doble

			fmt.Println("\n          ********** Iniciando bloque de apuntadores doble")
			//pointerblock := newBlock

			numBlock1 := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count

			var newPointer1 Pointerblock

			for m := 0; m < len(newPointer1.B_pointers); m++ {
				newPointer1.B_pointers[m] = -1
			}

			tempSuperblock.S_free_blocks_count -= 1

			for j, bloque1 := range newPointer1.B_pointers { // primer bloque de apuntadores

				indice = j

				//execute -path=/home/darkun/Escritorio/prueba.mia

				if bloque1 == -1 {

					newBlock = tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count

					newPointer1.B_pointers[j] = newBlock

					//pointerblock := newBlock

					numBlock2 := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count

					var newPointer2 Pointerblock

					for m := 0; m < len(newPointer2.B_pointers); m++ {
						newPointer2.B_pointers[m] = -1
					}

					tempSuperblock.S_free_blocks_count -= 1

					for k, bloque2 := range newPointer2.B_pointers { // segundo bloque de apuntadores
						//indice = k

						//execute -path=/home/darkun/Escritorio/prueba.mia

						if bloque2 == -1 {

							newBlock = tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count

							newPointer2.B_pointers[k] = newBlock

							var newFileblock Fileblock

							tempSuperblock.S_free_blocks_count -= 1

							var c int

							//data := "1,G,root\n1,U,root,root,123\n"
							for i := 0; i < len(newFileblock.B_content[:]); i++ {

								if newFileblock.B_content[i] == 0 { // si hay espacio

									if c < len(resto_cadena) {
										//fileblock.B_content [2,U,usuarios,user2,    contra2sena]
										newFileblock.B_content[i] = byte(resto_cadena[c])

										c++

									}

								}
							}

							x = c

							fileblock_start := tempSuperblock.S_block_start + int32(newBlock)*int32(binary.Size(Fileblock{})) // bloque1

							//fmt.Println("\n Escribiendo newFileblock en el archivo..........")

							err = EscribirObjeto(file, newFileblock, int64(fileblock_start)) //Bloque 1

							if err != nil {
								fmt.Println("Error: ", err)
							}

							// Escribiendo en bitmap de bloques
							err = EscribirObjeto(file, byte(1), int64(tempSuperblock.S_bm_block_start+int32(newBlock)))
							if err != nil {
								fmt.Println("Error: ", err)
							}

							var espacios int

							for i := 0; i < len(newFileblock.B_content); i++ {

								if newFileblock.B_content[i] == 0 {
									espacios++
								}
							}

							if espacios > 0 {

								//fmt.Println("\n Todavia sobra espacio despues de escribir la cadena en el slice")

								break

							} else { // si ya no hay espacios en el slice para ingresar la cadena
								//fmt.Println("\nEl valor de c es: ", c)

								if len(resto_cadena[c:]) != 0 { // si todavia hay caracteres en la cadena para seguir ingresando en el slice

									resto_cadena = resto_cadena[c:] //si todavia hay caracteres en newCadena para seguir ingresando en slice de fileblock.Bcontent

								} else { // si ya no hay caracteres en la cadena
									break
								}
							}

						}

					}

					// escribiendo blocks
					err = EscribirObjeto(file, &newPointer2, int64(tempSuperblock.S_block_start+numBlock2*int32(binary.Size(Pointerblock{})))) //Bloque 0
					if err != nil {
						fmt.Println("Error: ", err)
					}

					/*

						// escribiendo y actualizando el superblock
						err = EscribirObjeto(file, tempSuperblock, int64(superblock_start))
						if err != nil {
							fmt.Println("Error: ", err)
						}*/

				}

			}

			// escribiendo blocks
			err = EscribirObjeto(file, &newPointer1, int64(tempSuperblock.S_block_start+numBlock1*int32(binary.Size(Pointerblock{})))) //Bloque 0
			if err != nil {
				fmt.Println("Error: ", err)
			}

			/*
				// escribiendo y actualizando el superblock
				err = EscribirObjeto(file, tempSuperblock, int64(superblock_start))
				if err != nil {
					fmt.Println("Error: ", err)
				}*/

		}

		//execute -path=/home/darkun/Escritorio/prueba.mia

		if indice == 15 && len(resto_cadena[x:]) != 0 {

			if block == 12 {

				//fmt.Println("\nCreando Bloque de apuntadores Doble")

				return CreateNewIndirectBlock(resto_cadena, x, crrInode, num_inodo, tempSuperblock, file, superblock_start, block+1, 2)
			}

			if block == 13 {

				fmt.Println("\nCreando Bloque de apuntadores Triple")

				//CreateNewIndirectBlock(resto_cadena, x, crrInode, num_inodo, tempSuperblock, file, superblock_start, block+1, 3)
			}

		}

	} else if crrInode.I_block[block] == -1 && block == 13 { //si el bloque de inodo no esta vacio, pasa al siguiente

		fmt.Println("\nBloque de apuntadores doble")

		return CreateNewIndirectBlock(resto_cadena, x, crrInode, num_inodo, tempSuperblock, file, superblock_start, block+1, 2)

	} else if crrInode.I_block[block] == -1 && block == 14 {

		fmt.Println("\nBloque de apuntadores triple")

		//CreateNewIndirectBlock(resto_cadena, x, crrInode, num_inodo, tempSuperblock, file, superblock_start, block+1, 2)
	}

	fmt.Println("\n\n========================= Fin CreateNewIndirectBlock ===========================")

	return tempSuperblock
}

func VerificaTipoArchivo(carpeta []string, cont_folder int, crrFolderBlock *Folderblock, tempSuperblock Superblock, file *os.File, block int32, cont string, superblock_start int32, k int) error {

	fmt.Println("\n\n========================= Iniciando VerificaTipoArchivo ===========================")
	//execute -path=/home/darkun/Escritorio/prueba.mia

	fmt.Println("\n         Es un archivo")
	fmt.Println("\nEl valor de k es: ", k)

	copy(crrFolderBlock.B_content[k].B_name[:], []byte(carpeta[cont_folder]))

	newInodo := tempSuperblock.S_inodes_count - tempSuperblock.S_free_inodes_count

	crrFolderBlock.B_content[k].B_inodo = newInodo

	var newInode Inode
	// escribiendo blocks
	tempSuperblock.S_free_inodes_count -= 1

	userid, _ := strconv.Atoi(user_.Uid)
	newInode.I_uid = int32(userid)

	groupid, _ := strconv.Atoi(user_.Gid)
	newInode.I_gid = int32(groupid)

	copy(newInode.I_type[:], "1")
	copy(newInode.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		newInode.I_block[i] = -1 //-1 no han sido utilizados
	}

	//execute -path=/home/darkun/Escritorio/prueba.mia

	for k, block := range newInode.I_block {

		if block == -1 {

			newBlock := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count
			newInode.I_block[k] = int32(newBlock)

			var newCadena string

			if cont != "" { // si hay parametro cont

				archivo_, err := os.Open(cont)

				if err != nil {
					fmt.Println("Error al abrir el archivo:", err)
					return err
				}

				scanner := bufio.NewScanner(archivo_)

				for scanner.Scan() {

					linea := scanner.Text()

					if linea != "" {

						for i := 0; i < len(linea); i++ {

							if string(linea[i]) == "#" {

								break
							}

							newCadena += string(linea[i])
						}

					}

				}
			} else { // si no hay parametro cont
				var cont int

				var cont_num int

				for cont < size_ {

					newCadena += strconv.Itoa(cont_num)
					cont_num++

					if cont_num == 10 {
						cont_num = 0
					}

					cont++
				}
			}

			var newFile Fileblock

			tempSuperblock.S_free_blocks_count -= 1

			var c int

			for m := 0; m < len(newFile.B_content[:]); m++ {

				if newFile.B_content[m] == 0 { // si hay todavia espacio

					if c < len(newCadena) {

						newFile.B_content[m] = byte(newCadena[c])
						//fmt.Printf("agregando letra:  %s   ", string(newCadena[c]))
						c++

					} else {
						break
					}

				}
			}

			var espacios int

			for i := 0; i < len(newFile.B_content); i++ {

				if newFile.B_content[i] == 0 {
					espacios++
				}
			}

			err := EscribirObjeto(file, newInode, int64(tempSuperblock.S_inode_start+newInodo*int32(binary.Size(Inode{})))) //Inode 0 carpeta raiz

			if err != nil {
				fmt.Println("Error: ", err)
			}

			//Escribiendo en bitmap inodos

			err = EscribirObjeto(file, byte(1), int64(tempSuperblock.S_bm_inode_start+newInodo))
			if err != nil {
				fmt.Println("Error: ", err)
			}

			fmt.Println("\n\n ****Escribiendo objeto FILEBLOCK en el archivo *****")

			if err := EscribirObjeto(file, newFile, int64(tempSuperblock.S_block_start+newBlock*int32(binary.Size(Fileblock{})))); err != nil { //aqui solo escribi el primer EBR
				return err
			}

			// Escribiendo en bitmap de bloques
			err = EscribirObjeto(file, byte(1), int64(tempSuperblock.S_bm_block_start+newBlock))
			if err != nil {
				fmt.Println("Error: ", err)
			}

			// escribiendo y actualizando el superblock
			err = EscribirObjeto(file, tempSuperblock, int64(superblock_start))
			if err != nil {
				fmt.Println("Error: ", err)
			}

			//data := "1,G,root\n1,U,root,root,123\n"
			if espacios > 0 {

				fmt.Println("\n Todavia sobra espacio despues de escribir la cadena en el slice")

				return nil

			} else { // si ya no hay espacios en el slice para ingresar la cadena

				if len(newCadena[c:]) != 0 { //si todavia hay caracteres en newCadena para seguir ingresando en slice de fileblock.Bcontent

					//fmt.Println("\n Enviando cadena:\n", newCadena)
					fmt.Println("\nEnviando c: ", c, " Enviando inodos: ", newInodo)
					CreateNewBlock(newCadena, c, newInode, newInodo, &tempSuperblock, file, superblock_start)

					// escribiendo y actualizando el superblock
					err = EscribirObjeto(file, tempSuperblock, int64(superblock_start))
					if err != nil {
						fmt.Println("Error: ", err)
					}

					numBlocks := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count

					fmt.Println("\nEl nuevo numero de bloques ocupados es(despues de createnewblock): ", numBlocks)

					fmt.Println("\n\n========================= Finalizando VerificaTipoArchivo ===========================")

					return nil
				}

			}

		}

	}

	//execute -path=/home/darkun/Escritorio/prueba.mia

	return nil
}
