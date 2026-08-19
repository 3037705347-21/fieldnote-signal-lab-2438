package lab

import "strings"

func scoreObservation(item Observation) (int, []string) {
	score := 0
	reasons := []string{}
	if len(strings.Fields(item.Description)) >= 8 {
		score += 25
		reasons = append(reasons, "description contains sufficient field detail")
	} else {
		reasons = append(reasons, "description is brief")
	}
	score += minInt(len(item.Labels)*12, 30)
	if len(item.Labels) > 0 {
		reasons = append(reasons, "controlled tags were attached")
	}
	if item.Review != nil {
		score += reviewScore(*item.Review)
		reasons = append(reasons, "review confidence was incorporated")
	}
	if score > 100 {
		score = 100
	}
	return score, reasons
}
func reviewScore(review Review) int {
	if review.Verdict == "confirmed" {
		return int(30 * review.Confidence)
	}
	if review.Verdict == "needs-followup" {
		return int(20 * review.Confidence)
	}
	return 0
}
func scoreBand(score int) string {
	if score >= 70 {
		return "strong"
	}
	if score >= 40 {
		return "developing"
	}
	return "early"
}
func minInt(left, right int) int {
	if left < right {
		return left
	}
	return right
}
