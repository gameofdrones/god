package main

import (
  "code.google.com/p/gorest"
  "net/http"
  "log"
)

type Coord struct{
  X int
  Y int
}

func NewCoord(x int, y int) *Coord{
  c := new(Coord)
  c.X = x
  c.Y = y
  return c
}

func Register(ctx *Thunder){
  gorest.RegisterService( NewPositionService(ctx) )
  gorest.RegisterService( NewRocketService(ctx) )
  gorest.RegisterService( NewActionService(ctx) )
  http.Handle("/",gorest.Handle())
  http.ListenAndServe(":9000",nil)
}

func allowCross(rb *gorest.ResponseBuilder) *gorest.ResponseBuilder {
  rb.AddHeader("Access-Control-Allow-Origin", "*")
  rb.AddHeader("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE")
  rb.AddHeader("Access-Control-Allow-Headers", "accept, origin, x-requested-with, content-type")
  return rb
}

//Position Definition
type PositionService struct {
  ctx *Thunder
  gorest.RestService `root:"/position/" consumes:"application/json" produces:"application/json"`
  position  gorest.EndPoint `method:"GET" path:"/" output:"Coord"`
  set       gorest.EndPoint `method:"PUT" path:"/" postdata:"Coord"`
}
func NewPositionService(ctx *Thunder) *PositionService{
  s := new(PositionService)
  s.ctx = ctx
  return s
}
func(serv PositionService) Position() Coord{
  allowCross(serv.ResponseBuilder())
  return Coord{123, 12}
}
func(serv PositionService) Set(data Coord) {
  log.Printf("SET")
  log.Printf("%+v", data)
  serv.ctx.SetPosition(data.X, data.Y, false)
  allowCross(serv.ResponseBuilder())
}

//Rocket Definition
type RocketService struct {
  ctx *Thunder
  gorest.RestService `root:"/rocket/" consumes:"application/json"`
  rocketAllowCross    gorest.EndPoint `method:"OPTIONS" path:"/"`
  rocket    gorest.EndPoint `method:"PUT" path:"/" postdata:"Coord"`
}
func NewRocketService(ctx *Thunder) *RocketService{
  s := new(RocketService)
  s.ctx = ctx
  return s
}
func(serv RocketService) Rocket(data Coord) {
    log.Printf("ROCKET")
    log.Printf("%+v", data)
    serv.ctx.SetPosition(data.X, data.Y, true)
    allowCross(serv.ResponseBuilder())
}

// Actions Definition
type ActionService struct {
    ctx *Thunder
    gorest.RestService `root:"/actions/"`
    putAction     gorest.EndPoint `method:"PUT"     path:"/{action:string}" postdata:"string"`
    deleteAction  gorest.EndPoint `method:"DELETE"  path:"/{action:string}"`
}
func NewActionService(ctx *Thunder) *ActionService{
  s := new(ActionService)
  s.ctx = ctx
  return s
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

func(serv RocketService) RocketAllowCross() {
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
