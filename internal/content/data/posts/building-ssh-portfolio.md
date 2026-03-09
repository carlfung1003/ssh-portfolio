---
title: "Building an SSH Portfolio with Go"
date: "2026-03-09"
tags: "go, tui, ssh, charmbracelet"
readingTime: 8
---

# Building an SSH Portfolio with Go

What if your portfolio lived in the terminal instead of a browser?

That's the question that kicked off this project. Inspired by the growing movement
of SSH-based apps and the incredible Charmbracelet ecosystem, I decided to build
a personal portfolio that anyone can access with a single command:

```
ssh hi.carlfung.dev
```

## The Stack

- **Bubble Tea** — TUI framework using the Elm Architecture
- **Wish** — SSH server middleware (no OpenSSH needed)
- **Lip Gloss** — CSS-like styling for the terminal
- **Glamour** — Markdown rendering with ANSI colors

## Why SSH?

Browsers are great, but there's something delightfully minimal about the terminal.
No JavaScript bundles, no CSS frameworks, no loading spinners. Just text, color,
and keyboard shortcuts.

Plus, it's a conversation starter. When you put `ssh hi.carlfung.dev` on your
resume, people actually try it.

## Lessons Learned

1. **Go is surprisingly approachable** coming from TypeScript
2. **The Elm Architecture** (Model/Update/View) feels natural for TUI apps
3. **Terminal color support** varies wildly — test in multiple terminals
4. **embed.FS** is magic — bake all content into a single binary
