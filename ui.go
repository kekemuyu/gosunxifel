package main

import (
	"fmt"
	"gosunxifel/config"
	"gosunxifel/util"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"

	log "github.com/donnie4w/go-logger/logger"

	"github.com/zserge/lorca"
)

type Myweb struct {
	UI lorca.UI
}

var Defaultweb Myweb

func New(width, height int) Myweb {
	var myweb Myweb
	var err error
	myweb.UI, err = lorca.New("", "", width, height)
	if err != nil {
		log.Fatal(err)
	}

	return myweb
}

func (m *Myweb) Run() {

	ui := m.UI
	defer ui.Close()

	// A simple way to know when UI is ready (uses body.onload even in JS)
	ui.Bind("start", func() {
		log.Debug("UI is ready")
	})

	ui.Bind("getdisk", getdisk)
	ui.Bind("browseclientpath", Browseclientpath)
	ui.Bind("browseclientuppage", Browseclientuppage)
	ui.Bind("installdriver", installdriver)
	ui.Bind("flashburn", flashburn)
	// Load HTML.
	// You may also use `data:text/html,<base64>` approach to load initial HTML,
	// e.g: ui.Load("data:text/html," + url.PathEscape(html))

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go http.Serve(ln, http.FileServer(FS))
	ui.Load(fmt.Sprintf("http://%s", ln.Addr()))

	//	go recieve(ui)
	// You may use console.log to debug your JS code, it will be printed via
	// log.Println(). Also exceptions are printed in a similar manner.
	// ui.Eval(`
	// 	console.log("Hello, world!");
	// 	console.log('Multiple values:', [1, false, {"x":5}]);
	// `)

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}

	log.Debug("exiting...")
}

//更新win根目录
func getdisk() {
	log.Debug("GetClientDisk", runtime.GOOS)
	if runtime.GOOS != "windows" {
		return
	}
	dinfo := util.GetDiskInfo()
	if len(dinfo) <= 0 {
		return
	}
	log.Debug(dinfo)
	jsStr := `$('#disk').find("option").remove()`
	Defaultweb.UI.Eval(jsStr)
	for k, v := range dinfo {
		jsStr := fmt.Sprintf(`$('#disk').append("<option value=%s>%s</option>")`, strconv.Itoa(k), v)
		Defaultweb.UI.Eval(jsStr)
	}
}

func Browseclientpath(bpath string) byte {
	curpath := config.Cfg.Section("file").Key("clientpath").MustString(config.GetRootdir())
	log.Debug(curpath)

	if bpath != "" {
		if string(curpath[len(curpath)-1]) != "/" {
			curpath += `/` + bpath
		} else {
			curpath += bpath
		}

	}
	if len(bpath) > 4 {
		if bpath[:4] == "disk" {
			curpath = bpath[4:]
		}
	}

	log.Debug(curpath)

	s, err := os.Stat(curpath)
	if err != nil {
		log.Error(err)
		log.Debug(config.GetRootdir())

		return 3
	}
	if s.IsDir() {
		files, _ := ioutil.ReadDir(curpath)
		jsStr1 := fmt.Sprintf(`$('#clientpath').val("%s");$("#clientfiles").find("li").remove()`, curpath)
		Defaultweb.UI.Eval(jsStr1)

		for _, f := range files {
			log.Debug(f.Name())
			if f.IsDir() {
				jsStr1 = fmt.Sprintf(`$('#clientfiles').append("<li><span><img style=\"width:20px;height:20px;\" src='./img/folder.ico'></span>%s</li>")`, f.Name())
			} else {
				jsStr1 = fmt.Sprintf(`$('#clientfiles').append("<li><span><img style=\"width:20px;height:20px;\" src='./img/doc.ico'></span>%s</li>")`, f.Name())
			}

			Defaultweb.UI.Eval(jsStr1)
		}
		config.Cfg.Section("file").Key("clientpath").SetValue(curpath)

		config.Save()
		return 0
	} else {
		// files, _ := ioutil.ReadDir(curpath)
		// jsStr1 := `$("#filesgroup").find("li").remove()`
		// Defaultweb.UI.Eval(jsStr1)
		// for _, f := range files {
		// 	log.Debug(f.Name())

		// 	jsStr := fmt.Sprintf(`$('#filesgroup').append("<li>%s</li>")`, f.Name())
		// 	Defaultweb.UI.Eval(jsStr)
		// }
		return 1
	}

}

func Browseclientuppage() {
	curpath := config.Cfg.Section("file").Key("clientpath").MustString(config.GetRootdir())
	log.Debug(curpath)

	curpath = util.GetParentDirectory(curpath)

	_, err := os.Stat(curpath)
	if err != nil {
		log.Error(err)
		return
	}

	files, _ := ioutil.ReadDir(curpath)
	jsStr1 := fmt.Sprintf(`$('#clientpath').val("%s");$("#clientfiles").find("li").remove()`, curpath)

	Defaultweb.UI.Eval(jsStr1)
	for _, f := range files {
		log.Debug(f.Name())
		if f.IsDir() {
			jsStr1 = fmt.Sprintf(`$('#clientfiles').append("<li><span><img style=\"width:20px;height:20px;\" src='./img/folder.ico'></span>%s</li>")`, f.Name())
		} else {
			jsStr1 = fmt.Sprintf(`$('#clientfiles').append("<li><span><img style=\"width:20px;height:20px;\" src='./img/doc.ico'></span>%s</li>")`, f.Name())
		}

		Defaultweb.UI.Eval(jsStr1)
	}
	config.Cfg.Section("file").Key("clientpath").SetValue(curpath)

	config.Save()

}

func installdriver() {

	go func() {
		rootdir := config.GetRootdir()
		cmd := exec.Command("cmd.exe", "/c", rootdir+"/sunxi-tools-win32support_f1c100s/zadig-2.3.exe")
		err := cmd.Run()

		if err != nil {
			fmt.Println(err)
		}
	}()

}

func flashburn(addr string, file string) string {
	if addr == "" || file == "" {
		return "请检查地址或选择的文件是否正确"
	}
	rootdir := config.GetRootdir()

	log.Debug(rootdir + "/sunxi-tools-win32support_f1c100s/sunxi-fel.exe -p spiflash-write " + addr + " " + file)
	cmd := exec.Command("cmd.exe", "/c", rootdir+"/sunxi-tools-win32support_f1c100s/sunxi-fel.exe -p spiflash-write "+addr+" "+file)
	out, err := cmd.CombinedOutput()
	err = cmd.Run()
	log.Debug(err, string(out))
	if err != nil {
		fmt.Println(err)
	}
	return "烧写完成"
}
