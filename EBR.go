package main

type EBR struct { //extended boot record

	Part_mount bool // indica si la particion esta montada o no

	Part_fit [16]byte // indica el tipo de ajuste de la particion(B mejor ajuste F primer ajuste W peor ajuste)

	Part_start int32 // indica en que byte del disco inicia la particion

	Part_s int32 //contiene el tamano total de la particion en bytes

	Part_next int32 // byte en el que esta el proximo EBR . -1 si no hay siguiente

	Part_name [16]byte // nombre de la particion
}
