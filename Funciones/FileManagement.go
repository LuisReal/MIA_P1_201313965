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
		ext2(n, TempMBR.Mbr_partitions[index], newSuperblock, "06/03/2024", file)
	} else {
		fmt.Println("EXT3")
	}

	// Close bin file
	defer file.Close()

	fmt.Println("\n\n=========================Finalizando MKFS===========================")
}

func ext2(n int32, partition Partition, newSuperblock Superblock, date string, file *os.File) {

	fmt.Println("\n\n=========================Creando ext2===========================")
	fmt.Println("\nDate: ", date)
	fmt.Println("\nFile: ", file)

	newSuperblock.S_filesystem_type = 2
	newSuperblock.S_bm_inode_start = partition.Part_start + int32(binary.Size(Superblock{}))
	newSuperblock.S_bm_block_start = newSuperblock.S_bm_inode_start + n
	newSuperblock.S_inode_start = newSuperblock.S_bm_block_start + 3*n
	newSuperblock.S_block_start = newSuperblock.S_inode_start + n*int32(binary.Size(Inode{}))

	newSuperblock.S_free_inodes_count -= 1
	newSuperblock.S_free_blocks_count -= 1
	newSuperblock.S_free_inodes_count -= 1
	newSuperblock.S_free_blocks_count -= 1

	fmt.Println("\n\n=========================Finalizando ext2===========================")
}
