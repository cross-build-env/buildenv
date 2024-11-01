package io

import "testing"

func TestDownload(t *testing.T) {
	url := "http://192.168.100.25:8083/repository/build_resource/rootfs/ti-processor-sdk-rtos-j721e-evm-07_03_00_07.tar.gz"
	path, err := Download(url, "testdata/temp", func(percent int) {
		t.Logf("progress: %v", percent)
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("path: %v", path)
}
