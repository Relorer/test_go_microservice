FROM alpine:3.14 as builder  

RUN apk update && apk add ca-certificates tzdata  

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/  
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo/  

COPY bin/main /
COPY config.yaml /
CMD ["/main"]