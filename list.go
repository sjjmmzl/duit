package duit

import (
	"image"

	"9fans.net/go/draw"
)

type ListValue struct {
	Label    string
	Value    interface{}
	Selected bool
}

type List struct {
	Values   []*ListValue
	Multiple bool
	Changed  func(index int, result *Result)
	Click    func(index, buttons int, r *Result)
	Keys     func(index int, m draw.Mouse, k rune, r *Result)

	size image.Point
}

var _ UI = &List{}

func (ui *List) Layout(env *Env, size image.Point) image.Point {
	font := env.Display.DefaultFont
	ui.size = image.Pt(size.X, len(ui.Values)*(4*font.Height/3))
	return ui.size
}

func (ui *List) Draw(env *Env, img *draw.Image, orig image.Point, m draw.Mouse) {
	font := env.Display.DefaultFont
	r := rect(ui.size).Add(orig)
	img.Draw(r, env.Normal.Background, nil, image.ZP)
	lineR := r
	lineR.Max.Y = lineR.Min.Y + 4*font.Height/3

	for _, v := range ui.Values {
		colors := env.Normal
		if v.Selected {
			colors = env.Inverse
		}
		img.Draw(lineR, colors.Background, nil, image.ZP)
		img.String(lineR.Min.Add(image.Pt(font.Height/4, font.Height/4)), colors.Text, image.ZP, font, v.Label)
		lineR = lineR.Add(image.Pt(0, 4*font.Height/3))
	}
}

func (ui *List) Mouse(env *Env, m draw.Mouse) (result Result) {
	result.Hit = ui
	if !m.In(rect(ui.size)) {
		return
	}
	font := env.Display.DefaultFont
	index := m.Y / (4 * font.Height / 3)
	if m.Buttons != 0 && ui.Click != nil {
		ui.Click(index, m.Buttons, &result)
	}
	if !result.Consumed && m.Buttons == 1 {
		v := ui.Values[index]
		v.Selected = !v.Selected
		if v.Selected && !ui.Multiple {
			for _, vv := range ui.Values {
				if vv != v {
					vv.Selected = false
				}
			}
		}
		if ui.Changed != nil {
			ui.Changed(index, &result)
		}
		result.Redraw = true
		result.Consumed = true
	}
	return
}

func (ui *List) selectedIndices() (l []int) {
	for i, lv := range ui.Values {
		if lv.Selected {
			l = append(l, i)
		}
	}
	return
}

func (ui *List) Key(env *Env, orig image.Point, m draw.Mouse, k rune) (result Result) {
	result.Hit = ui
	if !m.In(rect(ui.size)) {
		return
	}
	if ui.Keys != nil {
		// xxx what should "index" be? especially for multiple: true...
		sel := ui.selectedIndices()
		index := -1
		if len(sel) == 1 {
			index = sel[0]
		}
		ui.Keys(index, m, k, &result)
		if result.Consumed {
			return
		}
	}
	switch k {
	case draw.KeyUp, draw.KeyDown:
		if len(ui.Values) == 0 {
			return
		}
		sel := ui.selectedIndices()
		oindex := -1
		nindex := -1
		switch k {
		case draw.KeyUp:
			result.Consumed = true
			if len(sel) == 0 {
				nindex = len(ui.Values) - 1
			} else {
				oindex = sel[0]
				nindex = (sel[0] - 1 + len(ui.Values)) % len(ui.Values)
			}
		case draw.KeyDown:
			result.Consumed = true
			if len(sel) == 0 {
				nindex = 0
			} else {
				oindex = sel[len(sel)-1]
				nindex = (sel[len(sel)-1] + 1) % len(ui.Values)
			}
		}
		if oindex >= 0 {
			ui.Values[oindex].Selected = false
			result.Redraw = true
		}
		if nindex >= 0 {
			ui.Values[nindex].Selected = true
			result.Redraw = true
			if ui.Changed != nil {
				ui.Changed(nindex, &result)
			}
			// xxx orig probably should not be a part in this...
			font := env.Display.DefaultFont
			p := orig.Add(image.Pt(m.X, nindex*(4*font.Height/3)+font.Height/2))
			result.Warp = &p
		}
	}
	return
}

func (ui *List) FirstFocus(env *Env) *image.Point {
	return &image.ZP
}

func (ui *List) Focus(env *Env, o UI) *image.Point {
	if o != ui {
		return nil
	}
	return ui.FirstFocus(env)
}

func (ui *List) Print(indent int, r image.Rectangle) {
	uiPrint("List", indent, r)
}