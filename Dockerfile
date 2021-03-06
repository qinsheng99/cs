FROM golang:latest as BUILDER

# build binary
RUN mkdir -p /go/src/openeuler/cve-sa-backend
COPY . /go/src/openeuler/cve-sa-backend
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN cd /go/src/openeuler/cve-sa-backend && go mod tidy && CGO_ENABLED=1 go build -v -o ./cve-sa-backend main.go

# copy binary config and utils
FROM openeuler/openeuler:21.03
RUN mkdir -p /opt/app/cve-sa-backend/conf/ && mkdir /opt/app/cve-sa-backend/webapp
COPY ./conf/app_test.ini /opt/app/cve-sa-backend/conf/app.ini
COPY ./webapp/manager.html /opt/app/cve-sa-backend/webapp/manager.html
COPY --from=BUILDER /go/src/openeuler/cve-sa-backend/cve-sa-backend /opt/app/cve-sa-backend
ENV TimeZone=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TimeZone /etc/localtime && echo $TimeZone > /etc/timezone
WORKDIR /opt/app/cve-sa-backend/
ENTRYPOINT ["/opt/app/cve-sa-backend/cve-sa-backend"]