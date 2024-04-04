package Funciones

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func Mkdir(path string, r string) error { //mkdir -path=/bin        // permisos (rw- rw- r--) = 664

	fmt.Println("\n\n=========================Creando Carpeta (MKDIR)===========================")

	//fmt.Println("\nUsuario con sesion actual: ", user_.Nombre, " ID: ", user_.Id, " r: ", r)

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

	carpetas = carpetas[1:]

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

				copy(journaling.Contenido[k].Operation[:], "mkdir")
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

	if r == "0" { //el parametro r no esta incluido

		for i := 0; i < len(carpetas); i++ {

			AddingFolderRoot(carpetas[i], inodo0, file, tempSuperblock, superblock_start)

		}

	} else { // si esta incluido el parametro r

		var cont_folder int

		AddingNewFolder(carpetas, inodo0, file, tempSuperblock, -1, superblock_start, cont_folder)

	}

	fmt.Println("\n\n=======================Finalizando Creacion Carpeta (MKDIR)===========================")

	return nil
}

func AddingNewFolder(carpeta []string, Inodo Inode, file *os.File, tempSuperblock Superblock, Inodo_actual int32, superblock_start int32, cont_folder int) error {

	//mkdir -r -path="/home/archivos/archivos"19"

	fmt.Println("\n\n========================= Iniciando AddingNewFolder ===========================")

	var folder_bytes [12]byte

	fmt.Println("\nEl valor de cont_folder es: ", cont_folder)

	if cont_folder == len(carpeta) {
		fmt.Println("\n return : debido a que cont_folder es igual que la longitud de la carpeta arrreglos")
		return nil
	}

	//execute -path=/home/darkun/Escritorio/prueba.mia

	copy(folder_bytes[:], []byte(carpeta[cont_folder]))

	indice := int32(0)

	fmt.Println("\nIngresando Carpeta: ", carpeta[cont_folder])

	//fmt.Println("El tipo de inodo es: ", string(Inodo.I_type[:]))

	// Iterate over i_blocks from Inode
	for j, block := range Inodo.I_block {
		if block != -1 {
			if indice < 13 {

				//CASO DIRECTO

				//// carpeta -> 0   archivo ->1

				if string(Inodo.I_type[:]) == "0" { // si es carpeta

					var crrFolderBlock Folderblock

					/*
						for i := int32(0); i < 4; i++ {
							crrFolderBlock.B_content[i].B_inodo = -1
						}*/

					//grafo += `Bloque` + strconv.Itoa(int(block)) + `:0;` + "\n"
					// Read object from bin file                                       // un Inodo y un bloque las estructuras miden lo mismo
					if err := LeerObjeto(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))); err != nil {
						return err
					}

					//mkdir -r -path="/home/archivos/archivos"19"

					var exist_folder int

					for _, folder := range crrFolderBlock.B_content {

						//fmt.Println("\nEl valor de k es: ", k)

						if string(folder.B_name[:]) == string(folder_bytes[:]) {

							//fmt.Println("\nLa carpeta "+carpeta[cont_folder]+" si existe en el bloque: ", block)

							//fmt.Println("\n ======= NextInode ======")
							var NextInode Inode
							// Read object from bin file
							if err := LeerObjeto(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(Inode{})))); err != nil {
								return err
							}

							fmt.Println("Enviando el INODO: " + strconv.Itoa(int(folder.B_inodo)))
							//grafo, _ = EnlazandoNodos(path, NextInode, file, tempSuperblock, disco, grafo, folder.B_inodo)
							exist_folder++

							AddingNewFolder(carpeta, NextInode, file, tempSuperblock, folder.B_inodo, superblock_start, cont_folder+1)

							return nil

						}

					} //execute -path=/home/darkun/Escritorio/prueba.mia

					if exist_folder == 0 { // si no existe la carpeta en el bloque de carpetas
						//fmt.Println("\nComo NO existe la carpeta en el bloque: ", block)
						//mkdir -r -path="/home/archivos/archivos"19"

						for k, folder := range crrFolderBlock.B_content { // creo la carpeta si todavia existe espacio en el bloque

							if folder.B_inodo == -1 {
								//fmt.Println("\nIngresando la carpeta " + carpeta[cont_folder] + " en el bloque " + strconv.Itoa(int(block)))

								inodo := tempSuperblock.S_inodes_count - tempSuperblock.S_free_inodes_count

								//fmt.Println("\nEl valor de inodo es: ", inodo)

								crrFolderBlock.B_content[k].B_inodo = inodo

								copy(crrFolderBlock.B_content[k].B_name[:], []byte(carpeta[cont_folder]))

								err := EscribirObjeto(file, crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))) //Inode 0 carpeta raiz

								if err != nil {
									fmt.Println("Error: ", err)
								}

								numBlocks := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count

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

								//fmt.Println("\nEl valor de numBlocks es: ", numBlocks)

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

								//execute -path=/home/darkun/Escritorio/prueba.mia

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

								//fmt.Println("\n ======= NextInode ======")
								var NextInode Inode
								// Read object from bin file
								if err := LeerObjeto(file, &NextInode, int64(tempSuperblock.S_inode_start+inodo*int32(binary.Size(Inode{})))); err != nil {
									return err
								}

								AddingNewFolder(carpeta, NextInode, file, tempSuperblock, inodo, superblock_start, cont_folder+1)

								//return nil
								return nil
							}
						}
					}

				}

			}

		} else { //execute -path=/home/darkun/Escritorio/prueba.mia

			///home/archivos/carpeta1/carpeta2/carpeta3/carpeta4/carpeta5
			/*
				fmt.Println("\nLeyendo el siguiente bloque " + strconv.Itoa(int(block)) + " vacio del inodo")
				fmt.Println("El indice J del bloque en el inodo es: ", j)
				fmt.Println("El cont_folder es: ", cont_folder)
				fmt.Println("El nuevo cont_folder es: ", cont_folder+1)

				fmt.Println("*******La carpeta a ingresar es: ", carpeta[cont_folder])*/

			newBlock := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count
			Inodo.I_block[j] = int32(newBlock)
			//fmt.Println("\nEl nuevo bloque es: ", Inodo.I_block[j])

			//fmt.Println("\nEl valor de Inodo actual es: ", Inodo_actual)
			//fmt.Println("\nEl Inodo "+strconv.Itoa(int(Inodo_actual))+" en su posicion "+strconv.Itoa(j)+" apunta al bloque : ", Inodo.I_block[j])

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
			//fmt.Println("\nEl nuevo numero de bloque es: ", newInode.I_block[0])

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

			//fmt.Println("\n ======= NextInode ======")
			var NextInode Inode
			// Read object from bin file
			if err := LeerObjeto(file, &NextInode, int64(tempSuperblock.S_inode_start+inodos_ocupados*int32(binary.Size(Inode{})))); err != nil {
				return err
			}

			AddingNewFolder(carpeta, NextInode, file, tempSuperblock, inodos_ocupados, superblock_start, cont_folder+1)

			return nil

			//execute -path=/home/darkun/Escritorio/prueba.mia

			///home/archivos/carpeta1/carpeta2/carpeta3/carpeta4/carpeta5
		}

		indice++
	}

	fmt.Println("\n\n========================= Finalizando AddingNewFolder ===========================")

	return nil
}

func AddingFolderRoot(carpeta string, Inodo Inode, file *os.File, tempSuperblock Superblock, superblock_start int32) error {

	fmt.Println("\n\n========================= Iniciando AddingFolderRoot ===========================")

	fmt.Println("\nIngresando Carpeta: ", carpeta)

	indice := int32(0)

	//fmt.Println("El tipo de inodo es: ", string(Inodo.I_type[:]))

	// Iterate over i_blocks from Inode
	for j, block := range Inodo.I_block {
		if block != -1 {
			if indice < 13 {

				//CASO DIRECTO

				//// carpeta -> 0   archivo ->1

				if string(Inodo.I_type[:]) == "0" { // si es carpeta
					//fmt.Println("Inodo de tipo carpeta")
					var crrFolderBlock Folderblock

					// Read object from bin file                                       // un Inodo y un bloque las estructuras miden lo mismo
					if err := LeerObjeto(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))); err != nil {
						return err
					}

					for k, folder := range crrFolderBlock.B_content {

						//execute -path=/home/darkun/Escritorio/prueba.mia

						if folder.B_inodo == -1 { // si es -1 es porque hay espacio para ingresar la nueva carpeta

							inodo := crrFolderBlock.B_content[k-1].B_inodo + int32(1)

							crrFolderBlock.B_content[k].B_inodo = inodo

							copy(crrFolderBlock.B_content[k].B_name[:], []byte(carpeta))

							err := EscribirObjeto(file, crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))) //Inode 0 carpeta raiz

							if err != nil {
								fmt.Println("Error: ", err)
							}

							numBlocks := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count
							//fmt.Println("\nEl numero de bloques ocupados es: ", numBlocks)

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

							newInode.I_block[0] = numBlocks

							err = EscribirObjeto(file, newInode, int64(tempSuperblock.S_inode_start+inodo*int32(binary.Size(Inode{})))) //Inode 0 carpeta raiz

							if err != nil {
								fmt.Println("Error: ", err)
							}

							//Escribiendo en bitmap de inodos
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

							return nil // detiene la recursividad

						}

					}

				} else {
					break // si no es de tipo carpeta se rompe el ciclo
				}

			} else {
				fmt.Println("CASO INDIRECTO")
			}
		} else { //creando un nuevo bloque en el Inodo0   carpeta raiz porque no habia espacio en el bloque anterior

			//execute -path=/home/darkun/Escritorio/prueba.mia

			newBlock := (tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count)
			Inodo.I_block[j] = int32(newBlock)

			tempSuperblock.S_free_blocks_count -= 1

			err := EscribirObjeto(file, Inodo, int64(tempSuperblock.S_inode_start)) //Inode 0 carpeta raiz

			if err != nil {
				fmt.Println("Error: ", err)
			}

			inodos_ocupados := tempSuperblock.S_inodes_count - tempSuperblock.S_free_inodes_count

			var newFolder Folderblock //Bloque 0 -> carpetas

			for i := int32(0); i < 4; i++ {
				newFolder.B_content[i].B_inodo = -1
			}

			newFolder.B_content[0].B_inodo = inodos_ocupados
			copy(newFolder.B_content[0].B_name[:], []byte(carpeta))

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

			// llenando con -1 los primeros 15 bloques
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
			//fmt.Println("\nEl nuevo numero de bloque es: ", newInode.I_block[0])

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

			return nil

		}

		indice++
	}

	fmt.Println("\n\n========================= Finalizando AddingFolderRoot ===========================")

	return nil

	//execute -path=/home/darkun/Escritorio/basico.mia
	//execute -path=/home/darkun/Escritorio/avanzado.mia
}
