package main

import (
    "fmt"
    "github.com/nsf/termbox-go"
    "os"
    "os/exec"
    "time"
)

func draw_image(width, height, a_paddle_pos, b_paddle_pos, 
    ball_pos_x, ball_pos_y, paddle_size, paddle_dist, 
    a_score, b_score, b_paddle_direction int) {
    c := exec.Command("clear")
    c.Stdout = os.Stdout
    c.Run()
    output := ""
    for h := 0; h < height; h++ {
        for w := 0; w < width; w++ {
            switch {
            // Draw left wall
            case w == 0:
                output += "|"
            // Draw Right
            case w == width-1:
                output += "|\n"
            // Draw top and bottom wall
            case h == 0 || h == height-1:
                output += "-"
            // Draw A Paddle
            case h == height - paddle_dist - 1 && 
            w >= a_paddle_pos - paddle_size && 
            w <= a_paddle_pos + paddle_size:
                output += "_"

            case h == height - paddle_dist - 1 && 
            w >= a_paddle_pos - paddle_size  && 
            w <= a_paddle_pos + paddle_size + 4:
                output += "~"

            case h == height - paddle_dist - 1 && 
            w <= a_paddle_pos &&
            w >= a_paddle_pos - paddle_size - 4 :
                output += "~"
            // Draw B Paddle
            case h == paddle_dist && 
            w >= b_paddle_pos - paddle_size && 
            w <= b_paddle_pos+ paddle_size:
                output += "_"
            //Draw Ball
            case h == ball_pos_y && 
            w == ball_pos_x:
                output += "o"
            default:
                output += " "
            }
        }

    }
    fmt.Println(output)
    fmt.Println("Player A Score ", a_score)
    fmt.Println("Player B Score ", b_score)
    fmt.Println(b_paddle_direction)
}


func computer(ch_b_paddle_direction chan int, height int ,width int,
    b_paddle_pos int, ball_pos_x int, ball_pos_y int,
    ball_direction_x int, ball_direction_y int) {
        findpos:
        for {
            if ball_pos_y >= height {
                ball_direction_y = ball_direction_y * -1
            }
            if ball_pos_x >= width-1 || ball_pos_x == 1 {
                ball_direction_x = ball_direction_x * -1
            }
            if ball_pos_y == 1 {
                ball_direction_y = ball_direction_y * -1
                break findpos
            }
            ball_pos_x += ball_direction_x
            ball_pos_y += ball_direction_y
        }
        switch {
        case b_paddle_pos > ball_pos_x:
            ch_b_paddle_direction <- 1
        case b_paddle_pos < ball_pos_x:
            ch_b_paddle_direction <- -1
        default:
            ch_b_paddle_direction <- 0
        }

}

func main() {
    width := 100
    height := 30
    a_paddle_pos := 40 
    b_paddle_pos := 40
    paddle_size := 5
    paddle_dist := 3
    ball_pos_x := 50
    ball_pos_y := 15
    ball_direction_x := 1
    ball_direction_y := -1
    a_score := 0
    b_score := 0
    var ch_b_paddle_direction chan int = make(chan int)
    err := termbox.Init()
    if err != nil {
            panic(err)
    }
    defer termbox.Close()

    event_queue := make(chan termbox.Event)
    go func() {
            for {
                    event_queue <- termbox.PollEvent()
            }
    }()
loop:
    for {
        go computer(ch_b_paddle_direction, height, width,
        b_paddle_pos, ball_pos_x, ball_pos_y,
        ball_direction_x,ball_direction_y)
        select {
        case ev := <-event_queue:
                if ev.Type == termbox.EventKey && 
                ev.Key == termbox.KeyEsc {
                        break loop
                }
                if ev.Type == termbox.EventKey && 
                ev.Key == termbox.KeyArrowLeft {
                        a_paddle_pos -= 2
                }
                if ev.Type == termbox.EventKey && 
                ev.Key == termbox.KeyArrowRight {
                        a_paddle_pos += 2
                }
        default:
            termbox.Flush()
            }

        switch {
        // paddle A ball interaction
        case ball_pos_y == height - paddle_dist - 1 && 
                ball_pos_x > a_paddle_pos - paddle_size && 
                ball_pos_x < a_paddle_pos + paddle_size:
            ball_direction_y = ball_direction_y * -1
        case ball_pos_y == height - paddle_dist - 1 && 
        ball_pos_x >= a_paddle_pos - paddle_size  && 
        ball_pos_x <= a_paddle_pos + paddle_size + 4:
            ball_direction_y = ball_direction_y * -1
            ball_direction_x = ball_direction_x * -1

        case ball_pos_y == height - paddle_dist - 1 && 
        ball_pos_x <= a_paddle_pos &&
        ball_pos_x >= a_paddle_pos - paddle_size - 4 :
            ball_direction_y = ball_direction_y * -1
            ball_direction_x = ball_direction_x * -1

        // paddle B ball interaction
        case ball_pos_y == paddle_dist && 
                ball_pos_x > b_paddle_pos - paddle_size && 
                ball_pos_x < b_paddle_pos + paddle_size:
            ball_direction_y = ball_direction_y * -1
        // Top Rebound
        case ball_pos_y == 1:
            ball_direction_y = ball_direction_y * -1
            a_score += 1
        // Bottom Rebound
        case ball_pos_y >= height:
            ball_direction_y = ball_direction_y * -1
            b_score += 1
        // X Axis Rebound
        }
        if ball_pos_x >= width-1 || ball_pos_x == 1 {
            ball_direction_x = ball_direction_x * -1
        }
        
        ball_pos_x += ball_direction_x
        ball_pos_y += ball_direction_y
        b_paddle_direction := <- ch_b_paddle_direction
        b_paddle_pos = b_paddle_pos - b_paddle_direction
        // b_paddle_pos =+ b_paddle_direction
        draw_image(width, height, 
            a_paddle_pos, b_paddle_pos, 
            ball_pos_x, ball_pos_y, 
            paddle_size, paddle_dist,
            a_score, b_score, b_paddle_direction)
        time.Sleep(60 * time.Millisecond)

}
}
