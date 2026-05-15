package providerschema

// DocStringFormat specifies the format used for an associated documentation
// string.
type DocStringFormat int

//go:generate go tool golang.org/x/tools/cmd/stringer -type=DocStringFormat -trimprefix=DocString

const (
	DocStringUnsupported DocStringFormat = 0

	// DocStringPlain means that the associated documentation string is intended
	// as plain text without any markup, although it may still contain
	// meaningful newline characters separating multiple paragraphs.
	DocStringPlain DocStringFormat = 1

	// DocStringMarkdown means that the associated documentation contains
	// Markdown-like formatting markup.
	//
	// Unfortunately no specific Markdown implementation has been specified,
	// so the string may contain markdown extensions that only apply to
	// certain implementations. It's the caller's responsibility to somehow
	// deal with unsupported or invalid Markdown text.
	DocStringMarkdown DocStringFormat = 2
)
