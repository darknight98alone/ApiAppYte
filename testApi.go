package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Hypotheses struct {
	Confidence int    `json:"confidence"`
	Utterance  string `json:"utterance"`
}

type myjsonstruct struct {
	Status     int          `json:"status"`
	Hypotheses []Hypotheses `json:"hypotheses"`
	Id         string       `json:"id"`
}

type myjsstruct struct {
	Content string `json:"content"`
}

type myresponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func (myres *myresponse) returnerr(k int) string {
	myres.Status = k
	myres.Result = ""
	switch k {
	case 1:
		myres.Message = "Don't have sound"
	case 2:
		myres.Message = "Be destroyed"
	case 9:
		myres.Message = "Server is busy"
	case 10:
		myres.Message = "Don't have result"
	}
	b, _ := json.Marshal(myres)
	return string(b)
}

var listAPIKey []string
var startindex int

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var myres myresponse
		//fmt.Println(base64Decode(r.Header.Get("data")))
		file := r.Body

		url := "https://api.fpt.ai/hmi/asr/general"

		req, err5 := http.NewRequest("POST", url, file)
		if err5 != nil {
			fmt.Fprintf(w, myres.returnerr(2))
			return
		}
		var locs myjsonstruct

		for {
			index := startindex
			if startindex == len(listAPIKey) {
				fmt.Fprintf(w, "Api rate limited")
				return
			}
			key := listAPIKey[startindex]
			req.Header.Add("api_key", key)
			req.Header.Add("cache-control", "no-cache")
			req.Header.Add("postman-token", "8845e06b-a5d5-50e5-2738-7fee3b251d68")
			res, err6 := http.DefaultClient.Do(req)
			if res.StatusCode == 429 {
				startindex = startindex + 1
				continue
			}
			if err6 != nil { //9
				if index != len(listAPIKey)-1 {
					continue
				}
				fmt.Fprintf(w, myres.returnerr(9))
				return
			}
			defer res.Body.Close()
			body, err7 := ioutil.ReadAll(res.Body)
			if err7 != nil {
				fmt.Fprintf(w, myres.returnerr(2))
				return
			}
			json.Unmarshal(body, &locs)
			fmt.Println("USE:", index)
			break
		}

		if len(locs.Hypotheses) != 0 {
			myres.Status = 0
			myres.Message = locs.Hypotheses[0].Utterance
			myres.Result = ""
			b, err8 := json.Marshal(myres)
			if err8 != nil {
				fmt.Fprintf(w, myres.returnerr(2))
				return
			}
			w.Header().Set("content-type", "application/json")
			fmt.Fprintf(w, string(b))
		} else {
			fmt.Fprintf(w, myres.returnerr(1))
		}
	}
}

func base64Decode(str string) string {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(data)
}

func main() {
	listAPIKey = []string{"Pldlpwx4sLK5za9LJRpbbSVg2wcVCMtD", "YANs1Rb1TyXZh9uXupI0NDk24qXCIGtZ", "7dRvL89TCyEtr8Zz6WQS3Sikv123Mtba", "ktmLa3q7wqsNgZX2PSvDI3IVxY5zXBuX", "cg1OjiLWYCI9G6GVMvduKl5jA0BwyyPo", "Flj3tR4sCWMlG7AHgN0yyBe6nWjQp94h", "QMltUM5ulCxYjfvomNgJhSieeEBpbQhB", "3GRw9uuRbMmQIBQra98MVZ7CJHaSdEYS", "WBme28eANXFEdvFCCXHC868LPZrOIX05", "Im4Qwznz450eKCXTp6p4srUYVHPCgr4D", "dMEk9pcQgD74FkVTM91tPwQYj3KGI8Qu", "86IF4IQejkU5d1je8gjyPRo3QCvhleD9", "887gP6n2JJj6a4GaX5NELczEjlIZNWbw", "D5IhpTwv8ckoNtfMdp7vlJP4KUE7OQwh", "4AOFLyJsFhadvXYboJCMmk6nWMEOj4WI", "GoxclRGQoBrAbt2qoOZlOziQrWx9gIeP", "t5ecZv4lMtrAfdtTZTDj0vzNQfgjS4dW", "69NsuGtzRPodRegNT9Yedjee6tLvYNzI", "U0hsq2cyefjaDgRR3RFCGfzGoGdwCMPJ", "gPbZTLJ5L77GyYmbfBw7pXDL5ILPV9Ae"}
	startindex = 0
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":9091", nil)
}
