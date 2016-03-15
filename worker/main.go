package main

import (
    "github.com/miekg/dns"
    "os"
    "log"
    "fmt"
    "io/ioutil"
    "net"
    "net/http"
    "strconv"

    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)


// use Consul DNS to resolve
func Resolve(q string) (ip net.IP, port uint16, target string, err error) {
    c := new(dns.Client)
    m := new(dns.Msg)
    m.SetQuestion(dns.Fqdn(q), dns.TypeSRV)
    m.RecursionDesired = true

    dns_server := "127.0.0.1:8600"
    if len(os.Args) > 1 {
        dns_server = os.Args[1]
    }
    fmt.Printf("Using dns server: %v\n", dns_server)

    r, _, err := c.Exchange(m, dns_server)
    if r == nil {
        log.Fatalf("error: %s\n", err.Error())
    }

    if r.Rcode != dns.RcodeSuccess {
        log.Fatalf("dns lookup failed\n")
    }

    for _, srv := range r.Answer {
        port = srv.(*dns.SRV).Port
        target = srv.(*dns.SRV).Target

        fmt.Printf("%v %v\n", port, target)

        for _, a := range r.Extra {
            if target != a.(*dns.A).Hdr.Name {
                continue
            }
            ip = a.(*dns.A).A
            fmt.Printf("%v %v\n", target, ip)
            return
        }
    }

    log.Fatalf("no DNS record found\n")
    return
}


func Workers(w http.ResponseWriter) {
    ip, port, _, err := Resolve("shoveling-worker.service.consul.")

    fmt.Fprintf(w, "http://"+ip.String()+":"+strconv.Itoa(int(port))+"/ping\n")
    resp, err := http.Get("http://"+ip.String()+":"+strconv.Itoa(int(port))+"/ping")
    if err != nil {
        log.Printf("GET failed\n")
        return
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Fprintf(w, "%v\n", string(body))
}


type Db struct {
    conn *sql.DB
}

func (db *Db) Get() (*sql.DB) {
    if db.conn != nil {
        return db.conn
    }

    ip, port, _, err := Resolve("name-service.service.consul.")
    d, err := sql.Open("mysql",
        "root:@tcp("+ip.String()+":"+strconv.Itoa(int(port))+")/nodes")
    if err != nil {
        log.Fatal(err)
    }

    db.conn = d
    return d
}

func Nodes(w http.ResponseWriter, db *Db) {
    err := db.Get().Ping()
    if err != nil {
        log.Fatal(err)
    }

    var availabledisk int
    err = db.Get().QueryRow("select availabledisk from nodes where id = ?", 1).Scan(&availabledisk)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(availabledisk)
}


func main() {
    fmt.Println("Starting worker.")
    var db Db

    http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "pong")
    })

    http.HandleFunc("/workers", func(w http.ResponseWriter, r *http.Request) {
        Workers(w)
    })

    http.HandleFunc("/nodes", func(w http.ResponseWriter, r *http.Request) {
        Nodes(w, &db)
    })

    log.Fatal(http.ListenAndServe(":5060", nil))
}
