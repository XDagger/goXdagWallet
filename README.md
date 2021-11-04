# Cross Platform XDAG GUI Wallet

This is a cross-platform XDAG GUI wallet, especially for macOS and Linux, powered by [fyne](https://github.com/fyne-io/fyne).

Fyne is a cross-platform GUI in Go inspired by Material Design.

The wallet can run on Windows, Linux, Mac now.

Mobile version is coming.
## repo structure
 - clib - a wrapper of XDAG Wallet C library
   - xDagWallet - XDAG wallet C library
 - wallet - golang XDAG wallet app 
   - i18n - international strings
   - data - i18n config json, fonts
   - images - image and icon bundled in components/resource.go
   - component - ui of wallet window
   - config - wallet config
   - wallet_state - wallet state
   - xlog - wallet log

## build
enter /clib

build runtime library with CMakeLists.txt

need MingW64 in Windows

enter /wallet

`go mod tidy`

`CGO_ENABLED=1 go build`

in Windows

`CGO_ENABLED=1 go build -ldflags -H=windowsgui`

in Mac

if library not found , go to /wallet/components/wallet_cgo.go

change the library path in code: 

`//#cgo darwin LDFLAGS: ....`

## features

 - frequently transferred addresses list 
 - history query filter and pagination