package sftp

import (
	"fmt"
	"time"
)


const (
	SSH_FILEXFER_ATTR_SIZE              = 0x00000001
	SSH_FILEXFER_ATTR_PERMISSIONS       = 0x00000004
	SSH_FILEXFER_ATTR_ACCESSTIME        = 0x00000008
	SSH_FILEXFER_ATTR_CREATETIME        = 0x00000010
	SSH_FILEXFER_ATTR_MODIFYTIME        = 0x00000020
	SSH_FILEXFER_ATTR_ACL               = 0x00000040
	SSH_FILEXFER_ATTR_OWNERGROUP        = 0x00000080
	SSH_FILEXFER_ATTR_SUBSECOND_TIMES   = 0x00000100
	SSH_FILEXFER_ATTR_BITS              = 0x00000200	// don't care
	SSH_FILEXFER_ATTR_ALLOCATION_SIZE   = 0x00000400	// don't care
	SSH_FILEXFER_ATTR_TEXT_HINT         = 0x00000800	// don't care
	SSH_FILEXFER_ATTR_MIME_TYPE         = 0x00001000	// don't care
	SSH_FILEXFER_ATTR_LINK_COUNT        = 0x00002000	// don't care
	SSH_FILEXFER_ATTR_UNTRANSLATED_NAME = 0x00004000	// don't care
	SSH_FILEXFER_ATTR_CTIME             = 0x00008000
	SSH_FILEXFER_ATTR_EXTENDED          = 0x80000000

	SSH_FILEXFER_TYPE_REGULAR           = 1
	SSH_FILEXFER_TYPE_DIRECTORY         = 2
	SSH_FILEXFER_TYPE_SYMLINK           = 3
	SSH_FILEXFER_TYPE_SPECIAL           = 4
	SSH_FILEXFER_TYPE_UNKNOWN           = 5
	SSH_FILEXFER_TYPE_SOCKET            = 6
	SSH_FILEXFER_TYPE_CHAR_DEVICE       = 7
	SSH_FILEXFER_TYPE_BLOCK_DEVICE      = 8
	SSH_FILEXFER_TYPE_FIFO              = 9
)


func DumpP6Attrs(b []byte) string {
	var s string

	// FLAGS
	flags, b := unmarshalUint32(b)
	fmt.Printf(".flags: %08X\n", flags)

	// TYPE
	typ, b := unmarshalByte(b)
	fmt.Print(".type:  ")
	switch typ {
	case SSH_FILEXFER_TYPE_REGULAR:
		fmt.Print("SSH_FILEXFER_TYPE_REGULAR")
	case SSH_FILEXFER_TYPE_DIRECTORY:
		fmt.Printf("SSH_FILEXFER_TYPE_DIRECTORY")
	}
	fmt.Printf(" (%01X)", typ)
	fmt.Println()

	// SIZE
	if flags&SSH_FILEXFER_ATTR_SIZE == SSH_FILEXFER_ATTR_SIZE {
		var size uint64
		size, b = unmarshalUint64(b)
		fmt.Printf(".size:  %d (%016X)\n", size, size)
	}

	// OWNER/GROUP
	if flags&SSH_FILEXFER_ATTR_OWNERGROUP == SSH_FILEXFER_ATTR_OWNERGROUP {
		var owner, group string
		owner, b = unmarshalString(b)
		group, b = unmarshalString(b)
		fmt.Printf(".owner: %s (%08X)\n.group: %s (%08X)\n", owner, len(owner), group, len(group))
	}

	// PERM
	if flags&SSH_FILEXFER_ATTR_PERMISSIONS == SSH_FILEXFER_ATTR_PERMISSIONS {
		var perm uint32
		perm, b = unmarshalUint32(b)
		fmt.Printf(".perm:  %s (%08X)\n", toFileMode(perm), perm)
	}

	// ATIME
	if flags&SSH_FILEXFER_ATTR_ACCESSTIME == SSH_FILEXFER_ATTR_ACCESSTIME {
		b = dumpTime(b, flags, "atime")
	}

	// CRTIME
	if flags&SSH_FILEXFER_ATTR_CREATETIME == SSH_FILEXFER_ATTR_CREATETIME {
		b = dumpTime(b, flags, "crtim")
	}

	// MTIME
	if flags&SSH_FILEXFER_ATTR_MODIFYTIME == SSH_FILEXFER_ATTR_MODIFYTIME {
		b = dumpTime(b, flags, "mtime")
	}

	// CTIME
	if flags&SSH_FILEXFER_ATTR_CTIME == SSH_FILEXFER_ATTR_CTIME {
		b = dumpTime(b, flags, "ctime")
	}

	// ACL
	if flags&SSH_FILEXFER_ATTR_ACL == SSH_FILEXFER_ATTR_ACL {
		var acl string
		acl, b = unmarshalString(b)
		fmt.Printf(".acl:  %s (%08X)\n", acl, len(acl))
	}

	// EXTENDED
	if flags&SSH_FILEXFER_ATTR_EXTENDED == SSH_FILEXFER_ATTR_EXTENDED {
		var count uint32
		count, b = unmarshalUint32(b)
		ext := make(map[string]string, count)
		for i := uint32(0); i < count; i++ {
			var key, data string
			key, b = unmarshalString(b)
			data, b = unmarshalString(b)
			ext[key] = data
		}
		fmt.Printf(".xattr: (%08X) %s\n", count, ext)
	}

	return s + "\n"
}

func dumpTime(b []byte, flags uint32, name string) []byte {
	var xtime uint64
	xtime, b = unmarshalUint64(b)
	fmt.Printf(".%s: %v (%016X)", name, time.Unix(int64(xtime), 0), xtime)

	if flags&SSH_FILEXFER_ATTR_SUBSECOND_TIMES == SSH_FILEXFER_ATTR_SUBSECOND_TIMES {
		var stime uint32
		stime, b = unmarshalUint32(b)
		fmt.Printf(", subsec %d (%08X)", stime, stime)
	}
	fmt.Println()
	return b
}
