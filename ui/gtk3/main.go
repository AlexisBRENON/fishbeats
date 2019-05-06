package gtk3

import (
  "log"
  "math"
  "time"

  "github.com/gotk3/gotk3/cairo"
  "github.com/gotk3/gotk3/gdk"
  "github.com/gotk3/gotk3/gtk"

  "github.com/AlexisBRENON/fishbeats/engine"
  "github.com/AlexisBRENON/fishbeats/utils"
)

type ApplicationData struct {
  e *engine.Engine
  surface *cairo.Surface
  dots []utils.Point
  previousDots []utils.Point
}

func NewApplicationData(e *engine.Engine) *ApplicationData {
  dots := []utils.Point{utils.Pt(0.5, 0.5)}
  previousDots := []utils.Point{}
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

    for i := 0; i < len(data.dots); i++ {
      switch keyEvent.KeyVal() {
      case 65361: // Left
        data.dots[i].X -= 0.01
        break
      case 65362: // Top
        data.dots[i].Y -= 0.01
        break
      case 65363: // Right
        data.dots[i].X += 0.01
        break
      case 65364: // Bottom
        data.dots[i].Y += 0.01
        break
      }
      data.dots[i].X = math.Min(math.Max(data.dots[i].X, 0.0), 1.0)
      data.dots[i].Y = math.Min(math.Max(data.dots[i].Y, 0.0), 1.0)
    }

    data.e.Update(timestamp, data.dots)

    cr := cairo.Create(data.surface);
    cr.Save()
    cr.SetOperator(cairo.OPERATOR_CLEAR)
    for i := 0; i < len(data.previousDots); i++ {
      point := data.previousDots[i]
      x := (point.X * float64(data.surface.GetWidth())) - 3.0
      y := (point.Y * float64(data.surface.GetHeight())) - 3.0
      w := 6
      h := 6
      cr.Rectangle(float64(x), float64(y), float64(w), float64(h))
      cr.Fill()
    }
    cr.Restore()
    cr.SetSourceRGBA(0, 0, 0, 1)
    for i := 0; i < len(data.dots); i++ {
      point := data.dots[i]
      x := (point.X * float64(data.surface.GetWidth())) - 3.0
      y := (point.Y * float64(data.surface.GetHeight())) - 3.0
      w := 6
      h := 6
      cr.Rectangle(float64(x), float64(y), float64(w), float64(h))
      cr.Fill()
    }
    widget.QueueDraw()
    data.previousDots = data.dots
    data.dots = make([]utils.Point, len(data.previousDots))
    copy(data.dots, data.previousDots)
    return true
}
