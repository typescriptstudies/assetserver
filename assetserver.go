package main

import(
	"os"
	"fmt"
	"strings"
	"path/filepath"
	"net/http"
	"github.com/gorilla/mux"
	"mime"
	"io"
)

var(	
	ASSETS_PATH string = os.Getenv("ASSETERVER_ASSETS_PATH")	
)

///////////////////////////////////////////////////////////////////////////

func pathfromparts(parts []string) string {
	return strings.Join(parts,string(filepath.Separator))
}

func assetpathfromparts(parts []string) string {
	return ASSETS_PATH+strings.Join(parts,string(filepath.Separator))
}

func rootpath(assetname string) string {
	return assetpathfromparts([]string{assetname})
}

func assetpath(assettype string,assetname string) string {
	return assetpathfromparts([]string{"assets",assettype,assetname})
}

func subassetpath(assettype string,assetsubtype string,assetname string) string {
	return assetpathfromparts([]string{"assets",assettype,assetsubtype,assetname})
}

///////////////////////////////////////////////////////////////////////////

func servAssets(w http.ResponseWriter, r *http.Request, path string) {
	ext := filepath.Ext(path)
	mimetype := mime.TypeByExtension(ext)	

	content, err := os.Open(path)

    if err == nil {
        defer content.Close()
	    w.Header().Set("Content-Type", mimetype)
	    io.Copy(w, content)
    } else {
    	fmt.Fprintf(w, "Not Found 404 [ info : path %v ext %v mine %v ]",path,ext,mimetype)
    }
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	servAssets(w,r,rootpath(vars["assetname"]))
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	servAssets(w,r,assetpath(vars["assettype"],vars["assetname"]))
}

func subassetsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	servAssets(w,r,subassetpath(vars["assettype"],vars["assetsubtype"],vars["assetname"]))
}

///////////////////////////////////////////////////////////////////////////

func main() {
	fmt.Printf("assetserver serving assets at %v",ASSETS_PATH)

	r := mux.NewRouter()

	r.HandleFunc("/{assetname}", rootHandler).Methods("GET")
	r.HandleFunc("/assets/{assettype}/{assetname}", assetsHandler).Methods("GET")
	r.HandleFunc("/assets/{assettype}/{assetsubtype}/{assetname}", subassetsHandler).Methods("GET")

	http.Handle("/",r)

	http.ListenAndServe(":9002", nil)
}