# gsnake

Snake in your terminal written in go

## Install

### From releases

Download latest release from the [releases page](https://github.com/evad1n/gsnake/releases).

### From source

`go install github.com/evad1n/gsnake@latest`

## Usage

### Controls

-   WASD or arrow keys for movement (In double mode (`-d`) WASD controls one snake and arrow keys the other)
-   `p` to pause
-   `r` to reset
-   `n` to go frame by frame while paused
-   `esc`, `ctrl + c` or `q` to quit

### Flags

```
-d    Play with 2 independent snakes on the same board
-double
    Play with 2 independent snakes on the same board
-s float
    Base speed multiplier (default 1)
-size int
    Optional max board size (default 40)
-speed float
    Base speed multiplier (default 1)
-w    Wrap around screen
-wrap
    Wrap around screen
```
