package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
)

type Params struct {
	Count  int    `json:"count"`
	File   string `json:"file"`
	Height int    `json:"height"`
	Url    string `json:"url"`
	Width  int    `json:"width"`
}

type Yvonne struct {
	Thumb string `json:"thumb"`
	Link  string `json:"link"`
}

func checkErr(e error, msg string) {
	if e != nil {
		fmt.Println(fmt.Sprintf("ERROR: %s", msg))
		os.Exit(1)
	}
}

func getParams() Params {
	const h int = 640
	const w int = 480

	var p Params

	flag.IntVar(&p.Count, "count", 100, "the number of images to retrieve")
	flag.StringVar(&p.File, "file", os.Getenv("IMAGE_FILE"), "the number of images to retrieve")
	flag.IntVar(&p.Height, "height", h, "the height of the image")
	flag.StringVar(&p.Url, "url", os.Getenv("RANDOM_URL"), "the url to get a random image")
	flag.IntVar(&p.Width, "width", w, "the width of the image")
	flag.Parse()

	return p
}

// TODO: Put a mechanism in this function so it will retry on failures n times then stop.
func fetchImageFromLink(url string) http.Response {
	resp1, err := http.Get(url)
	checkErr(err, fmt.Sprintf("Could not acquire content from url: %s", url))

	jsonBody, _ := ioutil.ReadAll(resp1.Body)
	resp1.Body.Close()

	yvonneStruct := Yvonne{}
	json.Unmarshal(jsonBody, &yvonneStruct)

	resp2, err2 := http.Get(yvonneStruct.Thumb)
	checkErr(err2, fmt.Sprintf("Could not acquire content from url: %s", yvonneStruct.Thumb))

	return *resp2
}

func main() {
	params := getParams()

	var X int = 0
	var Y int = 0

	next_pt := []image.Point{}
	next_pt = append(next_pt, image.Point{X, Y})

	im := image.NewRGBA(image.Rect(X, Y, params.Width, params.Height))

	for loopCount := 0; loopCount < params.Count; loopCount++ {
		resp := fetchImageFromLink(params.Url).Body
		m, _, err := image.Decode(resp)
		checkErr(err, "Could not read body")
		resp.Close()

		n := next_pt[0]
		next_pt = next_pt[1:]
		r := image.Rect(n.X, n.Y, m.Bounds().Dx()+n.X, m.Bounds().Dy()+n.Y)
		draw.Draw(im, r, m, image.Point{0, 0}, draw.Over)

		tmpX := n.X + m.Bounds().Dx()
		tmpY := n.Y + m.Bounds().Dy()

		if tmpX < params.Width {
			next_pt = append(next_pt, image.Point{tmpX, n.Y})
		}

		if tmpY < params.Height {
			next_pt = append(next_pt, image.Point{n.X, tmpY})
		}

		// When there's no available pixels within the bounds end the loop.
		if tmpX >= params.Width && tmpY >= params.Height {
			break
		}
	}

	filename := fmt.Sprintf("%s/%s", os.Getenv("BASE_DIR"), params.File)

	if _, err2 := os.Stat(filename); !os.IsNotExist(err2) {
		err3 := os.Remove(filename)
		checkErr(err3, fmt.Sprintf("Could not delete file: %s", filename))
	}

	newImg, err4 := os.Create(filename)
	checkErr(err4, fmt.Sprintf("Could not create file %s", filename))

	err5 := png.Encode(newImg, im)
	checkErr(err5, "Could not convert image to PNG.")
	newImg.Close()
}
