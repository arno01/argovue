.PHONY: all argovue ui skaffold

all: argovue ui

argovue:
	cd src && GOOS=linux go build

ui:
	cd ui && yarn build

skaffold: ui argovue
	cp src/argovue skaffold/argovue
	rm -rf skaffold/ui
	mkdir skaffold/ui
	cp -a ui/dist skaffold/ui/

