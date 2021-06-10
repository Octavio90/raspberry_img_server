package main

import (
	"os"
	"log"
	"fmt"
	"flag"
	"image"
	"strconv"
	"os/exec"
	"net/http"
	"image/jpeg"
	"image/color"
	"html/template"
)

var first  = true
var width  = 1640
var height = 1232
var port   = "8080"
var imgs   = "imgs"
var green  = color.RGBA{144, 209, 107, 1}
var fname  = fmt.Sprintf("./%v/test.jpg",imgs)
var args   = []string{"-w","1640","-h","1232","-o",fname}

type Changeable interface { Set(x, y int, c color.Color) }
type Img struct {
    image.Image
    custom map[image.Point]color.Color
}

func NewImg(img image.Image) *Img { return &Img{img, map[image.Point]color.Color{}} }
func (m *Img) Set(x, y int, c color.Color) { m.custom[image.Point{x, y}] = c }
func (m *Img) At(x, y int) color.Color {
    if c := m.custom[image.Point{x, y}]; c != nil {return c }
    return m.Image.At(x, y)
}


func main(){
	portFlag := flag.String("port", "", "Specify server port")
	flag.Parse()

	if *portFlag != ""  { port   = *portFlag }
	fmt.Println(fmt.Sprintf("Serving on localhost:%v/view",port))

	http.Handle(fmt.Sprintf("/%v/",imgs), http.StripPrefix(fmt.Sprintf("/%v/",imgs), http.FileServer(http.Dir(fmt.Sprintf("./%v",imgs)))))
	http.HandleFunc("/view" , view )
	http.HandleFunc("/capture", capture)
	if err := http.ListenAndServe(":"+port, nil); err != nil { log.Panic(err) }
}

func view(w http.ResponseWriter, r *http.Request){
	if first{
		if _, err := os.Stat(fname); err == nil { os.Remove(fname) }
		first = false
	}
	t,err := template.ParseFiles("index.html")
	if err != nil { log.Panic(err) }
	t.Execute(w,"Align Server")
}

func capture(w http.ResponseWriter, r *http.Request){
	vpixel := r.FormValue("vpixel")
	hpixel := r.FormValue("hpixel")

	vlimit, verr := strconv.Atoi(vpixel)
	if verr != nil { vlimit = 0 }

	hlimit, herr := strconv.Atoi(hpixel)
	if herr != nil { hlimit = 0}

	if _, err     := os.Stat(fname); err == nil { os.Remove(fname) }
	cmd     := exec.Command("raspistill", args...)

	if err  := cmd.Start(); err != nil {log.Panic(err)}
	err     := cmd.Wait()
	if err  != nil { log.Fatal(err) }

	draw_line(vlimit,"height")
	draw_line(hlimit,"width")
	http.Redirect(w,r,"/view",http.StatusSeeOther)
}


func draw_line(pixel_number int, mode string){
	fmt.Println("Drawing line at",pixel_number,"mode:",mode)
	img, err := os.Open(fname)
	if err   != nil { log.Panic(err) }
	defer img.Close()

	jpeg_img, err := jpeg.Decode(img)
	if err  != nil { log.Panic(err) }

	ofile, err := os.Create(fname)
	if err != nil { log.Fatal(err) }
	defer ofile.Close()

	n_img := NewImg(jpeg_img)
        if mode == "height" {
		for x := 0; x < height; x += 1 {
			n_img.Set(pixel_number-2, x, green)
			n_img.Set(pixel_number-1, x, green)
			n_img.Set(pixel_number  , x, green)
			n_img.Set(pixel_number+1, x, green)
			n_img.Set(pixel_number+2, x, green)
		}
	}else{
		for x := 0; x < width; x += 1{
			n_img.Set(x, pixel_number-2, green)
			n_img.Set(x, pixel_number-1, green)
			n_img.Set(x, pixel_number  , green)
			n_img.Set(x, pixel_number+1, green)
			n_img.Set(x, pixel_number+2, green)
		}
	}
	jpeg.Encode(ofile, n_img, nil)
}

