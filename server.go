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

//Position Definition
type PositionService struct {
  gorest.RestService `root:"/position/" consumes:"application/json" produces:"application/json"`
  position  gorest.EndPoint `method:"GET" path:"/" output:"Coord"`
  set       gorest.EndPoint `method:"PUT" path:"/" postdata:"Coord"`
}
func(serv PositionService) Position() Coord{

    return Coord{123, 12}
}
func(serv PositionService) Set(data Coord) {
    log.Printf("SET")
    log.Printf("%+v", data)
}

//Rocket Definition
type RocketService struct {
  gorest.RestService `root:"/rocket/" consumes:"application/json"`
  rocket    gorest.EndPoint `method:"PUT" path:"/" postdata:"Coord"`
}
func(serv RocketService) Rocket(data Coord) {
    log.Printf("ROCKET")
    log.Printf("%+v", data)
}

// Actions Definition
type ActionService struct {
    gorest.RestService `root:"/actions/"`
    stop  gorest.EndPoint `method:"PUT"  path:"/stop" postdata:"string"`
    up    gorest.EndPoint `method:"POST" path:"/up" postdata:"string"`
    down  gorest.EndPoint `method:"POST" path:"/down" postdata:"string"`
    left  gorest.EndPoint `method:"POST" path:"/left" postdata:"string"`
    right gorest.EndPoint `method:"POST" path:"/right" postdata:"string"`
    fire  gorest.EndPoint `method:"PUT"  path:"/fire" postdata:"string"`

}
func(serv ActionService) Stop(data string) {
  log.Printf("STOP")
}
func(serv ActionService) Up(data string) {
  log.Printf("UP")
}
func(serv ActionService) Down(data string) {
  log.Printf("DOWN")
}
func(serv ActionService) Left(data string) {
  log.Printf("LEFT")
}
func(serv ActionService) Right(data string) {
  log.Printf("RIGHT")
}
func(serv ActionService) Fire(data string) {
  log.Printf("FIRE")
}