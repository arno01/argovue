#!/bin/bash -ae
TAG=$(git describe --tags --abbrev=0)
MINOR=$(echo $TAG | sed -E 's|^.*\.([0-9]+)$|\1|')
MAJOR=$(echo $TAG | sed -E 's|^v(.*)\.[0-9]+$|\1|')
NEWVER="$MAJOR.$((MINOR+1))"
git tag v$NEWVER

CHART=$(yq r helm/argovue/Chart.yaml version)
CHARTMINOR=$(echo $TAG | sed -E 's|^.*\.([0-9]+)$|\1|')
CHARTMAJOR=$(echo $TAG | sed -E 's|(.*)\.[0-9]+$|\1|')
NEWCHART="$CHARTMAJOR.$((CHARTMINOR+1))"
yq w -i helm/argovue/Chart.yaml appVersion $NEWVER
yq w -i helm/argovue/Chart.yaml version $NEWCHART
git commit helm/argovue/Chart.yaml -m "Bump chart version to app:v$NEWVER chart:$NEWCHART"
make helm
git commit docs -m "Bump chart release to v$NEWVER"
git push --tags
git push
