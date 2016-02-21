package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"text/template"

	"github.com/toma63/parse"
)

type PB struct {
	Path string
	Msgs []Msg
}
type Msg struct {
	Name         string
	Comm         string
	Code         int64
	IsEnum       bool
	IsMessage    bool
	IsNotRootMsg bool
	IsRootMsg    bool
	Lines        []MsgLine
	HasHandlerS  bool
	HasHandler   bool
}
type MsgLine struct {
	Type1 string
	Type2 string
	Name  string
	Comm  string
}

var last1 string
var last2 string
var commReg *regexp.Regexp = regexp.MustCompile(`/*(.*)*/`)
var codeReg *regexp.Regexp = regexp.MustCompile(`\[\d+]`)
var HandlerSReg *regexp.Regexp = regexp.MustCompile(`(.*)S`)

func main() {
	templatePath := flag.String("templatePath", "./template/", " -templatePath")
	protoPath := flag.String("protoPath", "./testdata.txt", " -proto")
	asDicPath := flag.String("asDicPath", "./out/ProtoDic.as", " -asDicPath")
	javaDicPath := flag.String("javaDicPath", "./out/ProtoDic.java", " -javaDicPath")
	asHandlerPath := flag.String("asHandlerPath", "./out/", " -asHandlerPath")
	javaHandlerPath := flag.String("javaHandlerPath", "./out/", " -javaHandlerPath")
	flag.Parse()
	fmt.Println(flag.Args())

	pb := &PB{}
	lines := make(chan string, 10)
	tokens := make(chan string, 10)
	sp := regexp.MustCompile(`\s+`)
	go ReadLines(*protoPath, lines)

	go parse.SplitTokenizer(sp, lines, tokens)
	for t := range tokens {
		last1 = last2
		last2 = t
		switch t {
		case "package":
			parsePackage(pb, tokens)
		case "message":
			parseMessage(pb, tokens)
		case "enum":
			parseEnum(pb, tokens)
		}
	}
	time.Sleep(time.Millisecond * 500)
	t1, e1 := template.ParseFiles(*templatePath+"as3_template_dic.txt")
	if e1 != nil {
		fmt.Println(e1.Error())
	}
	t2, e2 := template.ParseFiles(*templatePath+"java_template_dic.txt")
	if e2 != nil {
		fmt.Println(e2.Error())
	}
	f, _ := os.OpenFile(*asDicPath, os.O_CREATE|os.O_TRUNC, 0777)
	t1.Execute(f, pb)
	f.Close()
	f, _ = os.OpenFile(*javaDicPath, os.O_CREATE|os.O_TRUNC, 0777)
	t2.Execute(f, pb)
	f.Close()

	t1, e1 = template.ParseFiles(*templatePath+"as3_template_handler.txt")
	if e1 != nil {
		fmt.Println(e1.Error())
	}
	t2, e2 = template.ParseFiles(*templatePath+"java_template_handler.txt")
	if e2 != nil {
		fmt.Println(e2.Error())
	}
	for i := 0; i < len(pb.Msgs); i++ {
		msg := pb.Msgs[i]
		if msg.IsMessage && msg.HasHandlerS {
			f, _ = os.OpenFile(*asHandlerPath+msg.Name+"_handler.as", os.O_CREATE|os.O_EXCL, 0777)
			t1.Execute(f, msg)
		}
		if msg.IsMessage && msg.HasHandler && msg.IsNotRootMsg {
			f, _ = os.OpenFile(*javaHandlerPath+msg.Name+"Handler.java", os.O_CREATE|os.O_EXCL, 0777)
			t2.Execute(f, msg)
		}
	}
	fmt.Println("\n\nprotodic done  \n\n:)\n")
	return
}
func parsePackage(pb *PB, tokens <-chan string) {
	x := parse.TakeN(1, tokens)
	pb.Path = strings.Replace(x[1], ";", "", -1)
}
func parseMessage(pb *PB, tokens <-chan string) {
	msg := Msg{}
	msg.IsMessage = true
	x := parse.TakeUntil("{", tokens)
	
	msg.Name = x[0]
	if HandlerSReg.MatchString(msg.Name) {
		msg.HasHandlerS = true
	} else {
		msg.HasHandler = true
	}
	if msg.Name == "Msg" {
		msg.IsRootMsg = true
	} else {
		msg.IsNotRootMsg = true
	}
	
	fmt.Println("-> "+msg.Name)
	
	msg.Comm, msg.Code = getCommAndCode(last1)

	arr := parse.TakeUntil("}", tokens)
	toks := make(chan string, len(arr))
	go array2chan(arr, toks)
	lineNum := 0
	for i := 0; i < len(arr); i++ {
		if strings.Contains(arr[i], "=") {
			lineNum++
		}
	}

	for i := 0; i < lineNum; i++ {
		x = parse.TakeUntilRE(commReg, toks)
		msgLine := MsgLine{}
		msgLine.Name = x[2]
		msgLine.Type1 = x[0]
		msgLine.Type2 = x[1]
		msgLine.Comm = strings.TrimSpace(rm(x[len(x)-1], "/*", "*/"))
		msg.Lines = append(msg.Lines, msgLine)
	}
	pb.Msgs = append(pb.Msgs, msg)
}
func parseEnum(pb *PB, tokens <-chan string) {
	msg := Msg{}
	msg.IsEnum = true
	x := parse.TakeUntil("{", tokens)
	msg.Name = x[0]
	msg.IsNotRootMsg = true
	msg.Comm, msg.Code = getCommAndCode(last1)

	arr := parse.TakeUntil("}", tokens)
	toks := make(chan string, len(arr))
	go array2chan(arr, toks)
	lineNum := 0
	for i := 0; i < len(arr); i++ {
		if strings.Contains(arr[i], "=") {
			lineNum++
		}
	}

	for i := 0; i < lineNum; i++ {
		x = parse.TakeUntilRE(commReg, toks)
		msgLine := MsgLine{}
		msgLine.Name = x[0]
		msgLine.Comm = strings.TrimSpace(rm(x[len(x)-1], "/*", "*/"))
		msg.Lines = append(msg.Lines, msgLine)
	}
	pb.Msgs = append(pb.Msgs, msg)
}
func array2chan(strArray []string, toks chan string) {
	for _, e := range strArray {
		toks <- e
	}
}
func getComm(str string) string {
	if commReg.Match([]byte(str)) {
		str = strings.Replace(str, "/*", "", 1)
		str = strings.Replace(str, "*/", "", 1)
		return str
	}
	return ""
}
func getCommAndCode(str string) (comm string, code int64) {
	if commReg.Match([]byte(str)) {
		str = strings.Replace(str, "/*", "", 1)
		comm = strings.Replace(str, "*/", "", 1)
		c := codeReg.FindString(comm)
		if c != "" {
			comm = strings.Replace(comm, c, "", 1)
			c = rm(c, "[", "]")
			code, _ = strconv.ParseInt(c, 10, 32)
		}
		return
	}
	return "", 0
}
func rp(s string, a1 string, a2 string) string {
	return strings.Replace(s, a1, a2, 1)
}
func rm(s string, a1 string, a2 string) string {
	return rp(rp(s, a1, ""), a2, "")
}
func ReadLines(filepath string, lines chan<- string) {

	fd, err := os.Open(filepath)

	if err != nil {
		panic(err)
	}

	defer fd.Close()

	scanln := bufio.NewScanner(fd)

	for scanln.Scan() {
		line := scanln.Text()
		if line == "" {
			continue
		}
		lines <- line
	}
	if err := scanln.Err(); err != nil {
		panic(err)
	}

	close(lines)
}
