package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var fileStorage string = "/home/"+os.Getenv("USER")+"/Documents/shell-bookmarks/bookmarks.txt"
var bookmarkPath string = ""
var file *os.File

func main(){
    app := setUpApp()
    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

func setUpApp() *cli.App {
    return &cli.App {
        Name: "shell-bookmark",
        Flags: []cli.Flag {
            &cli.StringFlag{ 
                Name: "file",
                Aliases: []string{"f"},
                Usage: "File to bookmark",
            },
        },
        Action: func(context *cli.Context) error {
            if context.NArg() > 0 {
                args := context.Args()
                command := args.Get(0)
                bookmarkPath = args.Get(1)
                var err error

                err = openOrCreateFile()
                defer file.Close()

                if err != nil {
                    return err
                }

                if command == "add" {
                    if err := addBookMark(); err != nil {
                        return err
                    }
                    return nil
                } 

                if command == "list" {
                    lines, err := readLines();
                    if  err != nil {
                        return err
                    }
                    
                    for line := range lines {
                        fmt.Printf("%ṣ\n",lines[line])
                    }
                    return nil
                }
            }

            return nil
        },
    }
}

func openOrCreateFile() error {
    openFile, err := os.OpenFile(fileStorage, os.O_RDWR, os.ModeAppend)
    if err != nil {
        if errors.Is(err, fs.ErrNotExist) {
            err = nil
            openFile, err = os.Create(fileStorage)
            if err != nil {
                log.Fatalf("createFile: %s", err)
                return err
            }
        } else {
            log.Fatalf("openFile: %s", err)
            return err
        }
    }
    file = openFile
    return err
}

func addBookMark() error {
    allBookMarksInFile, err := readLines()
    if err != nil {
        log.Fatalf("readLines: %s", err)
        return err
    }

    allBookMarksInFile = append(allBookMarksInFile, bookmarkPath)
    if err := writeLines(allBookMarksInFile); err != nil {
        log.Fatalf("writeLines: %s", err)
        return err
    }

    return nil
}

func readLines() ([]string, error) {
    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func writeLines(lines []string) error {
    writer := bufio.NewWriter(file)
    for _, line := range lines {
        fmt.Fprintln(writer, line)
    }

    return writer.Flush()
}
