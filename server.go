package daisy

import (
	"git.sr.ht/~sircmpwn/tokidoki/storage"
	"github.com/1f349/cardcaldav"
	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/carddav"
	"net/http"
	"path/filepath"
	"time"
)

type Conf struct {
	Listen string `json:"listen"`
	DB     string `json:"db"`
}

type daisyHandler struct {
	auth    AuthProvider
	backend carddav.Backend
}

func (d *daisyHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	principlePath, err := d.auth.CurrentUserPrincipal(req.Context())
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var homeSets []webdav.BackendSuppliedHomeSet
	path, err := d.backend.AddressBookHomeSetPath(req.Context())
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	homeSets = append(homeSets, carddav.NewAddressBookHomeSet(path))

	if req.URL.Path == principlePath {
		opts := webdav.ServePrincipalOptions{
			CurrentUserPrincipalPath: principlePath,
			HomeSets:                 homeSets,
			Capabilities: []webdav.Capability{
				carddav.CapabilityAddressBook,
			},
		}
		webdav.ServePrincipal(rw, req, &opts)
		return
	}

	if req.URL.Path == "/" {
		http.Error(rw, "Daisy API Endpoint", http.StatusOK)
		return
	}

	http.NotFound(rw, req)
}

type AuthProvider interface {
	cardcaldav.ProviderMiddleware
	webdav.UserPrincipalBackend
}

func NewHttpServer(conf Conf, wd string) *http.Server {
	cardcaldav.SetupLogger(Logger)
	principle := NullAuth(cardcaldav.NewAuth(conf.DB, Logger))

	_, cardStorage, err := storage.NewFilesystem(filepath.Join(wd, "storage"), "/calendar/", "/contacts/", principle)
	if err != nil {
		Logger.Fatal("Failed to load storage backend", "err", err)
	}
	cardHandler := &carddav.Handler{Backend: cardStorage}

	handler := &daisyHandler{
		auth:    principle,
		backend: cardStorage,
	}

	r := http.NewServeMux()
	r.Handle("/", handler)
	r.Handle("GET /ok", http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		http.Error(rw, "Daisy API Endpoint", http.StatusOK)
	}))
	r.Handle("/.well-known/carddav", cardHandler)
	r.Handle("/{user}/contacts/", cardHandler)

	r2 := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t := time.Now()
		r.ServeHTTP(rw, req)
		td := time.Since(t)
		Logger.Debug("Request", "method", req.Method, "url", req.URL.String(), "remote", req.RemoteAddr, "dur", td.String())
	})

	return &http.Server{
		Addr:              conf.Listen,
		Handler:           principle.Middleware(r2),
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		WriteTimeout:      time.Minute,
		IdleTimeout:       time.Minute,
		MaxHeaderBytes:    2500,
	}
}
