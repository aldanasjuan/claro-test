FROM golang:alpine as build
RUN apk --update add ca-certificates
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server . 

FROM scratch
COPY --from=build /app/server /
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# RUN mkdir /app/certs
ENTRYPOINT [ "./server" ]