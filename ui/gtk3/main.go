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

  mainLayout, err := gtk.OverlayNew()
  if err != nil {
    log.Fatal("Unable to create main layout:", err)
  }
  win.Add(mainLayout)

  scaleArea, err := gtk.DrawingAreaNew()
  scaleSurface := cairo.CreateImageSurface(cairo.FORMAT_INVALID, 0, 0)
  if err != nil {
    log.Fatal("Unable to create scale area:", err)
  }
  scaleArea.SetSizeRequest(100, 100)
  mainLayout.AddOverlay(scaleArea)
  scaleArea.Connect(
    "configure-event",
    func(widget *gtk.DrawingArea, event *gdk.Event) bool {
      return onScaleAreaConfigure(widget, &scaleSurface, e)
    })
  scaleArea.Connect("draw", onAreaDraw, &scaleSurface)

  fishArea, err := gtk.DrawingAreaNew()
  if err != nil {
    log.Fatal("Unable to create fish area:", err)
  }
  fishArea.SetSizeRequest(100, 100)
  mainLayout.AddOverlay(fishArea)
  /* Signals used to handle the backing surface */
  fishArea.Connect("configure-event", onAreaConfigure, &appData.surface)
  fishArea.Connect("draw", onAreaDraw, &appData.surface)
  /* Event signals */
  win.Connect(
    "key-press-event",
    func(widget *gtk.Window, event *gdk.Event, data *ApplicationData) bool {
      return onAreaKeyPressed(fishArea, event, data)
    }, appData)

  win.ShowAll()
  gtk.Main()
}

func mainQuit() {
  gtk.MainQuit()
}

func clearSurface(surface *cairo.Surface) {
  cr := cairo.Create(surface);
  cr.SetOperator(cairo.OPERATOR_CLEAR)
  cr.Paint();
}

func onAreaDraw(
  widget *gtk.DrawingArea,
  cr *cairo.Context,
  surface **cairo.Surface) bool {
    cr.SetSourceSurface(*surface, 0, 0)
    cr.Paint()
    return true
}

func onAreaConfigure(
  widget *gtk.DrawingArea,
  event *gdk.Event,
  surface **cairo.Surface) bool {
    *surface = cairo.CreateImageSurface(
      cairo.FORMAT_ARGB32,
      widget.GetAllocatedWidth(),
      widget.GetAllocatedHeight())

  /* Initialize the surface */
  clearSurface(*surface);

  /* We've handled the configure event, no need for further processing. */
  return true;
}

func onScaleAreaConfigure(
  widget *gtk.DrawingArea,
  surface **cairo.Surface,
  e *engine.Engine,
) bool {
  *surface = cairo.CreateImageSurface(
    cairo.FORMAT_ARGB32,
    widget.GetAllocatedWidth(),
    widget.GetAllocatedHeight())
  clearSurface(*surface);

  cr := cairo.Create(*surface);
  cr.SetSourceRGBA(0.3, 0.3, 0.3, 0.6)
  cr.SetLineWidth(1)
  cr.Save()
  cr.Scale(
    float64(widget.GetAllocatedWidth()),
    float64(widget.GetAllocatedHeight()))
  for i := 1; i < e.NumOctaves*8; i++ {
    yValue := float64(i) / float64(e.NumOctaves * 8)
    cr.MoveTo(0, yValue);
    cr.LineTo(1, yValue);
  }
  for i := 0.0 ; i < 1 ; i += e.LegatoValue {
    cr.MoveTo(i, 0)
    cr.LineTo(i, 1)
  }
  cr.Restore()
  cr.Stroke()

  return true;
}

func onAreaKeyPressed(
  widget *gtk.DrawingArea,
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
    cr.Save()
    cr.SetOperator(cairo.OPERATOR_CLEAR)
    for i := 0; i < len(data.previousDots); i++ {
      point := data.previousDots[i]
      x := (point.X * data.surface.GetWidth() / 100.0) - 3.0
      y := (point.Y * data.surface.GetHeight() / 100.0) - 3.0
      w := 6
      h := 6
      cr.Rectangle(float64(x), float64(y), float64(w), float64(h))
      cr.Fill()
    }
    cr.Restore()
    cr.SetSourceRGBA(0, 0, 0, 1)
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
