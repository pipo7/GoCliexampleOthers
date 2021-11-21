# Using dev as initial development 
FROM golang:1.17-alpine3.14 as dev
WORKDIR /work

# Building the image with  videos binary using docker multi-stage . here building videos binary image with go lang SDK
FROM golang:1.17-alpine3.14 as build 
WORKDIR /videos
COPY ./videos/* /videos/
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o videos

# Now building runtime even smaller image and just using the binary from build stage as executable
FROM alpine as runtime
COPY --from=build /videos/videos /usr/local/bin/videos
COPY ./videos/videos.json /
# Check the contents of run.sh in folder
RUN cd ..
COPY run.sh /
RUN chmod +x run.sh
ENTRYPOINT ["./run.sh" ]
## Build this image as "" docker build -t videos .  "" Now you observe that final image videos is far less than previous ones
