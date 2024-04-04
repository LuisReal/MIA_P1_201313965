package Reportes

import (
	"os"
	//Funciones "MIA_P1_201313965/Funciones"
	"fmt"
	"log"
	"os/exec"
)

func ReporteTree() error {

	/*grafo := `digraph H {
	graph [pad="0.5", nodesep="0.5", ranksep="1"];
	node [shape=plaintext];
	rankdir=LR;`*/

	grafo := `digraph H {
		graph [pad="0.5", nodesep="0.5", ranksep="1"];
		node [shape=plaintext]
		 rankdir=LR;
		 
	  Inodo0 [
	   label=<
		 <table  border="0" cellborder="1" cellspacing="0">
		   <tr><td colspan="3" port='0'>Inodo 0</td></tr>
		   <tr><td>AD1</td><td port='1'>1</td></tr>
		   <tr><td>AD2</td><td port='2'>-1</td></tr>
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
	  
	  Inodo1 [
	   label=<
		 <table  border="0" cellborder="1" cellspacing="0">
		   <tr><td colspan="3" port='0'>Inodo 1</td></tr>
		   <tr><td>AD1</td><td port='1'>-1</td></tr>
		   <tr><td>AD2</td><td port='2'>-1</td></tr>
		   <tr><td>AD3</td><td port='3'>-1</td></tr>
		   <tr><td>AD4</td><td port='4'>-1</td></tr>
			<tr><td>AD5</td><td port='5'>-1</td></tr>
		   <tr><td>AD6</td><td port='6'>-1</td></tr>
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
	`
	grafo += `}`

	dot := "articulos.dot"

	file, err := os.Create(dot)

	if err != nil {
		fmt.Println(err)
		return err
	}

	file.WriteString(grafo)

	file.Close()

	result := "tree.png"
	comando := "dot -Tpng " + dot + " -o " + result

	fmt.Println("\nEl comando es:", comando)
	//system(comando.c_str());

	//exec.Command("dot", "-Tpng", dot, "-o", result)
	out, err := exec.Command("dot", "-Tpng", dot, "-o", result).Output()

	if err != nil {

		log.Fatal(err)
	}

	fmt.Println(string(out))
	return nil
}
