package main

import (
	"syscall"
	"time"
	"unsafe"
)

type Process struct {
	ID              uint32
	ParentProcessID uint32
	CreateTime      *time.Time
	ExecName        string
}

func GetAllProcess() ([]*Process, error) {
	handle, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	defer syscall.CloseHandle(handle)
	if err != nil {
		return nil, err
	}
	var entry syscall.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	err = syscall.Process32First(handle, &entry)
	if err != nil {
		return nil, err
	}
	r := make([]*Process, 0)
	for {
		process := newProcess(&entry)
		if process.CreateTime != nil {
			r = append(r, process)
		}

		err := syscall.Process32Next(handle, &entry)
		if err != nil {
			// windows sends ERROR_NO_MORE_FILES on last process
			if err == syscall.ERROR_NO_MORE_FILES {
				return r, nil
			}
			return nil, err
		}
	}
}

func newProcess(e *syscall.ProcessEntry32) *Process {
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}
	p := &Process{
		ID:              e.ProcessID,
		ParentProcessID: e.ParentProcessID,
		ExecName:        syscall.UTF16ToString(e.ExeFile[:end]),
	}
	var creationTime, exitTime, kernelTime, userTime syscall.Filetime
	process, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, e.ProcessID)
	defer syscall.CloseHandle(process)
	if err == nil && syscall.GetProcessTimes(process, &creationTime, &exitTime, &kernelTime, &userTime) == nil {
		unix := time.Unix(0, creationTime.Nanoseconds())
		p.CreateTime = &unix
	}
	return p
}
