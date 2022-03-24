FROM golang:latest as BUILDER

# build binary
RUN mkdir -p /go/src/openeuler/cve-sa-backend
COPY . /go/src/openeuler/cve-sa-backend
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN cd /go/src/openeuler/cve-sa-backend && go mod tidy && CGO_ENABLED=1 go build -v -o ./cve-sa-backend main.go

# copy binary config and utils
FROM openeuler/openeuler:21.03
RUN mkdir -p /opt/app/cve-sa-backend/conf/
COPY ./conf/app_test.ini /opt/app/cve-sa-backend/conf/app.ini
COPY --from=BUILDER /go/src/openeuler/cve-sa-backend/cve-sa-backend /opt/app/cve-sa-backend
WORKDIR /opt/app/cve-sa-backend/
ENTRYPOINT ["/opt/app/cve-sa-backend/cve-sa-backend"]