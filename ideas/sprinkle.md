# Sprinkle
A transparent, borderless Webkit frame (and not much else)

## Overview

A "sprinkle" is nothing more than a very tiny web browser designed to load a single-page web application.  These applications will, in turn, talk to the Sprinkles API for performing desktop and system management tasks.  Collectively, these tasks form the necessary interactions and behaviors of using a modern Linux graphical environment.

## Command Line Usage

```
sprinkle [options] APPNAME
```

### `[options]`

* `--hide`:        Hide the window on startup, leaving it up to the application being launched to show it when it is ready.
* `--layer`:       Which layer of the window stacking order the window should be ordered in (desktop, below, **normal**, above)
* `-w`/`--width`:  The initial width of the window (in pixels [250] or percent of screen width [75%])
* `-h`/`--height`: The initial height of the window (in pixels [32] or percent of screen width [5%])
* `-X`:            The X-coordinate at which the window should be placed initially 
* `-Y`:            The Y-coordinate at which the window should be placed initially
* `-D`/`--dock`:   A shortcut for pinning the window to a particular edge of the screen (top, left, bottom, right)