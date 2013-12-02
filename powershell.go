package main

import (
    "flag"
    "fmt"
    "os"
    "strings"
)

type Characters struct {
    separator      string
    separator_thin string
    ln             string
    branch         string
    padlock        string
}

var chars Characters

type Segment struct {
    fg       int
    bg       int
    content  string
    callback func(Segment) string
}

var segments map[string]Segment = make(map[string]Segment)

func InitSegments() {
    segments["username"] = Segment{250, 240, "\\u", nil}
    segments["hostname"] = Segment{250, 238, "\\h", nil}
    segments["path"] = Segment{fg: 250, bg: 237, callback: GetPathSegment}
    segments["prompt"] = Segment{15, 236, "\\$", nil}
}

type Options struct {
    shell    string
    line     string
    segments string
}

func InitCli() Options {
    opt := Options{}
    flag.StringVar(&opt.shell, "shell", "bash", "The type of the shell (zsh or bash [default])")
    flag.StringVar(&opt.line, "line", "powerline", "The type of line characters (powerline [default], compatible or flat)")
    flag.StringVar(&opt.segments, "segments", "", "The segments of the prompt")
    flag.Parse()
    return opt
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

func PrintSegment(segment Segment, next Segment) {
    var (
        content  string
        sepColor string
    )
    if segment.callback == nil {
        content = fmt.Sprintf(" %s ", segment.content)
    } else {
        content = segment.callback(segment)
    }

    if next.bg == 0 && next.fg == 0 {
        sepColor = reset + FgColor(segment.bg)
    } else {
        sepColor = FgColor(segment.bg) + BgColor(next.bg)
    }
    fmt.Print(FgColor(segment.fg), BgColor(segment.bg), content, sepColor, chars.separator)
}

func main() {
    InitSegments()
    options := InitCli()

    switch options.line {
    default:
    case "powerline":
        chars = Characters{"\uE0B0", "\uE0B1", "\uE0A1", "\uE0A0", "\uE0A2"}
    case "compatible":
        chars = Characters{"\u25B6", "\u25B7", "", "", ""}
    }

    segs := strings.Split(options.segments, ",")
    for i := 0; i < len(segs); i++ {
        _, exists := segments[segs[i]]
        if exists {
            var next Segment
            if len(segs) > i+1 {
                next, exists = segments[segs[i+1]]
                if !exists {
                    next = Segment{fg: 0, bg: 0}
                }
            } else {
                next = Segment{fg: 0, bg: 0}
            }
            PrintSegment(segments[segs[i]], next)
        }
    }
    fmt.Print(reset)
}

func GetPathSegment(seg Segment) string {
    home := os.Getenv("HOME")
    path, e := os.Getwd()
    if e != nil {
        path = os.Getenv("PWD")
    }
    // Replace home directory with a tilde
    if strings.HasPrefix(path, home) {
        path = "~" + path[len(home):]
    }

    var (
        segString string
        split     []string
    )
    split = strings.Split(path, "/")
    // First element might be empty, ignore it
    if split[0] == "" {
        split = split[1:]
    }
    // Keep part short
    if len(split) > 4 {
        tmp := make([]string, 5)
        tmp[0] = split[0]
        tmp[1] = split[1]
        tmp[2] = "\u2026"
        tmp[3] = split[len(split)-2]
        tmp[4] = split[len(split)-1]
        split = tmp
    }

    for i, dir := range split {
        segString += " " + dir + " "
        if i < len(split)-1 {
            segString += chars.separator_thin
        }
    }
    return segString
}
