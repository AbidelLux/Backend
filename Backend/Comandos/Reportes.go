package Comandos

import (
	"Backend/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unsafe"
)

var contadorBloques int
var contadorArchivos int
var bloquesUsados []int64

func ValidarDatosREP(context []string) string {
	contadorBloques = 0
	contadorArchivos = 0
	bloquesUsados = []int64{}
	name := ""
	path := ""
	id := ""
	path_file_ls := ""
	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Comparar_Cadenas(tk[0], "path") {
			path = strings.ReplaceAll(tk[1], "\"", "")
		} else if Comparar_Cadenas(tk[0], "name") {
			name = tk[1]
		} else if Comparar_Cadenas(tk[0], "id") {
			id = tk[1]
		} else if Comparar_Cadenas(tk[0], "path_file_ls") {
			path_file_ls = tk[1]
		}
	}

	if name == "" || path == "" || id == "" {
		return Error("REP", "Se esperan parámetros obligatorios.")

	}

	if Comparar_Cadenas(name, "DISK") {
		return dks(path, id)
	} else if Comparar_Cadenas(name, "mbr") {
		return Mensaje("REP", "Reporte MBR")
	} else if Comparar_Cadenas(name, "Inodo") {
		return Mensaje("REP", "Reporte Inodo")
	} else if Comparar_Cadenas(name, "ls") {
		return Mensaje("REP", "Reporte LS")
	} else if Comparar_Cadenas(name, "block") {
		return Mensaje("REP", "Reporte LS")
	} else if Comparar_Cadenas(name, "bm_inode") {
		return Mensaje("REP", "Reporte BM INODE")
	} else if Comparar_Cadenas(name, "bm block") {
		return Mensaje("REP", "Reporte BM BLOCK")
	} else if Comparar_Cadenas(name, "sb") {
		return Mensaje("REP", "Reporte SB")
	} else if Comparar_Cadenas(name, "FILE") {
		if path_file_ls == "" {
			return Error("REP", "Se espera el parámetro ruta.")

		}
		return fileR(path, id, path_file_ls)
	} else if Comparar_Cadenas(name, "ls") {
		return Mensaje("REP", "Reporte LS")
	} else {
		return Error("REP", name+", no es un reporte válido.")
	}
}

func dks(p string, id string) string {
	var pth string
	GetMount("REP", id, &pth)

	file, err := os.Open(strings.ReplaceAll(pth, "\"", ""))

	if err != nil {
		return Error("REP", "No se ha encontrado el disco.")
	}
	var disk Structs.MBR
	file.Seek(0, 0)

	data := leerBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &disk)
	if err_ != nil {
		return Error("REP", "Error al leer el archivo")
	}
	file.Close()

	aux := strings.Split(p, ".")
	if len(aux) > 2 {
		return Error("REP", "No se admiten nombres de archivos que contengan punto (.)")
	}
	pd := aux[0] + ".dot"

	carpeta := ""
	direccion := strings.Split(pd, "/")

	fileaux, _ := os.Open(strings.ReplaceAll(pd, "\"", ""))
	if fileaux == nil {
		for i := 0; i < len(direccion); i++ {
			carpeta += "/" + direccion[i]
			if _, err_2 := os.Stat(carpeta); os.IsNotExist(err_2) {
				os.Mkdir(carpeta, 0777)
			}
		}
		os.Remove(pd)
	} else {
		fileaux.Close()
	}

	partitions := GetParticiones(disk)
	var extended Structs.Particion
	ext := false
	for i := 0; i < 4; i++ {
		if partitions[i].Part_status == '1' {
			if partitions[i].Part_type == "E"[0] || partitions[i].Part_type == "e"[0] {
				ext = true
				extended = partitions[i]
			}
		}
	}

	content := ""
	content = "digraph G{\n rankdir=TB;\n forcelabels= true;\n graph [ dpi = \"600\" ]; \n node [shape = plaintext];\n nodo1 [label = <<table>\n <tr>\n"
	var positions [5]int64
	var positionsii [5]int64
	positions[0] = disk.Mbr_partition_1.Part_start - (1 + int64(unsafe.Sizeof(Structs.MBR{})))
	positions[1] = disk.Mbr_partition_2.Part_start - disk.Mbr_partition_1.Part_start + disk.Mbr_partition_1.Part_s
	positions[2] = disk.Mbr_partition_3.Part_start - disk.Mbr_partition_2.Part_start + disk.Mbr_partition_2.Part_s
	positions[3] = disk.Mbr_partition_4.Part_start - disk.Mbr_partition_3.Part_start + disk.Mbr_partition_3.Part_s
	positions[4] = disk.Mbr_tamano + 1 - disk.Mbr_partition_4.Part_start + disk.Mbr_partition_4.Part_s

	copy(positionsii[:], positions[:])

	logic := 0
	tmplogic := ""
	if ext {
		tmplogic = "<tr>\n"
		auxEbr := Structs.NewEBR()
		//file, err := os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
		file, err = os.Open(strings.ReplaceAll(pth, "\"", ""))

		if err != nil {
			return Error("REP", "No se ha encontrado el disco.")
		}

		file.Seek(extended.Part_start, 0)
		data = leerBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &auxEbr)
		if err_ != nil {
			return Error("REP", "Error al leer el archivo")
		}
		file.Close()
		var tamGen int64 = 0
		for auxEbr.Part_next != -1 {
			tamGen += auxEbr.Part_s
			res := float64(auxEbr.Part_s) / float64(disk.Mbr_tamano)
			res = res * 100
			tmplogic += "<td>\"EBR\"</td>"
			s := fmt.Sprintf("%.2f", res)
			tmplogic += "<td>\"Logica\n " + s + "% de la partición extendida\"</td>\n"

			resta := float64(auxEbr.Part_next) - (float64(auxEbr.Part_start) + float64(auxEbr.Part_s))
			resta = resta / float64(disk.Mbr_tamano)
			resta = resta * 10000.00
			resta = math.Round(resta) / 100.00
			if resta != 0 {
				s = fmt.Sprintf("%f", resta)
				tmplogic += "<td>\"Logica\n " + s + "% libre de la partición extendida\"</td>\n"
				logic++
			}
			logic += 2
			file, err = os.Open(strings.ReplaceAll(pth, "\"", ""))

			if err != nil {
				return Error("REP", "No se ha encontrado el disco.")
			}

			file.Seek(auxEbr.Part_next, 0)
			data = leerBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &auxEbr)
			if err_ != nil {
				return Error("REP", "Error al leer el archivo")
			}
			file.Close()
		}
		resta := float64(extended.Part_s) - float64(tamGen)
		resta = resta / float64(disk.Mbr_tamano)
		resta = math.Round(resta * 100)
		if resta != 0 {
			s := fmt.Sprintf("%.2f", resta)
			tmplogic += "<td>\"Libre \n " + s + "% de la partición extendida.\"</td>\n"
			logic++
		}
		tmplogic += "</tr>\n\n"
		logic += 2

	}
	var tamPrim int64
	for i := 0; i < 4; i++ {
		if partitions[i].Part_type == 'E' {
			tamPrim += partitions[i].Part_s
			res := float64(partitions[i].Part_s) / float64(disk.Mbr_tamano)
			res = math.Round(res*10000.00) / 100.00
			s := fmt.Sprintf("%.2f", res)
			content += "<td COLSPAN='" + strconv.Itoa(logic) + "'> Extendida \n" + s + "% del disco</td>\n"
		} else if partitions[i].Part_start != -1 {
			tamPrim += partitions[i].Part_s
			res := float64(partitions[i].Part_s) / float64(disk.Mbr_tamano)
			res = math.Round(res*10000.00) / 100.00
			s := fmt.Sprintf("%.2f", res)
			content += "<td ROWSPAN='2'> Primaria \n" + s + "% del disco</td>\n"
		}
	}

	if tamPrim != 0 {
		libre := disk.Mbr_tamano - tamPrim
		res := float64(libre) / float64(disk.Mbr_tamano)
		res = math.Round(res * 100)
		s := fmt.Sprintf("%.2f", res)
		content += "<td ROWSPAN='2'> Libre \n" + s + "% del disco</td>"

	}
	content += "</tr>\n\n"
	content += tmplogic
	content += "</table>>];\n}\n"

	//CREAR IMAGEN
	b := []byte(content)
	err_ = ioutil.WriteFile(pd, b, 0644)
	if err_ != nil {
		log.Fatal(err_)
	}

	terminacion := strings.Split(p, ".")

	path, _ := exec.LookPath("dot")
	cmd, _ := exec.Command(path, "-T"+terminacion[1], pd).Output()
	mode := int(0777)
	ioutil.WriteFile(p, cmd, os.FileMode(mode))
	disco := strings.Split(pth, "/")
	return Mensaje("REP", "Reporte tipo DISK del disco "+disco[len(disco)-1]+", creado correctamente.")
}

func fileR(p string, id string, ruta string) string {

	carpeta := ""
	direccion := strings.Split(p, "/")

	fileaux, _ := os.Open(strings.ReplaceAll(p, "\"", ""))
	if fileaux == nil {
		for i := 0; i < len(direccion); i++ {
			carpeta += "/" + direccion[i]
			if _, err_2 := os.Stat(carpeta); os.IsNotExist(err_2) {
				os.Mkdir(carpeta, 0777)
			}
		}
		os.Remove(p)
	} else {
		fileaux.Close()
	}

	var path string
	particion := GetMount("MKDIR", id, &path)
	tmp := GetPath(ruta)
	data := getDataFile(tmp, particion, path)
	b := []byte(data)
	err_ := ioutil.WriteFile(p, b, 0644)
	if err_ != nil {
		log.Fatal(err_)
	}

	archivo := strings.Split(ruta, "/")
	return Mensaje("REP", "Reporte tipo FILE del archivo  "+archivo[len(archivo)-1]+", creado correctamente.")
}

func arregloString(arreglo [16]byte) string {
	reg := ""
	for i := 0; i < 16; i++ {
		if arreglo[i] != 0 {
			reg += string(arreglo[i])
		}
	}
	return reg
}

func existeEnArreglo(arreglo []int64, busqueda int64) int {
	regresa := 0
	for _, numero := range arreglo {
		if numero == busqueda {
			regresa++
		}
	}
	return regresa
}
