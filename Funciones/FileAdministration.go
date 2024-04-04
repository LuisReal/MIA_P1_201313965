package Funciones

import (
	//"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Cat(path string) error {

	//execute -path=/home/darkun/Escritorio/prueba.mia

	fmt.Println("\n\n========================= Iniciando (Cat)===========================")

	array := strings.Split(path, "/")

	array = array[1:]

	fmt.Println("\nEl nuevo array es: ", array)

	//cat -file=/user.txt

	id := strings.ToUpper(user_.Id)

	file, tempSuperblock, _, err := getSuperBlock(id)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	Inodo_start := tempSuperblock.S_inode_start

	var inodo0 Inode

	if err := LeerObjeto(file, &inodo0, int64(Inodo_start)); err != nil {
		return err
	}
	var cont_folder int

	indexInode, _ := SearchingFile(array, inodo0, file, tempSuperblock, -1, cont_folder)

	fmt.Println("\nindexInode el valor que devuelve SearchingFile: ", indexInode)

	var crrInode Inode //Inodo 1

	if err := LeerObjeto(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Inode{})))); err != nil {
		return err
	}

	var bloque int
	var fileblock Fileblock
	var cadena string
	var fileblock_start int32

	for i := 0; i < len(crrInode.I_block); i++ { //iterando bloques de inodo1

		if crrInode.I_block[i] != -1 {

			bloque = int(crrInode.I_block[i]) //obtiene el numero del ultimo bloque de archivos creado

			fileblock_start = tempSuperblock.S_block_start + int32(bloque)*int32(binary.Size(Fileblock{}))

			if err := LeerObjeto(file, &fileblock, int64(fileblock_start)); err != nil { //bloque1
				return err
			}

			cadena += string(fileblock.B_content[:])

		}
	}

	fmt.Println("\n                         *************************** Mostrando el contenido del archivo *****************************\n\n", cadena)

	fmt.Println("\n\n========================= Finalizando (Cat)===========================")

	//execute -path=/home/darkun/Escritorio/prueba.mia

	return nil
}

func SearchingFile(carpeta []string, Inodo Inode, file *os.File, tempSuperblock Superblock, Inodo_actual int32, cont_folder int) (int32, int32) {

	fmt.Println("\n\n========================= Iniciando (SearchingFile)===========================")

	var folder_bytes [12]byte

	fmt.Println("\nEl valor de cont_folder es: ", cont_folder)

	//execute -path=/home/darkun/Escritorio/prueba.mia

	copy(folder_bytes[:], []byte(carpeta[cont_folder]))

	indice := int32(0)

	// Iterate over i_blocks from Inode
	for _, block := range Inodo.I_block {
		if block != -1 {
			if indice < 13 {

				// carpeta -> 0   archivo ->1

				if string(Inodo.I_type[:]) == "0" { // si es carpeta

					var crrFolderBlock Folderblock

					// Read object from bin file                                       // un Inodo y un bloque las estructuras miden lo mismo
					if err := LeerObjeto(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))); err != nil {
						return -1, -1
					}

					//home/archivos/user/docs/Tarea2.txt

					for _, folder := range crrFolderBlock.B_content {

						//fmt.Println("\nEl valor de k es: ", k)

						if string(folder.B_name[:]) == string(folder_bytes[:]) {

							//execute -path=/home/darkun/Escritorio/prueba.mia

							fmt.Println("\nLa carpeta "+carpeta[cont_folder]+" si existe en el bloque: ", block)

							if len(carpeta) == cont_folder+1 {
								fmt.Println("\nEstoy if len(carpet) == cont_folder+1")
								fmt.Println("\n\n========================= Finalizando (SearchingFile)===========================")
								return folder.B_inodo, block
							} else {
								fmt.Println("\n ======= NextInode ======")
								var NextInode Inode
								// Read object from bin file
								if err := LeerObjeto(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(Inode{})))); err != nil {
									return -1, -1
								}

								fmt.Println("Enviando el INODO: " + strconv.Itoa(int(folder.B_inodo)))

								return SearchingFile(carpeta, NextInode, file, tempSuperblock, folder.B_inodo, cont_folder+1)

							}

						}

					}

				}

			}

		}

	}

	return 0, 0
}

func Remove(path string) error {

	//execute -path=/home/darkun/Escritorio/prueba.mia

	fmt.Println("\n\n========================= Iniciando (Remove)===========================")

	array := strings.Split(path, "/")

	array = array[1:]

	fmt.Println("\nEl nuevo array es: ", array)

	//cat -file=/user.txt

	id := strings.ToUpper(user_.Id)

	file, tempSuperblock, superblock_start, err := getSuperBlock(id)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	Inodo_start := tempSuperblock.S_inode_start

	var inodo0 Inode

	if err := LeerObjeto(file, &inodo0, int64(Inodo_start)); err != nil {
		return err
	}
	var cont_folder int

	indexInode, bloque_anterior := SearchingFile(array, inodo0, file, tempSuperblock, -1, cont_folder)

	fmt.Println("\nindexInode el valor que devuelve SearchingFile: ", indexInode)

	fmt.Println("\nEl bloque anterior es: ", bloque_anterior)

	//Recuperando y leyendo el bloque anterior(bloque_anterior) del inodo indexInode

	var bloque Folderblock

	if err := LeerObjeto(file, &bloque, int64(tempSuperblock.S_block_start+bloque_anterior*int32(binary.Size(Folderblock{})))); err != nil {
		return err
	}

	var folder_bytes [12]byte

	//execute -path=/home/darkun/Escritorio/prueba.mia

	copy(folder_bytes[:], []byte(array[len(array)-1]))

	for k, folder := range bloque.B_content {

		if string(folder.B_name[:]) == string(folder_bytes[:]) {
			fmt.Println("\nCarpeta " + string(folder_bytes[:]) + " Encontrada")
			bloque.B_content[k].B_inodo = -1
			copy(bloque.B_content[k].B_name[:], "            ") //ingresando 12 espacios vacios
		}

	}

	//Recuperando y Escribiendo el bloque anterior(bloque_anterior) del inodo indexInode
	err = EscribirObjeto(file, bloque, int64(tempSuperblock.S_block_start+bloque_anterior*int32(binary.Size(Folderblock{}))))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	var crrInode Inode //Inodo 1

	if err := LeerObjeto(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Inode{})))); err != nil {
		return err
	}

	Removing(crrInode, indexInode, file, &tempSuperblock)

	// escribiendo y actualizando el superblock
	err = EscribirObjeto(file, tempSuperblock, int64(superblock_start))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("\n\n========================= Finalizando (Remove)===========================")

	//execute -path=/home/darkun/Escritorio/prueba.mia

	return nil
}

func Removing(Inodo Inode, numInodo int32, file *os.File, tempSuperblock *Superblock) error {

	fmt.Println("\n\n========================= Iniciando (Removing)===========================")

	indice := int32(0)

	for _, block := range Inodo.I_block {
		if block != -1 {
			if indice < 13 {

				// carpeta -> 0   archivo ->1

				if string(Inodo.I_type[:]) == "0" { // si es carpeta

					fmt.Println("\nInodo de tipo carpeta")

					var crrFolderBlock Folderblock

					// Read object from bin file                                       // un Inodo y un bloque las estructuras miden lo mismo
					if err := LeerObjeto(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))); err != nil {
						return err
					}

					//home/archivos/user/docs/Tarea2.txt

					for _, folder := range crrFolderBlock.B_content {

						if folder.B_inodo > numInodo {

							fmt.Println("\n ======= NextInode ======")
							var NextInode Inode
							// Read object from bin file
							if err := LeerObjeto(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(Inode{})))); err != nil {
								return err
							}

							fmt.Println("Enviando el INODO: " + strconv.Itoa(int(folder.B_inodo)))

							Removing(NextInode, folder.B_inodo, file, tempSuperblock)
						}

					}

					//Empieza a eliminar cada bloque del inodo, ya sea de archivo o carpeta ya que miden lo mismo 64bytes

					fmt.Println("\n\n        ****** Eliminando el bloque: ", block)
					fmt.Println()

					var tempFolderblock Folderblock

					err := EscribirObjeto(file, tempFolderblock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{}))))
					if err != nil {
						fmt.Println("Error: ", err)
					}

					// Escribiendo en bitmap de bloques
					err = EscribirObjeto(file, byte(0), int64(tempSuperblock.S_bm_block_start+block))
					if err != nil {
						fmt.Println("Error: ", err)
					}

				} else { // si es archivo
					fmt.Println("\nInodo de tipo archivo")

					fmt.Println("\n\n        ****** Eliminando el bloque: ", block)
					fmt.Println()

					var tempFileblock Fileblock

					err := EscribirObjeto(file, tempFileblock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Fileblock{}))))
					if err != nil {
						fmt.Println("Error: ", err)
					}

					// Escribiendo en bitmap de bloques
					err = EscribirObjeto(file, byte(0), int64(tempSuperblock.S_bm_block_start+block))
					if err != nil {
						fmt.Println("Error: ", err)
					}
				}
			}
		}
	}

	fmt.Println("\n\n        ****** Eliminando el Inodo: ", numInodo)
	fmt.Println()

	var tempInodo Inode

	err := EscribirObjeto(file, tempInodo, int64(tempSuperblock.S_inode_start+numInodo*int32(binary.Size(Inode{}))))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	//Escribiendo en bitmap inodos

	err = EscribirObjeto(file, byte(0), int64(tempSuperblock.S_bm_inode_start+numInodo))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("\n\n========================= Finalizando (Removing)===========================")
	return nil

	//execute -path=/home/darkun/Escritorio/prueba.mia

}

func Move(path string, dest string) error {

	//execute -path=/home/darkun/Escritorio/prueba.mia

	fmt.Println("\n\n========================= Iniciando (Move)===========================")

	fmt.Println("\n ******* COMANDO NO REALIZADO ***********")

	fmt.Println("\n\n========================= Finalizando (Move)===========================")

	//execute -path=/home/darkun/Escritorio/prueba.mia

	return nil
}

func getSuperBlock(id string) (*os.File, Superblock, int32, error) {

	driveletter := string(id[0])

	// Open bin file
	filepath := "./archivos/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := AbrirArchivo(filepath)
	if err != nil {
		return nil, Superblock{}, 0, err
	}

	var TempMBR MBR
	// Read object from bin file
	if err := LeerObjeto(file, &TempMBR, 0); err != nil {
		return nil, Superblock{}, 0, err
	}

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
					return nil, Superblock{}, 0, err
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
		return nil, Superblock{}, 0, err
	}

	var tempSuperblock Superblock

	if err := LeerObjeto(file, &tempSuperblock, int64(TempMBR.Mbr_partitions[index].Part_start)); err != nil {
		return nil, Superblock{}, 0, err
	}

	superblock_start := TempMBR.Mbr_partitions[index].Part_start
	// initSearch /user.txt -> regresa no Inodo
	// initSearch -> 1

	// getInodeFileData -> Iterate the I_Block n concat the data

	return file, tempSuperblock, superblock_start, nil

}
