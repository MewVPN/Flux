package wg

import (
	"log"
	"os/exec"
	"time"

	"flux/internal/config"
	"flux/internal/docker"
	"flux/internal/util"
)

const (
	image = "ghcr.io/wg-easy/wg-easy:15.2.2"
	name  = "wg-easy"
)

func Ensure(cfg *config.Config) error {
	log.Println("[wg] ensure wg-easy")

	if !docker.Available() {
		log.Println("[wg] docker not available, skipping wg-easy management")
		return nil
	}

	if running() {
		log.Println("[wg] wg-easy already running")
		return nil
	}

	log.Println("[wg] wg-easy not running, preparing to start")

	if cfg.WGEasyUser == "" || cfg.WGEasyPassword == "" {
		log.Println("[wg] wg-easy credentials not found, generating")

		cfg.WGEasyUser = "admin"
		cfg.WGEasyPassword = util.Secret(16)

		if err := config.Save(cfg); err != nil {
			log.Printf("[wg] failed to save wg-easy credentials: %v\n", err)
			return err
		}

		log.Println("[wg] wg-easy credentials generated and saved")
	} else {
		log.Println("[wg] wg-easy credentials already present in config")
	}

	initHost := util.DetectPublicIP()
	log.Printf("[wg] detected IP for INIT_HOST: %s\n", initHost)

	log.Printf("[wg] starting wg-easy container (%s)\n", image)

	cmd := exec.Command(
		"docker", "run", "-d",
		"--name", name,
		"--restart", "unless-stopped",

		"--cap-add", "NET_ADMIN",
		"--cap-add", "SYS_MODULE",

		"--device", "/dev/net/tun",

		"-p", "51820:51820/udp",
		"-p", "51821:51821/tcp",

		"-e", "INIT_ENABLED=true",
		"-e", "INIT_USERNAME="+cfg.WGEasyUser,
		"-e", "INIT_PASSWORD="+cfg.WGEasyPassword,
		"-e", "INIT_HOST="+initHost,
		"-e", "INIT_PORT=51820",
		"-e", "EXPERIMENTAL_AWG=true",
		"-e", "INSECURE=true",

		"-v", "/etc/wireguard:/etc/wireguard",
		"-v", "/lib/modules:/lib/modules:ro",

		"--sysctl", "net.ipv4.conf.all.src_valid_mark=1",
		"--sysctl", "net.ipv4.ip_forward=1",

		image,
	)

	if err := cmd.Run(); err != nil {
		log.Printf("[wg] failed to start wg-easy container: %v\n", err)
		return err
	}

	log.Println("[wg] wg-easy container started, waiting for API")

	for i := 1; i <= 10; i++ {
		if running() {
			log.Printf("[wg] wg-easy is ready (attempt %d/10)\n", i)
			return nil
		}

		log.Printf("[wg] waiting for wg-easy API... (%d/10)\n", i)
		time.Sleep(time.Second)
	}

	log.Println("[wg] wg-easy did not become ready after 10 seconds")
	return nil
}
