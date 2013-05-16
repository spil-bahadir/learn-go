package main

import ( 
  "fmt"
	"http"
)

func add( x ,y int) int{
    return x + y;
}

func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return
}

var x, y, z int = 1, 2, 3
var c, python, java = true, false, "no!"

func array_logic(){
	p := []int{2, 3, 5, 7, 11, 13}
    fmt.Println("p ==", p)

    for i := 0; i < len(p); i++ {
        fmt.Printf("p[%d] == %d\n",
            i, p[i])
    }
}

func main(){
	i := 1;
	http_main()
	for {
	    // fmt.Printf("hello, world\n");
	    add(4,5)
	    split(17)   
	    array_logic()

	    if i % 10000000 == 0{
	    	fmt.Println(i/10000000,x,y,z,c,python,java);
	    }
	    i++
	}
}

func http_main(){
    //Process the http commands
    fmt.Printf("Starting http Server ... ")
    http.Handle("/", http.HandlerFunc(sayHello))
    err := http.ListenAndServe("0.0.0.0:8080", nil)
    if err != nil {
        fmt.Printf("ListenAndServe Error",err)
    }
}

func sayHello(c http.ResponseWriter, req *http.Request) {
    fmt.Printf("New Request\n")
    processRequest(c, req)
}

func processRequest(w http.ResponseWriter, req *http.Request){
    time.Sleep(time.Second*3)
    w.Write([]byte("Go Say’s Hello(Via http)"))
    fmt.Println("End")
}

type Vertex struct {
    Lat, Long float64
}

var m = map[string]Vertex{
    "Bell Labs": Vertex{
        40.68433, -74.39967,
    },
    "Google": Vertex{
        37.42202, -122.08408,
    },
}

