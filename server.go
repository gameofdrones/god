package main

import (
  "code.google.com/p/gorest"
  "net/http"
  "log"
)

type Coord struct{
  X int64
  Y int64
}

type Action byte

const (
  DIR_UP Action = 0x00
  DIR_DOWN Action = 0x01
  DIR_LEFT Action = 0x02
  DIR_RIGHT Action = 0x03

  FIRE Action = 0x04 // TODO

  WAIT Action = 0xFE
  STOP Action = 0xFF
)

func NewCoord(x int64, y int64) *Coord{
  c := new(Coord)
  c.X = x
  c.Y = y
  return c
}

func main() {
  Register()
}

func Register(){
  gorest.RegisterService(new(PositionService))
  gorest.RegisterService(new(RocketService))
  gorest.RegisterService(new(ActionService))
  http.Handle("/",gorest.Handle())
  http.ListenAndServe(":9000",nil)
}

func allowCross(rb *gorest.ResponseBuilder) *gorest.ResponseBuilder {
  rb.AddHeader("Access-Control-Allow-Origin", "*")
  rb.AddHeader("Access-Control-Allow-Method", "GET, PUT, POST, DELETE")
  rb.AddHeader("Access-Control-Allow-Headers", "accept, origin, x-requested-with, content-type")
  return rb
}

//Position Definition
type PositionService struct {
  gorest.RestService `root:"/position/" consumes:"application/json" produces:"application/json"`
  position  gorest.EndPoint `method:"GET" path:"/" output:"Coord"`
  set       gorest.EndPoint `method:"PUT" path:"/" postdata:"Coord"`
}
func(serv PositionService) Position() Coord{
  allowCross(serv.ResponseBuilder())
  return Coord{123, 12}
}
func(serv PositionService) Set(data Coord) {
  log.Printf("SET")
  log.Printf("%+v", data)
  allowCross(serv.ResponseBuilder())
}

//Rocket Definition
type RocketService struct {
  gorest.RestService `root:"/rocket/" consumes:"application/json"`
  rocket    gorest.EndPoint `method:"PUT" path:"/" postdata:"Coord"`
}
func(serv RocketService) Rocket(data Coord) {
    log.Printf("ROCKET")
    log.Printf("%+v", data)
    allowCross(serv.ResponseBuilder())
}

// Actions Definition
type ActionService struct {
    gorest.RestService `root:"/actions/"`
    putAction     gorest.EndPoint `method:"PUT"     path:"/{action:string}" postdata:"string"`
    deleteAction  gorest.EndPoint `method:"DELETE"  path:"/{action:string}"`
}

func(serv ActionService) PutAction(data string, actionStr string) {
  var action Action
  switch(actionStr){
    case "stop": action = STOP
    case "up": action = DIR_UP
    case "down": action = DIR_DOWN
    case "left": action = DIR_LEFT
    case "right": action = DIR_RIGHT
    case "fire": action = FIRE
    default: {
      allowCross(serv.ResponseBuilder()).SetResponseCode(404).Overide(true)
      return
    }
  }
  log.Printf("PUT %+v", action)
  allowCross(serv.ResponseBuilder())
}

func(serv ActionService) DeleteAction(actionStr string) {
  var action Action
  switch(actionStr){
    case "up": action = DIR_UP
    case "down": action = DIR_DOWN
    case "left": action = DIR_LEFT
    case "right": action = DIR_RIGHT
    default: {
      allowCross(serv.ResponseBuilder()).SetResponseCode(404).Overide(true)
      return
    }
  }
  log.Printf("DELETE %+v", action)
  allowCross(serv.ResponseBuilder())
}