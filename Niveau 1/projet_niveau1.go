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
	"os"
	"time"
)

var newImg = image.NewRGBA(image.Rect(0, 0, 10, 10))

func main() {

	nv_flou := 40

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

	// utilisation des goroutines pour faire plusieurs fois le floutage sur des portions de l'image.

	//création nvelle image qui sera l'image floue finale
	newImg = image.NewRGBA(image.Rect(0, 0, cat.Bounds().Size().X, cat.Bounds().Size().Y))

	for i := 0; i < (cat.Bounds().Size().X); i = i + nv_flou {
		for j := 0; j < (cat.Bounds().Size().Y); j = j + nv_flou {
			//lancer la goroutine avec la modification de la nouvelle image (globale) direct dans la fonction
			go box_blur(cat, nv_flou, i, j)
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

func box_blur(oldImg image.Image, nv_flou int, i int, j int) /* *image.RGBA*/ {

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

	for k := i; k < i+nv_flou; k++ {
		for l := j; l < j+nv_flou; l++ {

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
	for k := i; k < i+nv_flou; k++ {
		for l := j; l < j+nv_flou; l++ {
			newImg.Set(k, l, color.RGBA{newRedConv, newGreenConv, newBlueConv, newAlphaConv})
		}
	}
}
