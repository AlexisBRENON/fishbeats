package engine

import (
  "image"
  "log"

  "github.com/hajimehoshi/oto"
)

type Engine struct {
  player *oto.Player
}

func NewEngine() *Engine {
  var sampleRate = 44100
  var channelNum = 2
  var bitDepthInBytes = 1
  context, err := oto.NewContext(
    sampleRate,
    channelNum,
    bitDepthInBytes,
    4096,
  )
  if err != nil {
    return nil
  }

  return &Engine{
    player: context.NewPlayer(),
  }
}

func (e *Engine) Update(timestamp_nano int64, positions []image.Point) error {
  log.Println("Update at ", timestamp_nano, ":", positions)
  return nil
}

