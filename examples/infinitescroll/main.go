// Copyright 2016 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build example

package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

var (
	theViewport     = &viewport{}
	bgImage         *ebiten.Image
	repeatedBgImage *ebiten.Image
	groundImage     *ebiten.Image
)

type viewport struct {
	x16 int
	y16 int
}

func round(x float64) float64 {
	return math.Floor(x + 0.5)
}

func (p *viewport) Move() {
	w, h := bgImage.Size()
	mx := w * 16
	my := h * 16

	p.x16 += w / 32
	p.y16 += h / 32

	for mx <= p.x16 {
		p.x16 -= mx
	}
	for my <= p.y16 {
		p.y16 -= my
	}
	for p.x16 < 0 {
		p.x16 += mx
	}
	for p.y16 < 0 {
		p.y16 += my
	}
}

func (p *viewport) Position() (int, int) {
	return p.x16, p.y16
}

func updateGroundImage(ground *ebiten.Image) {
	ground.Clear()
	x16, y16 := theViewport.Position()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(-x16)/16, float64(-y16)/16)
	ground.DrawImage(repeatedBgImage, op)
}

func drawGroundImage(screen *ebiten.Image, ground *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(ground, op)
}

func update(screen *ebiten.Image) error {
	theViewport.Move()
	if ebiten.IsRunningSlowly() {
		return nil
	}
	updateGroundImage(groundImage)
	drawGroundImage(screen, groundImage)

	msg := fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS())
	ebitenutil.DebugPrint(screen, msg)
	return nil
}

func main() {
	var err error
	bgImage, _, err = ebitenutil.NewImageFromFile("_resources/images/tile.png", ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	w, h := bgImage.Size()
	const repeat = 5
	repeatedBgImage, _ = ebiten.NewImage(w*repeat, h*repeat, ebiten.FilterNearest)
	for j := 0; j < repeat; j++ {
		for i := 0; i < repeat; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(w*i), float64(h*j))
			repeatedBgImage.DrawImage(bgImage, op)
		}
	}
	groundImage, _ = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterNearest)

	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Infinite Scroll (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}
}
