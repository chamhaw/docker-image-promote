package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"

	"github.com/drone-plugins/image-promote"
)

var build = "0" // build number set at compile-time

func main() {
	// Load env-file if it exists first
	if env := os.Getenv("PLUGIN_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	app := cli.NewApp()
	app.Name = "docker image promote plugin"
	app.Usage = "docker image promote plugin"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "dry-run",
			Usage:  "dry run disables docker push",
			EnvVar: "PLUGIN_DRY_RUN",
		},
		cli.StringFlag{
			Name:   "daemon.mirror",
			Usage:  "docker daemon registry mirror",
			EnvVar: "PLUGIN_MIRROR",
		},
		cli.StringFlag{
			Name:   "daemon.storage-driver",
			Usage:  "docker daemon storage driver",
			EnvVar: "PLUGIN_STORAGE_DRIVER",
		},
		cli.StringFlag{
			Name:   "daemon.storage-path",
			Usage:  "docker daemon storage path",
			Value:  "/var/lib/docker",
			EnvVar: "PLUGIN_STORAGE_PATH",
		},
		cli.StringFlag{
			Name:   "daemon.bip",
			Usage:  "docker daemon bride ip address",
			EnvVar: "PLUGIN_BIP",
		},
		cli.StringFlag{
			Name:   "daemon.mtu",
			Usage:  "docker daemon custom mtu setting",
			EnvVar: "PLUGIN_MTU",
		},
		cli.StringSliceFlag{
			Name:   "daemon.dns",
			Usage:  "docker daemon dns server",
			EnvVar: "PLUGIN_CUSTOM_DNS",
		},
		cli.StringSliceFlag{
			Name:   "daemon.dns-search",
			Usage:  "docker daemon dns search domains",
			EnvVar: "PLUGIN_CUSTOM_DNS_SEARCH",
		},
		cli.StringSliceFlag{
			Name:   "daemon.insecureRegistries",
			Usage:  "docker daemon dns insecure registries",
			EnvVar: "PLUGIN_INSECURE_REGISTRIES",
		},
		cli.BoolFlag{
			Name:   "daemon.insecure",
			Usage:  "docker daemon allows insecure registries",
			EnvVar: "PLUGIN_INSECURE",
		},
		cli.BoolFlag{
			Name:   "daemon.ipv6",
			Usage:  "docker daemon IPv6 networking",
			EnvVar: "PLUGIN_IPV6",
		},
		cli.BoolFlag{
			Name:   "daemon.experimental",
			Usage:  "docker daemon Experimental mode",
			EnvVar: "PLUGIN_EXPERIMENTAL",
		},
		cli.BoolFlag{
			Name:   "daemon.debug",
			Usage:  "docker daemon executes in debug mode",
			EnvVar: "PLUGIN_DEBUG,DOCKER_LAUNCH_DEBUG",
		},
		cli.BoolFlag{
			Name:   "daemon.off",
			Usage:  "don't start the docker daemon",
			EnvVar: "PLUGIN_DAEMON_OFF",
		},
		cli.StringSliceFlag{
			Name:     "tags",
			Usage:    "build tags",
			Value:    &cli.StringSlice{"latest"},
			EnvVar:   "PLUGIN_TAG,PLUGIN_TAGS",
			FilePath: ".tags",
		},
		cli.StringFlag{
			Name: "pull-repo",
			Usage: "docker pull repository",
			EnvVar: "PLUGIN_PULL_REPO",
		},
		cli.StringFlag{
			Name:   "push-repo",
			Usage:  "docker push repository",
			EnvVar: "PLUGIN_PUSH_REPO",
		},
		cli.StringFlag{
			Name:   "docker.push.registry",
			Usage:  "docker push registry",
			Value:  "https://index.docker.io/v1/",
			EnvVar: "PLUGIN_PUSH_REGISTRY,DOCKER_REGISTRY",
		},
		cli.StringFlag{
			Name:   "docker.push.username",
			Usage:  "docker push username",
			EnvVar: "PLUGIN_PUSH_USERNAME,DOCKER_USERNAME",
		},
		cli.StringFlag{
			Name:   "docker.push.password",
			Usage:  "docker push password",
			EnvVar: "PLUGIN_PUSH_PASSWORD,DOCKER_PASSWORD",
		},
		cli.StringFlag{
			Name:   "docker.push.email",
			Usage:  "docker push email",
			EnvVar: "PLUGIN_PUSH_EMAIL,DOCKER_EMAIL",
		},
		cli.StringFlag{
			Name:   "docker.pull.registry",
			Usage:  "docker pull registry",
			Value:  "https://index.docker.io/v1/",
			EnvVar: "PLUGIN_PULL_REGISTRY,DOCKER_REGISTRY",
		},
		cli.StringFlag{
			Name:   "docker.pull.username",
			Usage:  "docker pull username",
			EnvVar: "PLUGIN_PULL_USERNAME,DOCKER_USERNAME",
		},
		cli.StringFlag{
			Name:   "docker.pull.password",
			Usage:  "docker pull password",
			EnvVar: "PLUGIN_PULL_PASSWORD,DOCKER_PASSWORD",
		},
		cli.StringFlag{
			Name:   "docker.pull.email",
			Usage:  "docker pull email",
			EnvVar: "PLUGIN_PULL_EMAIL,DOCKER_EMAIL",
		},
		cli.BoolTFlag{
			Name:   "docker.purge",
			Usage:  "docker should cleanup images",
			EnvVar: "PLUGIN_PURGE",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := promote.Plugin{
		Dryrun:  c.Bool("dry-run"),
		Cleanup: c.BoolT("docker.purge"),
		PushLogin: promote.Login{
			Registry: c.String("docker.push.registry"),
			Username: c.String("docker.push.username"),
			Password: c.String("docker.push.password"),
			Email:    c.String("docker.push.email"),
		},
		PullLogin: promote.Login{
			Registry: c.String("docker.pull.registry"),
			Username: c.String("docker.pull.username"),
			Password: c.String("docker.pull.password"),
			Email:    c.String("docker.pull.email"),
		},
		Promote: promote.Promote{
			Tags:        c.StringSlice("tags"),
			PushRepo:    c.String("push-repo"),
			PullRepo:    c.String("pull-repo"),

		},
		Daemon: promote.Daemon{
			Registry:           c.String("docker.registry"),
			Mirror:             c.String("daemon.mirror"),
			StorageDriver:      c.String("daemon.storage-driver"),
			StoragePath:        c.String("daemon.storage-path"),
			Insecure:           c.Bool("daemon.insecure"),
			Disabled:           c.Bool("daemon.off"),
			IPv6:               c.Bool("daemon.ipv6"),
			Debug:              c.Bool("daemon.debug"),
			Bip:                c.String("daemon.bip"),
			DNS:                c.StringSlice("daemon.dns"),
			DNSSearch:          c.StringSlice("daemon.dns-search"),
			MTU:                c.String("daemon.mtu"),
			Experimental:       c.Bool("daemon.experimental"),
			InsecureRegistries: c.StringSlice("daemon.insecureRegistries"),
		},
	}
	if len(plugin.Promote.Tags) == 0 {
		plugin.Promote.Tags = []string{"latest"}
	}
	return plugin.Exec()
}
