//create Gaussian flou
//package image en go
//package png permet de décoder les images
// flouter une image = faire une moyenne entre des groupes de pixels

//fonction At(x, y)
//AlphaAt : avoir le taux de transparence d'un png

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"time"
)

var newImg = image.NewRGBA(image.Rect(0, 0, 10, 10))
var pourcentage_flou = 30

func main() {

	catFile, err := os.Open("/mnt/c/Users/eolia/Documents/INSA/3TC/ELP/3TC-GO-projet/test4.png")
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
	// ca marche entre 15 et 80
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

	//création nvelle image qui sera l'image floue finale à la taille de l'ancienne
	newImg = image.NewRGBA(image.Rect(0, 0, cat.Bounds().Size().X, cat.Bounds().Size().Y))
	for i := 0; i < (cat.Bounds().Size().X); i = i + nv_flou_x {
		for j := 0; j < (cat.Bounds().Size().Y); j = j + nv_flou_y {
			//lancer la goroutine avec la modification de la nouvelle image (globale) direct dans la fonction
			go box_blur(cat, nv_flou_x, nv_flou_y, i, j)
		}
	}
	// lui faire faire les 2 bandes restantes sur le bord au cas où la taille de l'image est pas divisible par le flou demandé
	reste_sur_x := cat.Bounds().Size().X % nv_flou_x
	reste_sur_y := cat.Bounds().Size().Y % nv_flou_y
	for i := 0; i < (cat.Bounds().Size().Y); i = i + nv_flou_y {
		go box_blur(cat, reste_sur_x, nv_flou_y, (cat.Bounds().Size().X - reste_sur_x), i)
	}
	for j := 0; j < (cat.Bounds().Size().X); j = j + nv_flou_x {
		go box_blur(cat, nv_flou_x, reste_sur_y, j, (cat.Bounds().Size().Y - reste_sur_y))
	}

	/*go box_blur(cat, nv_flou_x, (cat.Bounds().Size().X - reste_sur_x), 0, reste_sur_x)
	go box_blur(cat, (cat.Bounds().Size().Y - reste_sur_y), nv_flou_y, reste_sur_y, 0)*/

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
	var nbreElem uint32

	var newRedConv uint8
	var newGreenConv uint8
	var newBlueConv uint8

	newRed = 0
	newGreen = 0
	newBlue = 0

	nbreElem = 0

	for k := i; k < i+nv_flou_x; k++ {
		for l := j; l < j+nv_flou_y; l++ {

			//rester en uint32 ici

			r, g, b, _ := oldImg.At(k, l).RGBA()

			newRed = (nbreElem*newRed + r) / (nbreElem + 1)
			newGreen = (nbreElem*newGreen + g) / (nbreElem + 1)
			newBlue = (nbreElem*newBlue + b) / (nbreElem + 1)

			nbreElem = nbreElem + 1
		}
	}

	//convertir en uint8 ici avec 4 nvelles var
	newRedConv = uint8(newRed / 257)
	newGreenConv = uint8(newGreen / 257)
	newBlueConv = uint8(newBlue / 257)

	/*println(newBlueConv)
	println(newRedConv)
	println(newGreenConv)*/

	// au lieu d'écrire dans newImg, on écrit dans la grande newImg (var globale)
	for k := i; k < i+nv_flou_x; k++ {
		for l := j; l < j+nv_flou_y; l++ {
			newImg.Set(k, l, color.RGBA{newRedConv, newGreenConv, newBlueConv, 255})
			//newImg.Set(k, l, color.RGBA{255, 0, 0, 255})
		}
	}
}
