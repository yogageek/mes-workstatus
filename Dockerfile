FROM golang:1.13-buster as build

WORKDIR /go/src/mes_workstatus
ADD . .

RUN go mod download
RUN go build -o /go/main

FROM gcr.io/distroless/base-debian10
WORKDIR /go/
COPY --from=build /go/main .
COPY .env .
#COPY docs/ docs

ENV POSTGRES_URL "host=42.159.86.191 port=5432 user=46d1a69b-6cd1-4b94-b009-537e2d575bba password=ssc8u7occfhqm3q6gkhm0gvcua dbname=ecd73592-abcd-4a8e-a7c9-26e1d5bab72c sslmode=disable"

EXPOSE 8080

CMD ["./main"]
