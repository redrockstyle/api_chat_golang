FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY ./ /app

# install go dependencies
RUN go mod download

# build api
RUN go build -o main cmd/main/main.go

FROM golang:1.22-alpine

WORKDIR /app
COPY --from=builder ./app ./

EXPOSE 5000

# RUN apt install openssl
# WORKDIR /app/certs
# RUN chmod +x gen_localhost.sh
# RUN bash gen_localhost.sh

CMD [ "./main" ]
# ENTRYPOINT [ "./main" ]