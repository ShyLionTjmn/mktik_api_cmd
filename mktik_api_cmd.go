package main

import "flag"
import "time"
import "fmt"
import "os"
import "mikrotik"
import "strconv"

const DEV_TIMEOUT = 10
const DEV_USER="admin"
const DEV_PASS=""
const RADAR_PORT="8728"

func myUsage() {
  fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] command [command options ...]\n", os.Args[0])
  fmt.Fprintf(os.Stderr, "\t-i\tDevice IP (mandatory)\n")
  fmt.Fprintf(os.Stderr, "\t-u\tUser (default %s, also MKTIK_API_USER environment variable used)\n", DEV_USER)
  fmt.Fprintf(os.Stderr, "\t-p\tPassword (default \"%s\", also MKTIK_API_PASS environment variable used)\n", DEV_PASS)
  fmt.Fprintf(os.Stderr, "\t-P\tPort (default %s)\n", RADAR_PORT)
  fmt.Fprintf(os.Stderr, "\t-d\tDebug\n")
  fmt.Fprintf(os.Stderr, "\nExample:\n")
  fmt.Fprintf(os.Stderr, "\t%s -i 10.100.26.160 '/interface/print' '=.proplist=name,type,disabled' '?type=ether'\n", os.Args[0])
}

func main() {

  var dev_ip string
  var dev_user string
  var dev_pass string
  var dev_port int
  var opt_d bool

  var def_user string=DEV_USER;
  var def_pass string=DEV_PASS;

  var env_user=os.Getenv("MKTIK_API_USER")

  if(env_user != "") {
    def_user=env_user
  }

  var env_pass=os.Getenv("MKTIK_API_PASS")

  if(env_pass != "") {
    def_pass=env_pass
  }

  flag.StringVar(&dev_ip, "i", "", "IP")
  flag.StringVar(&dev_user, "u", def_user, "User")
  flag.StringVar(&dev_pass, "p", def_pass, "Password")
  flag.IntVar(&dev_port, "P", 8728, "Port")
  flag.BoolVar(&opt_d, "d", false, "Debug")

  flag.Usage = myUsage
  flag.Parse()

  commands := flag.Args()

  if(dev_ip == "") {
    flag.Usage()
    os.Exit(1)
  }

  if(len(commands) == 0) {
    flag.Usage()
    os.Exit(1)
  }

  if(opt_d) {
    fmt.Printf("dev_ip:\t%v\n", dev_ip)
    fmt.Printf("dev_user:\t%v\n", dev_user)
    fmt.Printf("dev_pass:\t%v\n", dev_pass)
    fmt.Printf("dev_port:\t%v\n", dev_port)
    fmt.Printf("debug:\t%v\n", opt_d)

    for _, arg := range commands {
      fmt.Printf("commands:\t%v\n", arg)
    }
  }

  var err error
  var err_str string="No error"

  var ctrl_ch chan string
  ctrl_ch= make(chan string, 1)

  dev := mikrotik.Init(dev_ip, strconv.FormatInt(int64(dev_port), 10), DEV_TIMEOUT*time.Second, ctrl_ch)
  //dev.Debug=true
  defer dev.Close()

  var snt *mikrotik.MkSentence

  err = dev.Connect(dev_user, dev_pass)
  if(err != nil) {
    err_str = "ERROR: connect error: "+err.Error()
    fmt.Printf("%s\n", err_str)
    os.Exit(1)
  }


  err = dev.Send(commands...)
  if(err != nil) {
    err_str = "ERROR: comm error: "+err.Error()
    fmt.Printf("%s\n", err_str)
    os.Exit(1)
  }
  snt,err = dev.ReadSentence()
  for err == nil && snt.Answer == "!re" {
    snt.Dump()

    snt,err = dev.ReadSentence()
  }

  if(err != nil) {
    err_str = "ERROR: comm error: "+err.Error()
    fmt.Printf("%s\n", err_str)
    os.Exit(1)
  }

  if(snt.Answer != "!done") {
    err_str = "ERROR: no !done in last reply"
    fmt.Printf("%s\n", err_str)
    os.Exit(1)
  }

}
