package Comandos

import (
	"Backend/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

var DiscMont [99]DiscoMontado

type DiscoMontado struct {
	Path        [150]byte
	Estado      byte
	Particiones [26]ParticionMontada
}

type ParticionMontada struct {
	Letra  byte
	Estado byte
	Nombre [20]byte
}

var alfabeto = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

func ValidarDatosMOUNT(context []string) string {
	name := ""
	path := ""
	for i := 0; i < len(context); i++ {
		current := context[i]

		comando := strings.Split(current, "=")

		if Comparar_Cadenas(comando[0], "name") {
			name = comando[1]
		} else if Comparar_Cadenas(comando[0], "path") {
			path = strings.ReplaceAll(comando[1], "\"", "")
		}
	}
	if path == "" || name == "" {
		return Error("MOUNT", "El comando MOUNT requiere parámetros obligatorios")
	}
	salida := mount(path, name)
	fmt.Println(path)
	fmt.Println(name)
	return salida + listaMount()
}

func mount(p string, n string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recuperacion Archivo Inexistente:", r)
		}
	}()
	file, error_ := os.Open(p)
	if error_ != nil {
		return Error("MOUNT", "No se ha podido abrir el archivo.")
	}

	disk := Structs.NewMBR()
	file.Seek(0, 0)

	data := leerBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &disk)
	if err_ != nil {
		return Error("FDSIK", "Error al leer el archivo")
	}
	file.Close()

	particion := BuscarParticiones(disk, n, p)
	if particion.Part_type == 'E' || particion.Part_type == 'L' {
		var nombre [16]byte
		copy(nombre[:], n)
		if particion.Part_name == nombre && particion.Part_type == 'E' {
			return Error("MOUNT", "No se puede montar una partición extendida.")
		} else {
			ebrs := GetLogicas(*particion, p)
			encontrada := false
			if len(ebrs) != 0 {
				for i := 0; i < len(ebrs); i++ {
					ebr := ebrs[i]
					nombreebr := ""
					for j := 0; j < len(ebr.Part_name); j++ {
						if ebr.Part_name[j] != 0 {
							nombreebr += string(ebr.Part_name[j])
						}
					}

					if Comparar_Cadenas(nombreebr, n) && ebr.Part_mount == '1' {
						encontrada = true
						n = nombreebr
						break
					} else if nombreebr == n && ebr.Part_mount == '0' {
						return Error("MOUNT", "No se puede montar una partición Lógica eliminada.")
					}
				}
				if !encontrada {
					return Error("MOUNT", "No se encontró la partición Lógica.")
				}
			}
		}
	}
	for i := 0; i < 99; i++ {
		var ruta [150]byte
		copy(ruta[:], p)
		if DiscMont[i].Path == ruta {
			for j := 0; j < 26; j++ {
				var nombre [20]byte
				copy(nombre[:], n)
				if DiscMont[i].Particiones[j].Nombre == nombre {
					return Error("MOUNT", "Ya se ha montado la partición "+n)
				}
				if DiscMont[i].Particiones[j].Estado == 0 {
					DiscMont[i].Particiones[j].Estado = 1
					DiscMont[i].Particiones[j].Letra = alfabeto[j]
					copy(DiscMont[i].Particiones[j].Nombre[:], n)
					re := strconv.Itoa(i+1) + string(alfabeto[j])
					return Mensaje("MOUNT", "Se ha realizado correctamente el mount -id = 46"+re)

				}
			}
		}
	}
	for i := 0; i < 99; i++ {
		if DiscMont[i].Estado == 0 {
			DiscMont[i].Estado = 1
			copy(DiscMont[i].Path[:], p)
			for j := 0; j < 26; j++ {
				if DiscMont[i].Particiones[j].Estado == 0 {
					DiscMont[i].Particiones[j].Estado = 1
					DiscMont[i].Particiones[j].Letra = alfabeto[j]
					copy(DiscMont[i].Particiones[j].Nombre[:], n)

					re := strconv.Itoa(i+1) + string(alfabeto[j])
					return Mensaje("MOUNT", "se ha realizado correctamente el mount -id=46"+re)
				}
			}
		}
	}
	return ""
}

func GetMount(comando string, id string, p *string) Structs.Particion {
	if !(id[0] == '4' && id[1] == '6') {
		Error(comando, "El primer identificador no es válido.")
		return Structs.Particion{}
	}
	letra := id[len(id)-1]
	id = strings.ReplaceAll(id, "46", "")
	i, _ := strconv.Atoi(string(id[0] - 1))
	if i < 0 {
		Error(comando, "El primer identificador no es válido.")
		return Structs.Particion{}
	}
	for j := 0; j < 26; j++ {
		if DiscMont[i].Particiones[j].Estado == 1 {
			if DiscMont[i].Particiones[j].Letra == letra {

				path := ""
				for k := 0; k < len(DiscMont[i].Path); k++ {
					if DiscMont[i].Path[k] != 0 {
						path += string(DiscMont[i].Path[k])
					}
				}

				file, error := os.Open(strings.ReplaceAll(path, "\"", ""))
				if error != nil {
					Error(comando, "No se ha encontrado el disco")
					return Structs.Particion{}
				}
				disk := Structs.NewMBR()
				file.Seek(0, 0)

				data := leerBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
				buffer := bytes.NewBuffer(data)
				err_ := binary.Read(buffer, binary.BigEndian, &disk)

				if err_ != nil {
					Error("FDSIK", "Error al leer el archivo")
					return Structs.Particion{}
				}
				file.Close()

				nombreParticion := ""
				for k := 0; k < len(DiscMont[i].Particiones[j].Nombre); k++ {
					if DiscMont[i].Particiones[j].Nombre[k] != 0 {
						nombreParticion += string(DiscMont[i].Particiones[j].Nombre[k])
					}
				}
				*p = path
				return *BuscarParticiones(disk, nombreParticion, path)
			}
		}
	}
	return Structs.Particion{}
}

func listaMount() string {
	fmt.Println("\n<-------------------------- LISTADO DE MOUNTS -------------------------->")
	salidan_cmd := "\n<-------------------------- LISTADO DE MOUNTS -------------------------->\n"
	for i := 0; i < 99; i++ {
		for j := 0; j < 26; j++ {
			if DiscMont[i].Particiones[j].Estado == 1 {
				nombre := ""
				for k := 0; k < len(DiscMont[i].Particiones[j].Nombre); k++ {
					if DiscMont[i].Particiones[j].Nombre[k] != 0 {
						nombre += string(DiscMont[i].Particiones[j].Nombre[k])
					}
				}
				salidan_cmd += "\t id: 46" + strconv.Itoa(i+1) + string(alfabeto[j]) + ", Nombre: " + nombre + "\n"
			}
		}
	}
	return salidan_cmd
}
