FROM node:13.4.0-alpine3.10 as frontend
COPY . kubevue
RUN cd kubevue/ui && yarn install
RUN cd kubevue/ui && yarn build
RUN apk add git && \
	cd kubevue && \
	export VERSION=$(git describe --tags) && sed -i "s/_VERSION_/$VERSION/" ui/dist/config.js && \
	export COMMIT=$(git rev-parse --short HEAD) && sed -i "s/_COMMIT_/$COMMIT/" ui/dist/config.js && \
	export BUILDDATE=$(date +%Y%m%d%H%M%S) && sed -i "s/_BUILDDATE_/$BUILDDATE/" ui/dist/config.js

FROM golang:1.13-alpine as backend
COPY . /home/kubevue
RUN apk add git && \
	cd /home/kubevue && \
	export VERSION=$(git describe --tags) && \
	export COMMIT=$(git rev-parse --short HEAD) && \
	export BUILDDATE=$(date +%Y%m%d%H%M%S) && \
	cd src && go build -ldflags="-X main.version=$VERSION -X main.builddate=$BUILDDATE -X main.commit=$COMMIT"

FROM alpine:3.10
RUN apk update
COPY --from=backend /home/kubevue/src/kubevue kubevue
COPY --from=frontend kubevue/ui/dist ui/dist
CMD ["./kubevue"]
