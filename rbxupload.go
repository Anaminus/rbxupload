package main

import (
	"bufio"
	"encoding/json"
	flags "github.com/jessevdk/go-flags"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	path "path/filepath"
	"strconv"
	"strings"
)

// https://www.roblox.com/Services/Secure/LoginService.asmx/ValidateLogin
// {"userName":"<username>","password":"<password>","isCaptchaOn":false,"challenge":"<challenge id>","captchaResponse":"<user input>"}
// https://www.google.com/recaptcha/api/challenge?k=6Lc9HdsSAAAAAI8CUt1wkYPL8nZaWYkn9fLe3ApF

const (
	BaseURL   = `www.roblox.com`
	LoginURL  = `https://` + BaseURL + `/Services/Secure/LoginService.asmx/ValidateLogin`
	UploadURL = `http://` + BaseURL + `/Data/Upload.ashx`
)

func main() {

	var opts struct {
		Help     bool              `short:"h" long:"help"     optional:"true" description:"Shows this help message."`
		File     string            `short:"f" long:"file"     value-name:"[path]" description:"The location of the file to be uploaded."`
		Username string            `short:"u" long:"username" value-name:"[string]" optional:"false" description:"Username for logging in."`
		Password string            `short:"p" long:"password" value-name:"[string]" optional:"false" description:"Password for logging in."`
		Type     string            `short:"t" long:"type"     value-name:"[string]" optional:"true" default:"Model" description:"The type of file to upload. May be Model or Place."`
		Asset    uint              `short:"a" long:"asset"    value-name:"[id]" optional:"true" default:"0" description:"Asset ID to upload to. 0 creates a new asset. Places may only be updated, not created."`
		Info     map[string]string `short:"i" long:"info"     value-name:"[key]:[value]" optional:"true" description:"If uploading a new model, this sets information about the model. i.e. '-i name:string -i description:string' sets the name and description."`
		Skip     bool              `short:"s" long:"skip"     optional:"true" description:"Skip prompts."`
	}

	parser := flags.NewParser(&opts, flags.PrintErrors)
	parser.Parse()

	if opts.Help {
		parser.WriteHelp(os.Stderr)
		return
	}

	read := bufio.NewReader(os.Stdin)

	if opts.File == "" {
		if opts.Skip {
			print("File (-f) required.\n")
			return
		}
		print("File: ")
		file, _, _ := read.ReadLine()
		opts.File = string(file)
	}

	file, err := os.Open(opts.File)
	if err != nil {
		print("Invalid file.\n")
		return
	}
	defer file.Close()

	if opts.Username == "" {
		if opts.Skip {
			print("Username (-u) required.\n")
			return
		}
		print("Username: ")
		username, _, _ := read.ReadLine()
		opts.Username = string(username)
	}

	if opts.Password == "" {
		if opts.Skip {
			print("Password (-p) required.\n")
			return
		}
		print("\nIn order to upload assets, you must log in to a valid Roblox account.\nDO NOT ENTER YOUR PASSWORD INTO UNTRUSTED PROGRAMS!\nEnter a blank password to cancel.\nPassword: ")
		password, _, _ := read.ReadLine()
		opts.Password = string(password)
		if opts.Password == "" {
			print("Operation canceled.\n")
			return
		}
	}

	client := http.DefaultClient
	client.Jar, _ = cookiejar.New(&cookiejar.Options{})

	// Login
	print("\nLogging in...\n")

	req, _ := http.NewRequest("POST", LoginURL, strings.NewReader(`{"userName":"`+opts.Username+`","password":"`+opts.Password+`","isCaptchaOn":false,"challenge":"","captchaResponse":""}`))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		print(err, "\n")
		return
	}
	if resp.StatusCode != 200 {
		print("Login failed. Status code: ", resp.StatusCode, "\n")
		return
	}

	// Check response data
	// {"d":{"sl_translate":"Message","IsValid":true,"Message":"","ErrorCode":""}}
	respData := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		print("Login failed. JSON decode failed. ", err.Error(), "\n")
		return
	}
	d := respData["d"].(map[string]interface{})
	if !d["IsValid"].(bool) {

		print("Login failed. Error code ", d["ErrorCode"].(string), ": \"", d["Message"].(string), "\"\n")
		return
	}
	resp.Body.Close()
	print("Login succeeded.\n")

	// Upload file
	base := path.Base(opts.File)
	params := url.Values{
		"assetid":       {strconv.FormatUint(uint64(opts.Asset), 10)},
		"type":          {opts.Type},
		"name":          {base[:len(base)-len(path.Ext(base))]},
		"description":   {""},
		"genreTypeId":   {"1"},
		"isPublic":      {"False"},
		"allowComments": {"False"},
	}
	for k, v := range opts.Info {
		params.Set(k, v)
	}

	req, _ = http.NewRequest("POST", UploadURL+"?"+params.Encode(), file)
	req.Header.Set("User-Agent", "Roblox/WinInet")
	if stat, err := file.Stat(); err == nil {
		req.ContentLength = stat.Size()
	}
	//	req.TransferEncoding = []string{"identity"}
	//	req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))

	print("Uploading...\n")
	resp, err = client.Do(req)
	if err != nil {
		print(err.Error(), "\n")
		return
	}
	if resp.StatusCode != 200 {
		print("Upload failed. Status Code: ", resp.StatusCode, "\n")
		return
	}
	defer resp.Body.Close()
	return
}
