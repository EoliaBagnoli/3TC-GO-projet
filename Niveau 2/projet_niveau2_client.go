package main

import (
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

	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)
	if err != nil {
		fmt.Printf("Resolve failed")
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		fmt.Printf("Dial failed")
	}

	file, err := os.Open("/Users/eolia/Documents/INSA/3TC/ELP/3TC-GO-projet/test1.png") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // make sure to close the file even if we panic.

	buf := make([]byte, 1024)
	a_envoyer, err := io.CopyBuffer(conn, file, buf)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write(a_envoyer)
	if err != nil {
		fmt.Printf("Write failed")
	}

	//fmt.Println(n, "bytes sent")

	/*//charger l'image cherchée
	catFile, err := os.Open("/Users/eolia/Documents/INSA/3TC/ELP/3TC-GO-projet/test2.png")
	if err != nil {
		log.Fatal(err) // trouver comment enlever le fatal pour pas shutdown tout le programme
	}
	defer catFile.Close()

	cat, err := png.Decode(catFile)
	if err != nil {
		log.Fatal(err)
	}*/

	/*//mettre cat dans un buffer :
	buf := new(bytes.Buffer)
	png.Encode(buf, cat)
	a_envoyer := buf.Bytes()*/

	//****************************************************************************
	/*fileBuffer := make([]byte, 1024)
	file, err := os.Open("/Users/eolia/Documents/INSA/3TC/ELP/3TC-GO-projet/test3.png")
	var currentByte int64 = 0
	if err != nil {
		log.Fatal(err)
	}
	//conn.Write([]byte("send " + ))
	//read file until there is an error
	for err == nil || err != io.EOF {
		_, err = file.ReadAt(fileBuffer, currentByte)
		currentByte += 1024
		//fmt.Println(fileBuffer)
		conn.Write(fileBuffer)
	}*/
	//******************************************************************************

	/*//envoyer un msg
	_, err = conn.Write(a_envoyer)
	if err != nil {
		fmt.Printf("Write failed")
	}*/

	//recevoir un msg
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		println("Read data failed:", err.Error())
		os.Exit(1)
	}

	println("Received message:", string(received))

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

	// cette fois, le niveau de flou dépend du pourcentage donné (100% = moyenne de tous les pixels, 0% = image initiale) ça march po :((

	/*nv_flou_x := int(pourcentage_flou * float64(cat.Bounds().Size().X))
	nv_flou_y := int(pourcentage_flou * float64(cat.Bounds().Size().Y))
	fmt.Println(nv_flou_x)
	fmt.Println(nv_flou_y)*/
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

	/*sans go routines : tps moyen d'execution = 40 ms sur test3.png (1280 x 800 px)
	avec go routines : tps moyen d'execution = entre 1 et 11 ms sur même fichier pour flou de 30
	avec go routines : le temps dépend du niveau de flou que l'on veut et donc du nbre de go routines à créer : très efficace pour un flou
	elevé mais pas pour un flou petit. En dessous de 7, c'est mieux d'utiliser la version sans go routines. */

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

			//rester en uint32 ici

			r, g, b, a := oldImg.At(k, l).RGBA()

			newRed = (nbreElem*newRed + r) / (nbreElem + 1)
			newGreen = (nbreElem*newGreen + g) / (nbreElem + 1)
			newBlue = (nbreElem*newBlue + b) / (nbreElem + 1)
			newAlpha = (nbreElem*newAlpha + a) / (nbreElem + 1)

			nbreElem = nbreElem + 1
		}
	}

	//convertir en uint8 ici avec 4 nvelles var
	newRedConv = uint8(newRed / 257)
	newGreenConv = uint8(newGreen / 257)
	newBlueConv = uint8(newBlue / 257)
	newAlphaConv = uint8(newAlpha / 257)

	// au lieu d'écrire dans newImg, on écrit dans la grande newImg (var globale)
	for k := i; k < i+nv_flou_x; k++ {
		for l := j; l < j+nv_flou_y; l++ {
			newImg.Set(k, l, color.RGBA{newRedConv, newGreenConv, newBlueConv, newAlphaConv})
		}
	}
}