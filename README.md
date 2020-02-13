# go-console-tetris
Test project to get familiar with the basic principles of Go.

## Debugging with Visual Studio Code
Since the Visual Studio Code Go debugger plugin does only support the debug console, while termbox needs to read/write from/to /dev/tty, the only quick solution to debug the project with Visual Studio Code, I came up with, is to start a debug server with Delve.
* Install delve (https://github.com/go-delve/delve)
* Start debug server: 
```sh
$ ~/go/bin/dlv debug --headless --listen :2345 ~/go/src/github.com/gkuehn001/go-console-tetris
```
* Add the following configuration to your launch.json in VS Code:
```sh
{
    "name": "Connect to server",
    "type": "go",
    "request": "attach",
    "mode": "remote",
    "remotePath": "${workspaceFolder}",
    "port": 2345,
    "host": "127.0.0.1",
    "apiVersion": 1
}
```
* Connect to the delve server with VS Code and start debugging

## Tetromino Roation
This version of the game uses the Tetris Super Rotation System (current Tetris Guideline standard); wall kicks not yet included
(https://tetris.fandom.com/wiki/SRS)
