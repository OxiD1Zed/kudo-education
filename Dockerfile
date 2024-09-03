FROM golang:1.23

WORKDIR /app

COPY . .

RUN go build -o kudo-go-app cmd/main.go

CMD [ "./kudo-go-app" ]
