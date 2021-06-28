# Cross Platform XDAG GUI Wallet

This is a cross-platform XDAG GUI wallet, especially for macOS and Linux, powered by [fyne](https://github.com/fyne-io/fyne).

Fyne is a cross-platform GUI in Go inspired by Material Design.

##repo structure
 - xDagWallet - XDAG wallet C library
 - src - XDAG wallet C runtime wrapper
 - lib - XDAG wallet C runtime static library, built by CmakeLists.txt
 - data - i18n config json, images, fonts
 - component - tab pages ui of main window

##environment variable
CGO_ENABLED=1;FYNE_FONT=data/msyh.ttf;FYNE_SCALE=1.2