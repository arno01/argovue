FROM node:13.4.0-alpine3.10 as frontend
ARG VERSION
COPY ui ui
RUN cd ui && yarn install
RUN cd ui && yarn build
FROM golang:1.13-alpine as backend
COPY . /home/kubevue
RUN apk add git && \
	cd /home/kubevue && \
	export COMMIT=$(git rev-parse --short HEAD) && \
	export BUILDDATE=$(date +%Y%m%d%H%M%S) && \
	cd src && go build -ldflags="-X main.version=$VERSION -X main.builddate=$BUILDDATE -X main.commit=$COMMIT"

FROM alpine:3.10
RUN apk update
COPY --from=backend /home/kubevue/src/kubevue kubevue
COPY --from=frontend ui/dist ui/dist
CMD ["./kubevue"]
