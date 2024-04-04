package Funciones

import (
	//"encoding/binary"

	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Reportes(name string, path string, id string, ruta string) error {

	name = strings.ToLower(name)

	fmt.Println("\n\n========================= Inicio REPORTES ===========================")

	fmt.Printf("\nName: %s, Path: %s, Id: %s, Ruta: %s\n", name, path, id, ruta)

	driveletter := string(id[0])

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
				fmt.Println("\n****Particion Encontrada*****")

				index = i
			}
		}
	}

	if index != -1 {
		ImprimirParticion(TempMBR.Mbr_partitions[index])
		fmt.Println()
	} else {
		fmt.Println("\n*****Particion NO encontrada******")
		return nil
	}

	//crea las carpetas donde se guardara el archivo

	err2 := CrearArchivo(path)

	if err2 != nil {
		fmt.Println(err2)
		return err2
	}

	if name == "tree" {
		ReporteTree(index, path, file, TempMBR, driveletter)
		//ReporteTree(path, Inode0, file, tempSuperblock, driveletter)
	} else if name == "mbr" {
		ReporteMbr(path, file, driveletter)
	} else if name == "disk" {
		ReporteDisk(path, file, driveletter)
	} else if name == "bm_inode" {
		ReporteBitmap_inodos(index, path, file, TempMBR, driveletter)
	} else if name == "bm_block" {
		ReporteBitmap_bloques(index, path, file, TempMBR, driveletter)
	} else if name == "sb" {
		ReporteSuperbloque(index, path, file, TempMBR, driveletter)
	} else if name == "inode" {
		ReporteInode(index, path, file, TempMBR, driveletter)
	} else if name == "block" {
		ReporteBlock(index, path, file, TempMBR, driveletter)
	} else if name == "file" {
		ReporteFile(index, path, file, TempMBR, driveletter, ruta)
	} else if name == "journaling" {
		ReporteJournaling(index, path, file, TempMBR, driveletter)
	}

	//rep -id=191a -path="/home/serchiboi/archivos/reportes/reporte5_bm_inode.txt" -name=bm_inode

	fmt.Println("\n\n========================= Fin REPORTES ===========================")

	return nil

}

func ReporteDisk(path string, file *os.File, disco string) error {

	fmt.Println("\n\n========================= Iniciando Reporte DISK ===========================")

	var TempMBR MBR

	if err := LeerObjeto(file, &TempMBR, int64(0)); err != nil {
		return err
	}

	PrintMBR(TempMBR)

	var tamano_total float64

	tamano_disco := float64(TempMBR.Mbr_tamano) //tamano en bytes

	/*node [shape=record];
	  struct3 [label="MBR &#92;n 20%|
	  Libre
	  |{ Extendida |{EBR|Logica1|EBR |Logica2 |EBR}}|
	  Primaria | Libre

	  "];
	*/
	tamano_MBR := float64(binary.Size(TempMBR))

	tamano_total += tamano_MBR

	grafo := `digraph G {
		labelloc="t";
        label="Disco ` + disco + `";
        fontsize="50"
		node [shape=record];`

	//MBR &#92;n 20%|
	grafo += `struct1 [label="MBR|`

	var tamano_particion float64
	var tamano_partLibre float64
	var porcentaje float64

	var num_part_libres float64

	for i := 0; i < 4; i++ {

		if TempMBR.Mbr_partitions[i].Part_size == 0 {
			num_part_libres++

		}

	}

	for i := 0; i < 4; i++ {

		if TempMBR.Mbr_partitions[i].Part_size == 0 {

			if TempMBR.Mbr_partitions[i].Part_start != 0 {

				tamano_partLibre = float64(TempMBR.Mbr_partitions[i+1].Part_start) - float64(TempMBR.Mbr_partitions[i].Part_start)
				porcentaje = (100 * tamano_partLibre) / tamano_disco

				grafo += `Libre&#92;n` + strconv.FormatFloat(porcentaje, 'f', 2, 64) + `%|`

			} else {
				//fmt.Println("\nTamano total es: ", tamano_total)

				tamano_partLibre = tamano_disco - tamano_total

				tamano_partLibre = tamano_partLibre / num_part_libres

				porcentaje = (100 * tamano_partLibre) / tamano_disco

				if num_part_libres > 1 {
					grafo += `Libre&#92;n` + strconv.FormatFloat(porcentaje, 'f', 2, 64) + `%|`
				} else {
					grafo += `Libre&#92;n` + strconv.FormatFloat(porcentaje, 'f', 2, 64) + `%`
				}

			}

		} else {
			if string(TempMBR.Mbr_partitions[i].Part_type[:]) == "p" {

				tamano_particion = float64(TempMBR.Mbr_partitions[i].Part_size)

				tamano_total += tamano_particion

				porcentaje = (float64(100) * tamano_particion) / tamano_disco

				grafo += `Primaria&#92;n ` + strconv.FormatFloat(porcentaje, 'f', 2, 64) + `%|`

			} else if string(TempMBR.Mbr_partitions[i].Part_type[:]) == "e" {

				var tam_tot_logicas int32

				tamano_particion = float64(TempMBR.Mbr_partitions[i].Part_size)
				tamano_total += tamano_particion

				porcentaje = (float64(100) * tamano_particion) / tamano_disco

				inicio := TempMBR.Mbr_partitions[i].Part_start

				grafo += `{ Extendida  ` + strconv.FormatFloat(porcentaje, 'f', 2, 64) + `%|`

				var tempEBR EBR

				if err := LeerObjeto(file, &tempEBR, int64(inicio)); err != nil { // obtiene el primer ebr
					return err
				}

				tamano_ebr := float64(binary.Size(tempEBR))

				tam_tot_logicas += int32(tamano_ebr)

				grafo += `{`

				grafo += `EBR|`

				if tempEBR.Part_size != 0 {

					var part_name string // elimina los espacios en el slice para que pueda ser leido por graphviz
					for j := 0; j < len(tempEBR.Part_name); j++ {
						if tempEBR.Part_name[j] != 0 {
							part_name += string(tempEBR.Part_name[j])
						}

					}

					tamano_particion = float64(tempEBR.Part_size) // tamano particion logica

					tam_tot_logicas += int32(tamano_particion)

					porcentaje = (100 * tamano_particion) / tamano_disco

					grafo += part_name + ` &#92;n` + strconv.FormatFloat(porcentaje, 'f', 2, 64) + `%|`
				} else {

					tamano_libre_ebr := float64(tempEBR.Part_next) - float64(tempEBR.Part_start)

					tam_tot_logicas += int32(tamano_libre_ebr)

					porcentaje = (100 * tamano_libre_ebr) / tamano_disco

					grafo += `Libre&#92;n` + strconv.FormatFloat(porcentaje, 'f', 2, 64) + `%|`
				}

				part_next := tempEBR.Part_next

				for part_next != 0 { // obtiene los siguientes EBR analizando si existen por medio de su tamano

					grafo += `EBR|`

					if err := LeerObjeto(file, &tempEBR, int64(part_next)); err != nil { // obtiene el primer ebr
						return err
					}

					tam_tot_logicas += int32(binary.Size(EBR{}))

					if tempEBR.Part_size != 0 {

						var part_name string // elimina los espacios en el slice para que pueda ser leido por graphviz
						for j := 0; j < len(tempEBR.Part_name); j++ {
							if tempEBR.Part_name[j] != 0 {
								part_name += string(tempEBR.Part_name[j])
							}

						}

						tam_tot_logicas += tempEBR.Part_size

						tamano_particion = float64(tempEBR.Part_size) // tamano particion logica

						porcentaje = (100 * tamano_particion) / tamano_disco

						grafo += part_name + ` &#92;n` + strconv.FormatFloat(porcentaje, 'f', 2, 64) + `%|`

						//grafo += `Logica1 | EBR | Logica2 | EBR | Logica3}|`

					} else {

						if part_next != 0 {

							if tempEBR.Part_next != 0 {

								tamano_libre_ebr := float64(tempEBR.Part_next) - float64(tempEBR.Part_start)

								fmt.Println("\n El tamano_libre_ebr es: ", tamano_libre_ebr)

								tam_tot_logicas += int32(tamano_libre_ebr)

								porcentaje = (100 * tamano_libre_ebr) / tamano_disco

								grafo += `Libre&#92;n` + strconv.FormatFloat(porcentaje, 'f', 2, 64) + `%|`
							} else {

								tamano_libre_ebr := float64(TempMBR.Mbr_partitions[i].Part_size - tam_tot_logicas)

								fmt.Println("\n El tamano_libre_ebr es: ", tamano_libre_ebr)
								porcentaje = (100 * tamano_libre_ebr) / tamano_disco

								grafo += `Libre&#92;n` + strconv.FormatFloat(porcentaje, 'f', 2, 64) + `%`
							}

						}

					}

					part_next = tempEBR.Part_next

				}

				grafo += `}}|`

				//grafo += `{ Extendida | {EBR | Logica1 | EBR | Logica2 | EBR | Logica3}}|`
			}
		}

	}

	tamano_libre2 := tamano_disco - tamano_total

	if num_part_libres == float64(0) {

		if tamano_libre2 > 0 {

			porcentaje = (100 * tamano_libre2) / tamano_disco

			grafo += `Libre&#92;n` + strconv.FormatFloat(porcentaje, 'f', 2, 64) + `%`
		}
	}

	grafo += `"];`
	grafo += `}`

	dot := "disk.dot"

	file_dot, err := os.Create(dot)

	if err != nil {
		fmt.Println(err)
		return err
	}

	file_dot.WriteString(grafo)

	file_dot.Close()

	result := path

	//exec.Command("dot", "-Tpng", dot, "-o", result)
	out, err := exec.Command("dot", "-Tpdf", dot, "-o", result).Output()

	if err != nil {

		log.Fatal(err)
	}

	fmt.Println(string(out))

	fmt.Println("\n\n========================= Finalizando Reporte DISK ===========================")

	return nil

}

func ReporteMbr(path string, file *os.File, disco string) error {

	fmt.Println("\n\n========================= Iniciando Reporte MBR ===========================")
	fmt.Printf("\npath: %s", path)
	fmt.Println()

	var TempMBR MBR

	if err := LeerObjeto(file, &TempMBR, int64(0)); err != nil {
		return err
	}

	PrintMBR(TempMBR)

	grafo := `digraph H {
			labelloc="t";
			label="Disco ` + disco + `";
			fontsize="50"
			graph [pad="0.5", nodesep="0.5", ranksep="1"];
			node [shape=plaintext]
			rankdir=LR;`

	var contador int

	grafo += `label=<
				<table  border="0" cellborder="1" cellspacing="0">`
	contador++

	grafo += `<tr><td colspan="3" style="filled" bgcolor="#FFD700"  port='` + strconv.Itoa(contador) + `'>Reporte MBR</td></tr>`

	contador++

	grafo += `<tr><td>mbr_tamano</td><td port='` + strconv.Itoa(contador) + `'>` + strconv.Itoa(int(TempMBR.Mbr_tamano)) + `</td></tr>`

	contador++

	grafo += `<tr><td>mbr_fecha_creacion</td><td port='` + strconv.Itoa(contador) + `'>` + string(TempMBR.Mbr_fecha_creacion[:]) + `</td></tr>`

	contador++

	grafo += `<tr><td>mbr_disk_signature </td><td port='` + strconv.Itoa(contador) + `'>` + strconv.Itoa(int(TempMBR.Mbr_dsk_signature)) + `</td></tr>`

	//grafo += `<tr><td colspan="3" port='` + strconv.Itoa(contador) + `'>Particion</td></tr>`

	for i := 0; i < 4; i++ {

		if int(TempMBR.Mbr_partitions[i].Part_size) != 0 {

			contador++

			grafo += `

				<tr><td colspan="3" align="left" style="filled" bgcolor="lightblue" port='` + strconv.Itoa(contador) + `'>Particion</td></tr>`

			contador++

			grafo += `<tr><td>status</td><td port='` + strconv.Itoa(contador) + `'>` + strconv.FormatBool(TempMBR.Mbr_partitions[i].Part_status) + `</td></tr>`

			contador++

			grafo += `<tr><td>type</td><td port='` + strconv.Itoa(contador) + `'>` + string(TempMBR.Mbr_partitions[i].Part_type[:]) + `</td></tr>`

			//verifica si el slice de part_fit tiene elementos
			var elemento bool
			for _, v := range TempMBR.Mbr_partitions[i].Part_fit[:] {

				if v != 0 {
					elemento = true
				}
			}

			if elemento {

				contador++
				grafo += `<tr><td>fit</td><td port='` + strconv.Itoa(contador) + `'>` + string(TempMBR.Mbr_partitions[i].Part_fit[:]) + `</td></tr>`

			} else {

				contador++
				grafo += `<tr><td>fit</td><td port='` + strconv.Itoa(contador) + `'>` + "" + `</td></tr>`
			}

			contador++

			grafo += `<tr><td>start</td><td port='` + strconv.Itoa(contador) + `'>` + strconv.Itoa(int(TempMBR.Mbr_partitions[i].Part_start)) + `</td></tr>`

			contador++

			grafo += `<tr><td>size</td><td port='` + strconv.Itoa(contador) + `'>` + strconv.Itoa(int(TempMBR.Mbr_partitions[i].Part_size)) + `</td></tr>`

			var part_name string // elimina los espacios en el slice para que pueda ser leido por graphviz
			for j := 0; j < len(TempMBR.Mbr_partitions[i].Part_name); j++ {
				if TempMBR.Mbr_partitions[i].Part_name[j] != 0 {
					part_name += string(TempMBR.Mbr_partitions[i].Part_name[j])
				}

			}

			contador++
			grafo += `

				<tr><td>name</td><td port='` + strconv.Itoa(contador) + `'>` + part_name + `</td></tr>`

		}

	}

	//grafo += `</table> >`

	/*grafo += `label=<
	<table  border="0" cellborder="1" cellspacing="0">`*/

	fmt.Println("\nCreando REPORTE EBR")
	var cont int
	for j := 0; j < 4; j++ {

		if string(TempMBR.Mbr_partitions[j].Part_type[:]) == "e" {

			cont++

			grafo += `

				<tr><td colspan="3" style="filled" bgcolor="#FFD700" port='` + strconv.Itoa(cont) + `'> Reporte EBR</td></tr>`

			inicio := TempMBR.Mbr_partitions[j].Part_start

			var tempEBR EBR

			if err := LeerObjeto(file, &tempEBR, int64(inicio)); err != nil { // obtiene el primer ebr
				return err
			}

			cont++

			var part_name string // elimina los espacios en el slice para que pueda ser leido por graphviz
			for j := 0; j < len(tempEBR.Part_name); j++ {

				if tempEBR.Part_name[j] != 0 {
					part_name += string(tempEBR.Part_name[j])
				}

			}

			grafo += `<tr><td colspan="3" style="filled" bgcolor="lightblue" port='` + strconv.Itoa(cont) + `'>` + part_name + `</td></tr>`

			cont++
			grafo += `<tr><td>Status</td><td port='` + strconv.Itoa(cont) + `'>` + strconv.FormatBool(tempEBR.Part_mount) + `</td></tr>`

			var part_fit string // elimina los espacios en el slice para que pueda ser leido por graphviz
			for j := 0; j < len(tempEBR.Part_fit); j++ {

				if tempEBR.Part_fit[j] != 0 {
					part_name += string(tempEBR.Part_fit[j])
				}

			}

			cont++
			grafo += `<tr><td>Fit</td><td port='` + strconv.Itoa(cont) + `'>` + part_fit + `</td></tr>`

			cont++
			grafo += `<tr><td>Size</td><td port='` + strconv.Itoa(cont) + `'>` + strconv.Itoa(int(tempEBR.Part_size)) + `</td></tr>`

			cont++
			grafo += `<tr><td>Next</td><td port='` + strconv.Itoa(cont) + `'>` + strconv.Itoa(int(tempEBR.Part_next)) + `</td></tr>`

			cont++
			grafo += `<tr><td>Start</td><td port='` + strconv.Itoa(cont) + `'>` + strconv.Itoa(int(tempEBR.Part_start)) + `</td></tr>`

			part_next := tempEBR.Part_next

			for part_next != 0 { // obtiene los siguientes EBR analizando si existen por medio de su tamano

				if err := LeerObjeto(file, &tempEBR, int64(part_next)); err != nil { // obtiene el primer ebr
					return err
				}

				if tempEBR.Part_size != 0 {
					var part_name string // elimina los espacios en el slice para que pueda ser leido por graphviz
					for j := 0; j < len(tempEBR.Part_name); j++ {

						if tempEBR.Part_name[j] != 0 {
							part_name += string(tempEBR.Part_name[j])
						}

					}

					cont++

					grafo += `<tr><td colspan="3" style="filled" bgcolor="lightblue" port='` + strconv.Itoa(cont) + `'>` + part_name + `</td></tr>`

					cont++
					grafo += `<tr><td>Status</td><td port='` + strconv.Itoa(cont) + `'>` + strconv.FormatBool(tempEBR.Part_mount) + `</td></tr>`

					var part_fit string // elimina los espacios en el slice para que pueda ser leido por graphviz
					for j := 0; j < len(tempEBR.Part_fit); j++ {

						if tempEBR.Part_fit[j] != 0 {
							part_name += string(tempEBR.Part_fit[j])
						}

					}

					cont++
					grafo += `<tr><td>Fit</td><td port='` + strconv.Itoa(cont) + `'>` + part_fit + `</td></tr>`

					cont++
					grafo += `<tr><td>Size</td><td port='` + strconv.Itoa(cont) + `'>` + strconv.Itoa(int(tempEBR.Part_size)) + `</td></tr>`

					cont++
					grafo += `<tr><td>Next</td><td port='` + strconv.Itoa(cont) + `'>` + strconv.Itoa(int(tempEBR.Part_next)) + `</td></tr>`

					cont++
					grafo += `<tr><td>Start</td><td port='` + strconv.Itoa(cont) + `'>` + strconv.Itoa(int(tempEBR.Part_start)) + `</td></tr>`

				}

				part_next = tempEBR.Part_next

			}

		}

	}

	grafo += `</table>
				>`

	grafo += `}`

	//fmt.Println("\nImprimiendo grafo: ", grafo)
	//fmt.Println("\nImprimiendo grafo: ", grafo)
	dot := "mbr.dot"

	file_dot, err := os.Create(dot)

	if err != nil {
		fmt.Println(err)
		return err
	}

	file_dot.WriteString(grafo)

	file_dot.Close()

	//-path=/home/darkun/Escritorio/mbr.pdf
	//dot -Tpdf mbr.dot  -o mbr.pdf

	result := path

	//exec.Command("dot", "-Tpng", dot, "-o", result)
	out, err := exec.Command("dot", "-Tpdf", dot, "-o", result).Output()

	if err != nil {
		fmt.Println("\n Imprimiendo el error en ReportesMBR y EBR")
		log.Fatal(err)
	}

	fmt.Println(string(out))

	fmt.Println("\n\n========================= Finalizando Reporte MBR ===========================")

	return nil
}

func ReporteTree(index int, path string, file *os.File, TempMBR MBR, disco string) error {

	fmt.Println("\n\n========================= Iniciando Reporte Tree ===========================")

	var tempSuperblock Superblock

	if err := LeerObjeto(file, &tempSuperblock, int64(TempMBR.Mbr_partitions[index].Part_start)); err != nil {
		return err
	}

	Inodo_start := tempSuperblock.S_inode_start

	var inodo Inode

	grafo := `digraph H {
        label="Disco ` + disco + `";
        fontsize="50"
		graph [pad="0.5", nodesep="0.5", ranksep="1"];
		node [shape=plaintext]
		 rankdir=LR;`

	//fmt.Println("\n EL numero de inodos es: ", tempSuperblock.S_inodes_count)
	//fmt.Println("\n EL numero de bloques es: ", tempSuperblock.S_blocks_count)

	inodos_ocupados := tempSuperblock.S_inodes_count - tempSuperblock.S_free_inodes_count

	var inodo_vacio Inode

	for i := 0; i < int(inodos_ocupados); i++ {
		//fmt.Println("\nEstoy dentro del for de inodos")
		if err := LeerObjeto(file, &inodo, int64(Inodo_start+int32(i)*int32(binary.Size(Inode{})))); err != nil {
			return err
		}

		/*for i := int32(0); i < 15; i++ {
			Inode0.I_block[i] = -1
		}

		Inode0.I_block[0] = 0
		*/

		if inodo == inodo_vacio {
			continue
		}
		grafo += `Inodo` + strconv.Itoa(i) + ` [
			label=<
				<table  border="0" cellborder="1" cellspacing="0">
				<tr><td colspan="3" port='0'>Inodo` + strconv.Itoa(i) + `</td></tr>`

		grafo += `<tr><td>I_Uid</td><td port='0'>` + strconv.Itoa(int(inodo.I_uid)) + `</td></tr>`
		grafo += `<tr><td>I_Gid</td><td port='0'>` + strconv.Itoa(int(inodo.I_gid)) + `</td></tr>`
		grafo += `<tr><td>I_perm</td><td port='0'>` + string(inodo.I_perm[:]) + `</td></tr>`
		grafo += `<tr><td>I_type</td><td port='0'>` + string(inodo.I_type[:]) + `</td></tr>`

		for j := 0; j < 15; j++ {
			grafo += `<tr><td>AD` + strconv.Itoa(j+1) + `</td><td port='` + strconv.Itoa(j+1) + `'>` + strconv.Itoa(int(inodo.I_block[j])) + `</td></tr>`
		}

		grafo += `</table>
			>];
			
			`

		var bloque int32

		//var index int = 0

		for k := 0; k < 15; k++ {

			if inodo.I_block[k] != -1 {

				bloque = inodo.I_block[k] // esto contiene el numero de bloque

				//fmt.Println("\nEl bloque es: ", bloque)
				//fmt.Println("\nEl valor de k es: ", k)

				if k < 12 {

					//fmt.Println("\nEl bloque a crear es: ", bloque)

					// carpeta -> 0   archivo -> 1
					if string(inodo.I_type[:]) == "0" { // es un bloque de carpetas
						var folder Folderblock

						//fmt.Println("\nEstoy dentro del for de inodos")
						if err := LeerObjeto(file, &folder, int64(tempSuperblock.S_block_start+bloque*int32(binary.Size(Folderblock{})))); err != nil {
							return err
						}

						grafo += `Bloque` + strconv.Itoa(int(bloque)) + ` [
						label=<
						<table  border="0" cellborder="1" cellspacing="0">
						<tr><td colspan="3" port='0'>Bloque` + strconv.Itoa(int(bloque)) + `</td></tr>`

						var indice int

						for j := 0; j < 4; j++ { // si es un bloque de carpetas
							//fmt.Println("\nEl name de folder es: ", string(folder.B_content[j].B_name[:]))
							for k := 0; k < len(folder.B_content[j].B_name[:]); k++ {
								if folder.B_content[j].B_name[k] == 0 { //quitando espacios(los ceros restantes) al slice de B_name
									indice = k
									break
								}

							}
							//fmt.Println("\nEl indice es: ", indice)
							name := string(folder.B_content[j].B_name[:indice])
							inodo := folder.B_content[j].B_inodo

							//fmt.Println("\nInodo es: ", inodo)
							grafo += `<tr><td>` + name + `</td><td port='` + strconv.Itoa(j+1) + `'>` + strconv.Itoa(int(inodo)) + `</td></tr>`
						}

						grafo += `</table>
						>];	`

					} else if string(inodo.I_type[:]) == "1" { //es un bloque de archivos
						var file_block Fileblock

						//fmt.Println("\nEstoy dentro del for de inodos")
						if err := LeerObjeto(file, &file_block, int64(tempSuperblock.S_block_start+int32(bloque)*int32(binary.Size(Fileblock{})))); err != nil {
							return err
						}

						grafo += `Bloque` + strconv.Itoa(int(bloque)) + ` [
						label=<
						<table  border="0" cellborder="1" cellspacing="0">
						<tr><td colspan="3" port='0'>Bloque` + strconv.Itoa(int(bloque)) + `</td></tr>`

						var indice int

						for k := 0; k < len(file_block.B_content[:]); k++ {
							if file_block.B_content[k] == 0 { //quitando espacios(los ceros restantes) al slice de B_content
								indice = k
								break
							}

						}
						//contenido := string(folder.B_content[j].B_name[:indice])
						var contenido string
						if indice == 0 { // significa que el slice fileblock.B_content esta lleno
							contenido = string(file_block.B_content[:])
						} else { //el slice todavia tiene espacios vacios
							contenido = string(file_block.B_content[:indice])
						}

						//fmt.Println("\nEl contenido de fileblock es: ", contenido)

						grafo += `<tr><td port='` + strconv.Itoa(int(bloque)+1) + `'>` + contenido + `</td></tr>`

						grafo += `</table>
						>];`

					}

				} else { // bloques indirectos

					if k == 12 {

						var newPointerBlock Pointerblock

						if err := LeerObjeto(file, &newPointerBlock, int64(tempSuperblock.S_block_start+int32(bloque)*int32(binary.Size(Pointerblock{})))); err != nil {
							return err
						}

						//execute -path=/home/darkun/Escritorio/prueba.mia

						grafo += `Bloque` + strconv.Itoa(int(bloque)) + ` [
						label=<
						<table  border="0" cellborder="1" cellspacing="0">
						<tr><td colspan="3" port='0'>Bloque` + strconv.Itoa(int(bloque)) + `</td></tr>` + "\n"

						for i := 0; i < len(newPointerBlock.B_pointers); i++ {

							grafo += `<tr><td port='` + strconv.Itoa(int(i+1)) + `'>` + strconv.Itoa(int(newPointerBlock.B_pointers[i])) + `</td></tr>` + "\n"
						}

						grafo += `</table>
						>];` + "\n"

						for j, bloque := range newPointerBlock.B_pointers {

							if bloque != -1 {

								var file_block Fileblock

								//fmt.Println("\nEstoy dentro del for de inodos")
								if err := LeerObjeto(file, &file_block, int64(tempSuperblock.S_block_start+int32(bloque)*int32(binary.Size(Fileblock{})))); err != nil {
									return err
								}

								grafo += `Bloque` + strconv.Itoa(int(bloque)) + ` [
								label=<
								<table  border="0" cellborder="1" cellspacing="0">
								<tr><td colspan="3" port='0'>Bloque` + strconv.Itoa(int(bloque)) + `</td></tr>` + "\n"

								var indice int

								for k := 0; k < len(file_block.B_content[:]); k++ {
									if file_block.B_content[k] == 0 { //quitando espacios(los ceros restantes) al slice de B_content
										indice = k
										break
									}

								}

								//contenido := string(folder.B_content[j].B_name[:indice])
								var contenido string
								if indice == 0 { // significa que el slice fileblock.B_content esta lleno
									contenido = string(file_block.B_content[:])
								} else { //el slice todavia tiene espacios vacios
									contenido = string(file_block.B_content[:indice])
								}

								//fmt.Println("\nEl contenido de fileblock es: ", contenido)

								grafo += `<tr><td port='` + strconv.Itoa(j+1) + `'>` + contenido + `</td></tr>`

								grafo += `</table>
							>];` + "\n"
							}

						}
					} else if k == 13 {

						var newPointerBlock1 Pointerblock

						if err := LeerObjeto(file, &newPointerBlock1, int64(tempSuperblock.S_block_start+int32(bloque)*int32(binary.Size(Pointerblock{})))); err != nil {
							return err
						}

						//execute -path=/home/darkun/Escritorio/prueba.mia

						grafo += `Bloque` + strconv.Itoa(int(bloque)) + ` [
						label=<
						<table  border="0" cellborder="1" cellspacing="0">
						<tr><td colspan="3" port='0'>Bloque` + strconv.Itoa(int(bloque)) + `</td></tr>` + "\n"

						for i := 0; i < len(newPointerBlock1.B_pointers); i++ {

							grafo += `<tr><td port='` + strconv.Itoa(int(i+1)) + `'>` + strconv.Itoa(int(newPointerBlock1.B_pointers[i])) + `</td></tr>` + "\n"
						}

						grafo += `</table>
						>];` + "\n"

						for _, bloque := range newPointerBlock1.B_pointers {

							if bloque != -1 {

								var newPointerBlock2 Pointerblock

								//fmt.Println("\nEstoy dentro del for de inodos")
								if err := LeerObjeto(file, &newPointerBlock2, int64(tempSuperblock.S_block_start+int32(bloque)*int32(binary.Size(Pointerblock{})))); err != nil {
									return err
								}

								grafo += `Bloque` + strconv.Itoa(int(bloque)) + ` [
								label=<
								<table  border="0" cellborder="1" cellspacing="0">
								<tr><td colspan="3" port='0'>Bloque` + strconv.Itoa(int(bloque)) + `</td></tr>` + "\n"

								for i := 0; i < len(newPointerBlock2.B_pointers); i++ {

									grafo += `<tr><td port='` + strconv.Itoa(int(i+1)) + `'>` + strconv.Itoa(int(newPointerBlock2.B_pointers[i])) + `</td></tr>` + "\n"
								}

								grafo += `</table>
									>];` + "\n"

								for m, bloque := range newPointerBlock2.B_pointers {

									if bloque != -1 {

										var file_block Fileblock

										//fmt.Println("\nEstoy dentro del for de inodos")
										if err := LeerObjeto(file, &file_block, int64(tempSuperblock.S_block_start+int32(bloque)*int32(binary.Size(Fileblock{})))); err != nil {
											return err
										}

										grafo += `Bloque` + strconv.Itoa(int(bloque)) + ` [
										label=<
										<table  border="0" cellborder="1" cellspacing="0">
										<tr><td colspan="3" port='0'>Bloque` + strconv.Itoa(int(bloque)) + `</td></tr>` + "\n"

										var indice int

										for k := 0; k < len(file_block.B_content[:]); k++ {
											if file_block.B_content[k] == 0 { //quitando espacios(los ceros restantes) al slice de B_content
												indice = k
												break
											}

										}
										//contenido := string(folder.B_content[j].B_name[:indice])
										var contenido string
										if indice == 0 { // significa que el slice fileblock.B_content esta lleno
											contenido = string(file_block.B_content[:])
										} else { //el slice todavia tiene espacios vacios
											contenido = string(file_block.B_content[:indice])
										}

										//fmt.Println("\nEl contenido de fileblock es: ", contenido)

										grafo += `<tr><td port='` + strconv.Itoa(m+1) + `'>` + contenido + `</td></tr>`

										grafo += `</table>
									>];` + "\n"
									}

								}

							}

						}
					}

				}

			}
		}

	}

	var Inode0 Inode

	if err := LeerObjeto(file, &Inode0, int64(tempSuperblock.S_inode_start)); err != nil {
		return err
	}

	grafo, _ = EnlazandoNodos(path, Inode0, file, tempSuperblock, disco, grafo, 0) //enlazando los nodos

	grafo += `}`

	//fmt.Println("\nImprimiendo grafo: ", grafo)
	dot := "tree.dot"

	fmt.Println("\n Creando archivo tree.dot")

	file_dot, err := os.Create(dot)

	if err != nil {
		fmt.Println(err)
		return err
	}

	file_dot.WriteString(grafo)

	file_dot.Close()

	result := path

	//exec.Command("dot", "-Tpng", dot, "-o", result)
	out, err := exec.Command("dot", "-Tsvg", dot, "-o", result).Output()

	if err != nil {

		log.Fatal(err)
	}

	fmt.Println(string(out))

	fmt.Println("\n\n========================= Finalizando Reporte Tree ===========================")

	return nil
}

func EnlazandoNodos(path string, Inodo Inode, file *os.File, tempSuperblock Superblock, disco string, grafo string, Inodo_actual int32) (string, error) {

	//fmt.Println("\n\n========================= Iniciando EnlazandoNodos ===========================")

	// Iterate over i_blocks from Inode
	for k, block := range Inodo.I_block {
		if block != -1 {
			//fmt.Println("\nEl bloque al que apunta el inodo es: ", block)
			if k < 12 {
				//CASO DIRECTO

				//// carpeta -> 0   archivo ->1
				//fmt.Println("\nEl indice de block es: ", k)
				if string(Inodo.I_type[:]) == "0" { // si es carpeta

					//fmt.Println("\nInodo de tipo carpeta")

					var crrFolderBlock Folderblock

					// Read object from bin file                                       // un Inodo y un bloque las estructuras miden lo mismo
					if err := LeerObjeto(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))); err != nil {
						return "", err
					}

					//fmt.Println("\nEl valor de Inodo_actual es: ", Inodo_actual)
					if Inodo_actual == 0 {

						grafo += `Inodo` + strconv.Itoa(int(Inodo_actual)) + `:` + strconv.Itoa(k+1) + `:e->`
					} else {
						//fmt.Println("\nEl valor de k es: ", k)
						grafo += `Inodo` + strconv.Itoa(int(Inodo_actual)) + `:` + strconv.Itoa(k+1) + `:e->`
					}

					grafo += `Bloque` + strconv.Itoa(int(block)) + `:0;` + "\n"

					for j, folder := range crrFolderBlock.B_content {

						//execute -path=/home/darkun/Escritorio/prueba.mia

						var actual [12]byte
						copy(actual[:], ".")

						var padre [12]byte
						copy(padre[:], "..")

						if folder.B_inodo > 0 && (string(folder.B_name[:]) != string(actual[:])) { // apunta a un inodo (no puede ser 0 y tampoco -1)

							grafo += `Bloque` + strconv.Itoa(int(block)) + `:` + strconv.Itoa(j+1) + `:e->`

							grafo += `Inodo` + strconv.Itoa(int(folder.B_inodo)) + `:0;` + "\n"

							var NextInode Inode
							// Read object from bin file
							if err := LeerObjeto(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(Inode{})))); err != nil {
								return "", err
							}

							grafo, _ = EnlazandoNodos(path, NextInode, file, tempSuperblock, disco, grafo, folder.B_inodo)

						}

					}
				} else if string(Inodo.I_type[:]) == "1" { // si es archivo

					//fmt.Println("\nInodo de tipo archivo")

					//grafo += `Inodo` + strconv.Itoa(int(folder.B_inodo)) + `:` //1->
					//fmt.Println("\nEl valor de k es: ", k)

					//m := k+4

					grafo += `Inodo` + strconv.Itoa(int(Inodo_actual)) + `:` + strconv.Itoa(int(k+1)) + `:e->` //1->

					grafo += `Bloque` + strconv.Itoa(int(block)) + `:0;` + "\n"

					//execute -path=/home/darkun/Escritorio/prueba.mia

				}

			} else {

				if k == 12 { // apuntador simple

					//fmt.Println("CASO INDIRECTO SIMPLE")

					grafo += `Inodo` + strconv.Itoa(int(Inodo_actual)) + `:` + strconv.Itoa(int(k+1)) + `:e->` //1->

					grafo += `Bloque` + strconv.Itoa(int(block)) + `:0;` + "\n"

					var newPointerBlock Pointerblock

					//fmt.Println("\nEl numero de bloque es: ", block)
					//fmt.Println("\nEstoy dentro del for de inodos")
					if err := LeerObjeto(file, &newPointerBlock, int64(tempSuperblock.S_block_start+int32(block)*int32(binary.Size(Pointerblock{})))); err != nil {
						return "", err
					}

					for p, bloque := range newPointerBlock.B_pointers {

						if bloque != -1 {

							grafo += `Bloque` + strconv.Itoa(int(block)) + `:` + strconv.Itoa(int(p+1)) + `:e->` //1->

							grafo += `Bloque` + strconv.Itoa(int(bloque)) + `:0;` + "\n"
						}
					}

				} else if k == 13 {

					//fmt.Println("CASO INDIRECTO DOBLE")

					grafo += `Inodo` + strconv.Itoa(int(Inodo_actual)) + `:` + strconv.Itoa(int(k+1)) + `:e->` //1->

					grafo += `Bloque` + strconv.Itoa(int(block)) + `:0;` + "\n"

					var newPointerBlock1 Pointerblock

					fmt.Println("\nEl numero de bloque es: ", block)
					//fmt.Println("\nEstoy dentro del for de inodos")
					if err := LeerObjeto(file, &newPointerBlock1, int64(tempSuperblock.S_block_start+int32(block)*int32(binary.Size(Pointerblock{})))); err != nil {
						return "", err
					}

					for p, bloque1 := range newPointerBlock1.B_pointers {

						if bloque1 != -1 {

							grafo += `Bloque` + strconv.Itoa(int(block)) + `:` + strconv.Itoa(int(p+1)) + `:e->` //1->

							grafo += `Bloque` + strconv.Itoa(int(bloque1)) + `:0;` + "\n"

							var newPointerBlock2 Pointerblock

							//fmt.Println("\nEl numero de bloque es: ", block)
							//fmt.Println("\nEstoy dentro del for de inodos")
							if err := LeerObjeto(file, &newPointerBlock2, int64(tempSuperblock.S_block_start+int32(bloque1)*int32(binary.Size(Pointerblock{})))); err != nil {
								return "", err
							}

							for n, bloque2 := range newPointerBlock2.B_pointers {

								if bloque2 != -1 {

									grafo += `Bloque` + strconv.Itoa(int(bloque1)) + `:` + strconv.Itoa(int(n+1)) + `:e->` //1->

									grafo += `Bloque` + strconv.Itoa(int(bloque2)) + `:0;` + "\n"

								}

							}

						}
					}

				}

			}
		}

	}

	//fmt.Println("\n\n========================= Finalizando EnlazandoNodos ===========================")

	return grafo, nil
}

func ReporteBitmap_inodos(index int, path string, file *os.File, tempMBR MBR, driveletter string) error {
	//index, path, file, TempMBR, driveletter

	fmt.Println("\n\n========================= Iniciando Reporte Bitmap_inodos ===========================")
	//fmt.Printf("\nIndex: %d, path: %s, driveletter: %s", index, path, driveletter)

	var tempSuperblock Superblock

	if err := LeerObjeto(file, &tempSuperblock, int64(tempMBR.Mbr_partitions[index].Part_start)); err != nil {
		return err
	}

	//crea el archivo donde se mostrara el reporte de bitmap inodos

	bitmap_start := tempSuperblock.S_bm_inode_start

	bitmap_end := tempSuperblock.S_bm_block_start

	fmt.Println("\nTotal de binarios: ", bitmap_end-bitmap_start)

	archivo, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return err
	}

	file.Seek(int64(bitmap_start), 0)

	var buf [1]byte

	var contador int = 1

	for i := bitmap_start; i < bitmap_end; i++ {

		err := binary.Read(file, binary.LittleEndian, &buf)
		if err != nil {
			return err
		}
		//fmt.Print(buf)

		s := fmt.Sprintf("%b", buf) // convirtiendo binario a texto

		if contador <= 20 {

			archivo.WriteString(s)
			contador++

		} else {
			contador = 2
			archivo.WriteString("\n")
			archivo.WriteString(s)
		}
	}

	//execute -path=/home/darkun/Escritorio/prueba.mia

	archivo.Close()
	file.Close()
	//file.Close() //execute -path=/home/darkun/Escritorio/prueba.mia

	fmt.Println("\n\n========================= Finalizando Reporte Bitmap_inodos ===========================")

	return nil

}

func ReporteBitmap_bloques(index int, path string, file *os.File, tempMBR MBR, driveletter string) error {
	//index, path, file, TempMBR, driveletter

	fmt.Println("\n\n========================= Iniciando Reporte Bitmap_bloques ===========================")

	var tempSuperblock Superblock

	if err := LeerObjeto(file, &tempSuperblock, int64(tempMBR.Mbr_partitions[index].Part_start)); err != nil {
		return err
	}

	//crea el archivo donde se mostrara el reporte de bitmap inodos

	bitmap_start := tempSuperblock.S_bm_block_start

	bitmap_end := tempSuperblock.S_inode_start

	fmt.Println("\nTotal de binarios: ", bitmap_end-bitmap_start)

	archivo, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return err
	}

	file.Seek(int64(bitmap_start), 0)

	var buf [1]byte

	var contador int = 1

	for i := bitmap_start; i < bitmap_end; i++ {

		err := binary.Read(file, binary.LittleEndian, &buf)
		if err != nil {
			return err
		}
		//fmt.Print(buf)

		s := fmt.Sprintf("%b", buf) // convirtiendo binario a texto

		if contador <= 20 {
			archivo.WriteString(s)
			contador++

		} else {
			contador = 2
			archivo.WriteString("\n")
			archivo.WriteString(s)
		}

	}

	//execute -path=/home/darkun/Escritorio/prueba.mia

	archivo.Close()
	file.Close()
	//file.Close() //execute -path=/home/darkun/Escritorio/prueba.mia

	fmt.Println("\n\n========================= Finalizando Reporte Bitmap_bloques ===========================")

	return nil

}

func ReporteSuperbloque(index int, path string, file *os.File, tempMBR MBR, driveletter string) error {

	fmt.Println("\n\n========================= Iniciando Reporte Superbloque ===========================")

	superblock_start := tempMBR.Mbr_partitions[index].Part_start

	var superbloque Superblock

	if err := LeerObjeto(file, &superbloque, int64(superblock_start)); err != nil { // obtiene el primer ebr
		return err
	}

	grafo := `digraph H {
			labelloc="t";
			label="Disco ` + strings.ToUpper(driveletter) + `";
			fontsize="50"
			graph [pad="0.5", nodesep="0.5", ranksep="1"];
			node [shape=plaintext]
			rankdir=LR;`

	var contador int

	grafo += `label=<
				<table  border="0" cellborder="1" cellspacing="0">`
	contador++

	grafo += `<tr><td colspan="3" style="filled" bgcolor="#FFD700"  port='` + strconv.Itoa(contador) + `'>Reporte Superbloque</td></tr>`

	contador++

	grafo += `<tr><td>File System Type</td><td port='` + strconv.Itoa(contador) + `'>` + strconv.Itoa(int(superbloque.S_filesystem_type)) + `</td></tr>`

	contador++

	grafo += `<tr><td>Inodes count</td><td port='` + strconv.Itoa(contador) + `'>` + strconv.Itoa(int(superbloque.S_inodes_count)) + `</td></tr>`

	contador++

	grafo += `<tr><td>Blocks count </td><td port='` + strconv.Itoa(contador) + `'>` + strconv.Itoa(int(superbloque.S_blocks_count)) + `</td></tr>`

	contador++

	grafo += `<tr><td>Free Blocks count </td><td port='` + strconv.Itoa(contador) + `'>` + strconv.Itoa(int(superbloque.S_free_blocks_count)) + `</td></tr>`

	contador++

	grafo += `<tr><td>Free Inodes count </td><td port='` + strconv.Itoa(contador) + `'>` + strconv.Itoa(int(superbloque.S_free_inodes_count)) + `</td></tr>`

	contador++

	grafo += `<tr><td>Mnt count </td><td port='` + strconv.Itoa(contador) + `'>` + strconv.Itoa(int(superbloque.S_mnt_count)) + `</td></tr>`

	contador++

	grafo += `<tr><td>Inode size </td><td port='` + strconv.Itoa(contador) + `'>` + strconv.Itoa(int(superbloque.S_inode_size)) + `</td></tr>`

	contador++

	grafo += `<tr><td>Block size </td><td port='` + strconv.Itoa(contador) + `'>` + strconv.Itoa(int(superbloque.S_block_size)) + `</td></tr>`

	//execute -path=/home/darkun/Escritorio/prueba.mia

	contador++

	grafo += `<tr><td>Bitmap Inode Start </td><td port='` + strconv.Itoa(contador) + `'>` + strconv.Itoa(int(superbloque.S_bm_inode_start)) + `</td></tr>`

	contador++

	grafo += `<tr><td>Bitmap Block start </td><td port='` + strconv.Itoa(contador) + `'>` + strconv.Itoa(int(superbloque.S_bm_block_start)) + `</td></tr>`

	contador++

	grafo += `<tr><td>Inode Start </td><td port='` + strconv.Itoa(contador) + `'>` + strconv.Itoa(int(superbloque.S_inode_start)) + `</td></tr>`

	contador++

	grafo += `<tr><td>Block Start </td><td port='` + strconv.Itoa(contador) + `'>` + strconv.Itoa(int(superbloque.S_block_start)) + `</td></tr>`

	grafo += `</table>
				>`

	grafo += `}`

	dot := "sb.dot"

	file_dot, err := os.Create(dot)

	if err != nil {
		fmt.Println(err)
		return err
	}

	file_dot.WriteString(grafo)

	file_dot.Close()

	//-path=/home/darkun/Escritorio/mbr.pdf
	//dot -Tpdf mbr.dot  -o mbr.pdf

	//fmt.Println("\nEl path es: ", path)

	result := path

	//exec.Command("dot", "-Tpng", dot, "-o", result)
	out, err := exec.Command("dot", "-Tpdf", dot, "-o", result).Output()

	if err != nil {
		fmt.Println("\n Imprimiendo el error en Reporte Superbloque")
		log.Fatal(err)
	}

	fmt.Println(string(out))

	fmt.Println("\n\n========================= Finalizando Reporte Superbloque ===========================")

	return nil

}

func ReporteInode(index int, path string, file *os.File, tempMBR MBR, driveletter string) error {

	//execute -path=/home/darkun/Escritorio/prueba.mia

	//fmt.Println("\n\n========================= Iniciando Reporte Inode ===========================")

	var tempSuperblock Superblock

	if err := LeerObjeto(file, &tempSuperblock, int64(tempMBR.Mbr_partitions[index].Part_start)); err != nil {
		return err
	}

	Inodo_start := tempSuperblock.S_inode_start

	var inodo Inode

	grafo := `digraph H {
        label="Disco ` + strings.ToUpper(driveletter) + `";
        fontsize="50"
		graph [pad="0.5", nodesep="0.5", ranksep="1"];
		node [shape=plaintext]
		 rankdir=LR;`

	inodos_ocupados := tempSuperblock.S_inodes_count - tempSuperblock.S_free_inodes_count

	var inodo_vacio Inode

	for i := 0; i < int(inodos_ocupados); i++ {
		//fmt.Println("\nEstoy dentro del for de inodos")
		if err := LeerObjeto(file, &inodo, int64(Inodo_start+int32(i)*int32(binary.Size(Inode{})))); err != nil {
			return err
		}

		if inodo == inodo_vacio { // si hay un inodo vacio entre los que estan ocupados
			continue
		}

		grafo += `Inodo` + strconv.Itoa(i) + ` [
			label=<
				<table  border="0" cellborder="1" cellspacing="0">
				<tr><td colspan="3" port='0'>Inodo` + strconv.Itoa(i) + `</td></tr>`

		grafo += `<tr><td>I_Uid</td><td port='0'>` + strconv.Itoa(int(inodo.I_uid)) + `</td></tr>`
		grafo += `<tr><td>I_Gid</td><td port='0'>` + strconv.Itoa(int(inodo.I_gid)) + `</td></tr>`
		grafo += `<tr><td>I_perm</td><td port='0'>` + string(inodo.I_perm[:]) + `</td></tr>`
		grafo += `<tr><td>I_type</td><td port='0'>` + string(inodo.I_type[:]) + `</td></tr>`

		for j := 0; j < 15; j++ {
			grafo += `<tr><td>AD` + strconv.Itoa(j+1) + `</td><td port='` + strconv.Itoa(j+1) + `'>` + strconv.Itoa(int(inodo.I_block[j])) + `</td></tr>`
		}

		grafo += `</table>
			>];
			
			`

	}

	grafo, _ = EnlazandoInodos(file, tempSuperblock, driveletter, grafo) //enlazando los nodos

	grafo += `}`

	//fmt.Println("\nImprimiendo grafo: ", grafo)
	dot := "Inode.dot"

	file_dot, err := os.Create(dot)

	if err != nil {
		fmt.Println(err)
		return err
	}

	file_dot.WriteString(grafo)

	file_dot.Close()

	result := path

	//exec.Command("dot", "-Tpng", dot, "-o", result)
	out, err := exec.Command("dot", "-Tsvg", dot, "-o", result).Output()

	if err != nil {

		log.Fatal(err)
	}

	fmt.Println(string(out))

	//fmt.Println("\n\n========================= Finalizando Reporte Inode ===========================")

	//execute -path=/home/darkun/Escritorio/prueba.mia

	return nil
}

func EnlazandoInodos(file *os.File, tempSuperblock Superblock, disco string, grafo string) (string, error) {

	//fmt.Println("\n\n========================= Iniciando Enlazando INODOS ===========================")

	inodos_ocupados := tempSuperblock.S_inodes_count - tempSuperblock.S_free_inodes_count

	var inodo_vacio Inode

	var inodo Inode

	var cont int = 1

	for i := 0; i < int(inodos_ocupados-1); i++ {

		if err := LeerObjeto(file, &inodo, int64(tempSuperblock.S_inode_start+int32(i)*int32(binary.Size(Inode{})))); err != nil {
			return "", err
		}

		//execute -path=/home/darkun/Escritorio/prueba.mia

		if inodo == inodo_vacio {
			//grafo += `Inodo` + strconv.Itoa(int(i+1)) + `:0:w` + "\n"
			continue
		}

		grafo += `Inodo` + strconv.Itoa(int(i)) + `:0:e->`

		if err := LeerObjeto(file, &inodo, int64(tempSuperblock.S_inode_start+int32(i+1)*int32(binary.Size(Inode{})))); err != nil {
			return "", err
		}

		for inodo == inodo_vacio {
			fmt.Println("El Inodo" + strconv.Itoa(i+1) + " esta vacio")

			cont++

			if err := LeerObjeto(file, &inodo, int64(tempSuperblock.S_inode_start+int32(i+cont)*int32(binary.Size(Inode{})))); err != nil {
				return "", err
			}

		}

		if inodo != inodo_vacio {
			grafo += `Inodo` + strconv.Itoa(int(i+cont)) + `:0:w` + "\n"
			cont = 1
		}

	}

	//fmt.Println("\n\n========================= Finalizando Enlazando INODOS ===========================")

	return grafo, nil
}

func ReporteBlock(index int, path string, file *os.File, tempMBR MBR, driveletter string) error {

	//execute -path=/home/darkun/Escritorio/prueba.mia

	//fmt.Println("\n\n========================= Iniciando Reporte Block ===========================")

	var tempSuperblock Superblock

	if err := LeerObjeto(file, &tempSuperblock, int64(tempMBR.Mbr_partitions[index].Part_start)); err != nil {
		return err
	}

	Inodo_start := tempSuperblock.S_inode_start

	var inodo Inode

	grafo := `digraph H {
        label="Disco ` + strings.ToUpper(driveletter) + `";
        fontsize="50"
		graph [pad="0.5", nodesep="0.5", ranksep="1"];
		node [shape=plaintext]
		 rankdir=LR;`

	inodos_ocupados := tempSuperblock.S_inodes_count - tempSuperblock.S_free_inodes_count

	var inodo_vacio Inode

	for i := 0; i < int(inodos_ocupados); i++ {
		//fmt.Println("\nEstoy dentro del for de inodos")
		if err := LeerObjeto(file, &inodo, int64(Inodo_start+int32(i)*int32(binary.Size(Inode{})))); err != nil {
			return err
		}

		if inodo == inodo_vacio {
			continue
		}

		var bloque int32
		for k := 0; k < 15; k++ {

			if inodo.I_block[k] != -1 {

				bloque = inodo.I_block[k] // esto contiene el numero de bloque

				//fmt.Println("\nEl bloque a crear es: ", bloque)

				// carpeta -> 0   archivo -> 1
				if string(inodo.I_type[:]) == "0" { // es un bloque de carpetas
					var folder Folderblock

					//fmt.Println("\nEstoy dentro del for de inodos")
					if err := LeerObjeto(file, &folder, int64(tempSuperblock.S_block_start+bloque*int32(binary.Size(Folderblock{})))); err != nil {
						return err
					}

					grafo += `Bloque` + strconv.Itoa(int(bloque)) + ` [
					label=<
					<table  border="0" cellborder="1" cellspacing="0">
					<tr><td colspan="3" port='0'>Bloque` + strconv.Itoa(int(bloque)) + `</td></tr>`

					var indice int

					for j := 0; j < 4; j++ { // si es un bloque de carpetas
						//fmt.Println("\nEl name de folder es: ", string(folder.B_content[j].B_name[:]))
						for k := 0; k < len(folder.B_content[j].B_name[:]); k++ {
							if folder.B_content[j].B_name[k] == 0 { //quitando espacios(los ceros restantes) al slice de B_name
								indice = k
								break
							}

						}
						//fmt.Println("\nEl indice es: ", indice)
						name := string(folder.B_content[j].B_name[:indice])
						inodo := folder.B_content[j].B_inodo

						//fmt.Println("\nInodo es: ", inodo)
						grafo += `<tr><td>` + name + `</td><td port='` + strconv.Itoa(j+1) + `'>` + strconv.Itoa(int(inodo)) + `</td></tr>`
					}

					grafo += `</table>
						>];	
		
			`

				} else if string(inodo.I_type[:]) == "1" { //es un bloque de archivos
					var file_block Fileblock

					//fmt.Println("\nEstoy dentro del for de inodos")
					if err := LeerObjeto(file, &file_block, int64(tempSuperblock.S_block_start+int32(bloque)*int32(binary.Size(Fileblock{})))); err != nil {
						return err
					}

					grafo += `Bloque` + strconv.Itoa(int(bloque)) + ` [
						label=<
						<table  border="0" cellborder="1" cellspacing="0">
						<tr><td colspan="3" port='0'>Bloque` + strconv.Itoa(int(bloque)) + `</td></tr>`

					var indice int

					for k := 0; k < len(file_block.B_content[:]); k++ {
						if file_block.B_content[k] == 0 { //quitando espacios(los ceros restantes) al slice de B_content
							indice = k
							break
						}

					}
					//contenido := string(folder.B_content[j].B_name[:indice])
					var contenido string
					if indice == 0 { // significa que el slice fileblock.B_content esta lleno
						contenido = string(file_block.B_content[:])
					} else { //el slice todavia tiene espacios vacios
						contenido = string(file_block.B_content[:indice])
					}

					//fmt.Println("\nEl contenido de fileblock es: ", contenido)

					grafo += `<tr><td port='` + strconv.Itoa(int(bloque)+1) + `'>` + contenido + `</td></tr>`

					grafo += `</table>
					>];
			
				`

				}
			}
		}

	}

	grafo, _ = EnlazandoBloques(file, tempSuperblock, driveletter, grafo) //enlazando los nodos

	grafo += `}`

	//fmt.Println("\nImprimiendo grafo: ", grafo)
	dot := "Block.dot"

	file_dot, err := os.Create(dot)

	if err != nil {
		fmt.Println(err)
		return err
	}

	file_dot.WriteString(grafo)

	file_dot.Close()

	result := path

	//exec.Command("dot", "-Tpng", dot, "-o", result)
	out, err := exec.Command("dot", "-Tsvg", dot, "-o", result).Output()

	if err != nil {

		log.Fatal(err)
	}

	fmt.Println(string(out))

	//fmt.Println("\n\n========================= Finalizando Reporte Block ===========================")

	//execute -path=/home/darkun/Escritorio/prueba.mia

	return nil
}

func EnlazandoBloques(file *os.File, tempSuperblock Superblock, disco string, grafo string) (string, error) {

	fmt.Println("\n\n========================= Iniciando Enlazando BLOQUES ===========================")

	bloques_ocupados := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count

	var bloque_vacio Folderblock

	var bloque Folderblock

	var cont int = 1

	for i := 0; i < int(bloques_ocupados-1); i++ {

		if err := LeerObjeto(file, &bloque, int64(tempSuperblock.S_block_start+int32(i)*int32(binary.Size(Folderblock{})))); err != nil {
			return "", err
		}

		//execute -path=/home/darkun/Escritorio/prueba.mia

		if bloque == bloque_vacio {
			//grafo += `Bloque` + strconv.Itoa(int(i+1)) + `:0:w` + "\n"
			continue
		}

		grafo += `Bloque` + strconv.Itoa(int(i)) + `:0:e->`

		if err := LeerObjeto(file, &bloque, int64(tempSuperblock.S_block_start+int32(i+1)*int32(binary.Size(Folderblock{})))); err != nil {
			return "", err
		}

		for bloque == bloque_vacio {
			//fmt.Println("El bloque" + strconv.Itoa(i+1) + " esta vacio")
			cont++

			if err := LeerObjeto(file, &bloque, int64(tempSuperblock.S_block_start+int32(i+cont)*int32(binary.Size(Folderblock{})))); err != nil {
				return "", err
			}

		}

		if bloque != bloque_vacio {
			grafo += `Bloque` + strconv.Itoa(int(i+cont)) + `:0:w` + "\n"
			cont = 1
		}

	}

	fmt.Println("\n\n========================= Finalizando Enlazando BLOQUES ===========================")

	return grafo, nil
}

func ReporteFile(index int, path string, file *os.File, tempMBR MBR, driveletter string, ruta string) error {
	//index, path, file, TempMBR, driveletter

	fmt.Println("\n\n========================= Iniciando Reporte File ===========================")

	var tempSuperblock Superblock

	if err := LeerObjeto(file, &tempSuperblock, int64(tempMBR.Mbr_partitions[index].Part_start)); err != nil {
		return err
	}

	//crea el archivo donde se mostrara el reporte de bitmap inodos

	archivo, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return err
	}

	array := strings.Split(ruta, "/")

	array = array[1:]

	//fmt.Println("\nEl nuevo array es: ", array)

	Inodo_start := tempSuperblock.S_inode_start

	var inodo0 Inode

	if err := LeerObjeto(file, &inodo0, int64(Inodo_start)); err != nil {
		return err
	}

	var cont_folder int

	indexInode, _ := SearchingFile(array, inodo0, file, tempSuperblock, -1, cont_folder)

	//fmt.Println("\nindexInode el valor que devuelve SearchingFile: ", indexInode)

	//fmt.Println("\nEl bloque anterior es: ", bloque_anterior)

	//execute -path=/home/darkun/Escritorio/prueba.mia

	var crrInode Inode //Inodo 1

	if err := LeerObjeto(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Inode{})))); err != nil {
		return err
	}

	cadena, _ := ReadingFile(crrInode, indexInode, file, &tempSuperblock)

	archivo.WriteString(cadena)

	archivo.Close()
	file.Close()
	//file.Close() //execute -path=/home/darkun/Escritorio/prueba.mia

	fmt.Println("\n\n========================= Finalizando Reporte File ===========================")

	return nil

}

func ReadingFile(Inodo Inode, indexInode int32, file *os.File, tempSuperblock *Superblock) (string, error) {

	fmt.Println("\n\n========================= Iniciando ReadingFile ===========================")

	indice := int32(0)

	var cadena string

	for _, block := range Inodo.I_block {
		if block != -1 {
			if indice < 13 {

				// carpeta -> 0   archivo ->1

				if string(Inodo.I_type[:]) == "1" { // si es inodo de tipo archivo

					var archivo Fileblock

					// Read object from bin file                                       // un Inodo y un bloque las estructuras miden lo mismo
					if err := LeerObjeto(file, &archivo, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Fileblock{})))); err != nil {
						return "", err
					}

					for i := 0; i < len(archivo.B_content[:]); i++ {

						if archivo.B_content[i] != 0 {
							cadena += string(archivo.B_content[i])
						}

					}

					//execute -path=/home/darkun/Escritorio/prueba.mia

					cadena += "\n"

					//home/archivos/user/docs/Tarea2.txt

				}
			}
		}

		indice++
	}

	fmt.Println("\n\n========================= Finalizando ReadingFile ===========================")

	return cadena, nil
}

func ReporteJournaling(index int, path string, file *os.File, tempMBR MBR, driveletter string) error {

	fmt.Println("\n\n========================= Iniciando Reporte Journaling ===========================")

	var tempSuperblock Superblock

	if err := LeerObjeto(file, &tempSuperblock, int64(tempMBR.Mbr_partitions[index].Part_start)); err != nil {
		return err
	}

	if tempSuperblock.S_filesystem_type != 3 {
		fmt.Println("\n           ********* ERROR: La particion No tiene Journaling **********")
		fmt.Println()
		return nil
	}

	var journaling Journaling

	if err := LeerObjeto(file, &journaling, int64(tempMBR.Mbr_partitions[index].Part_start+int32(binary.Size(Superblock{})))); err != nil {
		return err
	}

	grafo := `digraph H {
			labelloc="t";
			label="Disco ` + strings.ToUpper(driveletter) + `";
			fontsize="50"
			graph [pad="0.5", nodesep="0.5", ranksep="1"];
			node [shape=plaintext]
			rankdir=LR;`

	var contador int

	grafo += `label=<
				<table  border="0" cellborder="1" cellspacing="0">`
	contador++

	grafo += `<tr><td colspan="4" style="filled" bgcolor="#FFD700"  port='` + strconv.Itoa(contador) + `'>Reporte Journaling</td></tr>`

	contador++

	grafo += `<tr><td>Operacion</td><td port='` + strconv.Itoa(contador) + `'>Path</td><td>Content</td><td>Date</td></tr>`

	for i := 0; i < len(journaling.Contenido); i++ {

		contador++

		var operation string

		for j := 0; j < len(journaling.Contenido[i].Operation[:]); j++ {

			if journaling.Contenido[i].Operation[j] != 0 {
				operation += string(journaling.Contenido[i].Operation[j])
			}

		}

		var path_ string

		for j := 0; j < len(journaling.Contenido[i].Path[:]); j++ {

			if journaling.Contenido[i].Path[j] != 0 {
				path_ += string(journaling.Contenido[i].Path[j])
			}

		}

		var content string

		for j := 0; j < len(journaling.Contenido[i].Content[:]); j++ {

			if journaling.Contenido[i].Content[j] != 0 {
				content += string(journaling.Contenido[i].Content[j])
			}

		}

		var date string

		for j := 0; j < len(journaling.Contenido[i].Date[:]); j++ {

			if journaling.Contenido[i].Date[j] != 0 {
				date += string(journaling.Contenido[i].Date[j])
			}

		}

		grafo += `<tr><td>` + operation + `</td><td port='` + strconv.Itoa(contador) + `'>` + path_ + `</td><td>` + content + `</td><td>` + date + `</td></tr>`

	}

	grafo += `</table>
					>;`

	grafo += `}`

	//fmt.Println("\nImprimiendo grafo: ", grafo)
	dot := "Journaling.dot"

	fmt.Println("\n Creando archivo Journaling.dot")

	file_dot, err := os.Create(dot)

	if err != nil {
		fmt.Println(err)
		return err
	}

	file_dot.WriteString(grafo)

	file_dot.Close()

	result := path

	//exec.Command("dot", "-Tpng", dot, "-o", result)
	out, err := exec.Command("dot", "-Tpdf", dot, "-o", result).Output()

	if err != nil {

		log.Fatal(err)
	}

	fmt.Println(string(out))

	fmt.Println("\n\n========================= Iniciando Reporte Journaling ===========================")

	return nil
}
