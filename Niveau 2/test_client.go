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
	connection, err := net.Dial("tcp", "localhost:27001")
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	answer(connection)
}

func answer(connection net.Conn) {
	fmt.Println("Connected to a server !")
	defer connection.Close()
	sendFileToServer(connection)
	getFileFromServer(connection)
	connection.Close()
}

func sendFileToServer(connection net.Conn) {
	fmt.Println("Let's send the picture we want to modify")
	defer connection.Close()
	file, err := os.Open("/mnt/c/Users/eolia/Documents/INSA/3TC/ELP/3TC-GO-projet/test1.png")
	if err != nil {
		fmt.Println(err)
		return
	}
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

	//*********************************************************************PROBLEME****************************************************************
	connection.Write(size)

	//connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File has been sent !")
}

func getFileFromServer(connection net.Conn) {
	fmt.Println("Receiving the modified file")
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	fmt.Println(" ")
	fmt.Println("Receiving file of size : ")
	fmt.Println(bufferFileSize)
	fmt.Println(" ")
	fmt.Println(" ")
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	connection.Read(bufferFileName)
	//fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create("test_reception_TCP.png")

	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
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