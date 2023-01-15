package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const (
	HOST = "localhost"
	PORT = "1213"
	TYPE = "tcp"
)

var newImg = image.NewRGBA(image.Rect(0, 0, 10, 10))
var pourcentage_flou = 0.005

func main() {
	//création du server et attente du client
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		fmt.Printf("problème au listen")
		log.Fatal(err) // pas fatal
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("problème à listen.Accept")
			log.Fatal(err) // pas fatal
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {

	var currentByte int64 = 0

	fileBuffer := make([]byte, 1024)

	var err error
	file, err := os.Create("test_reception_TCP.png")
	if err != nil {
		log.Fatal(err)
	}

	for err == nil || err != io.EOF {

		conn.Read(fileBuffer)

		cleanedFileBuffer := bytes.Trim(fileBuffer, "\x00")

		_, err = file.WriteAt(cleanedFileBuffer, currentByte)
		if len(string(fileBuffer)) != len(string(cleanedFileBuffer)) {
			break
		}
		currentByte += 1024

	}

	conn.Close()
	file.Close()
	return

	/*for {
		//arrivée de message
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}

		// conversion buffer to image.image
		img := image.NewGray(image.Rect(0, 0, 100, 100))
		img.Pix = buffer

		//enregistrer dans un fichier
		outputFile, err := os.Create("test_reception_TCP.png")
		if err != nil {
			fmt.Println("pas possible de créer le nv fichier")
		}
		png.Encode(outputFile, img)
		outputFile.Close()

	}*/
}

func answer() {

	catFile, err := os.Open("/mnt/c/Users/eolia/Documents/INSA/3TC/ELP/3TC-GO-projet/test3.png")
	if err != nil {
		log.Fatal(err) // trouver comment enlever le fatal pour pas shutdown tout le programme
	}
	defer catFile.Close()

	cat, err := png.Decode(catFile)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	nv_flou_x := 15
	nv_flou_y := 10

	//création nvelle image qui sera l'image floue finale à la taille de l'ancienne
	newImg = image.NewRGBA(image.Rect(0, 0, cat.Bounds().Size().X, cat.Bounds().Size().Y))

	for i := 0; i < (cat.Bounds().Size().X); i = i + nv_flou_x {
		for j := 0; j < (cat.Bounds().Size().Y); j = j + nv_flou_y {
			//lancer la goroutine avec la modification de la nouvelle image (globale) direct dans la fonction
			go box_blur(cat, nv_flou_x, nv_flou_y, i, j)
		}
	}

	end := time.Now()
	fmt.Println(end.Sub(start))

	//création du fichier image floutée
	outputFile, err := os.Create("test_flou_niveau1.png")
	if err != nil {
		fmt.Println("pas possible de créer le nv fichier")
	}
	png.Encode(outputFile, newImg)
	outputFile.Close()

}

// @param : image à flouter, niveau de flou, numéro de la portion d'image par rapport à l'image originale

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
