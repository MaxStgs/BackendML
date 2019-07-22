package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

var ip = flag.String("ip", "50341", "Ip for route")

func main() {
	flag.Parse()
	http.HandleFunc("/api/acceptFile", handleAcceptFile)
	http.HandleFunc("/", handleIndex)

	err := http.ListenAndServe(":"+*ip, nil)
	if err != nil {
		panic(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func setupResponse(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func handleAcceptFile(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	err := r.ParseMultipartForm(math.MaxInt16)
	//err := r.ParseForm()
	if err != nil {
		fmt.Printf("handleAcceptFile() err: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	remoteName := "NV_" + strings.Replace(r.RemoteAddr, ":", "_", -1)
	result, err := isFolderExists(remoteName)
	err = os.Chdir("data/")
	if err != nil {
		fmt.Printf("Error: %v, folder 'data' was not found, continue...\n", err)
	}
	// Is new user?
	if !result {
		err = os.Mkdir(remoteName, os.FileMode(0777))
	}
	err = os.Chdir(remoteName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Create new folder with data
	folderName := "NV_" + strings.Replace(time.Now().Format("02-Jan-2006 15:04:05"), ":", "_", -1)
	err = os.Mkdir(folderName, os.FileMode(0777))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	err = os.Chdir(folderName)

	// Write data
	err = ioutil.WriteFile("data.json", []byte(r.FormValue("data")), os.FileMode(0777))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Write photos
	err = os.Mkdir("images", os.FileMode(0777))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	err = os.Chdir("images")

	for _, v := range r.MultipartForm.File {
		for _, v1 := range v {
			file, err := v1.Open()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			defer file.Close()
			filedata, err := ioutil.ReadAll(file)
			data := make([]byte, 0)
			//count, err := file.Read(data)
			count, err := bufio.NewReader(file).Read(data)
			fmt.Printf("%s Readed bytes: %d\n", v1.Filename, count)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			err = ioutil.WriteFile(v1.Filename, filedata, os.FileMode(0777))
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}

	// Go back
	err = os.Chdir("../../../../")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func isFolderExists(folderName string) (result bool, err error) {
	arrInfo, err := ioutil.ReadDir("")
	if err != nil {
		return false, err
	}
	for _, v := range arrInfo {
		if v.IsDir() && v.Name() == folderName {
			return true, nil
		}
	}
	return false, nil
}
