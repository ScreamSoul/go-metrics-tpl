FROM golang:1.22 AS builder
LABEL stage=gobuilder
ENV CGO_ENABLED 0
RUN apt-get update && apt-get install tzdata
WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .


FROM builder as agent-build
RUN cd cmd/agent && go build -buildvcs=false  -o /app/main


FROM builder as server-build
RUN cd cmd/server && go build -buildvcs=false  -o /app/main


FROM scratch as app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Russia/Moskow
WORKDIR /app


FROM app as agent-app
COPY --from=agent-build /app/main /app/main
CMD ["./main"]


FROM app as server-app
COPY --from=server-build /app/main /app/main
CMD ["./main"]