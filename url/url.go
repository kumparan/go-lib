package url

import (
	"errors"
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

// Join is to join base URL with another URL/path/directory based on python lib urljoin
// e.g. Join(https://kumparan.com/trending , feed/category) => https://kumparan.com/feed/category
func Join(base, url string) (string, error) {
	if len(base) == 0 {
		return url, nil
	}
	if len(url) == 0 {
		return base, nil
	}

	// Split baseURL and url into 3 components
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

	// if url have a different scheme than baseURL return url
	if scheme != bscheme || !existInArrayString(usesRelative, scheme) {
		return url, nil
	}

	if existInArrayString(usesNetloc, scheme) {
		if len(netloc) > 0 {
			return urlUnsplit(scheme, netloc, path), nil
		}
		netloc = bnetloc
	}

	if len(path) == 0 {
		path = bpath
		return urlUnsplit(scheme, netloc, path), nil
	}

	// to get path from baseURL and remove the last one
	baseParts := strings.Split(bpath, "/")
	if len(baseParts[len(baseParts)-1]) > 0 {
		baseParts = baseParts[:len(baseParts)-1]
	}

	// append path to path from baseURL
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

	// to pop path if ".." occurs and ignore if "."  occurs
	resolvedPath := make([]string, 0)
	for _, v := range segments {
		if v == ".." {
			resolvedPath = resolvedPath[:len(resolvedPath)-1]
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

// urlSplit is to split url into 3 components (scheme, netloc, path)
// <scheme>://<netloc>/<path>
func urlSplit(url string) (scheme, netloc, path string, err error) {
	// check and get scheme
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

	// splitting between netloc and path
	if url[:2] == "//" {
		netloc, url = splitnetloc(url, 2)
		if (strings.Contains(netloc, "[") && !strings.Contains(netloc, "]")) ||
			(!strings.Contains(netloc, "[") && strings.Contains(netloc, "]")) {
			return "", "", "", errors.New("error when url split")
		}
	}

	// path is the rest of it
	if url == "/" {
		url = ""
	}

	return scheme, netloc, url, nil
}

// splitnetloc is to get netloc
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

// urlUnsplit is to join the url from 3 components
func urlUnsplit(scheme, netloc, path string) (url string) {
	if len(scheme) > 0 {
		url = scheme + "://" + netloc
	} else {
		url = netloc
	}

	if len(path) > 0 {
		if string(path[0]) == "/" {
			url += path
			return
		}
	}

	if url != "" {
		url = url + "/" + path
		return
	}

	url = path
	return
}

// existInArrayString is to check string exist in an array of string
func existInArrayString(arr []string, body string) bool {
	for _, v := range arr {
		if v == body {
			return true
		}
	}
	return false
}
