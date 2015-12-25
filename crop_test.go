/*
 * Copyright (c) 2014 Christian Muehlhaeuser
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 *	Authors:
 *		Christian Muehlhaeuser <muesli@gmail.com>
 *		Michael Wendland <michael@michiwend.com>
 */

package smartcrop

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var (
	testFile = "./samples/gopher.jpg"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func BenchmarkCrop(b *testing.B) {
	fi, err := os.Open(testFile)
	if err != nil {
		b.Fatal(err)
	}
	defer fi.Close()
	img, _, err := image.Decode(fi)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := SmartCrop(img, 250, 250); err != nil {
			b.Error(err)
		}
	}
}

func TestCrop(t *testing.T) {
	fi, err := os.Open(testFile)
	if err != nil {
		t.Fatal(err)
	}
	defer fi.Close()

	img, _, err := image.Decode(fi)
	if err != nil {
		t.Error(err)
	}

	topCrop, err := SmartCrop(img, 250, 250)
	if err != nil {
		t.Error(err)
	}
	want := image.Rect(59, 0, 486, 427)
	if topCrop != want {
		t.Fatalf("want %v, got %v", want, topCrop)
	}

	/*
		sub, ok := img.(SubImager)
		if ok {
			cropImage := sub.SubImage(topCrop)
			writeImageToJpeg(cropImage, "./smartcrop.jpg")
		} else {
			t.Error(errors.New("No SubImage support"))
		}
	*/
}

func BenchmarkEdge(b *testing.B) {
	fname := "24391757.jpg"
	fi, err := os.Open("./samples/" + fname)
	if err != nil {
		b.Fatal(err)
	}
	defer fi.Close()
	img, _, err := image.Decode(fi)
	if err != nil {
		b.Error(err)
	}
	rgbaImg := toRGBA(img)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o := image.NewRGBA(img.Bounds())
		edgeDetect(rgbaImg, o)
	}
}

func BenchmarkImageDir(b *testing.B) {

	b.StopTimer()

	files, err := ioutil.ReadDir("./samples/hatchet")
	if err != nil {
		b.Error(err)
	}

	b.StartTimer()
	for _, file := range files {
		if strings.Contains(file.Name(), ".jpg") {

			fi, _ := os.Open("./samples/hatchet/" + file.Name())
			defer fi.Close()

			img, _, err := image.Decode(fi)
			if err != nil {
				b.Error(err)
			}

			topCrop, err := SmartCrop(img, 900, 500)
			if err != nil {
				b.Error(err)
			}
			fmt.Printf("Top crop: %+v\n", topCrop)

			sub, ok := img.(SubImager)
			//sub, ok := img.(SubImager)
			if ok {
				cropImage := sub.SubImage(topCrop)
				writeImageToJpeg(cropImage, "/tmp/smartcrop/smartcrop-"+file.Name())
			} else {
				b.Error(errors.New("No SubImage support"))
			}
		}
	}
	//fmt.Println("average time/image:", b.t

}
