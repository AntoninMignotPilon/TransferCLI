package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Erreur :", err)
		os.Exit(1)
	}
	defer listener.Close()

	receiver, err := listener.Accept()
	if err != nil {
		fmt.Println("Erreur :", err)
		os.Exit(1)
	}
	defer receiver.Close()
	fmt.Println("Receiver connecté")

	sender, err := listener.Accept()
	if err != nil {
		fmt.Println("Erreur :", err)
		os.Exit(1)
	}
	defer sender.Close()
	fmt.Println("Sender connecté")

	n, err := io.Copy(receiver, sender)
	if err != nil {
		fmt.Println("Erreur :", err)
		os.Exit(1)
	}
	fmt.Println(n, "octets envoyés")

}
