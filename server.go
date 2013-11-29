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
  http.ListenAndServe(":9000", http.FileServer(http.Dir(".")))
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
  positionAllowCross    gorest.EndPoint `method:"OPTIONS" path:"/"`
  position              gorest.EndPoint `method:"GET" path:"/" output:"Coord"`
  set                   gorest.EndPoint `method:"PUT" path:"/" postdata:"Coord"`
  reset                 gorest.EndPoint `method:"DELETE" path:"/"`
}
func NewPositionService(ctx *Thunder) *PositionService{
  s := new(PositionService)
  s.ctx = ctx
  return s
}
func(serv PositionService) PositionAllowCross() {
  allowCross(serv.ResponseBuilder())
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
func(serv PositionService) Reset() {
  log.Printf("RESET")
  serv.ctx.ResetPosition()
  allowCross(serv.ResponseBuilder())
}

//Rocket Definition
type RocketService struct {
  ctx *Thunder
  gorest.RestService `root:"/rocket/" consumes:"application/json"`
  rocketAllowCross    gorest.EndPoint `method:"OPTIONS" path:"/"`
  rocket              gorest.EndPoint `method:"PUT" path:"/" postdata:"Coord"`
}
func NewRocketService(ctx *Thunder) *RocketService{
  s := new(RocketService)
  s.ctx = ctx
  return s
}
func(serv RocketService) RocketAllowCross() {
  allowCross(serv.ResponseBuilder())
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
    actionAllowCross    gorest.EndPoint `method:"OPTIONS" path:"/{action:string}"`
    postAction           gorest.EndPoint `method:"POST"    path:"/{action:string}" postdata:"string"`
    putAction           gorest.EndPoint `method:"PUT"     path:"/{action:string}" postdata:"string"`
    deleteAction        gorest.EndPoint `method:"DELETE"  path:"/{action:string}"`
}
func NewActionService(ctx *Thunder) *ActionService{
  s := new(ActionService)
  s.ctx = ctx
  return s
}
func(serv ActionService) ActionAllowCross(action string) {
  allowCross(serv.ResponseBuilder())
}
func(serv ActionService) PostAction(data string, actionStr string) {
  serv.PutAction(data, actionStr)
}

func(serv ActionService) PutAction(data string, actionStr string) {
  var action Action
  switch(actionStr){
    case "stop": action = STOP
    case "up": action = DIR_UP
    case "down": action = DIR_DOWN
    case "left": action = DIR_LEFT
    case "right": action = DIR_RIGHT
    case "fire": {
      serv.ctx.Fire()
      allowCross(serv.ResponseBuilder())
      return
    }
    default: {
      allowCross(serv.ResponseBuilder()).SetResponseCode(404).Overide(true)
      return
    }
  }
  log.Printf("PUT %+v", action)
  serv.ctx.Put(action)
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
  serv.ctx.Delete(action)
  allowCross(serv.ResponseBuilder())
}
