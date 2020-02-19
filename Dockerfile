FROM golang:1.12 AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    BINARY_NAME=gowebdis
RUN apt -y update && apt -y install git

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ${BINARY_NAME} .

WORKDIR /dist
RUN cp /build/${BINARY_NAME} ./${BINARY_NAME}
# RUN ldd gowebdis | tr -s '[:blank:]' '\n' | grep '^/' | \
#     xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;'
# RUN mkdir -p lib64 && cp /lib64/ld-linux-x86-64.so.2 lib64/
# RUN chown -R 1000:1000 /app

FROM scratch
COPY --chown=65534 --from=builder /dist /
USER 65534
ENV PORT 8080
ENV SENTINEL_ADDRESS localhost:26379
ENV MASTER_NAME mymaster
EXPOSE 8080
ENTRYPOINT [ "/gowebdis", "start" ]
