package util

import "github.com/ryomak/gameboys/common/math"

// Rect 矩形（整数座標）
type Rect struct {
	X, Y          int32
	Width, Height int32
}

// NewRect 矩形を作成
func NewRect(x, y, width, height int32) Rect {
	return Rect{X: x, Y: y, Width: width, Height: height}
}

// Intersects 矩形同士の衝突判定
func (r Rect) Intersects(other Rect) bool {
	return r.X < other.X+other.Width &&
		r.X+r.Width > other.X &&
		r.Y < other.Y+other.Height &&
		r.Y+r.Height > other.Y
}

// Contains 点が矩形内にあるか
func (r Rect) Contains(x, y int32) bool {
	return x >= r.X && x < r.X+r.Width &&
		y >= r.Y && y < r.Y+r.Height
}

// ContainsRect 矩形が別の矩形を完全に含むか
func (r Rect) ContainsRect(other Rect) bool {
	return r.X <= other.X &&
		r.Y <= other.Y &&
		r.X+r.Width >= other.X+other.Width &&
		r.Y+r.Height >= other.Y+other.Height
}

// Center 矩形の中心座標を取得
func (r Rect) Center() (x, y int32) {
	return r.X + r.Width/2, r.Y + r.Height/2
}

// Right 矩形の右端座標
func (r Rect) Right() int32 {
	return r.X + r.Width
}

// Bottom 矩形の下端座標
func (r Rect) Bottom() int32 {
	return r.Y + r.Height
}

// Circle 円（固定小数点座標）
type Circle struct {
	X, Y   math.Fixed
	Radius math.Fixed
}

// NewCircle 円を作成
func NewCircle(x, y, radius int32) Circle {
	return Circle{
		X:      math.NewFixed(x),
		Y:      math.NewFixed(y),
		Radius: math.NewFixed(radius),
	}
}

// NewCircleFixed 固定小数点座標で円を作成
func NewCircleFixed(x, y, radius math.Fixed) Circle {
	return Circle{X: x, Y: y, Radius: radius}
}

// Intersects 円同士の衝突判定
func (c Circle) Intersects(other Circle) bool {
	dx := c.X.Sub(other.X)
	dy := c.Y.Sub(other.Y)
	distSq := dx.Mul(dx).Add(dy.Mul(dy))
	radiusSum := c.Radius.Add(other.Radius)
	return distSq < radiusSum.Mul(radiusSum)
}

// Contains 点が円内にあるか
func (c Circle) Contains(x, y math.Fixed) bool {
	dx := c.X.Sub(x)
	dy := c.Y.Sub(y)
	distSq := dx.Mul(dx).Add(dy.Mul(dy))
	return distSq < c.Radius.Mul(c.Radius)
}

// IntersectsRect 円と矩形の衝突判定
func (c Circle) IntersectsRect(rect Rect) bool {
	// 矩形の最も近い点を見つける
	rectX := math.NewFixed(rect.X)
	rectY := math.NewFixed(rect.Y)
	rectW := math.NewFixed(rect.Width)
	rectH := math.NewFixed(rect.Height)

	closestX := c.X
	closestY := c.Y

	if c.X < rectX {
		closestX = rectX
	} else if c.X > rectX.Add(rectW) {
		closestX = rectX.Add(rectW)
	}

	if c.Y < rectY {
		closestY = rectY
	} else if c.Y > rectY.Add(rectH) {
		closestY = rectY.Add(rectH)
	}

	// 円の中心から最も近い点までの距離を計算
	dx := c.X.Sub(closestX)
	dy := c.Y.Sub(closestY)
	distSq := dx.Mul(dx).Add(dy.Mul(dy))

	return distSq < c.Radius.Mul(c.Radius)
}

// Point 点（固定小数点座標）
type Point struct {
	X, Y math.Fixed
}

// NewPoint 点を作成
func NewPoint(x, y int32) Point {
	return Point{
		X: math.NewFixed(x),
		Y: math.NewFixed(y),
	}
}

// NewPointFixed 固定小数点座標で点を作成
func NewPointFixed(x, y math.Fixed) Point {
	return Point{X: x, Y: y}
}

// DistanceTo 別の点までの距離
func (p Point) DistanceTo(other Point) math.Fixed {
	dx := p.X.Sub(other.X)
	dy := p.Y.Sub(other.Y)
	return math.FixedSqrt(dx.Mul(dx).Add(dy.Mul(dy)))
}

// DistanceToSq 別の点までの距離の二乗
func (p Point) DistanceToSq(other Point) math.Fixed {
	dx := p.X.Sub(other.X)
	dy := p.Y.Sub(other.Y)
	return dx.Mul(dx).Add(dy.Mul(dy))
}
