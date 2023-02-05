package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

var newImg = image.NewRGBA(image.Rect(0, 0, 10, 10))
var pourcentage_flou int64
var blur_group sync.WaitGroup

const BUFFERSIZE = 1024

var NUMBER_OF_CPUs = 12

func main() {
	var server_socket net.Listener
	var err error
	for {
		server_socket, err = net.Listen("tcp", "localhost:27001")
		if err != nil {
			fmt.Println("Error listetning: ", err)
		} else {
			break
		}
	}

	defer server_socket.Close()
	fmt.Println("Server started! Waiting for client_sockets...")

	for {
		var client_socket net.Conn
		var err error
		for {
			client_socket, err = server_socket.Accept()
			if err != nil {
				fmt.Println("Error: ", err)
			} else {
				break
			}
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
	fileSize := prepare_to_send(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := prepare_to_send(fileInfo.Name(), 64)
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
	bufferPercentage := make([]byte, 3)

	client_socket.Read(bufferPercentage)
	println("**************Pourcentage flou est :")
	//enlever les ":" inutiles de prepare_to_send
	p_string := strings.Trim(string(bufferPercentage), ":")
	//convertir en int64
	pourcentage_flou, _ = strconv.ParseInt(p_string, 10, 64)
	println(pourcentage_flou)

	client_socket.Read(bufferFileSize)
	fmt.Println(" ")
	fmt.Println("Receiving file of size : ")
	fmt.Println(bufferFileSize)
	fmt.Println(" ")
	fmt.Println(" ")
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	client_socket.Read(bufferFileName)

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

func prepare_to_send(retunString string, toLength int) string {
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

	// cette fois, le niveau de flou dépend du pourcentage donné (100% = moyenne de tous les pixels, 0% = image initiale)
	// ca marche bien entre 15 et 80

	if pourcentage_flou < 15 {
		pourcentage_flou = 15
	}
	if pourcentage_flou > 80 {
		pourcentage_flou = 80
	}
	x := float64(pourcentage_flou) / math.Log2(float64(cat.Bounds().Size().X))
	nv_flou_x := int(math.Pow(2, x))
	y := float64(pourcentage_flou) / math.Log2(float64(cat.Bounds().Size().Y))
	nv_flou_y := int(math.Pow(2, y))
	fmt.Println(nv_flou_x)
	fmt.Println(nv_flou_y)

	// création du channel
	numJobs := ((cat.Bounds().Size().X / nv_flou_x) + 1) * ((cat.Bounds().Size().Y / nv_flou_y) + 1)
	fmt.Println("numJobs :")
	fmt.Println(numJobs)
	jobs := make(chan [2]int, numJobs)

	newImg = image.NewRGBA(image.Rect(0, 0, cat.Bounds().Size().X, cat.Bounds().Size().Y))

	fmt.Println("WAITGROUP")
	fmt.Println(&blur_group)
	counter := 0
	fmt.Println(cat.Bounds().Size().X)
	fmt.Println(cat.Bounds().Size().Y)

	for i := 0; i < (cat.Bounds().Size().X); i = i + nv_flou_x {
		for j := 0; j < (cat.Bounds().Size().Y); j = j + nv_flou_y {
			jobs <- [2]int{i, j}
			counter++
		}
	}
	fmt.Println("counter :")
	fmt.Println(counter)

	for w := 1; w <= NUMBER_OF_CPUs; w++ {
		blur_group.Add(1)
		go box_blur(cat, nv_flou_x, nv_flou_y, jobs, &blur_group)
	}
	close(jobs)
	blur_group.Wait()

	outputFile, err := os.Create("image_temp.png")
	if err != nil {
		fmt.Println("pas possible de créer le nv fichier")
		return
	}
	png.Encode(outputFile, newImg)
	outputFile.Close()
}

func box_blur(oldImg image.Image, nv_flou_x int, nv_flou_y int, jobs <-chan [2]int, blur_group *sync.WaitGroup) {

	defer blur_group.Done()
	for index := range jobs {
		i := index[0]
		j := index[1]
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
}
