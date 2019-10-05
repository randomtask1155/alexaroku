# RokuRemote

this remote is incomplete and for personal use

## Usage

```
package main

import (
    "fmt"
    roku "github.com/randomtask1155/rokuremote"
)

func main() {
    rokuPlayer, err := roku.Connect(rokuIPAddress)
	if err != nil {
		fmt.Println(err)
	}

    rokuPlayer.StartChannel(roku.Netflix)
    rokuPlayer.Home()
}
```