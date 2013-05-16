package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "net/http"
)

import "time"

func main() {
    go download_url("http://gamesgames.com/")
    go download_url("http://agame.com/")
    go download_url("http://www.girlsgogames.com/")

    time.Sleep( time.Second * 1000000000 )
}

func download_url(url string){
  response, err := http.Get(url)
    if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        }
        fmt.Printf("%s\n", string(contents))
    }
}
