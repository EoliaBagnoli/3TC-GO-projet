package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

var newImg = image.NewRGBA(image.Rect(0, 0, 10, 10))
var pourcentage_flou = 0.005

const BUFFERSIZE = 1024

func main() {
	server_socket, err := net.Listen("tcp", "localhost:27001")
	if err != nil {
		fmt.Println("Error listetning: ", err)
		os.Exit(1)
	}
	defer server_socket.Close()
	fmt.Println("Server started! Waiting for client_sockets...")
	for {
		client_socket, err := server_socket.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println("Client connected")
		go answer(client_socket)
	}
}

func answer(client_socket net.Conn) {
	defer client_socket.Close()
	getFileFromClient(client_socket)
	do_box_blur()
	println("*************Box blur done**************")
	sendFileToClient(client_socket)
	client_socket.Close()
}

func sendFileToClient(client_socket net.Conn) {
	fmt.Println("Let's send the modified picture")
	file, err := os.Open("/mnt/c/Users/eolia/Documents/INSA/3TC/ELP/3TC-GO-projet/Niveau 2/image_temp.png")
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
	fileName := fillString(fileInfo.Name(), 64)
	//print("File has a size of " + fileSize)
	size := []byte(fileSize)
	println(" ")
	println("File has a size of : ")
	fmt.Println(size)
	println(" ")
	println(" ")
	client_socket.Write(size)
	client_socket.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		client_socket.Write(sendBuffer)
	}
	fmt.Println("File has been sent !")
}

func getFileFromClient(client_socket net.Conn) {
	fmt.Println("Receiving the file")
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	client_socket.Read(bufferFileSize)
	fmt.Println(" ")
	fmt.Println("Receiving file of size : ")
	fmt.Println(bufferFileSize)
	fmt.Println(" ")
	fmt.Println(" ")
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	client_socket.Read(bufferFileName)
	//fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create("image_temp.png")

	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, client_socket, (fileSize - receivedBytes))
			client_socket.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, client_socket, BUFFERSIZE)
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

func do_box_blur() {

	catFile, err := os.Open("/mnt/c/Users/eolia/Documents/INSA/3TC/ELP/3TC-GO-projet/Niveau 2/image_temp.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer catFile.Close()

	cat, err := png.Decode(catFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	nv_flou_x := 15
	nv_flou_y := 10

	newImg = image.NewRGBA(image.Rect(0, 0, cat.Bounds().Size().X, cat.Bounds().Size().Y))

	for i := 0; i < (cat.Bounds().Size().X); i = i + nv_flou_x {
		for j := 0; j < (cat.Bounds().Size().Y); j = j + nv_flou_y {
			//lancer la goroutine avec la modification de la nouvelle image (globale) direct dans la fonction
			go box_blur(cat, nv_flou_x, nv_flou_y, i, j)
		}
	}

	outputFile, err := os.Create("image_temp.png")
	if err != nil {
		fmt.Println("pas possible de créer le nv fichier")
		return
	}
	png.Encode(outputFile, newImg)
	outputFile.Close()
}

func box_blur(oldImg image.Image, nv_flou_x int, nv_flou_y int, i int, j int) /* *image.RGBA*/ {

	var newRed uint32
	var newGreen uint32
	var newBlue uint32
	var newAlpha uint32
	var nbreElem uint32

	var newRedConv uint8
	var newGreenConv uint8
	var newBlueConv uint8
	var newAlphaConv uint8

	newRed = 0
	newGreen = 0
	newBlue = 0
	newAlpha = 0

	nbreElem = 0

	for k := i; k < i+nv_flou_x; k++ {
		for l := j; l < j+nv_flou_y; l++ {
			r, g, b, a := oldImg.At(k, l).RGBA()

			newRed = (nbreElem*newRed + r) / (nbreElem + 1)
			newGreen = (nbreElem*newGreen + g) / (nbreElem + 1)
			newBlue = (nbreElem*newBlue + b) / (nbreElem + 1)
			newAlpha = (nbreElem*newAlpha + a) / (nbreElem + 1)

			nbreElem = nbreElem + 1
		}
	}
	newRedConv = uint8(newRed / 257)
	newGreenConv = uint8(newGreen / 257)
	newBlueConv = uint8(newBlue / 257)
	newAlphaConv = uint8(newAlpha / 257)
	for k := i; k < i+nv_flou_x; k++ {
		for l := j; l < j+nv_flou_y; l++ {
			newImg.Set(k, l, color.RGBA{newRedConv, newGreenConv, newBlueConv, newAlphaConv})
		}
	}
}
