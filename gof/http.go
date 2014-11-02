package gof

import (
	"fmt"
	"github.com/fuxiaohei/GoBlog/gof/log"
	"net/http"
)

type HttpServer struct {
	RouterInterface
	ConfigInterface
	*Injector
	mid      []RouterHandler
	notFound RouterHandler
}

func NewHttpServer(configFile string) *HttpServer {
	h := new(HttpServer)
	h.RouterInterface = NewRouter()
	h.ConfigInterface, _ = NewConfig(configFile)
	h.Injector = NewInjector()
	h.mid = make([]RouterHandler, 0)
	return h
}

func (hs *HttpServer) Use(fn ...RouterHandler) {
	hs.mid = append(hs.mid, fn...)
}

func (hs *HttpServer) NotFound(fn RouterHandler) {
	hs.notFound = fn
}

func (hs *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	params, fn := hs.RouterInterface.Find(p, r.Method)

	fn = append(hs.mid, fn...)

	ctx := NewContext(hs, w, r, fn, params, hs.Injector.Clone())

	ctx.Run()

	if ctx.Status > 0 {
		// force to send response
		ctx.SendResponse()
		if ctx.Status > 0 {
			if ctx.Status >= 400 {
				log.Warn("[%d] %s %s %s %s", ctx.Status, r.RemoteAddr, r.Proto, r.Method, p)
				return
			}
			log.Debug("[%d] %s %s %s %s", ctx.Status, r.RemoteAddr, r.Proto, r.Method, p)
			return
		}
	}

	if ctx.Status == 0 && hs.notFound != nil {
		hs.notFound(ctx)
		log.Debug("[%d] %s %s %s %s", ctx.Status, r.RemoteAddr, r.Proto, r.Method, p)
		return
	}

	w.WriteHeader(404)
	w.Write([]byte(http.StatusText(404)))
	log.Warn("[%d] %s %s %s %s", 404, r.RemoteAddr, r.Proto, r.Method, p)
}

func (hs *HttpServer) Listen(addr string, port int) error {
	log.Info("listen %s:%d", addr, port)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), hs)
}
