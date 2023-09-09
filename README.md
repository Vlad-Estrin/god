# What is it?
god (go download) is an wget-like application that allows downloading resources by URIs concurrently

# How to use:
1) Build
   - example: go build -o god ./cmd/main.go
2) Run (with added resources separated by space)
   - example: ./god https://update.flipperzero.one/builds/qFlipper/1.3.2/qFlipper-1.3.2.dmg https://update.flipperzero.one/builds/qFlipper/1.3.0-rc2/qFlipper-x86_64-1.3.0-rc2.AppImage
