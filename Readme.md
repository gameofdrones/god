# Game of Drones

Raspberry powered rocket launcher

## API

```
GET /position
Accept: application/json

{
  "x": 60,
  "y": 3
}
```

```
PUT /rocket
Content-Type: application/json

{
  "x": 100,
  "y": 6
}
```

```
PUT /actions/stop
```

```
POST /actions/up
```

```
POST /actions/down
```

```
POST /actions/left
```

```
POST /actions/right
```

```
PUT /actions/fire
```


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

