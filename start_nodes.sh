#!/bin/bash

#xterm --hold -e "cd bootnode; echo yeye"

if [[ "$OSTYPE" =~ "linux" ]]; then
    NEW_TERM=""
    if [ $TERM = "konsole" ]; then
        NEW_TERM="konsole --noclose -e"
    elif [ $TERM = "xterm" ]; then
        NEW_TERM="xterm --hold -e"
    else
        # TODO: Look into implementing gnome-terminal
        # Temporarily just use xterm if available
        if [ "$(which xterm)" =~ "\/xterm" ]; then
            NEW_TERM="xterm --hold -e"
        else
            echo "No GUI terminal found."
            exit 1
        fi
    fi
    cd bootnode
    $NEW_TERM "bootnode -nodekey boot.key"
    cd ../node1
    $NEW_TERM "./start.sh"
    cd ../node2
    $NEW_TERM "./start.sh"
    cd ../node3
    $NEW_TERM "./start.sh"
elif [[ "$OSTYPE" =~ "msys" || "$OSTYPE" =~ "win" ]]; then
    cd bootnode
    start powershell.exe -noexit -Command "bootnode -nodekey boot.key"
    cd ../node1
    start powershell.exe -noexit -Command ".\start.bat"
    cd ../node2
    start powershell.exe -noexit -Command ".\start.bat"
    cd ../node3
    start powershell.exe -noexit -Command ".\start.bat"
fi
