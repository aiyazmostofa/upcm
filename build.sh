#!/bin/bash
rm -rf build
cd frontend
npm run build
mv dist/ ..
cd ..
rm -rf backend/dist
mv dist/ backend/
cd backend
env GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 go build -o build/upcm.exe .
env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o build/upcm .
mv build ..
rm -rf dist
mkdir dist
touch dist/.gitkeep
cd ..
