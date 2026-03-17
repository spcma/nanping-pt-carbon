package ipfs

import "math"

// DistanceCalculator 距离计算器
type DistanceCalculator struct {
	earthRadius float64 // 地球半径（米）
}

// NewDistanceCalculator 创建距离计算器
func NewDistanceCalculator() *DistanceCalculator {
	return &DistanceCalculator{
		earthRadius: 6378137.0, // 地球平均半径 6378.137km
	}
}

// HaversineDistance 使用 Haversine 公式计算两点间的大圆距离
// lat1, lon1: 起点纬度和经度（十进制度）
// lat2, lon2: 终点纬度和经度（十进制度）
// 返回值：距离（米）
func (dc *DistanceCalculator) HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// 将角度转换为弧度
	dLat := dc.toRad(lat2 - lat1)
	dLon := dc.toRad(lon2 - lon1)

	lat1Rad := dc.toRad(lat1)
	lat2Rad := dc.toRad(lat2)

	// Haversine 公式
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return dc.earthRadius * c
}

// CalculateTotalDistance 计算所有记录点的总里程
// records: GPS 记录列表（按时间顺序排列）
// 返回值：总距离（米）、错误信息
func (dc *DistanceCalculator) CalculateTotalDistance(records []Record) (float64, error) {
	if len(records) < 2 {
		return 0, nil
	}

	totalDistance := 0.0

	for i := 1; i < len(records); i++ {
		dist := dc.HaversineDistance(
			records[i-1].Lat, records[i-1].Lon,
			records[i].Lat, records[i].Lon,
		)
		totalDistance += dist
	}

	return totalDistance, nil
}

// CalculateSegmentDistances 计算每段路程的距离
// records: GPS 记录列表
// 返回值：每段的距离列表（米）
func (dc *DistanceCalculator) CalculateSegmentDistances(records []Record) []SegmentDistance {
	if len(records) < 2 {
		return []SegmentDistance{}
	}

	segments := make([]SegmentDistance, 0, len(records)-1)

	for i := 1; i < len(records); i++ {
		dist := dc.HaversineDistance(
			records[i-1].Lat, records[i-1].Lon,
			records[i].Lat, records[i].Lon,
		)

		segment := SegmentDistance{
			From:       records[i-1],
			To:         records[i],
			Distance:   dist,
			Cumulative: 0, // 将在后面计算
		}

		segments = append(segments, segment)
	}

	// 计算累计距离
	cumulative := 0.0
	for i := range segments {
		cumulative += segments[i].Distance
		segments[i].Cumulative = cumulative
	}

	return segments
}

// toRad 将角度转换为弧度
func (dc *DistanceCalculator) toRad(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

// SegmentDistance 分段距离信息
type SegmentDistance struct {
	From       Record  // 起点记录
	To         Record  // 终点记录
	Distance   float64 // 这段的距离（米）
	Cumulative float64 // 累计距离（米）
}

// DistanceSummary 距离汇总信息
type DistanceSummary struct {
	TotalDistance   float64 `json:"total_distance"`    // 总距离（米）
	TotalDistanceKm float64 `json:"total_distance_km"` // 总距离（公里）
	PointCount      int     `json:"point_count"`       // 点数
	SegmentCount    int     `json:"segment_count"`     // 分段数
	AverageSpeed    float64 `json:"average_speed"`     // 平均速度（km/h）
	TimeSpanHours   float64 `json:"time_span_hours"`   // 时间跨度（小时）
}

// CalculateSummary 计算距离汇总信息
func (dc *DistanceCalculator) CalculateSummary(records []Record) DistanceSummary {
	if len(records) == 0 {
		return DistanceSummary{}
	}

	totalDist, _ := dc.CalculateTotalDistance(records)

	summary := DistanceSummary{
		TotalDistance:   totalDist,
		TotalDistanceKm: totalDist / 1000.0,
		PointCount:      len(records),
		SegmentCount:    max(0, len(records)-1),
	}

	// 计算时间跨度和平均速度
	if len(records) >= 2 {
		timeSpan := records[len(records)-1].Timestamp.Sub(records[0].Timestamp)
		summary.TimeSpanHours = timeSpan.Hours()

		if summary.TimeSpanHours > 0 {
			summary.AverageSpeed = summary.TotalDistanceKm / summary.TimeSpanHours
		}
	}

	return summary
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
