#!/bin/bash
version=$1
platforms=("darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm" "linux/arm64" "windows/amd64")

for platform in "${platforms[@]}"; do
  IFS='/' read -ra parts <<< "$platform"
  GOOS=${parts[0]}
  GOARCH=${parts[1]}

  output_name="static-exporter-${GOOS}-${GOARCH}"
  if [ $GOOS = "windows" ]; then
      output_name+='.exe'
  fi

  echo "Building release/$output_name"
  env GOOS=$GOOS GOARCH=$GOARCH go build -o release/$output_name cmd/main.go
  if [ $? -ne 0 ]; then
      echo 'An error has occurred! Aborting.'
      exit 1
  fi
done