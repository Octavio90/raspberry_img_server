package main

import (
	"os"
	"log"
	"fmt"
	"flag"
	"os/exec"
	"net/http"
	"html/template"
)
 
var port  = "8080"
var imgs  = "imgs"
var fname = fmt.Sprintf("./%v/test.jpg",imgs)
var args  = []string{"-w","1640","-h","1232","-o",fname}

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

func view(w http.ResponseWriter, r *http.Request) {
	t,err := template.ParseFiles("index.html")
	if err != nil { log.Panic(err) }
	t.Execute(w,"Welcome")
}

func capture(w http.ResponseWriter, r *http.Request){
	pixel_number := r.FormValue("pixel")
	if _, err    := os.Stat(fname); err == nil { 
		//fmt.Println("Removing file",fname) 
		os.Remove(fname)
	} 
	//fmt.Println("Taking picture on ",fname)
	cmd     := exec.Command("raspistill", args...)
	if err  := cmd.Start(); err != nil {log.Panic(err)}
	
	err     := cmd.Wait()
	if err  != nil { log.Fatal(err) }
	draw_line(pixel_number)
	http.Redirect(w,r,"/view",http.StatusSeeOther)
}


func draw_line(){
	file, err := os.Open(fname)
	if err != nil {log.Panic(err)}
	defer existingImageFile.Close()
}

