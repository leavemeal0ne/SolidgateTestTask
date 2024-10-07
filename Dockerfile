FROM golang:1.22 AS server_builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . ./

RUN go mod tidy

RUN  CGO_ENABLED=0 go build -o ./bin/main cmd/main.go

FROM scratch AS runner

COPY --from=server_builder /app/bin/main .

COPY /card_classification_data ./card_classification_data

CMD ["./main"]