package promote

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

type (
	// Daemon defines Docker daemon parameters.
	Daemon struct {
		Registry           string   // Docker registry
		Mirror             string   // Docker registry mirror
		Insecure           bool     // Docker daemon enable insecure registries
		StorageDriver      string   // Docker daemon storage driver
		StoragePath        string   // Docker daemon storage path
		Disabled           bool     // DOcker daemon is disabled (already running)
		Debug              bool     // Docker daemon started in debug mode
		Bip                string   // Docker daemon network bridge IP address
		DNS                []string // Docker daemon dns server
		DNSSearch          []string // Docker daemon dns search domain
		MTU                string   // Docker daemon mtu setting
		IPv6               bool     // Docker daemon IPv6 networking
		Experimental       bool     // Docker daemon enable experimental mode
		InsecureRegistries []string // Docker daemon insecure registries
	}

	// Login defines Docker login parameters.
	Login struct {
		Registry string // Docker registry address
		Username string // Docker registry username
		Password string // Docker registry password
		Email    string // Docker registry email
	}

	// Promote defines Docker build parameters.
	Promote struct {
		Tags        []string // Docker build tags
		PullRepo    string   // Docker pull repository
		PushRepo    string   // Docker build repository
	}

	// Plugin defines the Promote plugin parameters.
	Plugin struct {
		PullLogin   Login // Docker pull login configuration
		PushLogin   Login  // Docker login configuration
		Promote     Promote  // Docker image promote configuration
		Daemon      Daemon // Docker daemon configuration
		Dryrun      bool   // Docker push is skipped
		Cleanup     bool   // Docker purge is enabled
	}
)

// Exec executes the plugin step
func (p Plugin) Exec() error {
	// start the Docker daemon server
	if !p.Daemon.Disabled {
		cmd := commandDaemon(p.Daemon)
		if p.Daemon.Debug {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		} else {
			cmd.Stdout = ioutil.Discard
			cmd.Stderr = ioutil.Discard
		}
		go func() {
			trace(cmd)
			cmd.Run()
		}()
	}

	// poll the docker daemon until it is started. This ensures the daemon is
	// ready to accept connections before we proceed.
	for i := 0; i < 15; i++ {
		cmd := commandInfo()
		err := cmd.Run()
		if err == nil {
			break
		}
		time.Sleep(time.Second * 1)
	}

	// login to the docker pull registry
	if p.PullLogin.Password != "" {
		cmd := commandLogin(p.PullLogin)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("Pull login error authenticating: %s", err)
		}
	} else {
		fmt.Println("Pull login registry credentials not provided. Guest mode enabled.")
	}

	// login to the Docker push registry
	if p.PushLogin.Password != "" {
		cmd := commandLogin(p.PushLogin)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("Error authenticating: %s", err)
		}
	} else {
		fmt.Println("Registry credentials not provided. Guest mode enabled.")
	}


	var cmds []*exec.Cmd
	cmds = append(cmds, commandVersion()) // docker version
	cmds = append(cmds, commandInfo())    // docker info

	for _, tag := range p.Promote.Tags {
		
		cmds = append(cmds, commandPull(p.Promote, tag)) // docker pull
		cmds = append(cmds, commandTag(p.Promote, tag)) // docker tag
		
		if p.Dryrun == false {
			cmds = append(cmds, commandPush(p.Promote, tag)) // docker push
		}

		if p.Cleanup {
			cmds = append(cmds, commandRmi(fmt.Sprintf("%s:%s", p.Promote.PullRepo, tag)))// docker rmi
			cmds = append(cmds, commandRmi(fmt.Sprintf("%s:%s", p.Promote.PushRepo, tag))) // docker rmi			
		}
	}

	if p.Cleanup {
		cmds = append(cmds, commandPrune())           // docker system prune -f
	}

	// execute all commands in batch mode.
	for _, cmd := range cmds {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		trace(cmd)

		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

const dockerExe = "/usr/local/bin/docker"
const dockerdExe = "/usr/local/bin/dockerd"

// helper function to create the docker login command.
func commandLogin(login Login) *exec.Cmd {
	if login.Email != "" {
		return commandLoginEmail(login)
	}

	fmt.Printf("login insecure registry: %s\n", login.Registry)
	return exec.Command(
		dockerExe, "login",
		"-u", login.Username,
		"-p", login.Password,
		login.Registry,
	)
}

func commandLoginEmail(login Login) *exec.Cmd {
	return exec.Command(
		dockerExe, "login",
		"-u", login.Username,
		"-p", login.Password,
		"-e", login.Email,
		login.Registry,
	)
}

// helper function to create the docker info command.
func commandVersion() *exec.Cmd {
	return exec.Command(dockerExe, "version")
}

// helper function to create the docker info command.
func commandInfo() *exec.Cmd {
	return exec.Command(dockerExe, "info")
}

func dirExist(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}
	return false
}


// helper function to create the docker tag command.
func commandTag(build Promote, tag string) *exec.Cmd {
	var (
		source = fmt.Sprintf("%s:%s", build.PullRepo, tag)
		target = fmt.Sprintf("%s:%s", build.PushRepo, tag)
	)
	return exec.Command(
		dockerExe, "tag", source, target,
	)
}
func commandPull(build Promote, tag string) *exec.Cmd {
	target := fmt.Sprintf("%s:%s", build.PullRepo, tag)
	fmt.Println("Pull repo:", target)
	return exec.Command(dockerExe, "pull", target)
}

// helper function to create the docker push command.
func commandPush(build Promote, tag string) *exec.Cmd {
	target := fmt.Sprintf("%s:%s", build.PushRepo, tag)
	return exec.Command(dockerExe, "push", target)
}

// helper function to create the docker daemon command.
func commandDaemon(daemon Daemon) *exec.Cmd {
	args := []string{"-g", daemon.StoragePath}

	if daemon.StorageDriver != "" {
		args = append(args, "-s", daemon.StorageDriver)
	}
	if daemon.Insecure && daemon.Registry != "" {
		args = append(args, "--insecure-registry", daemon.Registry)
	}
	for _, registry := range daemon.InsecureRegistries {
		args = append(args, "--insecure-registry", registry)
	}
	if daemon.IPv6 {
		args = append(args, "--ipv6")
	}
	if len(daemon.Mirror) != 0 {
		args = append(args, "--registry-mirror", daemon.Mirror)
	}
	if len(daemon.Bip) != 0 {
		args = append(args, "--bip", daemon.Bip)
	}
	for _, dns := range daemon.DNS {
		args = append(args, "--dns", dns)
	}
	for _, dnsSearch := range daemon.DNSSearch {
		args = append(args, "--dns-search", dnsSearch)
	}
	if len(daemon.MTU) != 0 {
		args = append(args, "--mtu", daemon.MTU)
	}
	if daemon.Experimental {
		args = append(args, "--experimental")
	}
	return exec.Command(dockerdExe, args...)
}

func commandPrune() *exec.Cmd {
	return exec.Command(dockerExe, "system", "prune", "-f")
}

func commandRmi(tag string) *exec.Cmd {
	return exec.Command(dockerExe, "rmi", tag)
}

// trace writes each command to stdout with the command wrapped in an xml
// tag so that it can be extracted and displayed in the logs.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}
