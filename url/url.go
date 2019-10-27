package url

import (
	"strings"
)

// Join :nodoc:
func Join(base, url string) string {
	if len(base) == 0 {
		return url
	}

	if len(url) == 0 {
		return base
	}

	bscheme, bnetloc, bpath := urlsplit(base)
	scheme, netloc, path := urlsplit(url)

	// if url have a base
	if scheme != "" {
		return url
	}

	// if url is a directory
	if netloc == "" && len(path) > 0 {
		// if base have no scheme
		if bscheme == "" {
			return url
		}
		// delete path for base if url is a directory
		bpath = nil
	}

	// remove last path from base path
	if bpath != nil && len(bpath) > 0 {
		bpath = bpath[:len(bpath)-1]
	}

	// append netloc to base path
	if netloc != "" {
		bpath = append(bpath, netloc)
	}

	// append url path to base path
	if len(path) > 0 {
		for _, v := range path {
			bpath = append(bpath, v)
		}
	}

	// join base path into /./. format
	paths := strings.Join(bpath, "/")

	// if base scheme not empty
	if bscheme != "" {
		return bscheme + "://" + bnetloc + "/" + paths
	}

	return bnetloc + "/" + paths
}

func urlsplit(url string) (scheme, netloc string, path []string) {
	// <scheme>://<netloc>/<path>
	if strings.Contains(url, "://") {
		arr := strings.Split(url, "://")
		arr1 := strings.Split(arr[1], "/")
		return arr[0], arr1[0], arr1[1:]
	}

	arr := strings.Split(url, "/")
	return "", arr[0], arr[1:]
}
