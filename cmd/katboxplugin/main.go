/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/paypal/katbox/pkg/katbox"
)

func init() {
	flag.Set("logtostderr", "true")
}

var (
	endpoint          = flag.String("endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	driverName        = flag.String("drivername", "kbox.csi.paypal.com", "name of the driver")
	nodeID            = flag.String("nodeid", "", "node id")
	maxVolumesPerNode = flag.Int64("maxvolumespernode", 0, "limit of volumes per node")
	workdir           = flag.String("workdir", "/csi-data-dir", "Location where plugin will store configuration and directories")
	afterLifespan     = flag.Duration(
		"afterlifespan",
		time.Hour*12,
		"Length of time to keep a volume after a request for deletion",
	)
	pruneInterval = flag.Duration(
		"pruneinterval",
		time.Second*5,
		"Interval at which the background process looking to evict deleted volumes runs.",
	)
	headroom = flag.Float64(
		"headroom",
		0.1,
		"Value between 0.0 and 1.0 (inclusive) that determines the percentage of space that should be attempted to be kept free in the underlying storage device",
	)
	showVersion = flag.Bool("version", false, "Show version.")
	// Set by the build process
	version = ""
)

func main() {
	flag.Parse()

	if *showVersion {
		baseName := path.Base(os.Args[0])
		fmt.Println(baseName, version)
		return
	}

	handle()
	os.Exit(0)
}

func handle() {
	driver, err := katbox.NewKatboxDriver(
		*driverName,
		*nodeID,
		*endpoint,
		*workdir,
		*maxVolumesPerNode,
		*afterLifespan,
		*pruneInterval,
		*headroom,
		version,
	)
	if err != nil {
		fmt.Printf("Failed to initialize driver: %s", err.Error())
		os.Exit(1)
	}
	driver.Run()
}
