#!/bin/bash

install () {
    echo "Installing tasks cli"
    git clone https://github.com/Ayobami6/todo_cli
    chmod u+x todo_cli/bin/tasks
    sudo cp todo_cli/bin/tasks /bin
    # clean up
    rm -rf todo_cli
}

install