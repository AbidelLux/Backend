package main

import (
	"Backend/Comandos"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/rs/cors"
)

type Cmd_API struct {
	Cmd string `json:"CMD"`
}

var logued = false

var salida_cmd string = ""

func main() {
	fmt.Println("API backend")
	mux := http.NewServeMux()

	mux.HandleFunc("/analizar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var Content Cmd_API
		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Error reading body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &Content)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
			return
		}

		ejecuta_comando(Content.Cmd)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"result": salida_cmd})
		salida_cmd = ""
	})

	fmt.Println("Servidor en el puerto 5000")
	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":5000", handler))
}

func ejecuta_comando(cmd string) {
	arr_cmd := strings.Split(cmd, "\n")

	for i := 0; i < len(arr_cmd); i++ {
		if arr_cmd[i] != "" {
			texto := strings.TrimSpace(arr_cmd[i])
			if len(texto) > 0 {
				if texto[0] == '#' { // Si es un comentario
					salida_cmd += "# Comentario: " + texto[1:] + "\n"
				} else {
					comando := Comando(texto)
					texto = strings.TrimSpace(texto)
					texto = strings.TrimLeft(texto, comando)
					tokens := SepararTokens(texto)

					// Llamar a la función que ejecuta el comando basado en el token
					salida_cmd += ">> Procesando comando: " + comando + " " + texto + "\n"
					funciones(comando, tokens)
				}
			}
		}
	}
}
func Comando(eleccion string) string { //retorna solo el comando que recibimos de esa linea
	var token string
	bandera := false
	for i := 0; i < len(eleccion); i++ {
		if bandera {
			if string(eleccion[i]) == " " || string(eleccion[i]) == "-" {
				break
			}
			token += string(eleccion[i])

		} else if string(eleccion[i]) != " " && !bandera {
			if string(eleccion[i]) == "#" {
				token = eleccion
			} else {
				token += string(eleccion[i])
				bandera = true
			}
		}
	}
	return token
}

func SepararTokens(texto string) []string {
	var tokens []string
	if texto == "" {
		return tokens
	}
	texto += " " // Añadimos un espacio extra al final para manejar el último token
	var token string
	estado := 0

	for i := 0; i < len(texto); i++ {
		c := string(texto[i])

		if estado == 0 { // Estado base, esperando un "-"
			if c == "-" {
				estado = 1
				if token != "" { // Agregar el token si había algo antes
					tokens = append(tokens, token)
					token = ""
				}
			} else if c != " " {
				token += c // Si es un carácter, lo agregamos al token
			}
		} else if estado == 1 { // Estado esperando "=" o un valor
			if c == "=" {
				token += c
				estado = 2
			} else if c == " " {
				continue
			} else {
				token += c
			}
		} else if estado == 2 { // Estado esperando un valor después del "="
			if c == " " {
				if token != "" { // Si el token está completo, lo agregamos
					tokens = append(tokens, token)
					token = ""
				}
				estado = 0
			} else if c == "\"" {
				estado = 3 // Si encontramos comillas, entramos en estado 3
			} else {
				token += c
			}
		} else if estado == 3 { // Estado dentro de comillas
			if c == "\"" {
				estado = 4
			} else {
				token += c
			}
		} else if estado == 4 { // Estado después de comillas
			if c == " " {
				if token != "" { // Añadir el token completo
					tokens = append(tokens, token)
					token = ""
				}
				estado = 0
			}
		}
	}

	// Agregar el último token si no fue agregado
	if token != "" {
		tokens = append(tokens, token)
	}

	return tokens
}

func funciones(comando string, tokens []string) { //manda a llamar las funciones de cada comando
	if comando != "" {
		if Comandos.Comparar_Cadenas(comando, "mkdisk") {
			fmt.Println("============== Comando  \"MKDISK\" ==============")
			salida_cmd += "=============================================================\n"
			salida_cmd += "-->> Salida: \n" + Comandos.ValidarDatosMKDISK(tokens) + "\n"
			salida_cmd += "=============================================================\n"
		} else if Comandos.Comparar_Cadenas(comando, "rmdisk") {
			fmt.Println("============== Comando  \"RMDISK\" ==============")
			salida_cmd += "=============================================================\n"
			salida_cmd += "-->> Salida: \n" + Comandos.RMDISK(tokens) + "\n"
			salida_cmd += "=============================================================\n"
		} else if Comandos.Comparar_Cadenas(comando, "fdisk") {
			fmt.Println("=============== Comando \"FDISK\" ===============")
			salida_cmd += "=============================================================\n"
			salida_cmd += "-->> Salida: \n" + Comandos.ValidarDatosFDISK(tokens) + "\n"
			salida_cmd += "=============================================================\n"
		} else if Comandos.Comparar_Cadenas(comando, "mount") {
			fmt.Println("=============== Comando \"MOUNT\" ===============")
			salida_cmd += "=============================================================\n"
			salida_cmd += "-->> Salida: \n" + Comandos.ValidarDatosMOUNT(tokens) + "\n"
			salida_cmd += "=============================================================\n"
		} else if Comandos.Comparar_Cadenas(comando, "mkfs") {
			fmt.Println("=============== Comando  \"MKFS\" ===============")
			salida_cmd += "=============================================================\n"
			salida_cmd += "-->> Salida: \n" + Comandos.ValidarDatosMKFS(tokens) + "\n"
			salida_cmd += "=============================================================\n"
		} else if Comandos.Comparar_Cadenas(comando, "rep") {
			fmt.Println("================ Comando \"REP\" ================")
			salida_cmd += "=============================================================\n"
			salida_cmd += "-->> Salida: \n" + Comandos.ValidarDatosREP(tokens) + "\n"
			salida_cmd += "=============================================================\n"
		} else if Comandos.Comparar_Cadenas(comando, "cat") {
			fmt.Println("=============== Comando \"CAT\" ===============")
			salida_cmd += "=============================================================\n"

			salida_cmd += "=============================================================\n"
		} else if Comandos.Comparar_Cadenas(comando, "login") {
			fmt.Println("=============== Comando  \"LOGIN\" ===============")
			logued = true
			salida_cmd += "=============================================================\n"
			if logued {
				salida_cmd += Comandos.Error("LOGIN", "Usuario en linea.")
				salida_cmd += "=============================================================\n"
				return
			} else {
				logued = Comandos.ValidarDatosLOGIN(tokens)
				if logued == true {
					salida_cmd += Comandos.Mensaje("LOGIN", "Sesion Iniciada con Exito")

				}
				if logued == false {
					salida_cmd += Comandos.Error("LOGIN", "No se pudo iniciar sesion")
				}
			}
			salida_cmd += "=============================================================\n"
		} else if Comandos.Comparar_Cadenas(comando, "logout") {
			fmt.Println("============== Comando  \"LOGOUT\" ==============")
			salida_cmd += "=============================================================\n"
			if !logued {
				salida_cmd += Comandos.Error("LOGOUT", "No se ha iniciado sesion.")
				salida_cmd += "=============================================================\n"
				return
			} else {
				logued = Comandos.CerrarSesion()
				if logued == false {
					salida_cmd += Comandos.Mensaje("LOGOUT", "Sesion Terminada")
				}
			}
			salida_cmd += "=============================================================\n"
		} else if Comandos.Comparar_Cadenas(comando, "mkgrp") {
			fmt.Println("=============== Comando \"MKGRP\" ===============")
			salida_cmd += "=============================================================\n"
			if !logued {
				salida_cmd += Comandos.Error("MKGRP", "No se ha iniciado sesion.")
				salida_cmd += "=============================================================\n"
				return
			} else {
				salida_cmd += "Comando MKGRP \n"
			}
			salida_cmd += "=============================================================\n"
		} else if Comandos.Comparar_Cadenas(comando, "rmgrp") {
			fmt.Println("=============== Comando \"RMGRP\" ===============")
			salida_cmd += "=============================================================\n"
			if !logued {
				salida_cmd += Comandos.Error("RMGRP", "No se ha iniciado sesion.")
				salida_cmd += "=============================================================\n"
				return
			} else {
				salida_cmd += "Comando RMGRP \n"
			}
			salida_cmd += "=============================================================\n"
		} else if Comandos.Comparar_Cadenas(comando, "mkusr") {
			fmt.Println("=============== Comando \"MKUSR\" ===============")
			salida_cmd += "=============================================================\n"
			if !logued {
				salida_cmd += Comandos.Error("MKUSR", "No se ha iniciado sesion.")
				salida_cmd += "=============================================================\n"
				return
			} else {
				salida_cmd += "Comando MKUSR \n"
			}
			salida_cmd += "=============================================================\n"
		} else if Comandos.Comparar_Cadenas(comando, "rmusr") {
			fmt.Println("=============== Comando \"RMUSR\" ===============")
			salida_cmd += "=============================================================\n"
			if !logued {
				salida_cmd += Comandos.Error("RMUSR", "No se ha iniciado sesion.")
				salida_cmd += "=============================================================\n"
				return
			} else {
				salida_cmd += "Comando RMUSR \n"
			}
			salida_cmd += "=============================================================\n"
		} else if Comandos.Comparar_Cadenas(comando, "chgrp") {
			fmt.Println("=============== Comando \"CHGRP\" ===============")
			salida_cmd += "=============================================================\n"
			if !logued {
				salida_cmd += Comandos.Error("CHGRP", "Aún no se ha iniciado sesión.")
				salida_cmd += "=============================================================\n"
				return
			} else {
				salida_cmd += "Comando CHGRP \n"
			}
			salida_cmd += "=============================================================\n"
		} else if Comandos.Comparar_Cadenas(comando, "mkfile") {
			fmt.Println("============== Comando  \"MKFILE\" ==============")
			salida_cmd += "=============================================================\n"
			if !logued {
				salida_cmd += Comandos.Error("MKFILE", "No se ha iniciado sesion.")
				salida_cmd += "=============================================================\n"
				return

			} else {
				salida_cmd += "Comando MKFILE \n"
			}
			salida_cmd += "=============================================================\n"
		} else if Comandos.Comparar_Cadenas(comando, "mkdir") {
			fmt.Println("=============== Comando \"MKDIR\" ===============")
			salida_cmd += "=============================================================\n"
			if !logued {
				salida_cmd += Comandos.Error("MKDIR", "Aún no se ha iniciado sesión.")
				salida_cmd += "=============================================================\n"
				return
			} else {
				salida_cmd += "Comando MKDIR\n"
			}
			salida_cmd += "=============================================================\n"
		} else {
			salida_cmd += "=============================================================\n"
			salida_cmd += Comandos.Error("ANALIZADOR", "No se reconoce el comando \""+comando+"\"")
			salida_cmd += "=============================================================\n"
		}
	}
}
