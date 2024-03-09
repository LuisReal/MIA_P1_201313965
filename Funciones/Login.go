package Funciones

import (
	"encoding/binary"
	"fmt"
	"strings"
)

func Login(user string, pass string, id string) {
	fmt.Println("\n\n========================= LOGIN ===========================")

	fmt.Printf("\nUser: %s, pass: %s, id: %s\n", user, pass, id)

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
	fmt.Println("\nImprimiendo MBR")
	PrintMBR(TempMBR)
	fmt.Println("\nFinalizando Impresion de MBR")

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
					return
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
		return
	}

	var tempSuperblock Superblock

	if err := LeerObjeto(file, &tempSuperblock, int64(TempMBR.Mbr_partitions[index].Part_start)); err != nil {
		return
	}

	// initSearch /users.txt -> regresa no Inodo
	// initSearch -> 1

	indexInode := int32(1)

	var crrInode Inode

	if err := LeerObjeto(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Inode{})))); err != nil {
		return
	}

	// getInodeFileData -> Iterate the I_Block n concat the data

	var fileblock Fileblock

	if err := LeerObjeto(file, &fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(Fileblock{})))); err != nil {
		return
	}

	fmt.Println("Fileblock------------")
	//data := "1,G,root\n1,U,root,root,123\n"
	data := string(fileblock.B_content[:])

	lines := strings.Split(data, "\n")

	for _, line := range lines {
		// Imprimir cada línea
		fmt.Println(line)

	}

	fmt.Println("Inode", crrInode.I_block)

	// Close bin file
	defer file.Close()

	fmt.Println("\n\n========================= FIN LOGIN ===========================")
}
