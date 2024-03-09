package ftp

import (
	"crypto/tls"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/jmattheis/website/assets"
	"github.com/jmattheis/website/content"
	"github.com/jmattheis/website/util"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/acme/autocert"

	"github.com/fclairamb/ftpserver/server"
)

func getFiles(conf Config) []*virtualFileInfo {
	start := content.StartTXT(content.DnsSafeBanner, "ftp", conf.Port)
	files := []*virtualFileInfo{
		{
			dir:  "/",
			name: "blog",
			size: 4096,
			mode: 0666 | os.ModeDir,
		},
		{
			dir:     "/",
			name:    "start.txt",
			size:    int64(len(start)),
			mode:    0666,
			content: []byte(start),
		},
		{
			dir:     "/",
			name:    "projects.txt",
			size:    int64(len(content.ProjectsTXT)),
			mode:    0666,
			content: []byte(content.ProjectsTXT),
		},
		{
			dir:     "/",
			name:    "cat.txt",
			size:    int64(len(content.Cat)),
			mode:    0666,
			content: []byte(content.Cat),
		},
	}
	for i, entry := range assets.BlogList {
		value := assets.BlogContent[i]
		files = append(files, &virtualFileInfo{
			dir:     "/blog",
			name:    entry[2:],
			size:    int64(len(value)),
			mode:    0666,
			content: []byte(value),
		})
	}
	return files
}

type Config struct {
	Port string
}

func Listen(conf Config, manager *autocert.Manager, ip string) {
	files := getFiles(conf)
	drv := &MainDriver{
		Server: server.Settings{
			PublicHost: ip,
			ListenAddr: ":" + conf.Port,
			PassiveTransferPortRange: &server.PortRange{
				Start: 50000,
				End:   52999,
			},
		},
		Files: files,
		Dirs:  dirs(files),
		tlsConfig: &tls.Config{
			NextProtos: []string{"ftp"},
		},
	}
	if manager == nil {
		drv.tlsConfig.Certificates = []tls.Certificate{*util.NewUntrustedCert()}
	} else {
		drv.tlsConfig.GetCertificate = manager.GetCertificate
	}
	log.Info().Str("on", "init").Str("port", conf.Port).Msg("ftp")
	go func() {
		if err := server.NewFtpServer(drv).ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("ftp")
		}
	}()
}

func dirs(files []*virtualFileInfo) []string {
	m := map[string]struct{}{}
	for _, file := range files {
		m[file.dir] = struct{}{}
	}
	dirs := []string{}
	for dir := range m {
		dirs = append(dirs, dir)
	}
	return dirs
}

// MainDriver defines a very basic ftpserver driver
type MainDriver struct {
	tlsConfig *tls.Config     // TLS config (if applies)
	Server    server.Settings // Our settings
	nbClients int32           // Number of clients
	Files     []*virtualFileInfo
	Dirs      []string
}

// GetSettings returns some general settings around the server setup
func (driver *MainDriver) GetSettings() (*server.Settings, error) {
	return &driver.Server, nil
}

// GetTLSConfig returns a TLS Certificate to use
func (driver *MainDriver) GetTLSConfig() (*tls.Config, error) {
	return driver.tlsConfig, nil
}

// WelcomeUser is called to send the very first welcome message
func (driver *MainDriver) WelcomeUser(cc server.ClientContext) (string, error) {
	nbClients := atomic.AddInt32(&driver.nbClients, 1)
	if nbClients > 1000 {
		return "Cannot accept any additional client", fmt.Errorf(
			"too many clients: %d > %d",
			driver.nbClients,
			1000)
	}

	return "Welcome to jmattheis.de, you're on dir /", nil
}

// AuthUser authenticates the user and selects an handling driver
func (driver *MainDriver) AuthUser(server.ClientContext, string, string) (server.ClientHandlingDriver, error) {
	return &ClientDriver{BaseDir: "/", Dirs: driver.Dirs, Files: driver.Files}, nil
}

// UserLeft is called when the user disconnects, even if he never authenticated
func (driver *MainDriver) UserLeft(server.ClientContext) {
	atomic.AddInt32(&driver.nbClients, -1)
}
