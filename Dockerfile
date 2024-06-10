# syntax=docker/dockerfile:1

FROM golang:1.22.4-bookworm

WORKDIR /me/xboxbedrock/blockserver

RUN apt-get update

RUN apt-get install -y fontconfig

RUN apt-get install -y libvips-dev

RUN apt-get install -y libvips

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN cp /me/xboxbedrock/blockserver/blocks/minecraft.ttf /usr/local/share/fonts/minecraft.ttf

RUN chmod 644 /usr/local/share/fonts/minecraft.ttf

RUN CGO_ENABLED=1 GOOS=linux go build -o /blocksrv

EXPOSE 8000

CMD ["/blocksrv"]