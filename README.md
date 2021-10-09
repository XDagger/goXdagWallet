# Cross Platform XDAG GUI Wallet

This is a cross-platform XDAG GUI wallet, especially for macOS and Linux, powered by [fyne](https://github.com/fyne-io/fyne).

Fyne is a cross-platform GUI in Go inspired by Material Design.

## repo structure
 - clib - a wrapper of XDAG Wallet C library
   - xDagWallet - XDAG wallet C library
 - wallet - golang XDAG wallet app 
   - i18n - international strings
   - data - i18n config json, images, fonts
   - component - ui of wallet window
   - config - wallet config
   - wallet_state - wallet state
   - xlog - wallet log

## environment variable

CGO_ENABLED=1

FYNE_FONT=data/myFont.ttf

FYNE_SCALE=1.2

## features

 - frequent transfer address 