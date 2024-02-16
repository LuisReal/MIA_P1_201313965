package main

type Partition struct {
	Part_status bool // es de tipo bool(indica si la particion esta montada o no)

	Part_type [16]byte //(indica el tipo de particion: primaria(P) o extendida(E))

	Part_fit [16]byte // indica el tipo de ajuste(B mejor ajuste  F primer ajuste W peor ajuste)

	Part_start int32 // indica en que byte del disco inicia la particion

	Part_s int32 // contiene el tamano total de la particion en bytes

	Part_name [16]byte // contiene el nombre de la particion

	Part_correlative int32 // contiene el correlativo de la particion

	Part_id int32
}
