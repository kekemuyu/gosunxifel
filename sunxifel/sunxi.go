package sunxifel

import (
	"fmt"
	"gosunxifel/config"
	"os/exec"
	"strings"

	log "github.com/donnie4w/go-logger/logger"
	"github.com/zserge/lorca"
)

type fel struct{}

var Default fel

func flashburn(addr string, file string, ui lorca.UI) {
	if addr == "" || file == "" {
		return
	}
	var jsstr string
	jsstr = fmt.Sprintf(`$('#cmdout').append("<li>地址：%s开始烧写</li>");$('#cmdout').append("<li></li>");`, addr)
	ui.Eval(jsstr)

	rootdir := config.GetRootdir()

	log.Debug(rootdir + "/sunxi-tools-win32support_f1c100s/sunxi-fel.exe -p spiflash-write " + addr + " " + file)

	cmd := exec.Command("cmd.exe", "/c", rootdir+"/sunxi-tools-win32support_f1c100s/sunxi-fel.exe -p spiflash-write "+addr+" "+file)
	stdoutIn, _ := cmd.StdoutPipe()

	err := cmd.Start()
	if err != nil {
		log.Error("cmd.Start() failed with '%s'\n", err)
	}

	go func(ui lorca.UI) {
		for {
			bs := make([]byte, 1024, 1024)
			n, err := stdoutIn.Read(bs)
			if n > 0 {
				var re string
				re = string(bs[:n])

				re = re[strings.Index(re, "%")-3 : strings.LastIndex(re, "s")+1]

				jsstr = fmt.Sprintf(`$('#cmdout li:last').remove();$('#cmdout:last').append("<li>%s</li>");`, re)

				ui.Eval(jsstr)
				log.Debug(re)
			}
			if err != nil {
				jsstr = "地址：" + addr + "  烧写失败"

				ui.Eval(jsstr)
				log.Error(err)
				return
			}

		}
	}(ui)

	err = cmd.Wait()
	jsstr = fmt.Sprintf(`$('#cmdout').append("<li>地址：%s烧写完成</li>");`, addr)
	ui.Eval(jsstr)
}

func (f *fel) Burn(bs []Blockinfo, ui lorca.UI) {
	jsstr := `$('#cmdout').find("li").remove();`
	ui.Eval(jsstr)
	if len(bs) <= 0 {
		return
	}
	log.Debug("[burn]Blockinfos:", bs)
	for _, v := range bs {
		flashburn(v.Addr, v.Path, ui)
	}
}

func (f *fel) LoadBlockList(ui lorca.UI) []Blockinfo {
	jsstr := `$('#blocklist').find("li").remove();`

	bl := config.Cfg.Section("blocklist").KeysHash()
	bs := make([]Blockinfo, 0)
	var bi Blockinfo
	log.Debug("[LoadBlockList]blocklist:", bl)
	if len(bl) <= 0 {
		return bs
	}
	for k, v := range bl {
		str := "地址：" + k + "  " + v
		jsstr += fmt.Sprintf(`$('#blocklist').append("<li>%s</li>");`, str)
		bi.Addr = k
		bi.Path = v
		bs = append(bs, bi)
	}

	ui.Eval(jsstr)
	return bs
}

func (f *fel) AddOneBlock(bi Blockinfo) {
	log.Debug("[AddOneBlock]Blockinfo:", bi)
	config.Cfg.Section("blocklist").Key(bi.Addr).SetValue(bi.Path)
	config.Save()
}

func (f *fel) ClearBlockList(ui lorca.UI) {
	keys := config.Cfg.Section("blocklist").KeyStrings()
	if len(keys) <= 0 {
		return
	}
	for _, v := range keys {
		config.Cfg.Section("blocklist").DeleteKey(v)
	}
	config.Save()
}
