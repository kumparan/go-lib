package url

import (
	"errors"
	"fmt"
	"strings"
)

var usesRelative = []string{"", "ftp", "http", "gopher", "nntp", "imap",
	"wais", "file", "https", "shttp", "mms",
	"prospero", "rtsp", "rtspu", "sftp",
	"svn", "svn+ssh", "ws", "wss"}

var usesNetloc = []string{"", "ftp", "http", "gopher", "nntp", "telnet",
	"imap", "wais", "file", "mms", "https", "shttp",
	"snews", "prospero", "rtsp", "rtspu", "rsync",
	"svn", "svn+ssh", "sftp", "nfs", "git", "git+ssh",
	"ws", "wss"}

var schemeChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+-."

// Join :nodoc:
func Join(base, url string) (string, error) {
	if len(base) == 0 {
		return url, nil
	}

	if len(url) == 0 {
		return base, nil
	}

	bscheme, bnetloc, bpath, err := urlSplit(base)
	if err != nil {
		return "", err
	}
	scheme, netloc, path, err := urlSplit(url)
	if err != nil {
		return "", err
	}

	if scheme == "" {
		scheme = bscheme
	}

	if scheme != bscheme || !existInArray(usesRelative, scheme) {
		return url, nil
	}

	if existInArray(usesNetloc, scheme) {
		if len(netloc) > 0 {
			return urlUnsplit(scheme, netloc, path), nil
		}
		netloc = bnetloc
	}

	if len(path) == 0 {
		path = bpath
		return urlUnsplit(scheme, netloc, path), nil
	}

	baseParts := strings.Split(bpath, "/")
	if len(baseParts[len(baseParts)-1]) > 0 {
		baseParts = baseParts[:len(baseParts)-1]
	}

	var segments []string
	if path[:1] == "/" {
		segments = strings.Split(path, "/")
	} else {
		splitPath := strings.Split(path, "/")
		segments = baseParts
		for _, v := range splitPath {
			segments = append(segments, v)
		}
		//FILTER
	}

	resolvedPath := make([]string, 0)
	for _, v := range segments {
		if v == ".." {
			return "", errors.New("error when resolving path")
		} else if v == "." {
			continue
		} else {
			resolvedPath = append(resolvedPath, v)
		}
	}

	if segments[len(segments)-1] == "." || segments[len(segments)-1] == ".." {
		resolvedPath = append(resolvedPath, "")
	}

	path = strings.Join(resolvedPath, "/")
	return urlUnsplit(scheme, netloc, path), nil
}

func urlSplit(url string) (scheme, netloc, path string, err error) {
	if strings.Contains(url, ":") {
		var posColon int
		for i, v := range url {
			// 58 rune for ":"
			if v == 58 {
				posColon = i
			}
		}
		for _, v := range url[:posColon] {
			if !strings.Contains(schemeChars, string(v)) {
				break
			}
		}
		scheme = strings.ToLower(url[:posColon])
		url = url[posColon+1:]
	}

	if url[:2] == "//" {
		netloc, url = splitnetloc(url, 2)
		if (strings.Contains(netloc, "[") && !strings.Contains(netloc, "]")) ||
			(!strings.Contains(netloc, "[") && strings.Contains(netloc, "]")) {
			return "", "", "", errors.New("error when url split")
		}
	}

	if url == "/" {
		url = ""
	}

	return scheme, netloc, url, nil
}

func splitnetloc(url string, start int) (domain, rest string) {
	delim := len(url)
	c := "/?#"
	for _, v := range c {
		var wdelim int
		wdelim = strings.Index(url[2:], string(v))
		if wdelim >= 0 {
			if delim >= wdelim {
				delim = wdelim + 2
			}
		}
	}
	return url[start:delim], url[delim:]
}

func urlUnsplit(scheme, netloc, path string) (res string) {
	if len(scheme) > 0 {
		res = scheme + "://" + netloc
	} else {
		res = netloc
	}

	fmt.Println("dalem:", string(path[0]))
	if string(path[0]) == "/" {
		res += path
		return
	}

	if res != "" {
		res = res + "/" + path
		return
	}

	res = path
	return
}

func existInArray(arr []string, body string) bool {
	for _, v := range arr {
		if v == body {
			return true
		}
	}
	return false
}
