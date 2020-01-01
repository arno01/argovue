.PHONY: all argovue ui skaffold helm

all: argovue ui

argovue:
	cd src && GOOS=linux go build

ui:
	cd ui && yarn build

helm:
	helm package helm/argovue -d docs
	helm repo index docs --url https://jamhed.github.io/argovue/

skaffold: ui argovue
	cp src/argovue skaffold/argovue
	rm -rf skaffold/ui
	mkdir skaffold/ui
	cp -a ui/dist skaffold/ui/

