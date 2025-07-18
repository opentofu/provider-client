package providerschema

// DocStringFormat specifies the format used for an associated documentation
// string.
type DocStringFormat int

const (
	// DocStringPlain means that the associated documentation string is intended
	// as plain text without any markup, although it may still contain
	// meaningful newline characters separating multiple paragraphs.
	DocStringPlain DocStringFormat = 0

	// DocStringMarkdown means that the associated documentation contains
	// Markdown-like formatting markup.
	//
	// Unfortunately no specific Markdown implementation has been specified,
	// so the string may contain markdown extensions that only apply to
	// certain implementations. It's the caller's responsibility to somehow
	// deal with unsupported or invalid Markdown text.
	DocStringMarkdown DocStringFormat = 1
)
