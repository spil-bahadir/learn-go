// comment

package main

import (
    "os"
    "bufio"
    "bytes"
    "io"
    "fmt"
    "strings"
    "github.com/ziutek/mymysql/mysql"
    _ "github.com/ziutek/mymysql/thrsafe"
    // "time"
)

type Perf struct {
    ms string
	av string
}

type PerformanceStamp struct{
	ts string
	cdn map[string]Perf
	// availability map[string]string
	ip string
	url string
}

func lineParser(line string, obj_feed chan PerformanceStamp){
	parts := strings.Split(line,",")
	ts := parts[0]
	// fmt.Println(ts,line)
	rest := strings.Join(parts[1:len(parts)],",")
	args := strings.Split(rest,"&")
	performance := ""
	availability := ""
	ip := ""
	url := ""
	for i:= 0; i<len(args); i++{
		vars := strings.Split(args[i],"=")
		if vars[0] == "p"{
			performance = vars[1]	
		}
		if vars[0] == "a"{
			availability = vars[1]
		}
		if vars[0] == "ipaddr"{
			ip = vars[1]
			url_ip_part := strings.Split(ip, ";")
			ip = url_ip_part[0]
			url = url_ip_part[1]
		}
	}
	v := PerformanceStamp{ip: ip, url: url, ts: ts}
	v.cdn = make(map[string]Perf)
	
	values := strings.Split(performance,"|")
	for j := 0; j < len(values); j++ {
		perf_val := strings.Split(values[j],",")
		v.cdn[perf_val[0]] = Perf{ms:perf_val[1]}
	}

	values = strings.Split(availability,"|")
	for j := 0; j < len(values); j++ {
		avail_val := strings.Split(values[j],",")
		obj := v.cdn[avail_val[0]]
		// v.performance[avail_val[0]].av = avail_val[1]
		
		obj.av = avail_val[1]
		v.cdn[avail_val[0]] = obj
	}

	// fmt.Println(v)
    obj_feed <- v
}


func main() {
    

	line_feed := make(chan string)
    obj_feed := make(chan PerformanceStamp)

    dbConn(obj_feed)

	file := "/Users/bcambel/Downloads/log.bahadir.log"
	
	go readLines2(file, line_feed)

	for line := range line_feed {
        lineParser(line,obj_feed)
    }

    // time.Sleep( time.Second * 100 )
}

func printOK() {
    fmt.Println("OK")
}

func checkError(err error) {
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
func checkedResult(rows []mysql.Row, res mysql.Result, err error) ([]mysql.Row,
    mysql.Result) {
    checkError(err)
    return rows, res
}

func dbConn(obj_feed chan PerformanceStamp){
    user := "root"
    pass := ""
    dbname := "multicdn"
    //proto  := "unix"
    //addr   := "/var/run/mysqld/mysqld.sock"
    proto := "tcp"
    addr := "127.0.0.1:3306"

    db := mysql.New(proto, "", addr, user, pass, dbname)

    fmt.Printf("Connect to %s:%s... ", proto, addr)

    checkError(db.Connect())
    printOK()
    
    go func(){
    idx := 0
    for obj := range obj_feed{
        db.Query("insert cdn_logs values ('%s', '%s', '%s', '%s')", obj.ts, obj.cdn,obj.ip,obj.url)
        idx++;

        if idx % 1000 == 0{
            fmt.Printf("%d. records saved\n", idx)
        }
    }}()
}

func readLines2(path string, feed chan string) (err error) {
	var (
        file *os.File
        part []byte
        prefix bool
    )
    if file, err = os.Open(path); err != nil {
        return
    }
    defer file.Close()

    reader := bufio.NewReader(file)
    buffer := bytes.NewBuffer(make([]byte, 0))

     for {
        if part, prefix, err = reader.ReadLine(); err != nil {
            break
        }
        buffer.Write(part)
        if !prefix {
        	feed <- buffer.String()
            // lines = append(lines, buffer.String())
            buffer.Reset()
        }
    }
    if err == io.EOF {
        err = nil
    }
    close(feed)
    return
}

// Read a whole file into the memory and store it as array of lines
func readLines(path string) (lines []string, err error) {
    var (
        file *os.File
        part []byte
        prefix bool
    )
    if file, err = os.Open(path); err != nil {
        return
    }
    defer file.Close()

    reader := bufio.NewReader(file)
    buffer := bytes.NewBuffer(make([]byte, 0))
    for {
        if part, prefix, err = reader.ReadLine(); err != nil {
            break
        }
        buffer.Write(part)
        if !prefix {
            lines = append(lines, buffer.String())
            buffer.Reset()
        }
    }
    if err == io.EOF {
        err = nil
    }
    return
}

func writeLines(lines []string, path string) (err error) {
    var (
        file *os.File
    )

    if file, err = os.Create(path); err != nil {
        return
    }
    defer file.Close()

    //writer := bufio.NewWriter(file)
    for _,item := range lines {
        //fmt.Println(item)
        _, err := file.WriteString(strings.TrimSpace(item) + "\n"); 
        //file.Write([]byte(item)); 
        if err != nil {
            //fmt.Println("debug")
            fmt.Println(err)
            break
        }
    }
    /*content := strings.Join(lines, "\n")
    _, err = writer.WriteString(content)*/
    return
}
