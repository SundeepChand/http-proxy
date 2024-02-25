package core

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/SundeepChand/http-proxy/config"
)

type ProxyServer struct {
	// State for laod-balancing
	// TODO: Handle for concurrent writes to this map
	routesToCurOriginToUseMap map[string]int

	mux *http.ServeMux

	conf *config.Config
}

func NewProxyServer(conf *config.Config) *ProxyServer {
	proxy := &ProxyServer{
		routesToCurOriginToUseMap: make(map[string]int),
		mux:                       http.NewServeMux(),
		conf:                      conf,
	}

	proxy.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "proxy server is up and running")
	})

	for path := range conf.RoutesMapping {
		proxy.routesToCurOriginToUseMap[path] = 0
	}

	proxy.registerRoutesMap()

	return proxy
}

// getOriginServerUrl returns the exact URL to
// which the given request should be proxied to.
func (p *ProxyServer) getOriginServerUrl(originUrl string) string {
	if _, err := net.LookupHost(originUrl); err != nil {
		log.Println("error in finding host for origin:", originUrl, "error:", err)
		return fmt.Sprintf("%s:%v", "localhost", p.conf.Server.Port)
	}
	return originUrl
}

func (p *ProxyServer) registerRoutesMap() {
	client := http.Client{Timeout: time.Second * 1}
	_ = client

	for path, targets := range p.conf.RoutesMapping {
		p.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			log.Println("received req: ", r.URL, r.Method, r.RemoteAddr)

			originToUseIdx := p.routesToCurOriginToUseMap[path]
			// TODO: Fix this direct updation of map, which will lead to race conditions.
			p.routesToCurOriginToUseMap[path] = (p.routesToCurOriginToUseMap[path] + 1) % len(targets.Origins)

			originURL := targets.Origins[originToUseIdx]

			proxyReq, err := http.NewRequest(r.Method, originURL, r.Body)
			if err != nil {
				log.Println("error creating proxy request", err)
				http.Error(w, "error creating proxy request", http.StatusInternalServerError)
				return
			}

			for name, values := range r.Header {
				for _, value := range values {
					proxyReq.Header.Add(name, value)
				}
			}

			resp, err := client.Do(proxyReq)
			if err != nil {
				log.Println("error sending proxy request", err)
				http.Error(w, "error sending proxy request", http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			for name, values := range resp.Header {
				for _, value := range values {
					w.Header().Add(name, value)
				}
			}

			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
		})
	}
}

func (p *ProxyServer) ListenAndServe() error {
	return http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", p.conf.Server.Port), p.mux)
}
