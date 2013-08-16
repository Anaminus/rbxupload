# rbxupload

Upload Roblox models and places via command-line.

## Installation

1. [Download](http://code.google.com/p/go/downloads/list) and install Go
2. Make sure your [GOPATH](http://golang.org/doc/code.html#GOPATH) environment
   variable is set (i.e. `$HOME/Go`, `%USERPROFILE%\Go`)
3. Run `go get github.com/Anaminus/rbxupload` (requires
   [git](http://golang.org/s/gogetcmd))
5. rbxupload will be available in `$GOPATH/bin`

### Binaries

None at the moment.

## Usage

rbxupload can be run with the the following flags:

* `-h, --help`

	Displays a help message.

* `-f [path], --file=[path]`

	The location of the file to be uploaded.

* `-u [string], --username=[string]`

	Username for logging in.

* `-p [string], --password=[string]`

	Password for logging in.

	**DO NOT ENTER YOUR PASSWORD INTO UNTRUSTED PROGRAMS!**

* `-t [string], --type=[string]`

	The type of file to upload. May be "Model" or "Place". (Defaults to Model)

* `-a [id], --asset=[id]`

	Asset ID to upload to. 0 creates a new asset. Places may only be updated,
	not created. (Defaults to 0)

* `-i [key]:[value], --info=[key]:[value]`

	If uploading a new model, this sets information about the model. i.e. `-i
	name:string -i description:string` sets the name and description.

* `-s, --skip`

	Skips prompts. If omitted, and a required flag is missing, you will be
	prompted to input its value.
