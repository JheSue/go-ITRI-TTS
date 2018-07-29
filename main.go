package main

import (
	"os"
	"io"
	"net/http"
	"strings"
	"io/ioutil"
	"fmt"
	"encoding/xml"
	"github.com/aerth/playwav"
)

type Result struct {
	XMLName xml.Name `xml:"Envelope"`
	Result string `xml:"Body>ConvertSimpleResponse>Result"`
}

type Resultwav struct{
	XMLName xml.Name `xml:"Envelope"`
	Resultwav string `xml:"Body>GetConvertStatusResponse>Result"`
}

func main(){
	accountId:=os.Args[1]
	password:=os.Args[2]
	ttstext:=os.Args[3]
	reqdata :=`<?xml version="1.0" encoding="utf-8"?>
			<soap12:Envelope 
				xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
				xmlns:xsd="http://www.w3.org/2001/XMLSchema"
				xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
				<soap12:Body>
					<ConvertSimple xmlns="http://tts.itri.org.tw/">
						<accountID>`+accountId+`</accountID>
						<password>`+password+`</password>
						<TTStext>`+ttstext+`</TTStext>
					</ConvertSimple>
				</soap12:Body>
			</soap12:Envelope>`
	responsedata:=soapTts(reqdata)
	var result Result
	err2 := xml.Unmarshal(responsedata, &result)
	if err2 != nil {
		fmt.Println("error messages:", err2)
	}
	id := strings.Split(result.Result , "&")

	reqdata = `<?xml version="1.0" encoding="utf-8"?>
			<soap12:Envelope
				xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
				xmlns:xsd="http://www.w3.org/2001/XMLSchema"
				xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
				<soap12:Body>
					<GetConvertStatus xmlns="http://tts.itri.org.tw/">
						<accountID>`+accountId+`</accountID>
						<password>`+password+`</password>
						<convertID>`+id[2]+`</convertID>
					</GetConvertStatus>
				</soap12:Body>
			</soap12:Envelope>`
	responsedata = soapTts(reqdata)
	var resultwav Resultwav
	err := xml.Unmarshal(responsedata, &resultwav)
	if err != nil {
		fmt.Println("error messages:", err)
	}
	source_url := strings.Split(resultwav.Resultwav , "&")
	valid := true
	for valid {
		responsedata = soapTts(reqdata)
		err := xml.Unmarshal(responsedata, &resultwav)
		if err != nil {
			fmt.Println("error messages:", err)
		}
		source_url := strings.Split(resultwav.Resultwav , "&")
		if source_url[3] == "completed"{
			valid=false
		} 
	}
	source_url = strings.Split(resultwav.Resultwav , "&")
	wav, err := http.Get(source_url[4])
	out, err :=os.Create("TTS.wav")
	_,err =io.Copy(out, wav.Body) 
	playwav.FromFile("TTS.wav")
	
}

func soapTts(reqdata string) []byte {
	res, err :=http.Post("http://tts.itri.org.tw/TTSService/Soap_1_3.php?wsdl","application/soap+xml; charset=UTF-8",strings.NewReader(reqdata))
        if err != nil {
                fmt.Println("http post err",err)
		os.Exit(1)
        }
        if http.StatusOK !=res.StatusCode{
                fmt.Println("request failed", res.StatusCode)
		os.Exit(res.StatusCode)
        }

        responsedata, err := ioutil.ReadAll(res.Body)
        if err != nil{
                fmt.Println(err)
		os.Exit(2)
        }

        //fmt.Println(string(responsedata))
        return responsedata
}

