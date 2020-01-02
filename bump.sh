#!/bin/bash -ae
TAG=$(git describe --tags --abbrev=0)
MINOR=$(echo $TAG | sed -E 's|^.*\.([0-9]+)$|\1|')
MAJOR=$(echo $TAG | sed -E 's|^v(.*)\.[0-9]+$|\1|')
NEWVER="$MAJOR.$((MINOR+1))"
git tag v$NEWVER
yq w -i helm/argovue/Chart.yaml appVersion $NEWVER
git commit helm/argovue/Chart.yaml -m "Bump chart version to v$NEWVER"
make helm
git commit docs -m "Bump chart release to v$NEWVER"
git push
git push origin v$NEWVER
