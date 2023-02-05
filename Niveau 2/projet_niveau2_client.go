package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const BUFFERSIZE = 1024

func main() {
	var server_socket net.Conn
	var err error
	for {
		server_socket, err = net.Dial("tcp", "localhost:27001")
		if err != nil {
			fmt.Println(err)
		} else {
			break
		}
	}
	defer server_socket.Close()
	answer(server_socket)
}

func answer(server_socket net.Conn) {
	fmt.Println("Connected to a server !")
	defer server_socket.Close()
	sendFileToServer(server_socket)
	getFileFromServer(server_socket)
	server_socket.Close()
}

func sendFileToServer(server_socket net.Conn) {
	fmt.Println("Let's send the picture we want to modify")
	var i string
	var file *os.File
	var err error
	for {
		fmt.Printf("Enter the name of the picture you want to blur (png only) : ")
		// Taking input from user
		fmt.Scanln(&i)
		file, err = os.Open("../" + i)
		if err != nil {
			fmt.Println(err)
		} else {
			break
		}
	}
	var p string
	for {
		fmt.Printf("Enter the percentage at which you want to blur (between 1 and 50) : ")
		fmt.Scanln(&p)
		p_int, err := strconv.Atoi(p)
		if err != nil {
			fmt.Println("Please enter an integer")
		} else if p_int < 0 || p_int > 100 {
			fmt.Println("Please enter a percentage between 0 and 100")
		} else {
			p = strconv.Itoa(p_int)
			p = fillString(p, 3)
			break
		}
	}

	// on recup les stats du fichier demandé
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 2)
	p_string := fillString(p, 3)

	size := []byte(fileSize)
	println(" ")
	println("File has a size of : ")
	fmt.Println(size)
	println(" ")
	println(" ")

	println("en train d'envoyer p")
	server_socket.Write([]byte(p_string))
	println(p)
	println("p envoyé")

	server_socket.Write(size)

	server_socket.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	for {
		// de façon infinie, on met l'image dans le buffer, on regarde si on a atteint le EOF (end of file), si non on envoie
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		server_socket.Write(sendBuffer)
	}
	fmt.Println("File has been sent !")
}

func getFileFromServer(server_socket net.Conn) {
	fmt.Println("Receiving the modified file")
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	server_socket.Read(bufferFileSize)
	fmt.Println(" ")
	fmt.Println("Receiving file of size : ")
	fmt.Println(bufferFileSize)
	fmt.Println(" ")
	fmt.Println(" ")
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	server_socket.Read(bufferFileName)

	newFile, err := os.Create("test_reception_TCP.png")

	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, server_socket, (fileSize - receivedBytes))
			server_socket.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, server_socket, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Println("Received file completely!")
}

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}
