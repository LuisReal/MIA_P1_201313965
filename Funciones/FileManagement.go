package Funciones

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func Mkfs(id string, type_ string, fs_ string) {
	fmt.Println("\n\n=========================Iniciando MKFS===========================")
	fmt.Println()

	driveletter := string(id[0])

	// Open bin file
	filepath := "./archivos/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := abrirArchivo(filepath)
	if err != nil {
		return
	}

	var TempMBR MBR
	// Read object from bin file
	if err := LeerObjeto(file, &TempMBR, 0); err != nil {
		return
	}

	// Print object
	PrintMBR(TempMBR)

	var index int = -1
	// Iterate over the partitions
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_partitions[i].Part_size != 0 {
			if strings.Contains(string(TempMBR.Mbr_partitions[i].Part_id[:]), id) {
				fmt.Println("\n*********************Particion Encontrada****************")
				if TempMBR.Mbr_partitions[i].Part_status { // si la particion es true es porque esta montada
					fmt.Println("\n ********************La Particion esta montada**********************")
					index = i
				} else {
					fmt.Println("\n ********************La Particion no esta montada**********************")
					return
				}
				break
			}
		}
	}

	if index != -1 {
		ImprimirParticion(TempMBR.Mbr_partitions[index])
	} else {
		fmt.Println("\n*********************Particion NO Encontrada****************")
		return
	}

	numerador := int32(TempMBR.Mbr_partitions[index].Part_size - int32(binary.Size(Superblock{})))
	denominador_base := int32(4 + int32(binary.Size(Inode{})) + 3*int32(binary.Size(Fileblock{})))
	var temp int32 = 0
	if fs_ == "2fs" {
		temp = 0
	} else {
		temp = int32(binary.Size(Journaling{}))
	}
	denominador := denominador_base + temp
	n := int32(numerador / denominador)

	fmt.Println("\n*************************El numero de estructuras N es: ", n)

	// var newMRB Structs.MRB
	var newSuperblock Superblock
	newSuperblock.S_inodes_count = 0
	newSuperblock.S_blocks_count = 0

	newSuperblock.S_free_blocks_count = 3 * n
	newSuperblock.S_free_inodes_count = n

	//copy(newSuperblock.S_mtime[:], "06/03/2024")      Esto no se evaluara
	//copy(newSuperblock.S_umtime[:], "06/03/2024")		Esto no se evaluara
	//newSuperblock.S_mnt_count = 0                    (No se evaluara cuantas veces fue montado el sistema)

	if fs_ == "2fs" {
		ext2(n, TempMBR.Mbr_partitions[index], newSuperblock, file)
	} else {
		fmt.Println("EXT3")
	}

	// Close bin file
	defer file.Close()

	fmt.Println("\n\n=========================Finalizando MKFS===========================")
}

func ext2(n int32, partition Partition, newSuperblock Superblock, file *os.File) {

	fmt.Println("\n\n=========================Creando ext2===========================")
	fmt.Println("\nFile: ", file)

	newSuperblock.S_filesystem_type = 2
	newSuperblock.S_bm_inode_start = partition.Part_start + int32(binary.Size(Superblock{}))
	newSuperblock.S_bm_block_start = newSuperblock.S_bm_inode_start + n
	newSuperblock.S_inode_start = newSuperblock.S_bm_block_start + 3*n
	newSuperblock.S_block_start = newSuperblock.S_inode_start + n*int32(binary.Size(Inode{})) // n = numero de inodos

	// se resta dos veces -1 porque hay que crear dos inodos y dos bloques
	newSuperblock.S_free_inodes_count -= 1
	newSuperblock.S_free_blocks_count -= 1
	newSuperblock.S_free_inodes_count -= 1
	newSuperblock.S_free_blocks_count -= 1

	// escribiendo ceros en bitmap de inodos
	for i := int32(0); i < n; i++ {
		err := escribirObjeto(file, byte(0), int64(newSuperblock.S_bm_inode_start+i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	// escribiendo ceros en bitmap de bloques
	for i := int32(0); i < 3*n; i++ {
		err := escribirObjeto(file, byte(0), int64(newSuperblock.S_bm_block_start+i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	// llenando con -1 los primeros 15 bloques
	var newInode Inode
	for i := int32(0); i < 15; i++ {
		newInode.I_block[i] = -1 //-1 no han sido utilizados
	}

	// escribiendo los inodos en el archivo

	for i := int32(0); i < n; i++ {
		err := escribirObjeto(file, newInode, int64(newSuperblock.S_inode_start+i*int32(binary.Size(Inode{}))))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	// escribiendo los bloques en el archivo (los fileblock y folderblock tiene el mismo tamano de 64bytes)

	var newFileblock Fileblock
	for i := int32(0); i < 3*n; i++ {
		err := escribirObjeto(file, newFileblock, int64(newSuperblock.S_block_start+i*int32(binary.Size(Fileblock{}))))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	//creando el primer Inode en posicion 0
	var Inode0 Inode //Inode 0
	Inode0.I_uid = 1
	Inode0.I_gid = 1
	Inode0.I_size = 0
	copy(Inode0.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode0.I_block[i] = -1
	}

	Inode0.I_block[0] = 0 // el inode 0 apunta al bloque 0

	// . | 0
	// .. | 0
	// users.txt | 1
	//

	var Folderblock0 Folderblock //Bloque 0 -> carpetas
	Folderblock0.B_content[0].B_inodo = 0
	copy(Folderblock0.B_content[0].B_name[:], ".")
	Folderblock0.B_content[1].B_inodo = 0
	copy(Folderblock0.B_content[1].B_name[:], "..")
	Folderblock0.B_content[2].B_inodo = 1
	copy(Folderblock0.B_content[2].B_name[:], "users.txt")

	//creando el Inode 1

	var Inode1 Inode //Inode 1
	Inode1.I_uid = 1
	Inode1.I_gid = 1
	Inode1.I_size = int32(binary.Size(Folderblock{}))

	copy(Inode1.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode1.I_block[i] = -1
	}

	Inode1.I_block[0] = 1 // el Inode 1 apunta al bloque 1

	data := "1,G,root\n1,U,root,root,123\n"
	var Fileblock1 Fileblock //Bloque 1 -> archivo
	copy(Fileblock1.B_content[:], data)

	// Inodo 0 -> Bloque 0 -> Inodo 1 -> Bloque 1
	// Crear la carpeta raiz /
	// Crear el archivo users.txt "1,G,root\n1,U,root,root,123\n"

	// escribiendo el superblock
	err := escribirObjeto(file, newSuperblock, int64(partition.Part_start))

	// escribiendo bitmap inodes con unos

	err = escribirObjeto(file, byte(1), int64(newSuperblock.S_bm_inode_start))
	err = escribirObjeto(file, byte(1), int64(newSuperblock.S_bm_inode_start+1))

	// escribiendo bitmap blocks con unos
	err = escribirObjeto(file, byte(1), int64(newSuperblock.S_bm_block_start))
	err = escribirObjeto(file, byte(1), int64(newSuperblock.S_bm_block_start+1))

	fmt.Println("Inode 0:", int64(newSuperblock.S_inode_start))
	fmt.Println("Inode 1:", int64(newSuperblock.S_inode_start+int32(binary.Size(Inode{}))))

	// escribiendo inodes
	err = escribirObjeto(file, Inode0, int64(newSuperblock.S_inode_start))                             //Inode 0
	err = escribirObjeto(file, Inode1, int64(newSuperblock.S_inode_start+int32(binary.Size(Inode{})))) //Inode 1

	// escribiendo blocks
	err = escribirObjeto(file, Folderblock0, int64(newSuperblock.S_block_start))                               //Bloque 0
	err = escribirObjeto(file, Fileblock1, int64(newSuperblock.S_block_start+int32(binary.Size(Fileblock{})))) //Bloque 1

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("\n\n=========================Finalizando ext2===========================")
}
