#!/bin/bash
cd frontend
npm run build
mv dist/ ..
cd ..
rm -rf backend/dist
mv dist/ backend/
cd backend
env GOOS=windows GOARCH=amd64 go build -o build/upcm.exe .
env GOOS=linux GOARCH=amd64 go build -o build/upcm .
mv build ..
rm -rf dist
mkdir dist
touch dist/.gitkeep
cd ..
