package evaluator

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

type UIWidget struct {
	ID      string
	Type    string
	Text    string
	OnClick *Function
}

type UIWindowState struct {
	mu       sync.Mutex
	ID       string
	Title    string
	Width    int
	Height   int
	Widgets  []UIWidget
	nextID   int
	Values   map[string]string
	server   *http.Server
	closeCh  chan struct{}
	isClosed bool
}

var (
	uiMu     sync.Mutex
	windows  = make(map[string]*UIWindowState)
	winCount = 0
)

func uiWindowFunc(args ...Object) Object {
	if len(args) < 1 || len(args) > 3 {
		return &Error{Message: "jash_ui.window() requires 1-3 arguments: title, width, height"}
	}

	title, ok := args[0].(*String)
	if !ok {
		return &Error{Message: "window title must be a string"}
	}

	width, height := 600, 400
	if len(args) > 1 {
		if w, ok := args[1].(*Integer); ok { width = int(w.Value) }
	}
	if len(args) > 2 {
		if h, ok := args[2].(*Integer); ok { height = int(h.Value) }
	}

	uiMu.Lock()
	winCount++
	id := fmt.Sprintf("jwin_%d", winCount)
	win := &UIWindowState{
		ID:      id,
		Title:   title.Value,
		Width:   width,
		Height:  height,
		Widgets: []UIWidget{},
		nextID:  0,
		Values:  make(map[string]string),
		closeCh: make(chan struct{}),
	}
	windows[id] = win
	uiMu.Unlock()

	return &JSONObject{
		Pairs: map[string]Object{
			"add_label":  &Builtin{Fn: uiMakeAdder(id, "label")},
			"add_button": &Builtin{Fn: uiMakeButton(id)},
			"add_entry":  &Builtin{Fn: uiMakeAdder(id, "entry")},
			"add_text":   &Builtin{Fn: uiMakeAdder(id, "text-area")},
			"get_value":  &Builtin{Fn: uiMakeGetValue(id)},
			"run":        &Builtin{Fn: uiMakeRun(id)},
			"close":      &Builtin{Fn: uiMakeClose(id)},
		},
	}
}

func uiMakeAdder(winID, widgetType string) func(args ...Object) Object {
	return func(args ...Object) Object {
		uiMu.Lock()
		win, ok := windows[winID]
		uiMu.Unlock()
		if !ok {
			return &Error{Message: "window not found"}
		}
		if len(args) < 1 {
			return &Error{Message: fmt.Sprintf("add_%s() requires a text argument", widgetType)}
		}
		text, ok := args[0].(*String)
		if !ok {
			return &Error{Message: "text must be a string"}
		}

		win.mu.Lock()
		win.nextID++
		wid := fmt.Sprintf("%s_w%d", winID, win.nextID)
		win.Widgets = append(win.Widgets, UIWidget{ID: wid, Type: widgetType, Text: text.Value})
		win.mu.Unlock()

		if widgetType == "entry" || widgetType == "text-area" {
			return &JSONObject{Pairs: map[string]Object{"id": &String{Value: wid}}}
		}
		return NULL
	}
}

func uiMakeButton(winID string) func(args ...Object) Object {
	return func(args ...Object) Object {
		uiMu.Lock()
		win, ok := windows[winID]
		uiMu.Unlock()
		if !ok {
			return &Error{Message: "window not found"}
		}
		if len(args) < 2 {
			return &Error{Message: "add_button() requires 2 arguments: text and callback function"}
		}
		text, ok := args[0].(*String)
		if !ok {
			return &Error{Message: "button text must be a string"}
		}
		fn, ok := args[1].(*Function)
		if !ok {
			return &Error{Message: "button callback must be a function"}
		}

		win.mu.Lock()
		win.nextID++
		wid := fmt.Sprintf("%s_w%d", winID, win.nextID)
		win.Widgets = append(win.Widgets, UIWidget{ID: wid, Type: "button", Text: text.Value, OnClick: fn})
		win.mu.Unlock()
		return NULL
	}
}

func uiMakeGetValue(winID string) func(args ...Object) Object {
	return func(args ...Object) Object {
		if len(args) < 1 {
			return &Error{Message: "get_value() requires a widget ID"}
		}
		id, ok := args[0].(*String)
		if !ok {
			return &Error{Message: "widget ID must be a string"}
		}
		uiMu.Lock()
		win, ok := windows[winID]
		uiMu.Unlock()
		if !ok {
			return &Error{Message: "window not found"}
		}
		win.mu.Lock()
		val, exists := win.Values[id.Value]
		win.mu.Unlock()
		if !exists {
			return &String{Value: ""}
		}
		return &String{Value: val}
	}
}

func uiMakeRun(winID string) func(args ...Object) Object {
	return func(args ...Object) Object {
		uiMu.Lock()
		win, ok := windows[winID]
		uiMu.Unlock()
		if !ok {
			return &Error{Message: "window not found"}
		}

		addr := fmt.Sprintf("127.0.0.1:%d", findAvailablePort())
		html := generateUIHTML(win)
		mux := http.NewServeMux()

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(html))
		})

		mux.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				http.Error(w, "method not allowed", 405)
				return
			}
			var payload struct {
				Widget string            `json:"widget"`
				Values map[string]string `json:"values"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			win.mu.Lock()
			for k, v := range payload.Values {
				win.Values[k] = v
			}
			win.mu.Unlock()
			for _, wgt := range win.Widgets {
				if wgt.ID == payload.Widget && wgt.OnClick != nil {
					env := NewEnclosedEnvironment(wgt.OnClick.Env)
					for k, v := range payload.Values {
						env.Set(k, &String{Value: v})
					}
					vals := &JSONObject{Pairs: make(map[string]Object)}
					for k, v := range payload.Values {
						vals.Pairs[k] = &String{Value: v}
					}
					applyFunction(wgt.OnClick, []Object{vals})
					break
				}
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"ok":true}`))
		})

		mux.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
			win.mu.Lock()
			if !win.isClosed {
				win.isClosed = true
				close(win.closeCh)
			}
			win.mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"ok":true}`))
		})

		win.server = &http.Server{Addr: addr, Handler: mux}
		go win.server.ListenAndServe()
		time.Sleep(50 * time.Millisecond)
		openBrowser(fmt.Sprintf("http://%s", addr))
		<-win.closeCh
		win.server.Close()

		win.mu.Lock()
		resultPairs := make(map[string]Object)
		for k, v := range win.Values {
			resultPairs[k] = &String{Value: v}
		}
		win.mu.Unlock()
		return &JSONObject{Pairs: resultPairs}
	}
}

func uiMakeClose(winID string) func(args ...Object) Object {
	return func(args ...Object) Object {
		uiMu.Lock()
		win, ok := windows[winID]
		uiMu.Unlock()
		if !ok {
			return &Error{Message: "window not found"}
		}
		win.mu.Lock()
		if !win.isClosed {
			win.isClosed = true
			close(win.closeCh)
		}
		win.mu.Unlock()
		return NULL
	}
}

func findAvailablePort() int {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 8080
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port
}

func openBrowser(url string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		exec.Command("xdg-open", url).Start()
	}
}

func generateUIHTML(win *UIWindowState) string {
	var widgetsHTML string
	for _, w := range win.Widgets {
		switch w.Type {
		case "label":
			widgetsHTML += fmt.Sprintf(`<div class="label">%s</div>`, esc(w.Text))
		case "button":
			widgetsHTML += fmt.Sprintf(`<button class="btn" onclick="clickWidget('%s')">%s</button>`, w.ID, esc(w.Text))
		case "entry":
			widgetsHTML += fmt.Sprintf(`<input class="entry" id="%s" type="text" value="%s" oninput="updateValue('%s',this.value)">`, w.ID, esc(w.Text), w.ID)
		case "text-area":
			widgetsHTML += fmt.Sprintf(`<textarea class="text" id="%s" oninput="updateValue('%s',this.value)">%s</textarea>`, w.ID, w.ID, esc(w.Text))
		}
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>%s</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#f0f0f0;padding:20px;min-height:100vh}
.window{max-width:%dpx;margin:0 auto;background:#fff;border-radius:10px;box-shadow:0 2px 16px rgba(0,0,0,0.12);padding:24px;min-height:%dpx}
.label{font-size:14px;color:#333;margin:10px 0 4px}
.btn{display:block;width:100%%;padding:10px 16px;margin:8px 0;background:#007aff;color:#fff;border:none;border-radius:6px;font-size:14px;cursor:pointer;transition:background .15s}
.btn:hover{background:#005bbf}
.btn:disabled{opacity:.6;cursor:default}
.entry,.text{display:block;width:100%%;padding:8px 12px;margin:4px 0 10px;border:1px solid #ccc;border-radius:6px;font-size:14px;font-family:inherit;outline:none;transition:border .15s}
.entry:focus,.text:focus{border-color:#007aff}
.text{min-height:80px;resize:vertical}
</style>
</head>
<body>
<div class="window">%s</div>
<script>
const vals={};
function updateValue(id,v){vals[id]=v}
function clickWidget(id){
	var btn=event.target;btn.disabled=true
	fetch('/event',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({widget:id,values:vals})})
	.then(r=>r.json()).then(d=>{btn.disabled=false}).catch(e=>{btn.disabled=false})
}
window.addEventListener('beforeunload',function(){navigator.sendBeacon('/close','{}')})
</script>
</body>
</html>`, esc(win.Title), win.Width, win.Height, widgetsHTML)
}

func esc(s string) string {
	var out []byte
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '&':
			out = append(out, "&amp;"...)
		case '<':
			out = append(out, "&lt;"...)
		case '>':
			out = append(out, "&gt;"...)
		case '"':
			out = append(out, "&quot;"...)
		case '\'':
			out = append(out, "&#39;"...)
		default:
			out = append(out, s[i])
		}
	}
	return string(out)
}
