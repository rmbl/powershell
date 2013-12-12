package main

import (
    "os"
    "os/exec"
    "strings"
)

type Segment struct {
    fg       int
    bg       int
    content  string

    separator string
    separatorFg int
}

func AppendSegment(fg, bg int, content string) {
    AppendFullSegment(fg, bg, content, chars.separator, bg)
}

func AppendFullSegment(fg, bg int, content, separator string, separatorFg int) {
    segment := new(Segment)
    segment.fg = fg
    segment.bg = bg
    segment.content = content
    segment.separator = separator
    segment.separatorFg = separatorFg
    line = append(line, segment)
}

func InitSegments() map[string]func() {
    segments := make(map[string]func())
    segments["username"] = AddUserSegment
    segments["hostname"] = AddHostSegment
    segments["path"] = AddPathSegment
    segments["prompt"] = AddPromptSegment
    segments["git"] = AddGitSegment
    return segments
}

func AddUserSegment() {
    fg := colorScheme.UsernameFg
    bg := colorScheme.UsernameBg
    if os.Geteuid() == 0 {
        fg = colorScheme.RootFg
        bg = colorScheme.RootBg
    }

    AppendSegment(fg, bg, "\\u")
}

func AddPromptSegment() {
    AppendSegment(colorScheme.PromptFg, colorScheme.PromptBg, "\\$")
}

func AddHostSegment() {
    segment := "\\h"

    if os.Getenv("SSH_CLIENT") != "" && len(chars.padlock) > 0 {
        segment = chars.padlock + " " + segment
    }

    AppendSegment(colorScheme.HostnameFg, colorScheme.HostnameBg, segment)
}

func AddPathSegment() {
    home := os.Getenv("HOME")
    path, e := os.Getwd()
    if e != nil {
        path = os.Getenv("PWD")
    }
    // Replace home directory with a tilde
    if strings.HasPrefix(path, home) {
        path = "~" + path[len(home):]
    }

    split := strings.Split(path, "/")
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
        if dir == "~" {
            AppendSegment(colorScheme.HomeFg, colorScheme.HomeBg, dir)
        } else if i == len(split) - 1 {
            AppendSegment(colorScheme.CwdFg, colorScheme.PathBg, dir)
        } else {
            AppendFullSegment(colorScheme.PathFg, colorScheme.PathBg, dir, chars.separator_thin, colorScheme.SeparatorFg)
        }
    }
}

func AddGitSegment() {
    var (
        branch string
        filesAdded bool
        fg int = colorScheme.RepoCleanFg
        bg int = colorScheme.RepoCleanBg
    )

    cmd := exec.Command("git", "ls-files", "--exclude-standard", "--others")
    output, err := cmd.Output()
    if err != nil {
        return
    } else if len(output) > 0 {
        filesAdded = true
    }

    cmd = exec.Command("git", "diff-files", "--quiet")
    output, err = cmd.Output()
    if err != nil || filesAdded {
        fg = colorScheme.RepoDirtyFg
        bg = colorScheme.RepoDirtyBg
    }

    cmd = exec.Command("git", "rev-parse", "--symbolic-full-name", "--abbrev-ref", "HEAD")
    output, err = cmd.Output()
    branch = strings.TrimSpace(string(output))

    //Add branch character and print it
    if len(chars.branch) > 0 {
        branch = chars.branch + " " + branch
    }
    if filesAdded {
        branch += " +"
    }
    AppendSegment(fg, bg, branch)
}

