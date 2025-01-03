FROM golang:1.22 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://proxy.golang.org,direct make build


FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        curl \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app
COPY --from=builder /src/configs /app/conf

WORKDIR /app

EXPOSE 8008
EXPOSE 9008

CMD ["./nft_transfer", "-conf", "/app/conf"]