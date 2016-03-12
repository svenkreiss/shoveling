package main

import (
    "github.com/miekg/dns"
    "os"
    "log"
    "fmt"
    "io/ioutil"
    "net/http"
    "strconv"
)


func Workers(w http.ResponseWriter) {
    c := new(dns.Client)
    m := new(dns.Msg)
    m.SetQuestion(dns.Fqdn("shoveling-worker.service.consul."), dns.TypeSRV)
    m.RecursionDesired = true

    dns_server := "127.0.0.1:8600"
    if len(os.Args) > 1 {
        dns_server = os.Args[1]
    }
    fmt.Printf("Using dns server: ", dns_server)

    r, _, err := c.Exchange(m, dns_server)
    if r == nil {
        log.Fatalf("error: %s\n", err.Error())
    }

    if r.Rcode != dns.RcodeSuccess {
        log.Fatalf("dns lookup failed\n")
    }

    for _, srv := range r.Answer {
        fmt.Fprintf(w, "%v %v\n", srv.(*dns.SRV).Port, srv.(*dns.SRV).Target)

        for _, a := range r.Extra {
            if srv.(*dns.SRV).Target != a.(*dns.A).Hdr.Name {
                continue
            }
            fmt.Fprintf(w, "%v %v\n", a.(*dns.A).Hdr.Name, a.(*dns.A).A)

            resp, err := http.Get("http://"+a.(*dns.A).A.String()+":"+strconv.Itoa(int(srv.(*dns.SRV).Port))+"/ping")
            if err != nil {
                log.Printf("GET failed\n")
                return
            }
            defer resp.Body.Close()
            body, _ := ioutil.ReadAll(resp.Body)
            fmt.Fprintf(w, "%v\n", body)
        }
    }
}


func main() {
    fmt.Println("hello worker")

    http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "pong")
    })

    http.HandleFunc("/workers", func(w http.ResponseWriter, r *http.Request) {
        Workers(w)
    })

    log.Fatal(http.ListenAndServe(":5060", nil))
}
