package main

import (
  "github.com/kr/pty"
  "os"
//  "io"
  "log"
  "fmt"
  "time"
)

type MotionCommand struct {
  Motor byte
  Command byte
  Data byte
}

func Motion(ctx *Thunder) (err error) {
  log.Printf("motion iface")
  p, t, err := pty.Open()
  if err != nil {
    
    log.Printf("err", err)
    return err
  }
  log.Printf("motion iface")
  log.Printf(p.Name())
  log.Printf(t.Name())
  err = os.Symlink(t.Name(), "/dev/ttyGOD2")
  if err != nil {
    log.Printf("error creating symlink")
    return err
  }


  command := &MotionCommand {
    Motor: 0x00,
    Command: 0x00,
    Data: 0x00,
  }
  buf := make([]byte, 3)
  resp := make([]byte, 1)
  durationH := time.Duration(2)
  durationV := time.Duration(2)
  for {
    n, err := p.Read(buf)
    if err != nil {
      return err
    }
    if n != 3 {
      return nil
    }
    command.Motor = buf[0]
    command.Command = buf[1]
    command.Data = buf[2]

    fmt.Printf("read motor:%x command:%x data:%x\n", command.Motor, command.Command, command.Data)


    durationH = time.Duration(int(command.Data)*20)*time.Millisecond
    durationV = time.Duration(int(command.Data)*30)*time.Millisecond
    if (command.Command == 0x07) {
      // Set speed
    } else if (command.Command == 0x06) {
      ctx.Put(STOP)
    } else if (command.Motor == 0x00 && command.Command == 0x01) {
      ctx.RegisterMove(DIR_UP)
      ctx.RegisterWait(durationH)
      ctx.RegisterStop()
      resp[0] = 0x01
    } else if (command.Motor == 0x00 && command.Command == 0x02) {
      ctx.RegisterMove(DIR_DOWN)
      ctx.RegisterWait(durationH)
      ctx.RegisterStop()
      resp[0] = 0x02
    } else if (command.Motor == 0x01 && command.Command == 0x01) {
      ctx.RegisterMove(DIR_LEFT)
      ctx.RegisterWait(durationV)
      ctx.RegisterStop()
      resp[0] = 0x01
    } else if (command.Motor == 0x01 && command.Command == 0x02) {
      ctx.RegisterMove(DIR_RIGHT)
      ctx.RegisterWait(durationV)
      ctx.RegisterStop()
      resp[0] = 0x02
    }
    p.Write(resp)
    resp[0] = 0x00



  }

  _, _ = p, t

  return nil
}

