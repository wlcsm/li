package config

import "codeberg.org/wlcsm/li/core"

var Colorscheme = map[core.SyntaxHL]int{
	core.HLComment:   90,
	core.HLMlComment: 90,
	core.HLKeyword1:  94,
	core.HLKeyword2:  96,
	core.HLString:    36,
	core.HLNumber:    33,
	core.HLMatch:     32,
	core.HLNormal:    39,
}


func SyntaxConf(ext string) *core.EditorSyntax {
	switch ext {
	case "c", "h", "cpp", "cc":
		return &C
	case "go":
		return &Go
	case "js":
		return &JavaScript
	case "py":
		return &Python
	case "html", "htm":
		return &Html
	case "json":
		return &JSON
	default:
		return nil
	}
}

var (
	C = core.EditorSyntax{
		Filetype: "c",
		Keywords: map[core.SyntaxHL][]string{
			core.HLKeyword1: {
				"switch", "if", "while", "for", "break", "continue", "return",
				"else", "struct", "union", "typedef", "static", "enum", "class",
				"case",
			},
			core.HLKeyword2: {
				"int", "long", "double", "float", "char", "unsigned",
				"signed", "void",
			},
		},
		Scs:              "//",
		Mcs:              "/*",
		Mce:              "*/",
		HighlightStrings: true,
		HighlightNumbers: true,
	}

	Go = core.EditorSyntax{
		Filetype: "go",
		Keywords: map[core.SyntaxHL][]string{
			core.HLKeyword1: {
				"break", "default", "func", "interface", "select", "case", "defer",
				"go", "map", "struct", "chan", "else", "goto", "package", "switch",
				"const", "fallthrough", "if", "range", "type", "continue", "for",
				"import", "return", "var",
			},
			core.HLKeyword2: {
				"append", "bool", "byte", "cap", "close", "complex",
				"complex64", "complex128", "error", "uint16", "copy", "false",
				"float32", "float64", "imag", "int", "int8", "int16",
				"uint32", "int32", "int64", "iota", "len", "make", "new",
				"nil", "panic", "uint64", "print", "println", "real",
				"recover", "rune", "string", "true", "uint", "uint8",
				"uintptr",
			},
		},
		Scs:              "//",
		Mcs:              "/*",
		Mce:              "*/",
		HighlightStrings: true,
		HighlightNumbers: true,
	}

	JavaScript = core.EditorSyntax{
		Filetype: "javascript",
		Keywords: map[core.SyntaxHL][]string{
			core.HLKeyword1: {
				"abstract", "arguments", "await", "boolean", "break", "char",
				"debugger", "do", "double", "export", "final", "finally",
				"goto", "import", "in", "let", "null", "public",
				"super", "throw", "try", "volatile", "byte", "class",
				"else", "extends", "float", "if", "instance", "long",
			},
			core.HLKeyword2: {
				"package", "return", "switch", "throws", "typeof", "case",
				"const", "default", "enum", "for", "implement", "of",
				"native", "private", "short", "synchronized", "transien", "var",
				"while", "catch", "continue", "delete", "eval", "false",
				"function", "int", "this", "true", "yield", "interface",
				"new", "protected", "static", "void", "with",
			},
		},
		Scs:              "//",
		Mcs:              "/*",
		Mce:              "*/",
		HighlightStrings: true,
		HighlightNumbers: true,
	}

	Python = core.EditorSyntax{
		Filetype: "python",
		Keywords: map[core.SyntaxHL][]string{
			core.HLKeyword1: {
				"False", "None", "True", "and", "as", "assert",
				"break", "class", "continuepass", "def", "yield", "del",
				"elif", "else", "except", "finally", "for", "from",
				"print",
			},
			core.HLKeyword2: {
				"if", "import", "in", "is", "lambda", "nonlocal",
				"not", "or", "global", "raise", "return", "try",
				"while", "with",
			},
		},
		Scs:              "#",
		Mcs:              `"""`,
		Mce:              `"""`,
		HighlightStrings: true,
		HighlightNumbers: true,
	}

	Html = core.EditorSyntax{
		Filetype: "html",
		Keywords: map[core.SyntaxHL][]string{
			core.HLKeyword1: {
				"!DOC.cyPE", "html", "head", "meta", "link", "r",
				"title", "body", "script", "div",
			},
			core.HLKeyword2: {
				"rel", "name", "content", "href", "type", "id", "charset",
			},
		},
		Scs:              "",
		Mcs:              "<!--",
		Mce:              "-->",
		HighlightStrings: true,
		HighlightNumbers: true,
	}

	JSON = core.EditorSyntax{
		Filetype:         "json",
		HighlightStrings: true,
		HighlightNumbers: true,
	}
)
