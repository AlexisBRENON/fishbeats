package gtk3

import (
  "image"
  "log"
  "time"

  "github.com/gotk3/gotk3/cairo"
  "github.com/gotk3/gotk3/gdk"
  "github.com/gotk3/gotk3/gtk"

  "github.com/AlexisBRENON/fishbeats/engine"
)

type ApplicationData struct {
  e *engine.Engine
  surface *cairo.Surface
  dots []image.Point
  previousDots []image.Point
}

func NewApplicationData(e *engine.Engine) *ApplicationData {
  dots := []image.Point{image.Pt(50, 50)}
  previousDots := []image.Point{}
  return &ApplicationData{
    e: e,
    dots: dots,
    previousDots: previousDots,
  }
}

func Main(e *engine.Engine) {
  gtk.Init(nil)
  appData := NewApplicationData(e)
  win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
  if err != nil {
    log.Fatal("Unable to create window:", err)
  }
  win.SetTitle("FishBeats")
  win.SetDefaultSize(400, 400)
  win.Connect("destroy", mainQuit)

  mainFrame, err := gtk.FrameNew("Main frame")
  if err != nil {
    log.Fatal("Unable to create mainFrame:", err)
  }
  mainFrame.SetShadowType(gtk.SHADOW_IN)
  win.Add(mainFrame)

  drawingArea, err := gtk.DrawingAreaNew()
  if err != nil {
    log.Fatal("Unable to create drawing area:", err)
  }
  drawingArea.SetSizeRequest(100, 100)
  mainFrame.Add(drawingArea)
  /* Signals used to handle the backing surface */
  drawingArea.Connect("draw", onAreaDraw, appData)
  drawingArea.Connect("configure-event", onAreaConfigure, appData)
  /* Event signals */
  win.Connect("key-press-event", onWinKeyPressed, appData)

  win.ShowAll()
  gtk.Main()
}

func mainQuit() {
  gtk.MainQuit()
}

func clearSurface(appData *ApplicationData) {
  cr := cairo.Create(appData.surface);
  cr.SetSourceRGB(1, 1, 1);
  cr.Paint();
}

func onAreaDraw(
  widget *gtk.DrawingArea,
  cr *cairo.Context,
  appData *ApplicationData) bool {
    cr.SetSourceSurface(appData.surface, 0, 0)
    cr.Paint()
    return true
}
func onAreaConfigure(
  widget *gtk.DrawingArea,
  event *gdk.Event,
  appData *ApplicationData) bool {
    appData.surface = cairo.CreateImageSurface(
      cairo.FORMAT_RGB24,
      widget.GetAllocatedWidth(),
      widget.GetAllocatedHeight())

  /* Initialize the surface to white */
  clearSurface (appData);

  /* We've handled the configure event, no need for further processing. */
  return true;
}

func onWinKeyPressed(
  widget *gtk.Window,
  event *gdk.Event,
  data *ApplicationData) bool {
    timestamp := time.Now().UnixNano()
    keyEvent := gdk.EventKeyNewFromEvent(event)
    if keyEvent.KeyVal() < 65361 || keyEvent.KeyVal() > 65364 {
      return false
    }

    switch keyEvent.KeyVal() {
    case 65361: // Left
      log.Println("Left")
      for i := 0; i < len(data.dots); i++ {
        data.dots[i].X -= 1
      }
      break
    case 65362: // Top
      log.Println("Top")
      for i := 0; i < len(data.dots); i++ {
        data.dots[i].Y -= 1
      }
      break
    case 65363: // Right
      log.Println("Right")
      for i := 0; i < len(data.dots); i++ {
        data.dots[i].X += 1
      }
      break
    case 65364: // Bottom
      log.Println("Bottom")
      for i := 0; i < len(data.dots); i++ {
        data.dots[i].Y += 1
      }
      break
    }
    data.e.Update(timestamp, data.dots)

    cr := cairo.Create(data.surface);
    cr.SetSourceRGB(1, 1, 1)
    for i := 0; i < len(data.previousDots); i++ {
      point := data.previousDots[i]
      x := (point.X * data.surface.GetWidth() / 100.0) - 3.0
      y := (point.Y * data.surface.GetHeight() / 100.0) - 3.0
      w := 6
      h := 6
      cr.Rectangle(float64(x), float64(y), float64(w), float64(h))
      cr.Fill()
    }
    cr.SetSourceRGB(0, 0, 0)
    for i := 0; i < len(data.dots); i++ {
      point := data.dots[i]
      x := (point.X * data.surface.GetWidth() / 100.0) - 3.0
      y := (point.Y * data.surface.GetHeight() / 100.0) - 3.0
      w := 6
      h := 6
      cr.Rectangle(float64(x), float64(y), float64(w), float64(h))
      cr.Fill()
    }
    widget.QueueDraw()
    data.previousDots = data.dots
    data.dots = make([]image.Point, len(data.previousDots))
    copy(data.dots, data.previousDots)
    return true
}
