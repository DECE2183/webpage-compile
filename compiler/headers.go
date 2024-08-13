package compiler

type Header = string

const (
	HEADER_OK             Header = "HTTP/1.1 200 OK"
	HEADER_GZIP_ENCODING  Header = "Content-Encoding: gzip"
	HEADER_CONTENT_TYPE   Header = "Content-Type: "
	HEADER_CONTENT_LENGTH Header = "Content-Length: "
)

// MIME types
const (
	// text
	HEADER_TXT  Header = "text/plain; charset=utf-8"
	HEADER_HTML Header = "text/html; charset=utf-8"
	HEADER_CSS  Header = "text/css; charset=utf-8"
	HEADER_JS   Header = "text/javascript; charset=utf-8"
	HEADER_JSON Header = "application/json; charset=utf-8"
	// image
	HEADER_PNG  Header = "image/png"
	HEADER_JPG  Header = "image/jpeg"
	HEADER_GIF  Header = "image/gif"
	HEADER_SVG  Header = "image/svg+xml"
	HEADER_WEBP Header = "image/webp"
	// video
	HEADER_WEBM Header = "video/webm"
	HEADER_MP4  Header = "video/mp4"
	// audio
	HEADER_WEBA Header = "audio/weba"
	HEADER_MP3  Header = "audio/mpeg"
	HEADER_OGG  Header = "application/ogg"
	// font
	HEADER_TTF   Header = "font/ttf"
	HEADER_OTF   Header = "font/otf"
	HEADER_WOFF  Header = "font/woff"
	HEADER_WOFF2 Header = "font/woff2"
)

var defExtensionsMap = map[string]Header{
	".html":  HEADER_HTML,
	".htm":   HEADER_HTML,
	".css":   HEADER_CSS,
	".js":    HEADER_JS,
	".png":   HEADER_PNG,
	".jpg":   HEADER_JPG,
	".jpeg":  HEADER_JPG,
	".gif":   HEADER_GIF,
	".svg":   HEADER_SVG,
	".webp":  HEADER_WEBP,
	".webm":  HEADER_WEBM,
	".mp4":   HEADER_MP4,
	".weba":  HEADER_WEBA,
	".mp3":   HEADER_MP3,
	".ogg":   HEADER_OGG,
	".ttf":   HEADER_TTF,
	".otf":   HEADER_OTF,
	".woff":  HEADER_WOFF,
	".woff2": HEADER_WOFF2,
}
