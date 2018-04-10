// +build linux

package process

import (
	"os"
	"syscall"

	"github.com/opencontainers/runc/libcontainer/system"
	"github.com/opencontainers/runc/libcontainer/user"
)

func Su(userNameOrUidGid string) error {
	os.Unsetenv("HOME")
	currentUidGid := user.ExecUser{
		Uid:  syscall.Getuid(),
		Gid:  syscall.Getgid(),
		Home: "/",
	}
	passwdPath, err := user.GetPasswdPath()
	if err != nil {
		return err
	}
	groupPath, err := user.GetGroupPath()
	if err != nil {
		return err
	}
	effectiveUser, err := user.GetExecUserPath(userNameOrUidGid, &currentUidGid, passwdPath, groupPath)
	if err != nil {
		return err
	}
	if err := syscall.Setgroups(effectiveUser.Sgids); err != nil {
		return err
	}
	if err := system.Setgid(effectiveUser.Gid); err != nil {
		return err
	}
	if err := system.Setuid(effectiveUser.Uid); err != nil {
		return err
	}
	os.Setenv("HOME", effectiveUser.Home)
	return nil
}
