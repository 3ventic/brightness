## What's this

This is a small app I wrote so I can hook up keyboard shortcuts in i3 to control the brightness levels on my XPS15 when xbacklight failed to do so.

## Prerequisites

libnotify-dev

## Usage

```
./install
backlight # get current %
backlight -inc # increase by 2%
backlight -dec # decrease by 2%
```

add -notify for desktop notification via libnotify
