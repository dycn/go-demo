// Code generated by Thrift Compiler (0.14.2). DO NOT EDIT.

package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"github.com/apache/thrift/lib/go/thrift"
	"students"
)

var _ = students.GoUnusedProtection__

func Usage() {
  fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
  flag.PrintDefaults()
  fmt.Fprintln(os.Stderr, "\nFunctions:")
  fmt.Fprintln(os.Stderr, "  void Ping(i64 num)")
  fmt.Fprintln(os.Stderr, "  string GetClassName()")
  fmt.Fprintln(os.Stderr, "   List()")
  fmt.Fprintln(os.Stderr, "  void Add(Student s)")
  fmt.Fprintln(os.Stderr, "  bool IsNameExist(string name)")
  fmt.Fprintln(os.Stderr, "  i16 Count()")
  fmt.Fprintln(os.Stderr)
  os.Exit(0)
}

type httpHeaders map[string]string

func (h httpHeaders) String() string {
  var m map[string]string = h
  return fmt.Sprintf("%s", m)
}

func (h httpHeaders) Set(value string) error {
  parts := strings.Split(value, ": ")
  if len(parts) != 2 {
    return fmt.Errorf("header should be of format 'Key: Value'")
  }
  h[parts[0]] = parts[1]
  return nil
}

func main() {
  flag.Usage = Usage
  var host string
  var port int
  var protocol string
  var urlString string
  var framed bool
  var useHttp bool
  headers := make(httpHeaders)
  var parsedUrl *url.URL
  var trans thrift.TTransport
  _ = strconv.Atoi
  _ = math.Abs
  flag.Usage = Usage
  flag.StringVar(&host, "h", "localhost", "Specify host and port")
  flag.IntVar(&port, "p", 9090, "Specify port")
  flag.StringVar(&protocol, "P", "binary", "Specify the protocol (binary, compact, simplejson, json)")
  flag.StringVar(&urlString, "u", "", "Specify the url")
  flag.BoolVar(&framed, "framed", false, "Use framed transport")
  flag.BoolVar(&useHttp, "http", false, "Use http")
  flag.Var(headers, "H", "Headers to set on the http(s) request (e.g. -H \"Key: Value\")")
  flag.Parse()
  
  if len(urlString) > 0 {
    var err error
    parsedUrl, err = url.Parse(urlString)
    if err != nil {
      fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
      flag.Usage()
    }
    host = parsedUrl.Host
    useHttp = len(parsedUrl.Scheme) <= 0 || parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https"
  } else if useHttp {
    _, err := url.Parse(fmt.Sprint("http://", host, ":", port))
    if err != nil {
      fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
      flag.Usage()
    }
  }
  
  cmd := flag.Arg(0)
  var err error
  if useHttp {
    trans, err = thrift.NewTHttpClient(parsedUrl.String())
    if len(headers) > 0 {
      httptrans := trans.(*thrift.THttpClient)
      for key, value := range headers {
        httptrans.SetHeader(key, value)
      }
    }
  } else {
    portStr := fmt.Sprint(port)
    if strings.Contains(host, ":") {
           host, portStr, err = net.SplitHostPort(host)
           if err != nil {
                   fmt.Fprintln(os.Stderr, "error with host:", err)
                   os.Exit(1)
           }
    }
    trans, err = thrift.NewTSocket(net.JoinHostPort(host, portStr))
    if err != nil {
      fmt.Fprintln(os.Stderr, "error resolving address:", err)
      os.Exit(1)
    }
    if framed {
      trans = thrift.NewTFramedTransport(trans)
    }
  }
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error creating transport", err)
    os.Exit(1)
  }
  defer trans.Close()
  var protocolFactory thrift.TProtocolFactory
  switch protocol {
  case "compact":
    protocolFactory = thrift.NewTCompactProtocolFactory()
    break
  case "simplejson":
    protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
    break
  case "json":
    protocolFactory = thrift.NewTJSONProtocolFactory()
    break
  case "binary", "":
    protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid protocol specified: ", protocol)
    Usage()
    os.Exit(1)
  }
  iprot := protocolFactory.GetProtocol(trans)
  oprot := protocolFactory.GetProtocol(trans)
  client := students.NewClassMemberClient(thrift.NewTStandardClient(iprot, oprot))
  if err := trans.Open(); err != nil {
    fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
    os.Exit(1)
  }
  
  switch cmd {
  case "Ping":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "Ping requires 1 args")
      flag.Usage()
    }
    argvalue0, err21 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err21 != nil {
      Usage()
      return
    }
    value0 := argvalue0
    fmt.Print(client.Ping(context.Background(), value0))
    fmt.Print("\n")
    break
  case "GetClassName":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "GetClassName requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.GetClassName(context.Background()))
    fmt.Print("\n")
    break
  case "List":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "List requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.List(context.Background()))
    fmt.Print("\n")
    break
  case "Add":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "Add requires 1 args")
      flag.Usage()
    }
    arg22 := flag.Arg(1)
    mbTrans23 := thrift.NewTMemoryBufferLen(len(arg22))
    defer mbTrans23.Close()
    _, err24 := mbTrans23.WriteString(arg22)
    if err24 != nil {
      Usage()
      return
    }
    factory25 := thrift.NewTJSONProtocolFactory()
    jsProt26 := factory25.GetProtocol(mbTrans23)
    argvalue0 := students.NewStudent()
    err27 := argvalue0.Read(context.Background(), jsProt26)
    if err27 != nil {
      Usage()
      return
    }
    value0 := argvalue0
    fmt.Print(client.Add(context.Background(), value0))
    fmt.Print("\n")
    break
  case "IsNameExist":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "IsNameExist requires 1 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    fmt.Print(client.IsNameExist(context.Background(), value0))
    fmt.Print("\n")
    break
  case "Count":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "Count requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.Count(context.Background()))
    fmt.Print("\n")
    break
  case "":
    Usage()
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
  }
}
