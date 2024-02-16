package main

type MBR struct {
	Mbr_tamano int32

	Mbr_fecha_creacion [16]byte // de tipo time

	Mbr_dsk_signature int32

	Dsk_fit [16]byte // B (mejor ajuste)  F(primer ajuste) W(peor ajuste)

	Mbr_partitions [4]string // este arreglo simulara las 4 particiones

}
