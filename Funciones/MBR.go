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

		fmt.Println(fmt.Sprintf(" Particion: %d, Tipo de Particion:  %s, size: %d, start: %d ", i, string(data.Mbr_partitions[i].Part_type[:]), int(data.Mbr_partitions[i].Part_size), int(data.Mbr_partitions[i].Part_start)))
	}

	/*
		for i := 0; i < 4; i++ {
			fmt.Println(fmt.Sprintf("Partition %d: %s, %s, %d, %d", i, string(data.Mbr_partitions[i].Name[:]), string(data.Mbr_partitions[i].Type[:]), data.Mbr_partitions[i].Start, data.Mbr_partitions[i].Size))
		}*/
}
