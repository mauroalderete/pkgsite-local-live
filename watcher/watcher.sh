#!/bin/bash

NAME=watcher
VERSION_TAG=v1.0.0

# Colors
end="\033[0m"
black="\033[0;30m"
blackb="\033[1;30m"
white="\033[0;37m"
whiteb="\033[1;37m"
whitei="\033[3;37m"
red="\033[0;31m"
redb="\033[1;31m"
green="\033[0;32m"
greenb="\033[1;32m"
yellow="\033[0;33m"
yellowb="\033[1;33m"
lightyellow="\033[0;93m"
blue="\033[0;34m"
blueb="\033[1;34m"
purple="\033[0;35m"
purpleb="\033[1;35m"
lightblue="\033[0;36m"
lightblueb="\033[1;36m"

printcolor() {
    printf "$1$2${end}"
}

printverbose() {
    if [[ $VERBOSE == true ]]; then
        printcolor $whitei "$(printf "[$NAME] $1")\n";
    fi
}

printerror() {
    printcolor $redb "[$NAME] $@";
}

printwarn() {
    printcolor $lightyellow "[$NAME] $@";
}

printinfo() {
    printcolor $white "$@";
}

printnotify() {
    printcolor $lightblue "[inotifywait] $1";
}

printcommand() {
    printcolor $green "[command] $1";
}

show_usage() {
    printinfo "Watchers events on a folder and executes a command when something changes\n\n"
    printinfo "Usage: $NAME [ options ]\n\n"
    printinfo "Options\n"
    printinfo "\t-h|--help\t\tDisplay this help message.\n"
    printinfo "\t-e|--exec <command>\tComman to execute when an event occurs.\n\t\t\t\tYou can use '&' in the end to execute the command in background.\n"
    printinfo "\t-i|--inotify <args>\tList of arguments to pass to inotifywait.\n\t\t\t\tPlease execute 'inotifywait --help' to get more information.\n"
    printinfo "\t-k|--keeplive\t\tKeep alive when the command returns code != 0.\n"
    printinfo "\t--verbose\t\tMode verbose.\n"
    printinfo "\t-v|--version\t\tDisplay version.\n"
}

show_version() {
    printinfo "$NAME $VERSION_TAG"
}

exit_script() {
    printverbose "Exit with $1"
    if [[ $1 != 0 ]]; then
        printerror "Exit with $1\n";
    fi
    exit $1
}

OPT=$(getopt -o hve:i:k --long help,version,verbose,exec:,inotify:,keeplive -n "$NAME" -- "$@")

if [ $? != 0 ]; then
    show_usage
    exit $?
fi

eval set -- "$OPT"

HELP=false
VERSION=false
VERBOSE=false
COMMAND=""
INOTIFY=""
KEEPLIVE=false

while true; do
    case "$1" in
        -h | --help ) HELP=true; shift ;;
        -v | --version ) VERSION=true; shift ;;
        --verbose) VERBOSE=true; shift ;;
        -e | --exec ) COMMAND="$2"; shift 2 ;;
        -i | --inotify ) INOTIFY="$2"; shift 2 ;;
        -k | --keeplive ) KEEPLIVE=true; shift ;;
        -- ) shift; break ;;
        * ) break ;;
    esac
done

printverbose "Parse arguments..."

if [[ $HELP == true ]]; then
    show_usage;

    exit;
fi

if [[ $VERSION == true ]]; then
    show_version

    exit;
fi

if [[ $COMMAND == "" || $INOTIFY == "" ]]; then
    printerror "Argument missing\n"
    show_usage;

    exit 1;
fi

printverbose "Starting..."
printwarn "Press [CTRL+C] to stop...\n\n"
while true
do
    printverbose "Waiting event..."

    INOTIFY_STDOUT=$(inotifywait $INOTIFY)
    INOTIFY_CODE=$?
    printnotify "$INOTIFY_STDOUT\n"

    if [ $INOTIFY_CODE != 0 ]; then
        printwarn "inotifywait stopped, exited with code $INOTIFY_CODE\n"
        exit_script $INOTIFY_CODE
    fi

    printverbose "Command start '$COMMAND'"

    COMMAND_STDOUT=$($COMMAND)
    COMMAND_CODE=$?
    printcommand "$COMMAND_STDOUT\n"

    if [[ $COMMAND_CODE != 0 ]]; then
        if [[ $KEEPLIVE == false ]]; then
            printerror "Failed to execute the command, exited with code $COMMAND_CODE\n"
            printerror "Stopping...\n"

            exit_script $COMMAND_CODE
        else
            printwarn "Failed to execute the command, exited with code $COMMAND_CODE\n"
            printverbose "Command ended"
        fi
    else
        printverbose "Command ended"
    fi

    printinfo "\n"
done

printverbose "Exit"
