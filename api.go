package volumeapi

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

const (
	defaultDockerRootDirectory    = "/var/lib/docker/volumes"
	defaultContentTypeV1          = "appplication/vnd.docker.plugins.v1+json"
	defaultImplementationManifest = `{"Implements": ["VolumeDriver"]}`
)

type VolumeRequest struct {
	Root string `json:",omitempty"`
	Name string
}

type VolumeResponse struct {
	Mountpoint string
	Err        error
}

type VolumeDriver interface {
	Create(VolumeRequest) VolumeResponse
	Remove(VolumeRequest) VolumeResponse
	Path(VolumeRequest) VolumeResponse
	Mount(VolumeRequest) VolumeResponse
	Umount(VolumeRequest) VolumeResponse
}

type VolumeHandler struct {
	rootDirectory string
	handler       VolumeDriver
	mux           *http.ServeMux
}

type actionHandler func(VolumeRequest) VolumeResponse

func NewVolumeHandler(handler VolumeDriver) *VolumeHandler {
	return NewVolumeHandlerWithRoot(defaultDockerRootDirectory, handler)
}

func NewVolumeHandlerWithRoot(rootDirectory string, handler VolumeDriver) *VolumeHandler {
	h := &VolumeHandler{rootDirectory, handler, http.NewServeMux()}
	h.initMux()
	return h
}

func (h *VolumeHandler) initMux() {
	h.mux.HandleFunc("/Plugin.Activate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", defaultContentTypeV1)
		fmt.Fprintln(w, defaultImplementationManifest)
	})

	h.handle("/VolumeDriver.Create", func(req VolumeRequest) VolumeResponse {
		return h.handler.Create(req)
	})

	h.handle("/VolumeDriver.Remove", func(req VolumeRequest) VolumeResponse {
		return h.handler.Remove(req)
	})

	h.handle("/VolumeDriver.Path", func(req VolumeRequest) VolumeResponse {
		return h.handler.Path(req)
	})

	h.handle("/VolumeDriver.Mount", func(req VolumeRequest) VolumeResponse {
		return h.handler.Mount(req)
	})

	h.handle("/VolumeDriver.Umount", func(req VolumeRequest) VolumeResponse {
		return h.handler.Umount(req)
	})
}

func (h *VolumeHandler) handle(name string, actionCall actionHandler) {
	h.mux.HandleFunc(name, func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeRequest(w, r)
		if err != nil {
			return
		}

		req.Root = h.rootDirectory
		res := actionCall(req)

		encodeResponse(w, res)
	})
}

func (h *VolumeHandler) ListenAndServe(proto, addr, group string) error {
	server := http.Server{
		Addr:    addr,
		Handler: h.mux,
	}

	start := make(chan struct{})

	var l net.Listener
	var err error
	switch proto {
	case "tcp":
		l, err = newTcpSocket(addr, nil, start)
	case "unix":
		l, err = newUnixSocket(addr, group, start)
	}
	if err != nil {
		return err
	}

	close(start)
	return server.Serve(l)
}

func decodeRequest(w http.ResponseWriter, r *http.Request) (req VolumeRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return
}

func encodeResponse(w http.ResponseWriter, res VolumeResponse) {
	w.Header().Set("Content-Type", defaultContentTypeV1)
	if res.Err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(res)
}
