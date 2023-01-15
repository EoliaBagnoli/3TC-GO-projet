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
)

func main() {

	catFile, err := os.Open("/mnt/c/Users/eolia/Documents/INSA/3TC/ELM/test1.png")
	if err != nil {
		log.Fatal(err) // trouver comment enlever le fatal pour pas shutdown tout le programme
	}
	defer catFile.Close()

	cat, err := png.Decode(catFile)
	if err != nil {
		log.Fatal(err)
	}

	newImg := gaussian_blur(cat, 40)

	// outputFile is a File type which satisfies Writer interface
	outputFile, err := os.Create("test.png")
	if err != nil {
		fmt.Println("pas possible de créer le nv fichier")
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(outputFile, newImg)

	// Don't forget to close files
	outputFile.Close()

}

func gaussian_blur(oldImg image.Image, nv_flou int) *image.RGBA {

	// flou gaussien par association de 4 pixels.
	// créer nouvelle image taille img1.taille

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

					/*r1 := uint8(r / 257)
					g1 := uint8(g / 257)
					b1 := uint8(b / 257)
					a1 := uint8(a / 257)*/

					fmt.Println(nbreElem + 1)

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

			fmt.Println(("******************************************************************************"))

			/*r11, g11, b11, a11 := oldImg.At(i, j).RGBA()
			r12, g12, b12, a12 := oldImg.At(i, j+1).RGBA()
			r13, g13, b13, a13 := oldImg.At(i+1, j).RGBA()
			r14, g14, b14, a14 := oldImg.At(i+1, j+1).RGBA()

			newRed = uint8(((r11 + r12 + r13 + r14) / 4) / 25)
			newGreen = uint8(((g11 + g12 + g13 + g14) / 4) / 257)
			newBlue = uint8(((b11 + b12 + b13 + b14) / 4) / 257)
			newAlpha = uint8(((a11 + a12 + a13 + a14) / 4) / 257)

			newImg.Set(i, j, color.RGBA{newRed, newGreen, newBlue, 255})
			newImg.Set(i+1, j, color.RGBA{newRed, newGreen, newBlue, 255})
			newImg.Set(i, j+1, color.RGBA{newRed, newGreen, newBlue, 255})
			newImg.Set(i+1, j+1, color.RGBA{newRed, newGreen, newBlue, 255})

			fmt.Println(r11)
			fmt.Println(g11)
			fmt.Println(b11)
			fmt.Println(newAlpha)*/
		}
	}
	return (newImg)
}
