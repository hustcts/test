package main

import (
        "bytes"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "os"
        "runtime"
        "regexp"

        "github.com/astaxie/beego/validation"
        "github.com/codegangsta/cli"
        "github.com/opencontainers/specs"
        "errors"
        "io"
        "path/filepath"
        "strings"
)

const (
        // Path to config file inside the layout
        ConfigFile = "config.json"
        // Path to rootfs directory inside the layout
        RootfsDir = "rootfs"
)

var (
        ErrNoRootFS   = errors.New("no rootfs found in layout")
        ErrNoConfig = errors.New("no config json file found in layout")
)

type Config struct {
        configLinux specs.LinuxSpec
}

func validate(c *cli.Context) {
    args := c.String("json")

    if len(args) == 0 {
        args = c.String("layout")
        if len(args) == 0 {
            cli.ShowCommandHelp(c, "validate")
            return
        } else {
           err := validateLayout(args)
           if err != nil {
                                fmt.Printf("%s: invalid image layout: %v\n", args, err)
                        } else {
                                fmt.Printf("%s: valid image layout\n", args)
                        }
        }
    } else {
               validateConfigFile(args)
    }


}

func validateLayout(path string) error {
        fi, err := os.Stat(path)
        if err != nil {
                return fmt.Errorf("error accessing layout: %v", err)
        }
        if !fi.IsDir() {
                return fmt.Errorf("given path %q is not a directory", path)
        }
        var flist []string
        var imOK, rfsOK bool
        var im io.Reader
        walkLayout := func(fpath string, fi os.FileInfo, err error) error {
                rpath, err := filepath.Rel(path, fpath)
                if err != nil {
                        return err
                }
                switch rpath {
                case ".":
                case ConfigFile:
                        im, err = os.Open(fpath)
                        if err != nil {
                                return err
                        }
                        imOK = true
                case RootfsDir:
                        if !fi.IsDir() {
                                return errors.New("rootfs is not a directory")
                        }
                        rfsOK = true
                default:
                        flist = append(flist, rpath)
                }
                return nil
        }
        if err := filepath.Walk(path, walkLayout); err != nil {
                return err
        }
        return checkLayout(imOK, im, rfsOK, flist)
}

func checkLayout(imOK bool, im io.Reader, rfsOK bool, files []string) error {
        defer func() {
                if rc, ok := im.(io.Closer); ok {
                        rc.Close()
                }
        }()
        if !imOK {
                return ErrNoConfig
        }
        if !rfsOK {
                return ErrNoRootFS
        }
        _, err := ioutil.ReadAll(im)
        if err != nil {
                return fmt.Errorf("error reading the layout: %v", err)
        }

        for _, f := range files {
                if !strings.HasPrefix(f, "rootfs") {
                        return fmt.Errorf("unrecognized file path in layout: %q", f)
                }
        }
        return nil
}


func validateConfigFile(path string) {

        // Read json file and load into spec.LinuxSpec struct.
        // Validate as per rules and spec defination.
        config, err := NewConfig(path)
        if err != nil {
                fmt.Println("Error while opening File ", err)
                return
        }
        common := config.ValidateCommonSpecs()
        platform := config.ValidateLinuxSpecs()
        if !common || !platform {
                fmt.Println("\nNOTE: One or more errors found in", path)
        } else {
                fmt.Println("\n", path, "has Valid OC Format !!")
        }

        return

}

func dumpJSON(config Config) {
        b, err := json.Marshal(config.configLinux)
        if err != nil {
                fmt.Println(err)
                return
        }
        var out bytes.Buffer
        json.Indent(&out, b, "", "\t")
        out.WriteTo(os.Stdout)
}

func getOS() string {
        return runtime.GOOS
}

func NewConfig(path string) (Config, error) {

        data, err := ioutil.ReadFile(path)
        if err != nil {
                return Config{}, err
        }

        config := Config{}
        if getOS() == "linux" {
                json.Unmarshal(data, &config.configLinux)
        }

        return config, nil
}

func (conf *Config) ValidateCommonSpecs() bool {
        valid := validation.Validation{}

        //Validate mandatory fields.
        valid.Required(conf.configLinux.Version, "Version")
        //Version must complient with  SemVer v2.0.0
        valid.Match(conf.configLinux.Version, regexp.MustCompile("^(\\d+\\.)?(\\d+\\.)?(\\*|\\d+)$"),"Version")
        valid.Required(conf.configLinux.Platform.OS, "OS")
        valid.Required(conf.configLinux.Platform.Arch, "Platform.Arch")

        for _, env := range conf.configLinux.Process.Env {
                //If Process defined, env cannot be empty
                valid.Required(env, "Process.Env")
        }
        valid.Required(conf.configLinux.Process.User.UID, "User.UID")
        valid.Required(conf.configLinux.Process.User.GID, "User.GID")
        valid.Required(conf.configLinux.Root.Path, "Root.Path")
        //Iterate over Mount array
        for _, mount := range conf.configLinux.Mounts {
                //If Mount points defined, it must define these three.
                valid.Required(mount.Type, "Mount.Type")
                valid.Required(mount.Source, "Mount.Source")
                valid.Required(mount.Destination, "Mount.Destination")
        }

        if valid.HasErrors() {
                // validation does not pass
                for i, err := range valid.Errors {
                        fmt.Println(i, err.Key, err.Message)
                }
                return false
        }

        return true
}

func (conf *Config) ValidateLinuxSpecs() bool {
        valid := validation.Validation{}

        for _, namespace := range conf.configLinux.Linux.Namespaces  {
                valid.Required(namespace.Type, "Namespace.Type")
        }


        if valid.HasErrors() {
                // validation does not pass
                fmt.Println("\nLinux Specific config errors\n")

                for i, err := range valid.Errors {
                        fmt.Println(i, err.Key, err.Message)
                }
                return false
        }

        return true
}

func (conf *Config) Analyze() {
        fmt.Println("NOT-IMPLEMETED")
        return
}

func testOContainer(c *cli.Context) {
        fmt.Println("NOT-IMPLEMENTED")
        return
}
