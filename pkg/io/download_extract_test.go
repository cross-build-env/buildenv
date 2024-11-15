package io

import (
	"os"
	"testing"
)

func TestExtractTarGz(t *testing.T) {
	url := "http://192.168.100.25:8083/repository/build_resource/rootfs/ti-processor-sdk-rtos-j721e-evm-07_03_00_07.tar.gz"
	path, err := Download(url, "testdata/temp")
	if err != nil {
		t.Fatal(err)
	}

	if path != "testdata/temp/ti-processor-sdk-rtos-j721e-evm-07_03_00_07.tar.gz" {
		t.Logf("path: %v", path)
	}

	if err := Extract(path, "testdata/temp/ti-processor-sdk-rtos-j721e-evm-07_03_00_07"); err != nil {
		t.Fatal(err)
	}

	if err := os.RemoveAll("testdata/temp"); err != nil {
		t.Fatal(err)
	}
}
