package main

import (
	"fmt"

	raminfoprovider "dbs/client/internal/RAM_info_provider"
)

func main() {
	ramInfoProvider := raminfoprovider.RAMInfoProvider{}

	fmt.Println(ramInfoProvider.GetInfo())
}
