FROM golang:1.14-alpine AS backend
RUN apk add build-base g++ make protoc git protobuf-dev
WORKDIR /build
RUN git clone https://github.com/grpc/grpc-web.git
WORKDIR grpc-web
RUN make install-plugin
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make

FROM node:14-alpine AS frontend
RUN apk add build-base autoconf automake pngquant bash
WORKDIR /build
COPY frontend/yarn.lock .
COPY frontend/package.json .
COPY frontend/. .
RUN yarn install
RUN yarn run build

LABEL maintainer="Leigh MacDonald <leigh.macdonald@gmail.com>"
WORKDIR /app
COPY --from=backend /build/verimapcom .
COPY --from=backend /build/verimapcom/templates .
COPY --from=frontend /build/dist .
VOLUME uploads /app/uploads
EXPOSE 9090
EXPOSE 8001
CMD ["./verimapcom", "serve"]
