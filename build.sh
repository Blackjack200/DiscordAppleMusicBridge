
#!/bin/bash
hs=$(git rev-parse HEAD)

function build() {
  echo "[+] Build GOOS=${1} ARCH=${2}"
  NAME="music_${1}_${2}_$hs"
  CGO_ENABLED=1 GOOS=${1} GOARCH=${2} go build -ldflags="-w -s" -o "$NAME" main.go
  mv "$NAME" bin/"$NAME"
}
rm -rdf bin 2>/dev/null >/dev/null
mkdir bin 2>/dev/null >/dev/null
build darwin amd64 &
build darwin arm64 &
wait
