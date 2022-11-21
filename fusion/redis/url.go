package redis

import (
	"bytes"
	"net/url"
)

type RedisURL struct {
	url *url.URL
}

func ParseRedisURL(rawURL string) (*RedisURL, error) {
	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return &RedisURL{uri}, nil
}

func (u *RedisURL) Password() string {
	if user := u.url.User; user != nil {
		if pwd, ok := user.Password(); ok {
			if p, err := url.QueryUnescape(pwd); err == nil {
				return p
			}
		}
	}
	return ""
}

func (u *RedisURL) String() string {
	var buf bytes.Buffer
	if u.url.Host != "" {
		if h := u.url.Host; h != "" {
			buf.WriteString(escape(h, encodeHost))
		}
	}
	path := u.url.EscapedPath()
	if path != "" && path[0] != '/' && u.url.Host != "" {
		buf.WriteByte('/')
	}
	buf.WriteString(path)
	if u.url.RawQuery != "" {
		buf.WriteByte('?')
		buf.WriteString(u.url.RawQuery)
	}
	if u.url.Fragment != "" {
		buf.WriteByte('#')
		buf.WriteString(escape(u.url.Fragment, encodeFragment))
	}
	return buf.String()
}

type encoding int

const (
	encodeHost encoding = 1 + iota
	encodeQueryComponent
	encodeFragment
)

// Return true if the specified character should be escaped when
// appearing in a URL string, according to RFC 3986.
//
// Please be informed that for now shouldEscape does not check all
// reserved characters correctly. See golang.org/issue/5684.
func shouldEscape(c byte, mode encoding) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}

	if mode == encodeHost {
		// §3.2.2 Host allows
		//	sub-delims = "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
		// as part of reg-name.
		// We add : because we include :port as part of host.
		// We add [ ] because we include [ipv6]:port as part of host.
		// We add < > because they're the only characters left that
		// we could possibly allow, and Parse will reject them if we
		// escape them (because hosts can't use %-encoding for
		// ASCII bytes).
		switch c {
		case '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', ':', '[', ']', '<', '>', '"':
			return false
		}
	}

	switch c {
	case '-', '_', '.', '~': // §2.3 Unreserved characters (mark)
		return false

	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
		// Different sections of the URL allow a few of
		// the reserved characters to appear unescaped.
		switch mode {
		case encodeQueryComponent: // §3.4
			// The RFC reserves (so we must escape) everything.
			return true

		case encodeFragment: // §4.1
			// The RFC text is silent but the grammar allows
			// everything, so escape nothing.
			return false
		}
	}

	// Everything else must be escaped.
	return true
}

func escape(s string, mode encoding) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c, mode) {
			if c == ' ' && mode == encodeQueryComponent {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == ' ' && mode == encodeQueryComponent:
			t[j] = '+'
			j++
		case shouldEscape(c, mode):
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}
