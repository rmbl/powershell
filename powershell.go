package main

import (
    "fmt"
    "github.com/docopt/docopt.go"
)

type Characters struct {
    separator      string
    separator_thin string
    ln             string
    branch         string
    padlock        string
}

const (
    color_template string = "\\[\\e%s\\]"
    reset          string = "\\[\\e[0m\\]"
)

func Color(prefix int, code int) string {
    return fmt.Sprintf(color_template, fmt.Sprintf("[%v;5;%vm", prefix, code))
}
func FgColor(code int) string {
    return Color(38, code)
}
func BgColor(code int) string {
    return Color(48, code)
}

func PrintSegment(segment *Segment, next *Segment) {
    var sepColor string

    if next == nil {
        sepColor = reset + FgColor(segment.separatorFg)
    } else {
        sepColor = FgColor(segment.separatorFg) + BgColor(next.bg)
    }

    if segment.content[0] != ' ' {
        segment.content = " " + segment.content
    }
    if segment.content[len(segment.content) - 1] != ' ' {
        segment.content += " "
    }

    fmt.Print(FgColor(segment.fg), BgColor(segment.bg), segment.content, sepColor, segment.separator)
}

// Global variables
var chars Characters
var colorScheme ColorScheme
var line []*Segment

func main() {
    usage := `Go Powershell

Usage: powershell [options] <segment>...

Options:
    -h, --help           Show this screen.
    --version           Show the version.
    -c TYPE, --characters=TYPE  The type of line characters (powerline, compatible or flat) [default: powerline]

Segments:
    username
    hostname
    path
    git
    prompt
`
    arguments, _ := docopt.Parse(usage, nil, true, "Go Powershell 0.1", true)

    switch arguments["--characters"] {
    case "flat":
        chars = Characters{"", "", "", "", ""}
    case "compatible":
        chars = Characters{"\u25B6", "\u25B7", "", "", ""}
    case "powerline":
        chars = Characters{"\uE0B0", "\uE0B1", "\uE0A1", "\uE0A0", "\uE0A2"}
    default:
        chars = Characters{"\uE0B0", "\uE0B1", "\uE0A1", "\uE0A0", "\uE0A2"}
    }

    // Init colors and set used ones
    colors := InitColors()
    colorScheme = colors["default"]

    // Execute them and save them in an array before printing
    segments := InitSegments()
    line = make([]*Segment, 0, len(arguments["<segment>"].([]string)))
    for _, seg := range arguments["<segment>"].([]string) {
        if chosen, exists := segments[seg]; exists {
            chosen()
        }
    }

    // Print all chosen segments
    for i, seg := range line {
        if i < len(line) - 1 {
            PrintSegment(seg, line[i + 1])
        } else {
            PrintSegment(seg, nil)
        }
    }

    // Reset the cmdline before the user prompt
    fmt.Print(reset)
}
