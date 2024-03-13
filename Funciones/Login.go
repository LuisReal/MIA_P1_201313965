package Funciones

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Login(user string, pass string, id string) (string, error) {
	fmt.Println("\n\n========================= LOGIN ===========================")

	fmt.Printf("\nUser: %s, pass: %s, id: %s\n", user, pass, id)

	driveletter := string(id[0])

	// Open bin file
	filepath := "./archivos/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := abrirArchivo(filepath)
	if err != nil {
		return "", err
	}

	var TempMBR MBR
	// Read object from bin file
	if err := LeerObjeto(file, &TempMBR, 0); err != nil {
		return "", err
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
			if strings.Contains(string(TempMBR.Mbr_partitions[i].Part_id[:]), id) {
				fmt.Println("\n****Particion Encontrada*****")
				if TempMBR.Mbr_partitions[i].Part_status { // si la particion esta montada = true
					fmt.Println("\n*******La particion esta montada*****")
					index = i
				} else {
					fmt.Println("\n*******La particion NO esta montada*****")
					return "", err
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
		return "", err
	}

	var tempSuperblock Superblock

	if err := LeerObjeto(file, &tempSuperblock, int64(TempMBR.Mbr_partitions[index].Part_start)); err != nil {
		return "", err
	}

	// initSearch /users.txt -> regresa no Inodo
	// initSearch -> 1

	indexInode := int32(1)

	var crrInode Inode

	if err := LeerObjeto(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Inode{})))); err != nil {
		return "", err
	}

	// getInodeFileData -> Iterate the I_Block n concat the data

	var fileblock Fileblock

	fileblock_start := tempSuperblock.S_block_start + crrInode.I_block[0]*int32(binary.Size(Fileblock{}))

	if err := LeerObjeto(file, &fileblock, int64(fileblock_start)); err != nil {
		return "", err
	}

	fmt.Println("Fileblock------------")
	//data := "1,G,root\n1,U,root,root,123\n"
	data := string(fileblock.B_content[:])

	lines := strings.Split(data, "\n")

	datos := strings.Split(lines[1], ",")

	usuario := datos[3]
	contrasena := datos[4]

	fmt.Println("Inode", crrInode.I_block)

	// Close bin file
	defer file.Close()

	if usuario == user && pass == contrasena {
		fmt.Println("\n **********Usuario encontrado***********")

		fmt.Println("\n\n========================= FIN LOGIN ===========================")

		return usuario, err
	} else {
		fmt.Println("\n*********Usuario NO encontrado**********")

		fmt.Println("\n\n========================= FIN LOGIN ===========================")

		return "failed", err
	}

}

func Mkgrp(name string, id string) {
	fmt.Println("\n\n========================= Inicio MKGRP ===========================")

	fmt.Printf("El usuario a crear sera: %s, El id es: %s", name, id)
	fmt.Println()

	//return file, fileblock, fileblock_start, nil
	file, fileblock, start_fileblock, err := getUsersTXT(id)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("Fileblock------------")
	//data := "1,G,root\n1,U,root,root,123\n"

	var cadena string = " "

	cadena = string(fileblock.B_content[:])

	fmt.Println("\n Imprimiendo cadena: ", string(fileblock.B_content[:]))

	lines := strings.Split(cadena, "\n")

	if len(lines) > 0 {
		lines = lines[:len(lines)-1]
	}

	fmt.Println("\n\nContenido del arreglo lines: ", lines)
	fmt.Println("\nEl tamano del arreglo lines es: ", len(lines))

	fmt.Println("\nImprimiendo ultimo elemento de arreglo lines: ", lines[len(lines)-1])
	//2, G, usuarios, \n
	var contador int = 0
	var exist int = 0
	var datos []string

	for i := 0; i < len(lines); i++ {

		datos = strings.Split(lines[i], ",")

		contador_, _ := strconv.Atoi(datos[0])

		contador = contador_
		contador++

		if len(datos) != 0 {

			if string(datos[2]) == name {

				fmt.Println("\n\n      ********** El Grupo ya existe ************")

				fmt.Println("\n\n========================= Fin MKGRP ===========================")
				exist++
				return
			}
		}

	}

	if exist == 0 { // si el grupo a crear no existe

		newCadena := strconv.Itoa(contador) + ",G," + name + "\n"

		fmt.Println("\n ********datos de la variable newCadena: ", newCadena)

		var c int
		for i := 0; i < len(fileblock.B_content); i++ {
			//fmt.Println(fileblock[i])

			if fileblock.B_content[i] == 0 { // si hay todavia espacio

				if c < len(newCadena) {

					fileblock.B_content[i] = byte(newCadena[c])
					//fmt.Printf("agregando letra:  %s   ", string(newCadena[c]))
					c++

				} else {
					break
				}

			}
		}

		fmt.Println("\n El contenido nuevo de B_content es: ", string(fileblock.B_content[:]))

		fmt.Println("\n\n ********** Escribiendo objeto FILEBLOCK en el archivo ******************")
		if err := escribirObjeto(file, fileblock, int64(start_fileblock)); err != nil { //aqui solo escribi el primer EBR
			return
		}

	}

	var tempfileblock Fileblock

	fmt.Println("\n\n ********** Recuperando y Leyendo objeto FILEBLOCK del archivo binario ******************")
	if err := LeerObjeto(file, &tempfileblock, int64(start_fileblock)); err != nil {
		return
	}

	printFileblock(tempfileblock)

	//fmt.Println("\n\nLo que se guardo en fileblock.B_content es: ", string(fileblock.B_content[:]))

	fmt.Println("\n\n========================= Fin MKGRP ===========================")
}

func getUsersTXT(id string) (*os.File, Fileblock, int32, error) {

	driveletter := string(id[0])

	// Open bin file
	filepath := "./archivos/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := abrirArchivo(filepath)
	if err != nil {
		return nil, Fileblock{}, 0, err
	}

	var TempMBR MBR
	// Read object from bin file
	if err := LeerObjeto(file, &TempMBR, 0); err != nil {
		return nil, Fileblock{}, 0, err
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
			if strings.Contains(string(TempMBR.Mbr_partitions[i].Part_id[:]), id) {
				fmt.Println("\n****Particion Encontrada*****")
				if TempMBR.Mbr_partitions[i].Part_status { // si la particion esta montada = true
					fmt.Println("\n*******La particion esta montada*****")
					index = i
				} else {
					fmt.Println("\n*******La particion NO esta montada*****")
					return nil, Fileblock{}, 0, err
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
		return nil, Fileblock{}, 0, err
	}

	var tempSuperblock Superblock

	if err := LeerObjeto(file, &tempSuperblock, int64(TempMBR.Mbr_partitions[index].Part_start)); err != nil {
		return nil, Fileblock{}, 0, err
	}

	// initSearch /users.txt -> regresa no Inodo
	// initSearch -> 1

	indexInode := int32(1)

	var crrInode Inode

	if err := LeerObjeto(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Inode{})))); err != nil {
		return nil, Fileblock{}, 0, err
	}

	// getInodeFileData -> Iterate the I_Block n concat the data

	var fileblock Fileblock

	fileblock_start := tempSuperblock.S_block_start + crrInode.I_block[0]*int32(binary.Size(Fileblock{}))

	if err := LeerObjeto(file, &fileblock, int64(fileblock_start)); err != nil {
		return nil, Fileblock{}, 0, err
	}

	return file, fileblock, fileblock_start, nil

}

/*	 CODIGO PARA MANEJAR LOS SLICES DE BYTES DE TIPO [SIZE]BYTE

cadena := "1,U,root,123\n"
  //usuario := "2,U,user,dracker"
  var fileblock [32]byte
  copy(fileblock[:], []byte(cadena))

  //data := string(fileblock[:])
  //fmt.Println("\nLa data es: ",data)
  data := "2,U,usuario,562\n"
  //cadena += "3,U,user,002"
  fmt.Println("\nLa NUEVA data es: ",data)
  fmt.Println("la longitud de la cadena es: ", len(data))
  //Data := make([]byte,3)

  //fmt.Println(Data) //output is [0,0,0]

  var c int
  for i := 0; i < len(fileblock); i++ {
      //fmt.Println(fileblock[i])

      if fileblock[i] ==0 {

          if c < len(data){
              fileblock[i] = byte(data[c])
              fmt.Printf("letra:  %s   ", string(data[c]))
              c++

          }else{
              break
          }


      }
  }

  var contador int
  for i := 0; i < len(fileblock); i++ {
      if fileblock[i] ==0 {
        contador++
      }

  }


  fmt.Println("\nfileblock: ", string(fileblock[:]))
  fmt.Println("\nEl nuevo tamano de fileblock es: ", len(fileblock))

  if contador < len(data){
      fmt.Println("\nEl contador es: ", contador)
      fmt.Println("\nYa no hay suficiente espacio")
  }
*/
