FROM golang:latest

WORKDIR /app

ADD pkg/ ./pkg
ADD static/ ./static
ADD template/ ./template
ADD sql/ ./sql
COPY go.mod go.sum main.go ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

EXPOSE 3000

ENTRYPOINT [ "/main" ]