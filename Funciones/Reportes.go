package Funciones

import (
	//"encoding/binary"
	"encoding/binary"
	"fmt"
	"os"

	"log"
	"os/exec"
	"strconv"
	"strings"
)

func Reportes(name string, path string, id string, ruta string) error {

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
	fmt.Println("\n***********Imprimiendo MBR")
	fmt.Println()
	PrintMBR(TempMBR)
	fmt.Println("\n********Finalizando Impresion de MBR")

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

	if name == "tree" {
		ReporteTree(index, path, file, TempMBR)
	} else if name == "mbr" {
		ReporteMbr(path, file)
	} else if name == "disk" {
		ReporteDisk(path, file)
	}

	fmt.Println("\n\n========================= Fin REPORTES ===========================")

	return nil

}

func ReporteDisk(path string, file *os.File) error {

	fmt.Println("\n\n========================= Iniciando Reporte DISK ===========================")

	var TempMBR MBR

	if err := LeerObjeto(file, &TempMBR, int64(0)); err != nil {
		return err
	}

	PrintMBR(TempMBR)

	/*node [shape=record];
	  struct3 [label="MBR &#92;n 20%|
	  Libre
	  |{ Extendida |{EBR|Logica1|EBR |Logica2 |EBR}}|
	  Primaria | Libre

	  "];
	*/

	grafo := `digraph G {
		node [shape=record];`

	//MBR &#92;n 20%|
	grafo += `struct1 [label="MBR|`

	for i := 0; i < 4; i++ {

		if string(TempMBR.Mbr_partitions[i].Part_type[:]) == "p" {

			grafo += `Primaria|`

		} else if string(TempMBR.Mbr_partitions[i].Part_type[:]) == "e" {

			inicio := TempMBR.Mbr_partitions[i].Part_start

			grafo += `{ Extendida |`

			var tempEBR EBR

			if err := LeerObjeto(file, &tempEBR, int64(inicio)); err != nil { // obtiene el primer ebr
				return err
			}

			part_size := tempEBR.Part_size

			grafo += `{`

			grafo += `EBR |`

			var part_name string // elimina los espacios en el slice para que pueda ser leido por graphviz
			for j := 0; j < len(tempEBR.Part_name); j++ {
				if tempEBR.Part_name[j] != 0 {
					part_name += string(tempEBR.Part_name[j])
				}

			}

			grafo += part_name + ` |`

			for part_size != 0 { // obtiene los siguientes EBR analizando si existen por medio de su tamano

				grafo += `EBR |`

				part_start := tempEBR.Part_next

				if err := LeerObjeto(file, &tempEBR, int64(part_start)); err != nil { // obtiene el primer ebr
					return err
				}

				if tempEBR.Part_size != 0 {

					var part_name string // elimina los espacios en el slice para que pueda ser leido por graphviz
					for j := 0; j < len(tempEBR.Part_name); j++ {
						if tempEBR.Part_name[j] != 0 {
							part_name += string(tempEBR.Part_name[j])
						}

					}

					grafo += part_name + ` |`

					//grafo += `Logica1 | EBR | Logica2 | EBR | Logica3}|`

				}

				part_size = tempEBR.Part_size

			}

			grafo += `}}|`

			//grafo += `{ Extendida | {EBR | Logica1 | EBR | Logica2 | EBR | Logica3}}|`
		}
	}

	grafo += `Libre `
	grafo += `"];`
	grafo += `}`

	dot := "disk.dot"

	file, err := os.Create(dot)

	if err != nil {
		fmt.Println(err)
		return err
	}

	file.WriteString(grafo)

	file.Close()

	//-path=/home/darkun/Escritorio/mbr.pdf
	//dot -Tpdf mbr.dot  -o mbr.pdf

	fmt.Println("\nEl path es: ", path)

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

func ReporteMbr(path string, file *os.File) error {

	fmt.Println("\n\n========================= Iniciando Reporte MBR ===========================")
	fmt.Printf("\npath: %s", path)
	fmt.Println()

	var TempMBR MBR

	if err := LeerObjeto(file, &TempMBR, int64(0)); err != nil {
		return err
	}

	PrintMBR(TempMBR)

	grafo := `digraph H {
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

			contador++

			grafo += `<tr><td>fit</td><td port='` + strconv.Itoa(contador) + `'>` + string(TempMBR.Mbr_partitions[i].Part_fit[:]) + `</td></tr>`

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

			part_size := tempEBR.Part_size

			for part_size != 0 { // obtiene los siguientes EBR analizando si existen por medio de su tamano

				part_start := tempEBR.Part_next

				if err := LeerObjeto(file, &tempEBR, int64(part_start)); err != nil { // obtiene el primer ebr
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

				part_size = tempEBR.Part_size

			}

		}

	}

	grafo += `</table>
				>`

	grafo += `}`

	//fmt.Println("\nImprimiendo grafo: ", grafo)
	fmt.Println("\nImprimiendo grafo: ", grafo)
	dot := "mbr.dot"

	file, err := os.Create(dot)

	if err != nil {
		fmt.Println(err)
		return err
	}

	file.WriteString(grafo)

	file.Close()

	//-path=/home/darkun/Escritorio/mbr.pdf
	//dot -Tpdf mbr.dot  -o mbr.pdf

	fmt.Println("\nEl path es: ", path)

	result := path

	//exec.Command("dot", "-Tpng", dot, "-o", result)
	out, err := exec.Command("dot", "-Tpdf", dot, "-o", result).Output()

	if err != nil {

		log.Fatal(err)
	}

	fmt.Println(string(out))

	fmt.Println("\n\n========================= Finalizando Reporte MBR ===========================")

	return nil
}

func ReporteTree(index int, path string, file *os.File, TempMBR MBR) error {

	fmt.Println("\n\n========================= Iniciando Reporte Tree ===========================")
	fmt.Printf("\nIndex: %d, path: %s", index, path)
	fmt.Println()
	/*grafo := `digraph H {
	graph [pad="0.5", nodesep="0.5", ranksep="1"];
	node [shape=plaintext];
	rankdir=LR;`*/

	var tempSuperblock Superblock

	if err := LeerObjeto(file, &tempSuperblock, int64(TempMBR.Mbr_partitions[index].Part_start)); err != nil {
		return err
	}

	Inodo_start := tempSuperblock.S_inode_start

	var inodo Inode

	grafo := `digraph H {
		graph [pad="0.5", nodesep="0.5", ranksep="1"];
		node [shape=plaintext]
		 rankdir=LR;`

	fmt.Println("\n EL numero de inodos es: ", tempSuperblock.S_inodes_count)
	fmt.Println("\n EL numero de bloques es: ", tempSuperblock.S_blocks_count)

	for i := 0; i < int(tempSuperblock.S_inodes_count); i++ {
		//fmt.Println("\nEstoy dentro del for de inodos")
		if err := LeerObjeto(file, &inodo, int64(Inodo_start+int32(i)*int32(binary.Size(Inode{})))); err != nil {
			return err
		}

		/*for i := int32(0); i < 15; i++ {
			Inode0.I_block[i] = -1
		}

		Inode0.I_block[0] = 0
		*/
		grafo += `Inodo` + strconv.Itoa(i) + ` [
			label=<
				<table  border="0" cellborder="1" cellspacing="0">
				<tr><td colspan="3" port='0'>Inodo` + strconv.Itoa(i) + `</td></tr>`

		for j := 0; j < 15; j++ {
			grafo += `<tr><td>AD` + strconv.Itoa(j+1) + `</td><td port='` + strconv.Itoa(j+1) + `'>` + strconv.Itoa(int(inodo.I_block[j])) + `</td></tr>`
		}

		grafo += `</table>
			>];
			
			`

		var bloque int32
		for k := 0; k < 15; k++ {

			if inodo.I_block[k] != -1 {

				bloque = inodo.I_block[k] // esto contiene el numero de bloque

				// carpeta -> 0   archivo -> 1
				if string(inodo.I_type[:]) == "0" { // es un bloque de carpetas
					var folder Folderblock

					//fmt.Println("\nEstoy dentro del for de inodos")
					if err := LeerObjeto(file, &folder, int64(tempSuperblock.S_block_start+int32(bloque)*int32(binary.Size(Folderblock{})))); err != nil {
						return err
					}

					grafo += `Bloque` + strconv.Itoa(int(bloque)) + ` [
					label=<
					<table  border="0" cellborder="1" cellspacing="0">
					<tr><td colspan="3" port='0'>Bloque` + strconv.Itoa(int(bloque)) + `</td></tr>`

					var indice int

					for j := 0; j < 4; j++ { // si es un bloque de carpetas

						for k := 0; k < len(folder.B_content[j].B_name[:]); k++ {
							if folder.B_content[j].B_name[k] == 0 { //quitando espacios(los ceros restantes) al slice de B_name
								indice = k
								break
							}

						}
						//fmt.Println("\nEl indice es: ", indice)
						contenido := string(folder.B_content[j].B_name[:indice])

						grafo += `<tr><td>AD` + strconv.Itoa(j+1) + `</td><td port='` + strconv.Itoa(j+1) + `'>` + contenido + `</td></tr>`
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

					fmt.Println("\nEl contenido de fileblock es: ", contenido)

					grafo += `<tr><td port='` + strconv.Itoa(int(bloque)+1) + `'>` + contenido + `</td></tr>`

					grafo += `</table>
					>];
			
				`

				}
			}
		}

	}

	grafo += `}`

	//fmt.Println("\nImprimiendo grafo: ", grafo)
	dot := "tree.dot"

	file, err := os.Create(dot)

	if err != nil {
		fmt.Println(err)
		return err
	}

	file.WriteString(grafo)

	file.Close()

	result := path

	//exec.Command("dot", "-Tpng", dot, "-o", result)
	out, err := exec.Command("dot", "-Tpdf", dot, "-o", result).Output()

	if err != nil {

		log.Fatal(err)
	}

	fmt.Println(string(out))

	fmt.Println("\n\n========================= Finalizando Reporte Tree ===========================")

	return nil
}

/*	digraph H {
		graph [pad="0.5", nodesep="0.5", ranksep="1"];
		node [shape=plaintext]
		 rankdir=LR;

Inodo0 [
	   label=<
		 <table  border="0" cellborder="1" cellspacing="0">
		   <tr><td colspan="3" port='0'>Inodo 0</td></tr>
		   <tr><td>AD1</td><td port='1'>0</td></tr>
		   <tr><td>AD2</td><td port='2'>-1</td></tr>
		   <tr><td>AD3</td><td port='3'>-1</td></tr>
		   <tr><td>AD4</td><td port='4'>-1</td></tr>
			<tr><td>AD5</td><td port='5'>-1</td></tr>
		   <tr><td>AD6</td><td port='6'>-1</td></tr>
		 </table>
	  >];


Inodo1 [
	   label=<
		 <table  border="0" cellborder="1" cellspacing="0">
		   <tr><td colspan="3" port='0'>Inodo 1</td></tr>
		   <tr><td>AD1</td><td port='1'>1</td></tr>
		   <tr><td>AD2</td><td port='2'>2</td></tr>
		   <tr><td>AD3</td><td port='3'>-1</td></tr>
		   <tr><td>AD4</td><td port='4'>-1</td></tr>
			<tr><td>AD5</td><td port='5'>-1</td></tr>
		   <tr><td>AD6</td><td port='6'>-1</td></tr>
		 </table>
	  >];

Bloque0 [
	   label=<
		 <table  border="0" cellborder="1" cellspacing="0">
		   <tr><td colspan="3" port='0'>Bloque0</td></tr>
		   <tr><td>AD1</td><td port='1'>.</td></tr>
		   <tr><td>AD2</td><td port='2'>..</td></tr>
		   <tr><td>AD3</td><td port='3'>users.txt</td></tr>

		 </table>
	  >];

Bloque1 [
	   label=<
		 <table  border="0" cellborder="1" cellspacing="0">
		   <tr><td colspan="3" port='0'>Bloque1</td></tr>
		   <tr><td port='1'>1,G,root\n1,U,root,root,123</td></tr>
		 </table>
	  >]

	  Bloque2 [
	   label=<
		 <table  border="0" cellborder="1" cellspacing="0">
		   <tr><td colspan="3" port='0'>Bloque2</td></tr>
		   <tr><td port='1'>2,U,usuarios,user,123</td></tr>
		 </table>
	  >];


	Inodo0:1 -> Bloque0:0;
 	Bloque0:3 -> Inodo1:0;
	Inodo1:1 -> Bloque1:0;
	Inodo1:2 -> Bloque2:0;
}
*/
