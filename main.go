package main

import b64 "encoding/base64"

import (
	"fmt"
	"os/exec"
	"os"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"
	"strings"
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"bytes"
	"io"
)

type RetJSON struct {
	Error int
	Value string
	File string
	Tool_id string
	Good_id string
}




func mainAction(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	commandLine := "c:/cadv41/bin/cadlink"

	if req.Method == "GET" {
		retStr, _:=json.Marshal(RetJSON{Error: 1, Value:"Error in request method"})
		fmt.Fprintf(w, string(retStr))
		return
	} else {
		str:=req.FormValue("str")
		tool_id:=req.FormValue("tool_id")
		good_id:=req.FormValue("good_id")

		rand.Seed(time.Now().UnixNano())
		var sawFileName string = ""
		var tempFolder string = strconv.Itoa(rand.Int())
		var err error
		var path string = "d:\\temp\\optim\\"+tempFolder+"\\"
		var fileLines []string
		var exitName string
		var resSaw string

		defer os.RemoveAll(path)
		//decode file
		dec, err := b64.StdEncoding.DecodeString(str)
		if err!=nil {
			retStr, _ := json.Marshal(RetJSON{Error: 1, Value:"Error in base64 string"})
			fmt.Fprintf(w, string(retStr))
			return
		}
		fileLines=strings.Split(string(dec), "\n")
		for _, line := range (fileLines) {
			rows := strings.Split(line, ",")
			if (rows[0]=="JOBS") {
				exitName=rows[2]
				break
			}
		}
		r, err := charset.NewReader("windows-1251", bytes.NewBufferString(exitName))
		if err != nil {
			fmt.Println(err)
		}
		f, err := ioutil.ReadAll(r)
		if err != nil {
			fmt.Println(err)
		}
		sawFileName=string(f)

		//create random folder
		if err = os.Mkdir(path, 0777); err!=nil {
			retStr, _ := json.Marshal(RetJSON{Error: 1, Value:"Couldn't create temporary folder"})
			fmt.Fprintf(w, string(retStr))
			return
		}

		//create ptx file
		if err = ioutil.WriteFile(path  + "/file.ptx", dec, 0644); err!=nil {
			retStr, _:=json.Marshal(RetJSON{Error: 1, Value:"Couldn't create ptx file"})
			fmt.Fprintf(w, string(retStr))
			return
		}


		_, err = exec.Command(commandLine, path + "file.pt*").CombinedOutput()
		if err != nil {
			//			retStr, _:=json.Marshal(RetJSON{Error: 1, Value:"Couldn't start cadlink", Good_id:err.Error(), File:string(comm)})
			//			fmt.Fprint(w, string(retStr))
			//			return
		}

		//testing block
		/*
		TOTO remove on production
		 */
		err=genSaw(path, sawFileName)
		if err != nil {
			retStr, _ := json.Marshal(RetJSON{Error: 1, Value:"Couldn't gen saw file", File:""})
			fmt.Fprintf(w, string(retStr))
			return
		}

		resultf, err := ioutil.ReadFile(path + sawFileName+ ".rlt")
		if err == nil {
			retStr, _ := json.Marshal(RetJSON{Error: 1, Value:"Got rlt file, better try to repair known bugs", File:b64.StdEncoding.EncodeToString(resultf)})
			fmt.Fprintf(w, string(retStr))
			return
		}

		saw, err := ioutil.ReadFile(path + sawFileName + ".saw")
		if err == nil {
			resSaw = b64.StdEncoding.EncodeToString(saw)
			retStr, _ := json.Marshal(RetJSON{Error: 0, Value:"Success", File:string(resSaw), Tool_id:tool_id, Good_id:good_id})
			fmt.Fprintf(w, string(retStr))
			return
		}
		retStr, _ := json.Marshal(RetJSON{Error: 1, Value:"No file was found!", File:err.Error()})
		fmt.Fprintf(w, string(retStr))
		return
	}
	req.Body.Close()


}

func main() {
	http.HandleFunc("/ptxsaw", mainAction)
	http.ListenAndServe(":5684", nil)
}


func genSaw(path, sawFileName string) error {
	s, err := os.Open("d:\\temp\\optim\\02_06_2015 1.saw")
	defer s.Close()
	if err != nil {
		return err
	}

	d, err := os.Create(path+sawFileName +".saw")
	defer d.Close()
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		return err

	}

	return nil

}