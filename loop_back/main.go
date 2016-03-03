package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/melvinw/go-dpdk"
)

type loopCtx struct {
	portNum uint
}

func loop(arg unsafe.Pointer) int {
	ctx := (*loopCtx)(arg)
	n := uint(32)
	tbl := dpdk.GetCArray(n)

	for {
		nb_rx := dpdk.RteEthRxBurst(ctx.portNum, 0, tbl, 32)
		dpdk.RteEthTxBurst(ctx.portNum, 0, tbl, nb_rx)
	}
	return 0
}

func main() {
	ret := dpdk.RteEalInit(os.Args)
	if ret < 0 {
		fmt.Printf("Failed to init EAL: %s\n", dpdk.StrError(ret))
		return
	}

	numDevs := dpdk.RteEthDevCount()
	if numDevs < 1 {
		fmt.Println("No NICs found")
		return
	}

	/* Setup each NIC with 1 TX queue and 1 RX queue*/
	for i := uint(0); i < numDevs; i++ {
		mpName := fmt.Sprintf("p%d_rx0", i)
		rx_mp := dpdk.RtePktMbufPoolCreate(mpName, 512, 0, 0, 2048, 0)

		ret = dpdk.RteEthDevConfigure(i, 1, 1, &dpdk.RteEthConf{})
		if ret < 0 {
			fmt.Printf("Failed to setup eth dev: %s\n", dpdk.StrError(ret))
			return
		}

		ret = dpdk.RteEthRxQueueSetup(i, 0, 512, 0, nil, rx_mp)
		if ret < 0 {
			fmt.Println("Failed to setup rx queue")
			return
		}

		ret = dpdk.RteEthTxQueueSetup(i, 0, 512, 0, nil)
		if ret < 0 {
			fmt.Println("Failed to setup tx queue")
			return
		}

		ret = dpdk.RteEthDevStart(i)
		if ret < 0 {
			fmt.Println("Failed to setup tx queue")
			return
		}

		dpdk.RteEthPromiscuousEnable(i)

		c := loopCtx{i}
		dpdk.RteEalRemoteLaunch(loop, unsafe.Pointer(&c), i+1)
	}

	dpdk.RteEalMpWaitLCore()
}
