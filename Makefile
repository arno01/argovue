.PHONY: all kubevue ui skaffold

all: kubevue ui

kubevue:
	cd src && GOOS=linux go build

ui:
	cd ui && yarn build

skaffold: ui kubevue
	cp src/kubevue skaffold/kubevue
	rm -rf skaffold/ui
	mkdir skaffold/ui
	cp -a ui/dist skaffold/ui/

