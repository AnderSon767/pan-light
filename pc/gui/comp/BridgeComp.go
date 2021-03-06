package comp

import (
	"encoding/json"
	"github.com/peterq/pan-light/pc/gui/qt-rpc"
	"github.com/peterq/pan-light/qt/bindings/core"
	"github.com/peterq/pan-light/qt/bindings/gui"
	"github.com/peterq/pan-light/qt/bindings/quick"
	"github.com/peterq/pan-light/qt/bindings/widgets"
	"log"
)

func init() {
	BridgeComp_QmlRegisterType2("PanLight", 1, 0, "BridgeComp")
}

// 用来和qml通信的组件
type BridgeComp struct {
	quick.QQuickItem

	_ string `property:"someString"`
	_ func() `constructor:"init"`

	_ func(w *quick.QQuickItem, attr core.Qt__WidgetAttribute, b bool) `signal:"changeAttribute,auto"`
	_ func(msg string)                                                 `signal:"logMsg,auto"`
	_ func(data interface{})                                           `signal:"test,auto"`
	_ func(data string) string                                         `slot:"callSync,auto"`
	_ func(data string)                                                `signal:"callAsync,auto"`
	_ func(data string)                                                `signal:"goMessage"`
	_ func() *core.QPoint                                              `slot:"cursorPos,auto"`
	_ func(x, y int)                                                   `slot:"setCursorPos,auto"`
	//...
}

var Bridge *BridgeComp

func (t *BridgeComp) init() {
	Bridge = t
	qt_rpc.NotifyQml = func(event string, data *qt_rpc.Gson) {
		(*data)["event"] = event
		bin, _ := json.Marshal(*data)
		t.GoMessage(string(bin))
	}
}
func (t *BridgeComp) changeAttribute(w *quick.QQuickItem, attr core.Qt__WidgetAttribute, b bool) {
	widgets.NewQWidgetFromPointer(w.Pointer()).SetAttribute(attr, b)
}

func (t *BridgeComp) logMsg(msg string) {
	log.Println("bridge log:", msg)
}

func (t *BridgeComp) test(data interface{}) {
	//core.QVariant{}
	//log.Println(fmt.Sprintf("%#v", data))
}

// 同步调用go
func (t *BridgeComp) callSync(data string) string {
	var gson qt_rpc.Gson
	json.Unmarshal([]byte(data), &gson)
	result := qt_rpc.CallGoSync(&gson)
	bin, _ := json.Marshal(*result)
	return string(bin)
}

// 异步调用go
func (t *BridgeComp) callAsync(data string) {
	var gson qt_rpc.Gson
	json.Unmarshal([]byte(data), &gson)
	qt_rpc.CallGoAsync(&gson)
}

// 获取鼠标位置
func (t *BridgeComp) cursorPos() *core.QPoint {
	return gui.QCursor_Pos()
}

// 设置鼠标位置
func (t *BridgeComp) setCursorPos(x, y int) {
	gui.QCursor_SetPos(x, y)
}
