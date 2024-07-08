FROM golang:1.22 AS builder
ARG VERSION_NAME=github.com/screamsoul/go-metrics-tpl/internal/versions.buildVersion \
    COMMIT_NAME=github.com/screamsoul/go-metrics-tpl/internal/versions.buildCommit \
    DATE_NAME=github.com/screamsoul/go-metrics-tpl/internal/versions.buildDate

LABEL stage=gobuilder
ENV CGO_ENABLED 0
RUN apt-get update && apt-get install tzdata
WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .


FROM builder as agent-build
ARG AGENT_BUILD_VERSION=N/A \
    AGENT_BUILD_COMMIT=N/A
    
RUN cd cmd/agent && \ 
go build \ 
-buildvcs=false \ 
-o /app/main \ 
-ldflags "-X '${VERSION_NAME}=${AGENT_BUILD_VERSION}' -X '${COMMIT_NAME}=${AGENT_BUILD_COMMIT}' -X '${DATE_NAME}=$(date +'%Y/%m/%d %H:%M:%S')'"



FROM builder as server-build
ARG SERVER_BUILD_VERSION=N/A \
    SERVER_BUILD_COMMIT=N/A
  
RUN cd cmd/server && \
go build \ 
-buildvcs=false \ 
-o /app/main \ 
-ldflags "-X '${VERSION_NAME}=${SERVER_BUILD_VERSION}' -X '${COMMIT_NAME}=${SERVER_BUILD_COMMIT}' -X '${DATE_NAME}=$(date +'%Y/%m/%d %H:%M:%S')'"


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