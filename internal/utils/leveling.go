package utils

import (
	"log"
)

const MaxLevel = 80

type ExpTarget int

// Define exp targets for each level
var AllTargets = []ExpTarget{
	1500, 1815, 1980, 2200, 2420, 2640, 2915, 3190, 3520, 3905,
	4290, 4472, 4888, 5408, 5928, 6500, 7176, 7600, 8350, 9150,
	10100, 11100, 12200, 13450, 14750, 16250, 17900, 19650, 21650, 23800,
	26950, 27750, 28600, 29450, 30350, 31250, 32200, 33150, 34150, 35200,
	36250, 37300, 38450, 39600, 40800, 42000, 43250, 44550, 45900, 47250,
	48700, 50150, 51650, 53200, 54800, 56450, 58150, 59900, 61700, 63550,
	65450, 67400, 69400, 71500, 73650, 75850, 78150, 80500, 82900, 85400,
	87950, 90600, 93300, 96100, 99000, 101950, 105000, 108150, 111400, 114750,
}

// FindNextLevel finds the next level based on current exp
func FindNextLevel(exp int) int {
	log.Println(exp)

	for i, target := range AllTargets {
		if int(target) > exp {
			return i + 1
		}
	}

	return MaxLevel
}
