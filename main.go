package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

const (
	sendAction    = "send"
	receiveAction = "receive"
	bufferSize    = 64 * 1024
)

func main() {
	arg := os.Args[1:]

	if len(arg) < 2 {
		help()
	}

	switch arg[0] {
	case sendAction:
		send(arg[1])
	case receiveAction:
		receive(arg[1])
	default:
		help()
	}
}

func receive(code string) {

	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Erreur :", err)
		os.Exit(1)
	}
	defer conn.Close()

	// 1. lire la longueur du nom
	taille := make([]byte, 1)
	_, err = io.ReadFull(conn, taille)
	if err != nil {
		fmt.Println("Erreur de lecture :", err)
		os.Exit(1)
	}

	// 2. lire le nom
	nomBytes := make([]byte, taille[0])
	_, err = io.ReadFull(conn, nomBytes)
	if err != nil {
		fmt.Println("Erreur de lecture :", err)
		os.Exit(1)
	}

	// 3. convertir en texte, en gardant seulement le nom
	nom := filepath.Base(string(nomBytes))

	dest, err := os.Create(nom)
	if err != nil {
		fmt.Println("Erreur :", err)
		os.Exit(1)
	}
	defer dest.Close()

	// 4. lire tout le contenu du fichier
	buffer := make([]byte, bufferSize)
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Erreur de lecture :", err)
			os.Exit(1)
		}

		_, err = dest.Write(buffer[:n])
		if err != nil {
			fmt.Println("Erreur d'écriture :", err)
			os.Exit(1)
		}

		fmt.Println("Morceau reçu :", n, "octets")
	}

	fmt.Println("Fichier reçu :", nom)
}

func help() {
	fmt.Println("Utilisation :")
	fmt.Println("  transfert send <fichier>")
	fmt.Println("  transfert receive <code>")
	os.Exit(1)
}

func send(file string) {
	source, err := os.Open(file)
	if err != nil {
		fmt.Println("Erreur :", err)
		os.Exit(1)
	}
	defer source.Close()

	dest, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Erreur :", err)
		os.Exit(1)
	}
	defer dest.Close()

	// 1. garder seulement le nom, sans le dossier
	nom := filepath.Base(file)

	// 2. envoyer la longueur du nom
	_, err = dest.Write([]byte{byte(len(nom))})
	if err != nil {
		fmt.Println("Erreur d'envoi :", err)
		os.Exit(1)
	}

	// 3. envoyer le nom
	_, err = dest.Write([]byte(nom))
	if err != nil {
		fmt.Println("Erreur d'envoi :", err)
		os.Exit(1)
	}

	// 4. envoyer le contenu par chunk
	buffer := make([]byte, bufferSize)
	for {
		n, err := source.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Erreur de lecture :", err)
			os.Exit(1)
		}

		_, err = dest.Write(buffer[:n])
		if err != nil {
			fmt.Println("Erreur d'écriture :", err)
			os.Exit(1)
		}

		fmt.Println("Morceau envoyé :", n, "octets")
	}

	fmt.Println("Envoi terminé :", nom)
}
