FROM golang:1.22

WORKDIR /app

COPY . /app/

RUN go mod download
RUN go build 

EXPOSE 8000

ENTRYPOINT [ "/app/rssagg" ]
