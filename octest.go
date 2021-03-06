//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
        "github.com/codegangsta/cli"
        "os"
)

func main() {

        app := cli.NewApp()
        app.Name = "octest"
        app.Usage = "Tools for OCI specs"
        app.Version = "0.1.0"
        app.Commands = []cli.Command{
                {
                        Name:  "validate",
                        Usage: "validate container image / Json",
                        Flags: []cli.Flag{
                                cli.StringFlag{
                                        Name:  "json",
                                        Usage: "json config file to validate",
                                },
                                cli.StringFlag{
                                        Name: "layout",
                                        Usage: "directory layout to validate",
                                },
                        },
                        Action: validate,
                },
                {
                        Name:   "test",
                        Usage:  "Test the Container",
                        Action: testOContainer,
                },
        }

        app.Run(os.Args)
}
