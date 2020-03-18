package ftp

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/jmattheis/website/content"
	"github.com/jmattheis/website/util"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/acme/autocert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

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
	for _, entry := range content.BlogBox.List() {
		value, _ := content.BlogBox.FindString(entry)
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

func Listen(conf Config, manager *autocert.Manager) {
	files := getFiles(conf)
	drv := &MainDriver{
		Server: server.Settings{
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
	for dir, _ := range m {
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

	// This is the new IP loading change coming from Ray
	if driver.Server.PublicHost == "" {
		publicIP := ""

		log.Debug().Str("on", "fetchexternalip").Msg("ftp")

		if publicIP, err := externalIP(); err != nil {
			log.Warn().Str("on", "fetchexternalip").Err(err).Msg("ftp")
		} else {
			log.Debug().Str("on", "pubip").Str("ip", publicIP).Msg("ftp")
		}

		driver.Server.PublicIPResolver = func(cc server.ClientContext) (string, error) {
			if strings.HasPrefix(cc.RemoteAddr().String(), "127.0.0.1") {
				return "127.0.0.1", nil
			}
			return publicIP, nil
		}
	}

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

func externalIP() (string, error) {
	// If you need to take a bet, amazon is about as reliable & sustainable a service as you can get
	rsp, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		return "", err
	}

	defer func() {
		if errClose := rsp.Body.Close(); errClose != nil {
			fmt.Println("Problem closing checkip connection, err:", errClose)
		}
	}()

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(buf)), nil
}
