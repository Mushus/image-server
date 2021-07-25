FROM golang:1.16 AS builder

RUN apt-get update -y
RUN apt-get install -y wget dpkg-dev pkg-config libglib2.0-dev libexpat1-dev libjpeg62-turbo-dev libexif-dev libgif-dev libgsf-1-dev libtiff5-dev libfftw3-dev libpng-dev libimagequant-dev imagemagick liborc-0.4-dev libwebp-dev openexr libheif-dev ca-certificates
RUN wget -q https://github.com/libvips/libvips/releases/download/v8.11.2/vips-8.11.2.tar.gz && \
    tar xf vips-8.11.2.tar.gz  && \
    cd vips-8.11.2 && \
    ./configure && \
    make && \
    make install && \
    ldconfig && \
    rm -rf ./vips-8.11.2

WORKDIR /workspace

COPY go.mod .
COPY go.sum .
RUN ls -al && go mod download
COPY . .

RUN go build -o ./image-server ./cmd/app/main.go

FROM debian:buster-slim as runner

RUN apt-get update -y && \
    apt-get install -y wget dpkg-dev pkg-config libglib2.0-dev libexpat1-dev libjpeg62-turbo-dev libexif-dev libgif-dev libgsf-1-dev libtiff5-dev libfftw3-dev libpng-dev libimagequant-dev imagemagick liborc-0.4-dev libwebp-dev openexr libheif-dev ca-certificates && \
    wget -q https://github.com/libvips/libvips/releases/download/v8.11.2/vips-8.11.2.tar.gz && \
    tar xf vips-8.11.2.tar.gz  && \
    cd vips-8.11.2 && \
    ./configure && \
    make && \
    make install && \
    ldconfig && \
    rm -rf ./vips-8.11.2

WORKDIR /root/
COPY --from=builder /workspace/image-server image-server

CMD [ "/root/image-server" ]