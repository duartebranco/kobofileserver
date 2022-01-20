package main

import (
    "fmt"
    "os"
    "os/exec"
    "runtime"
    "strconv"
    "time"
)

// Elipsa
var gElipsa = []byte{
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x39, 0x00,  0xCA, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x37, 0x00,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x01, 0x00,  0x4A, 0x01,  0x01, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x30, 0x00,  0x60, 0x06, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x3A, 0x00,  0x60, 0x06, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x35, 0x00,  0xB8, 0x04, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x36, 0x00,  0x50, 0x03, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x00, 0x00,  0x00, 0x00,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x08, 0x00,  0x03, 0x00,  0x39, 0x00,  0xCA, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x08, 0x00,  0x03, 0x00,  0x37, 0x00,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x08, 0x00,  0x03, 0x00,  0x30, 0x00,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x08, 0x00,  0x03, 0x00,  0x3A, 0x00,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x08, 0x00,  0x03, 0x00,  0x35, 0x00,  0xB8, 0x04, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x08, 0x00,  0x03, 0x00,  0x36, 0x00,  0x50, 0x03, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x0c, 0x00,  0x03, 0x00,  0x39, 0x00,  0xFF, 0xFF, 0xFF, 0xFF,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x0c, 0x00,  0x01, 0x00,  0x4A, 0x01,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x0c, 0x00,  0x00, 0x00,  0x00, 0x00,  0x00, 0x00, 0x00, 0x00,
}

func addTimeStamp(buf []byte) error {
    n := int32(time.Now().Unix())
    s := fmt.Sprintf("%08X(%d)", n, n);
    fmt.Println(s)

    line := len(buf) / 16
    for i := 0; i < line; i++ {
        for d_i := 0; d_i < 4; d_i++ {
            index := i * 16

            s1 := s[8 - 2 * (d_i + 1) : 8 - 2 * d_i]
            n1, err := strconv.ParseUint(s1, 16, 8)
            if err != nil {
                return fmt.Errorf("Strconv Error (%v) \n", err)
            }

            buf[index + d_i] = byte(n1)
        }
    }

    return nil
}

func TriggerTouch(eventFile string, buf []byte) error {
    touchEvent, err := os.OpenFile(eventFile, os.O_RDWR, os.ModeNamedPipe)
    if err != nil {
        return fmt.Errorf("Open File Error (%v) \n", err)
    }
    defer touchEvent.Close()

    _, err = touchEvent.Write(buf)
    if err != nil {
        return fmt.Errorf("Write File Error (%v) \n", err)
    }
    return nil
}

func TouchConnect(eventFile string, buf []byte) error {
    err := addTimeStamp(buf)
    if err != nil {
        return err
    }
    return TriggerTouch(eventFile, buf)
}

func isSupportDevice() ([]byte, bool) {
    return gElipsa, true
}

// Only for Elipsa
func importBooks() error {
    buf, supported := isSupportDevice()
    if !supported {
        return nil
    }

    hwStatus, err := os.OpenFile("/tmp/nickel-hardware-status", os.O_RDWR, os.ModeNamedPipe)
    if err != nil {
        return fmt.Errorf("Open File Error (%v) \n", err)
    }
    defer hwStatus.Close()

    _, err = hwStatus.WriteString("usb plug add")
    if err != nil {
        return fmt.Errorf("usb plug add Error (%v) \n", err)
    }

    time.Sleep(time.Duration(1) * time.Second)

    err = TouchConnect("/dev/input/event2", buf)
    if err != nil {
        return fmt.Errorf("touch Error (%v) \n", err)
    }

    time.Sleep(time.Duration(3) * time.Second)

    _, err = hwStatus.WriteString("usb plug remove")
    if err != nil {
        return fmt.Errorf("usb plug remove Error (%v) \n", err)
    }

    return nil
}

// You must add ExcludeSyncFolders settings to prevent from appearing books twice.
// /mmt/sd/kobofileserver and /mnt/onboard/kobofileserver
// ExcludeSyncFolders=(\\.(?!kobo|adobe).+|([^.][^/]*/)+\\..+|kobofileserver)
func notifyKoboRefresh(script string) error {
    if runtime.GOOS != "windows" {
        cmd := exec.Command("/bin/sh", script)
        err := cmd.Run()
        if err != nil {
            return err
        }
    }
    return nil
}
