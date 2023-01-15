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

func main() {

	catFile, err := os.Open("/mnt/c/Users/eolia/Documents/INSA/3TC/ELP/3TC-GO-projet/test3.png")
	if err != nil {
		log.Fatal(err) // trouver comment enlever le fatal pour pas shutdown tout le programme
	}
	defer catFile.Close()

	cat, err := png.Decode(catFile)
	if err != nil {
		log.Fatal(err)
	}

	newImg := box_blur(cat, 30)

	// outputFile is a File type which satisfies Writer interface
	outputFile, err := os.Create("test_flou.png")
	if err != nil {
		fmt.Println("pas possible de créer le nv fichier")
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(outputFile, newImg)

	// Don't forget to close files
	outputFile.Close()

}

func box_blur(oldImg image.Image, nv_flou int) *image.RGBA {

	// flou gaussien par association de n pixels.

	//sans go routines : tps moyen d'execution = 40 ms sur test3.png (1280 x 800 px)

	start := time.Now()
	newImg := image.NewRGBA(image.Rect(0, 0, oldImg.Bounds().Size().X, oldImg.Bounds().Size().Y))

	var newRed uint32
	var newGreen uint32
	var newBlue uint32
	var newAlpha uint32
	var nbreElem uint32

	var newRedConv uint8
	var newGreenConv uint8
	var newBlueConv uint8
	var newAlphaConv uint8
	/*var r uint8
	var g uint8
	var b uint8
	var a uint8*/

	//for y := oldImg.Bounds().Min.Y; y < oldImg.Bounds().Max.Y; y++

	for i := 0; i < (oldImg.Bounds().Size().X); i = i + nv_flou {
		for j := 0; j < (oldImg.Bounds().Size().Y); j = j + nv_flou {

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

			for k := i; k < i+nv_flou; k++ {
				for l := j; l < j+nv_flou; l++ {
					newImg.Set(k, l, color.RGBA{newRedConv, newGreenConv, newBlueConv, newAlphaConv})
				}
			}

			//fmt.Println(("******************************************************************************"))
		}
	}
	end := time.Now()
	fmt.Println(end.Sub(start))
	return (newImg)
}
