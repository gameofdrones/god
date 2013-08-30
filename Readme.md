# Game of Drones

Raspberry powered rocket launcher

## API

### Position

#### Get current position

```
GET /position
Accept: application/json

{
  "x": 60,
  "y": 3
}
```

#### Reset position

```
DELETE /position
Accept: application/json

```

#### Set position

```
PUT /position
Content-Type: application/json

{
  "x": 60,
  "y": 3
}
```

### Rocket launcher


#### Launch !

```
PUT /rocket
Content-Type: application/json

{
  "x": 100,
  "y": 6
}
```

### Actions

#### Stop current move

```
PUT /actions/stop
```

#### Go up

```
POST /actions/up
```

#### Stop going up

```
DELETE /actions/up
```

#### Go down

```
POST /actions/down
```

#### Stop going down

```
DELETE /actions/down
```

#### Go left

```
POST /actions/left
```

#### Stop going left

```
DELETE /actions/left
```

#### Go right

```
POST /actions/right
```

#### Stop going right

```
DELETE /actions/right
```


#### Fire !!!!

```
PUT /actions/fire
```

### Authentication

None ;)


## Dev setup

### Debian

```
apt-get install libusb-1.0-0-dev
apt-get install golang
go get github.com/baloo/gousb/usb
go get code.google.com/p/gorest
```

### MacOS
```
brew install go
brew install libusb
go get code.google.com/p/gorest
go get github.com/baloo/gousb/usb
vi /usr/local/include/libusb-1.0/libusb.h
# Follow  https://gist.github.com/pjvds/4578277
```

