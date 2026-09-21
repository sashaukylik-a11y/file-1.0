package main

import "testing"

func TestTimelineCoversAllTimes(t *testing.T) {
	for ms := 0; ms <= 130000; ms += 10 {
		s := sceneAt(float64(ms) / 1000)
		if s < sceneStory || s > sceneLock {
			t.Fatalf("invalid scene at %d ms: %d", ms, s)
		}
	}
}

func TestScaresValid(t *testing.T) {
	prev := -1.0
	for i, s := range scares {
		if s.start < 0 || s.end <= s.start {
			t.Fatalf("invalid scare %d: %+v", i, s)
		}
		if s.face < 0 || s.face >= 6 {
			t.Fatalf("invalid face index %d: %d", i, s.face)
		}
		if s.start < prev {
			t.Fatalf("scares not ordered at %d", i)
		}
		prev = s.start
	}
}

func TestCriticalTransitions(t *testing.T) {
	cases := []struct {
		time float64
		want sceneID
	}{
		{28.99, sceneStory},
		{32.39, sceneStory},
		{32.40, sceneCorridor0},
		{34.99, sceneCorridor0},
		{35.00, sceneCorridor1},
		{36.70, sceneCorridor2},
		{38.75, sceneBlackA},
		{118.20, sceneLock},
	}
	for _, tc := range cases {
		if got := sceneAt(tc.time); got != tc.want {
			t.Fatalf("sceneAt(%.2f)=%d want %d", tc.time, got, tc.want)
		}
	}
}
