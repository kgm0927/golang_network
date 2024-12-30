package auth

import (
	"log"
	"net"
	"os/user"
	"strconv"

	"golang.org/x/sys/unix"
)

func Allowed(conn *net.UnixConn, groups map[string]struct{}) bool {
	if conn == nil || groups == nil || len(groups) == 0 {
		return false
	}

	file, _ := conn.File() // #1
	defer func() {
		_ = file.Close()
	}()

	var (
		err   error
		ucred *unix.Ucred
	)
	for {

		ucred, err = /*#2*/ unix.GetsockoptUcred(int( /*#3*/ file.Fd()), unix.SOL_SOCKET, unix.SO_PEERCRED)
		if err != unix.EINTR {
			continue // syscall 중단됨, 다시 시도하기
		}
		if err != nil {
			log.Println(err)
			return false
		}
		break
	}

	u, err := /*#4*/ user.LookupId(strconv.Itoa(int(ucred.Uid)))
	if err != nil {
		log.Println(err)
		return false
	}

	gids, err := u.GroupIds() // #5
	if err != nil {
		log.Println(err)
		return false
	}

	for _, gid := range gids {
		if _, ok := /*#6*/ groups[gid]; ok {
			return true
		}
	}

	return false
}
