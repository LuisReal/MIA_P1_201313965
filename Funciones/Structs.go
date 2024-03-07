package Funciones

import (
	"fmt"
)

type MBR struct {
	Mbr_tamano int32

	Mbr_fecha_creacion [16]byte // de tipo time

	Mbr_dsk_signature int32

	Dsk_fit [1]byte // B (mejor ajuste)  F(primer ajuste) W(peor ajuste)

	Mbr_partitions [4]Partition // este arreglo simulara las 4 particiones

}

func PrintMBR(data MBR) {
	fmt.Println(fmt.Sprintf("CreationDate: %s, fit: %s, size: %d", string(data.Mbr_fecha_creacion[:]), string(data.Dsk_fit[:]), data.Mbr_tamano))

	for i := 0; i < 4; i++ {

		fmt.Println(fmt.Sprintf(" Particion: %d, Tipo de Particion:  %s, Tamano de Particion: %d, start: %d, id: %s, correlativo: %d, status: %t", i, string(data.Mbr_partitions[i].Part_type[:]), int(data.Mbr_partitions[i].Part_size), int(data.Mbr_partitions[i].Part_start), string(data.Mbr_partitions[i].Part_id[:]), int(data.Mbr_partitions[i].Part_correlative), data.Mbr_partitions[i].Part_status))
	}

	/*
		for i := 0; i < 4; i++ {
			fmt.Println(fmt.Sprintf("Partition %d: %s, %s, %d, %d", i, string(data.Mbr_partitions[i].Name[:]), string(data.Mbr_partitions[i].Type[:]), data.Mbr_partitions[i].Start, data.Mbr_partitions[i].Size))
		}*/
}

type Partition struct {
	Part_status bool // es de tipo bool(indica si la particion esta montada o no)

	Part_type [1]byte //(indica el tipo de particion: primaria(P) o extendida(E))

	Part_fit [1]byte // indica el tipo de ajuste(B mejor ajuste  F primer ajuste W peor ajuste)

	Part_start int32 // indica en que byte del disco inicia la particion

	Part_size int32 //(part_s) contiene el tamano total de la particion en bytes (por defecto es cero)

	Part_name [16]byte // contiene el nombre de la particion

	Part_correlative int32 // contiene el correlativo de la particion

	Part_id [4]byte
}

func ImprimirParticion(data Partition) {
	fmt.Println(fmt.Sprintf("Name: %s, type: %s, start: %d, size: %d, status: %t, id: %s", string(data.Part_name[:]), string(data.Part_type[:]), data.Part_start, data.Part_size, data.Part_status, string(data.Part_id[:])))
}

type EBR struct { //extended boot record

	Part_mount bool // indica si la particion esta montada o no

	Part_fit [16]byte // indica el tipo de ajuste de la particion(B mejor ajuste F primer ajuste W peor ajuste)

	Part_start int32 // indica en que byte del disco inicia la particion

	Part_size int32 //contiene el tamano total de la particion en bytes

	Part_next int32 // byte en el que esta el proximo EBR . -1 si no hay siguiente

	Part_name [16]byte // nombre de la particion
}

func PrintEBR(data EBR) {

	fmt.Printf("mount: %t, fit: %s, start: %d, size: %d, next: %d, name: %s", data.Part_mount, string(data.Part_fit[:]), int(data.Part_start), int(data.Part_size), int(data.Part_next), string(data.Part_name[:]))
	fmt.Println()

}

type Superblock struct {
	S_filesystem_type   int32    //Guarda el número que identifica el sistema de archivos utilizado
	S_inodes_count      int32    //Guarda el número total de inodos
	S_blocks_count      int32    //Guarda el número total de bloques
	S_free_blocks_count int32    //Contiene el número de bloques libres
	S_free_inodes_count int32    //Contiene el número de inodos libres
	S_mtime             [17]byte //Última fecha en el que el sistema fue montado
	S_umtime            [17]byte //Última fecha en que el sistema fue desmontado
	S_mnt_count         int32    //Indica cuantas veces se ha montado el sistema
	S_magic             int32    //Valor que identifica al sistema de archivos, tendrá el valor 0xEF53
	S_inode_size        int32    //Tamaño del inodo
	S_block_size        int32    //Tamaño del bloque
	S_fist_inode        int32    //Primer inodo libre
	S_first_block       int32    //Primer bloque libre
	S_bm_inode_start    int32    //Guardará el inicio del bitmap de inodos
	S_bm_block_start    int32    //Guardará el inicio del bitmap de bloques
	S_inode_start       int32    //Guardará el inicio de la tabla de inodos
	S_block_start       int32    //Guardará el inicio de la tabla de bloques
}

type Inode struct {
	I_uid   int32    //UID del usuario propietario del archivo o carpeta
	I_gid   int32    //GID del grupo al que pertenece el archivo o carpeta.
	I_size  int32    //Tamaño del archivo en bytes
	I_atime [17]byte //Última fecha en que se leyó el inodo sin modificarlo
	I_ctime [17]byte //Fecha en la que se creó el inodo
	I_mtime [17]byte //Última fecha en la que se modifica el inodo
	I_block [15]int32
	I_type  [1]byte // Indica si es archivo o carpeta. 1 = Archivo, 0 = Carpeta
	I_perm  [3]byte
}

/*(I_block)Array en los que los primeros 12 registros son bloques directos.
El 13 será el número del bloque simple indirecto.
El 14 será el número del bloque doble indirecto.
El 15 será el número del bloque triple indirecto.
Si no son utilizados tendrá el valor: -1.
*/

/*(I_perm) Guardará los permisos del archivo o carpeta, Se trabajarán
usando los permisos UGO (User Group Other) en su forma octal.
Linux File Permission Cheatsheet
*/

type Folderblock struct { // bloque de carpetas
	B_content [4]Content
}

type Content struct { //content del bloque de carpetas
	B_name  [12]byte
	B_inodo int32
}

type Fileblock struct { // bloque de archivos
	B_content [64]byte
}

type Pointerblock struct { //bloque de apuntadores
	B_pointers [16]int32
}

type Content_J struct {
	Operation [10]byte
	Path      [100]byte
	Content   [100]byte
	Date      [17]byte
}

type Journaling struct {
	Size      int32
	Ultimo    int32
	Contenido [50]Content_J
}
