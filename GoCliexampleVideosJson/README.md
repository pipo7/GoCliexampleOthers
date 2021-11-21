Refrerence : https://www.youtube.com/watch?v=CODqM_rzwtk
> First download image
$ docker pull golang:1.17-alpine3.14
> Now using Dockerfile create customimage
$ docker build --target dev . -t go

 >Run the container with mounitng current Host PWD to /work in container
$ docker run  -it -v $(pwd):/work go sh
> Verify go version inside it
/work # go version
go version go1.17.3 linux/amd64
> Note since you mounted Host PWD to container ,so now you can see folders in vscode as well
> Runthe command : 
$ docker run  -it -v $(pwd):/work go sh

> Create a folder in /work and do following steps
/work: mkdir videos
cd videos
go mod init videos  ## creates a go.mod file

> create main.go ,videos.go in videos , if any error use sudo setfacl -Rm u:ps:rwx CLIExample1/
> go build .  ## Now you have ./videos binary ready

> We create logic at in videos folder in video.go and main.go

> Finally we build the final low MB image using Dockerfile multi-stage builds . Note in final image we jsut use go binary and go SDK is not installed 
> Finall we run its container to check get and all command
docker build . -t videos

> below should return get help
docker run -it videos get --help  

## See DOCKERFILE here  Build this image as "" docker build -t videos .  "" Now you observe that final image videos is far less than previous ones
> docker run -it videos get --all  OR docker run -it videos get --id "a"   ... should work
