/*
 * Copyright (c) 2017 Simon Schmidt
 * 
 * This software is provided 'as-is', without any express or implied
 * warranty. In no event will the authors be held liable for any damages
 * arising from the use of this software.
 * 
 * Permission is granted to anyone to use this software for any purpose,
 * including commercial applications, and to alter it and redistribute it
 * freely, subject to the following restrictions:
 * 
 * 1. The origin of this software must not be misrepresented; you must not
 *    claim that you wrote the original software. If you use this software
 *    in a product, an acknowledgment in the product documentation would be
 *    appreciated but is not required.
 * 2. Altered source versions must be plainly marked as such, and must not be
 *    misrepresented as being the original software.
 * 3. This notice may not be removed or altered from any source distribution.
 */
package main

import "image"
import "os"
import "image/png"
import "image/jpeg"
import "flag"
import "log"
import "github.com/nfnt/resize"

var dst = flag.String("dest","","Dest. Location")

var srcd = flag.String("srcd","","Src diffuse map")
var srcn = flag.String("srcn","","Src normal map")

var width = flag.Uint("width",0,"Output width")
var height = flag.Uint("height",0,"Output height")
var size = flag.Uint("size",0,"Output width and height")

var Bilinear = flag.Bool("bilinear",false,"Bilinear interpolation")
var Bicubic = flag.Bool("bicubic",false,"Bicubic interpolation")
var MitchellNetravali = flag.Bool("netravali",false,"Mitchell-Netravali interpolation")
var Lanczos = flag.Int("lanczos",0,"Lanczos interpolation a=?; possible 2 and 3")

var use_jpg = flag.Bool("jpg",false,"Use JPEG as output format")
var jpg_qual = flag.Int("jpgq",0,"JPEG quality (1-100")

var ipol = resize.NearestNeighbor 
func calc_ipol() {
	if *Bilinear { ipol = resize.Bilinear }
	if *Bicubic { ipol = resize.Bicubic }
	if *MitchellNetravali { ipol = resize.MitchellNetravali }
	switch *Lanczos {
	case 2: ipol = resize.Lanczos2
	case 3: ipol = resize.Lanczos3
	}
}
func load(str string) image.Image{
	f,e := os.Open(str)
	if e!=nil { log.Fatalln("Loading image: ",str," Error: ",e) }
	i,_,e := image.Decode(f)
	if e!=nil { log.Fatalln("Loading image: ",str," Error: ",e) }
	return i
}
func convert(i image.Image) image.Image {
	x := *width
	y := *height
	if x<1 { x = *size }
	if y<1 { y = *size }
	return resize.Resize(x,y,i,ipol)
}
func store(str string, i image.Image){
	if *use_jpg {
		f,e := os.Create(str+".jpg")
		if e!=nil { log.Fatalln("Storing image: ",str," Error: ",e) }
		if *jpg_qual>0 && *jpg_qual<=100 {
			e = jpeg.Encode(f,i,&jpeg.Options{Quality:*jpg_qual})
		}else{
			e = jpeg.Encode(f,i,nil)
		}
		if e!=nil { log.Fatalln("Storing image: ",str," Error: ",e) }
		return
	}
	{
		f,e := os.Create(str+".png")
		if e!=nil { log.Fatalln("Storing image: ",str," Error: ",e) }
		//e = png.Encode(f,i)
		e = (&png.Encoder{png.BestCompression}).Encode(f,i)
		if e!=nil { log.Fatalln("Storing image: ",str," Error: ",e) }
	}
}

func main(){
	flag.Parse()
	calc_ipol()
	if len(*dst)==0 {
		flag.PrintDefaults()
		return
	}
	if len(*srcd)>0 { store(*dst,convert(load(*srcd))) }
	if len(*srcn)>0 { store(*dst+".nm",convert(load(*srcn))) }
}
