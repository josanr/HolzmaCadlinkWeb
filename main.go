package main

import b64 "encoding/base64"
import "path/filepath"

import (
	"fmt"
	"os/exec"
	"os"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

type RetJSON struct {
	Error int
	Value string
	File string
	Tool_id string
	Good_id string
}


func hello(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	path:="d:/temp/"
	commandLine:="d:/cadv41/bin/cadlink"

	if req.Method == "GET" {
		retStr, _:=json.Marshal(RetJSON{Error: 1, Value:"Error in request method"})
		fmt.Fprintf(w, string(retStr))
		return
	} else {

		str:=req.FormValue("str")
		sawFileName:=req.FormValue("name")
		tool_id:=req.FormValue("tool_id")
		good_id:=req.FormValue("good_id")
		list, err:=filepath.Glob(path + "*.ptx");
		if err==nil{
			for x:=0; x<len(list);x=x+1{
				os.Remove(list[x])
			}
		}
		list, err=filepath.Glob(path + "*.saw");
		if err==nil{
			for x:=0; x<len(list);x=x+1{
				os.Remove(list[x])
			}
		}
		list, err=filepath.Glob(path + "*.rlt");
		if err==nil{
			for x:=0; x<len(list);x=x+1{
				os.Remove(list[x])
			}
		}

		dec, err:= b64.StdEncoding.DecodeString(str)
		if err!=nil{
			retStr, _:=json.Marshal(RetJSON{Error: 1, Value:"Error in base64 string"})
			fmt.Fprintf(w, string(retStr))

			return
		}

		err = ioutil.WriteFile(path + "file.ptx", dec, 0644)
		if err!=nil{
			retStr, _:=json.Marshal(RetJSON{Error: 1, Value:"Couldn't create ptx file"})
			fmt.Fprintf(w, string(retStr))
			os.Remove(path + "file.ptx")
			return
		}



		comm, err := exec.Command(commandLine, path + "file.pt*").CombinedOutput()
		if err != nil {

			retStr, _:=json.Marshal(RetJSON{Error: 1, Value:"Couldn't start cadlink", Good_id:err.Error(), File:string(comm)})
			fmt.Fprintf(w, string(retStr))
			return
		}

		//		fmt.Println(path +  sawFileName + ".saw")
		resSaw:=""
		saw, err := ioutil.ReadFile(path + sawFileName + ".saw")
		if err == nil {
			resSaw = b64.StdEncoding.EncodeToString(saw)
		}

		resultf, err := ioutil.ReadFile(path + "file.rlt")
		if err != nil {
			retStr, _:=json.Marshal(RetJSON{Error: 1, Value:"Couldn't read created rlt file"})
			fmt.Fprintf(w, string(retStr))
			os.Remove(path + "file.rlt")
			os.Remove(path + sawFileName + ".saw")
			return
		}



		retStr, _:=json.Marshal(RetJSON{Error: 0, Value:string(resultf), File:string(resSaw), Tool_id:tool_id, Good_id:good_id})
		fmt.Fprintf(w, string(retStr))
		os.Remove(path + "file.ptx")
		os.Remove(path + "file.rlt")
		os.Remove(path + sawFileName + ".saw")
	}
	req.Body.Close()

}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8000", nil)
}
