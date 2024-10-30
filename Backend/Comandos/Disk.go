package Comandos

import (
	"Backend/Structs"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
)

// exec -path=/home/daniel/Escritorio/ArchivosPrueba/ArchivoEjemplo2.script

func ValidarDatosMKDISK(tokens []string) string {
	size := ""
	fit := ""
	unit := ""
	path := ""
	error_ := false
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		tk := strings.Split(token, "=")
		if Comparar_Cadenas(tk[0], "fit") {
			if fit == "" {
				fit = tk[1]
			} else {
				return Error("MKDISK", "parametro f repetido en el comando: "+tk[0])

			}
		} else if Comparar_Cadenas(tk[0], "size") {
			if size == "" {
				size = tk[1]
			} else {
				return Error("MKDISK", "parametro SIZE repetido en el comando: "+tk[0])

			}
		} else if Comparar_Cadenas(tk[0], "unit") {
			if unit == "" {
				unit = tk[1]
			} else {
				return Error("MKDISK", "parametro U repetido en el comando: "+tk[0])

			}
		} else if Comparar_Cadenas(tk[0], "path") {
			if path == "" {
				path = tk[1]
			} else {
				return Error("MKDISK", "parametro PATH repetido en el comando: "+tk[0])

			}
		} else {

			error_ = true
			return Error("MKDISK", "no se esperaba el parametro "+tk[0])
		}
	}
	if fit == "" {
		fit = "FF"
	}
	if unit == "" {
		unit = "M"
	}
	if error_ {
		return "Error"
	}
	if path == "" && size == "" {
		return Error("MKDISK", "se requiere parametro Path y Size para este comando")
	} else if path == "" {
		return Error("MKDISK", "se requiere parametro Path para este comando")
	} else if size == "" {
		return Error("MKDISK", "se requiere parametro Size para este comando")
	} else if !Comparar_Cadenas(fit, "BF") && !Comparar_Cadenas(fit, "FF") && !Comparar_Cadenas(fit, "WF") {
		return Error("MKDISK", "valores en parametro fit no esperados")

	} else if !Comparar_Cadenas(unit, "k") && !Comparar_Cadenas(unit, "m") {
		return Error("MKDISK", "valores en parametro unit no esperados")
	} else {
		return makeFile(size, fit, unit, path)
	}
}

func makeFile(s string, f string, u string, path string) string {
	var disco = Structs.NewMBR()
	size, err := strconv.Atoi(s)
	if err != nil {
		return Error("MKDISK", "Size debe ser un número entero")
	}
	if size <= 0 {
		return Error("MKDISK", "Size debe ser mayor a 0")
	}
	if Comparar_Cadenas(u, "M") {
		size = 1024 * 1024 * size
	} else if Comparar_Cadenas(u, "k") {
		size = 1024 * size
	}
	f = string(f[0])

	disco.Mbr_tamano = int64(size)
	fecha := time.Now().String()
	copy(disco.Mbr_fecha_creacion[:], fecha)
	aleatorio, _ := rand.Int(rand.Reader, big.NewInt(999999999))
	entero, _ := strconv.Atoi(aleatorio.String())
	disco.Mbr_dsk_signature = int64(entero)
	copy(disco.Dsk_fit[:], string(f[0]))
	disco.Mbr_partition_1 = Structs.NewParticion()
	disco.Mbr_partition_2 = Structs.NewParticion()
	disco.Mbr_partition_3 = Structs.NewParticion()
	disco.Mbr_partition_4 = Structs.NewParticion()

	if ArchivoExiste(path) {
		_ = os.Remove(path)
	}

	if !strings.HasSuffix(path, "mia") {
		return Error("MKDISK", "Extensión de archivo no válida.")
	}
	carpeta := ""
	direccion := strings.Split(path, "/")

	for i := 0; i < len(direccion)-1; i++ {
		carpeta += "/" + direccion[i]
		if _, err_ := os.Stat(carpeta); os.IsNotExist(err_) {
			os.Mkdir(carpeta, 0777)
		}
	}

	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return Error("MKDISK", "No se pudo crear el disco.")
	}
	var vacio int8 = 0
	s1 := &vacio
	var num int64 = 0
	num = int64(size)
	num = num - 1
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s1)
	EscribirBytes(file, binario.Bytes())

	file.Seek(num, 0)

	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s1)
	EscribirBytes(file, binario2.Bytes())

	file.Seek(0, 0)
	disco.Mbr_tamano = num + 1

	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, disco)
	EscribirBytes(file, binario3.Bytes())
	file.Close()
	nombreDisco := strings.Split(path, "/")
	return Mensaje("MKDISK", "¡Disco \""+nombreDisco[len(nombreDisco)-1]+"\" creado correctamente!")
}

func RMDISK(tokens []string) string {
	if len(tokens) > 1 {
		return Error("RMDISK", "Solo se acepta el parámetro PATH.")

	}
	path := ""
	error_ := false
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		tk := strings.Split(token, "=")
		if Comparar_Cadenas(tk[0], "path") {
			if path == "" {
				path = tk[1]
			} else {
				return Error("RMDISK", "Parametro PATH repetido en el comando: "+tk[0])

			}
		} else {
			error_ = true
			return Error("RMDISK", "no se esperaba el parametro "+tk[0])
		}
	}
	if error_ {
		return ""
	}
	if path == "" {
		return Error("RMDISK", "se requiere parametro Path para este comando")
	} else {
		if !ArchivoExiste(path) {
			return Error("RMDISK", "No se encontró el disco en la ruta indicada.")
		}
		if !strings.HasSuffix(path, "mia") {
			return Error("RMDISK", "Extensión de archivo no válida.")
		}

		err := os.Remove(path)
		if err != nil {
			return Error("RMDISK", "Error al intentar eliminar el archivo. :c")
		}
		return Mensaje("RMDISK", "Disco ubicado en "+path+", ha sido eliminado exitosamente.")

	}

}
