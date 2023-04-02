package main

import (
	"flag"
	"os"
	"strings"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"github.com/titanous/json5"
)

func main() {
	keyboardDefFilename := flag.String("keyboardDef", "keyboard-right-def.json", "Filename for the keyboard deffinition to load and generate")
	flag.Parse()

	var kd KeyboardDefinition
	kdBytes, err := os.ReadFile(*keyboardDefFilename)
	if err != nil {
		panic(err)
	}
	err = json5.Unmarshal(kdBytes, &kd)
	if err != nil {
		panic(err)
	}

	var knp BubbleKeyNoduleProperties
	knpBytes, err := os.ReadFile(kd.BubbleKeyNodulePropertiesFile)
	if err != nil {
		panic(err)
	}

	err = json5.Unmarshal(knpBytes, &knp)
	if err != nil {
		panic(err)
	}

	points := make([]NoduleTypeAndPoint, 0)
	for _, col := range kd.FingerColumns {
		points = append(points, col.getKeyLocations()...)
	}

	for _, row := range kd.ThumbRows {
		points = append(points, row.getKeyLocations()...)
	}

	topNodules := make([]Nodule, len(points))
	bottomNodules := make([]Nodule, len(points))

	bubbleKeys := make(map[int64]KeyNodule)
	getBubbleKey := func(screwPossitionsBits int64) KeyNodule {
		k, ok := bubbleKeys[screwPossitionsBits]
		if !ok {
			k = knp.MakeBubbleKey(screwPossitionsBits)
			bubbleKeys[screwPossitionsBits] = k
		}
		return k
	}

	for i, p := range points {
		if p.noduleType == NoduleKey {
			topNodules[i] = getBubbleKey(p.screwPossitionsBits).Top.OrientAndMove(p.moveTo)
			bottomNodules[i] = getBubbleKey(p.screwPossitionsBits).Bottom.OrientAndMove(p.moveTo)
		} else if p.noduleType == NoduleDebug1 {
			topNodules[i] = MakeNoduleDebug1().OrientAndMove(p.moveTo)
			bottomNodules[i] = MakeNoduleDebug1().OrientAndMove(p.moveTo)
		} else if p.noduleType == NoduleDebug2 {
			topNodules[i] = MakeNoduleDebug2().OrientAndMove(p.moveTo)
			bottomNodules[i] = MakeNoduleDebug2().OrientAndMove(p.moveTo)

		} else if p.noduleType == NoduleDebug3 {
			topNodules[i] = MakeNoduleDebug3().OrientAndMove(p.moveTo)
			bottomNodules[i] = MakeNoduleDebug3().OrientAndMove(p.moveTo)
		}
	}
	top := NoduleCollection(topNodules).Combine()
	back := sdf.Difference3D(NoduleCollection(bottomNodules).Combine(), top)

	_ = top
	_ = back

	topOutputFileName := strings.TrimSuffix(*keyboardDefFilename, ".json") + "top.stl"
	backOutputFileName := strings.TrimSuffix(*keyboardDefFilename, ".json") + "back.stl"
	render.RenderSTL(top, 350, topOutputFileName)
	render.RenderSTL(back, 300, backOutputFileName)
	//render.RenderSTL(top, 300, "3x5plus2top.stl")
	//render.RenderSTL(back, 300, "3x5plus2back.stl")

}

/*
circomferance = 2 pi radius
circomferance = tau radius
arc = angle tau radius

crude slope over arc:
rise / run
run: delta angle * average radius
rise: detla radius
*/

/*
Distance calculation

angleMoved = (angle-startAngle)
radius = startRadius + radiusPerAngleIncrement*angleMoved
x2 = radius * COS(angle)
y2 = radius * SIN(angle)

distance = SQRT( (x2-x)^2 + (y2-y)^2 )

x2 = (startRadius + radiusPerAngleIncrement*(angle-startAngle)) * COS(angle)
y2 = (startRadius + radiusPerAngleIncrement*(angle-startAngle)) * SIN(angle)

distance = SQRT( ((startRadius + radiusPerAngleIncrement*(angle-startAngle)) * COS(angle)-x)^2 + ((startRadius + radiusPerAngleIncrement*(angle-startAngle)) * SIN(angle)-y)^2 )

*/
