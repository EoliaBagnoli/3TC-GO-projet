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
	server_socket, err := net.Dial("tcp", "localhost:27001")
	if err != nil {
		panic(err)
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
	/*var i string
	fmt.Printf("Enter the name of the picture you want to blur (png only) :")
	fmt.Scanf(i)
	file, err := os.Open("/mnt/c/Users/eolia/Documents/INSA/3TC/ELP/3TC-GO-projet/i")
	var p string
	fmt.Printf("Enter the name of the percentage at which you want to blur :")
	fmt.Scanf(p)*/
	//file, err := os.Open("/mnt/c/Users/eolia/Documents/INSA/3TC/ELP/3TC-GO-projet/test3.png")
	file, err := os.Open("/mnt/c/Users/eolia/Downloads/test3.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	// on recup les stats du fichier demandé
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	//print("File has a size of " + fileSize)
	fileName := fillString(fileInfo.Name(), 64)

	size := []byte(fileSize)
	println(" ")
	println("File has a size of : ")
	fmt.Println(size)
	println(" ")
	println(" ")

	server_socket.Write(size)

	//server_socket.Write([]byte(fileSize))
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
