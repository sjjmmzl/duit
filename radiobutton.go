package duit

import (
	"image"

	"9fans.net/go/draw"
)

type Radiobutton struct {
	Selected bool
	Value    interface{}
	Disabled bool
	Group    []*Radiobutton
	Changed  func(v interface{}, r *Result) // only the change function of the newly selected radiobutton in the group will be called

	m draw.Mouse
}

var _ UI = &Radiobutton{}

func (ui *Radiobutton) Layout(dui *DUI, self *Kid, sizeAvail image.Point, force bool) {
	dui.debugLayout("Radiobutton", self)

	hit := image.Point{0, 1}
	size := pt(2*BorderSize + 4*dui.Display.DefaultFont.Height/5).Add(hit)
	self.R = rect(size)
}

func (ui *Radiobutton) Draw(dui *DUI, self *Kid, img *draw.Image, orig image.Point, m draw.Mouse, force bool) {
	dui.debugDraw("Radiobutton", self)

	r := rect(pt(2*BorderSize + 4*dui.Display.DefaultFont.Height/5))
	hover := m.In(r)
	r = r.Add(orig)

	colors := dui.Regular.Normal
	color := colors.Text
	if ui.Disabled {
		colors = dui.Disabled
		color = colors.Border
	} else if hover {
		colors = dui.Regular.Hover
		color = colors.Border
	}

	hit := pt(0)
	if hover && m.Buttons&1 == 1 {
		hit = image.Pt(0, 1)
	}

	img.Draw(extendY(r, 1), colors.Background, nil, image.ZP)
	r = r.Add(hit)

	radius := r.Dx() / 2
	img.Arc(r.Min.Add(pt(radius)), radius, radius, 0, color, image.ZP, 0, 360)

	cr := r.Inset((4 * dui.Display.DefaultFont.Height / 5) / 5).Add(hit)
	if ui.Selected {
		radius = cr.Dx() / 2
		img.FillArc(cr.Min.Add(pt(radius)), radius, radius, 0, color, image.ZP, 0, 360)
	}
}

func (ui *Radiobutton) check(r *Result) {
	ui.Selected = true
	for _, r := range ui.Group {
		if r != ui {
			r.Selected = false
		}
	}
	if ui.Changed != nil {
		ui.Changed(ui.Value, r)
	}
}

func (ui *Radiobutton) Mouse(dui *DUI, self *Kid, m draw.Mouse, origM draw.Mouse, orig image.Point) (r Result) {
	if ui.Disabled {
		return
	}
	rr := rect(pt(2*BorderSize + 4*dui.Display.DefaultFont.Height/5))
	hover := m.In(rr)
	if hover != ui.m.In(rr) {
		self.Draw = Dirty
	}
	if hover && ui.m.Buttons&1 != m.Buttons&1 {
		self.Draw = Dirty
		if m.Buttons&1 == 0 {
			r.Consumed = true
			ui.check(&r)
		}
	}
	ui.m = m
	return
}

func (ui *Radiobutton) Key(dui *DUI, self *Kid, k rune, m draw.Mouse, orig image.Point) (r Result) {
	if k == ' ' {
		r.Consumed = true
		self.Draw = Dirty
		ui.check(&r)
	}
	return
}

func (ui *Radiobutton) FirstFocus(dui *DUI) *image.Point {
	return &image.ZP
}

func (ui *Radiobutton) Focus(dui *DUI, o UI) *image.Point {
	if o != ui {
		return nil
	}
	return ui.FirstFocus(dui)
}

func (ui *Radiobutton) Mark(self *Kid, o UI, forLayout bool, state State) (marked bool) {
	return self.Mark(o, forLayout, state)
}

func (ui *Radiobutton) Print(self *Kid, indent int) {
	PrintUI("Radiobutton", self, indent)
}
