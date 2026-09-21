package main

type scare struct {
	start, end float64
	face       int
}

var scares = []scare{
	{39.2, 40.15, 0}, {40.85, 41.25, 3}, {42.05, 43.3, 5}, {43.85, 44.2, 1},
	{62.4, 63.05, 2}, {64.15, 65.3, 4}, {67.05, 67.38, 0}, {69.1, 70.55, 5},
	{84.1, 84.45, 1}, {85, 85.9, 3}, {87.3, 87.62, 4}, {89, 90.35, 2}, {91.15, 91.78, 5},
	{105.4, 105.7, 0}, {106.2, 107, 4}, {108, 108.35, 1}, {109.2, 110.35, 5},
	{111.5, 112, 2}, {113, 114.35, 3},
}

type sceneID int

const (
	sceneStory sceneID = iota
	sceneCorridor0
	sceneCorridor1
	sceneCorridor2
	sceneBlackA
	sceneFalseEnd1
	sceneMonitor
	sceneBlackB
	sceneRecovered
	sceneAfterCameraText
	sceneBlackC
	sceneFalseEnd2
	sceneBlackD
	sceneSnapwave
	sceneLock
)

func scareAt(t float64) (scare, bool) {
	for _, s := range scares {
		if t >= s.start && t < s.end {
			return s, true
		}
	}
	return scare{}, false
}

func scareIndexAt(t float64) int {
	for i, s := range scares {
		if t >= s.start && t < s.end {
			return i
		}
	}
	return -1
}

func storySegment(t float64) int {
	switch {
	case t < 2.6:
		return 0
	case t < 5.6:
		return 1
	case t < 9:
		return 2
	case t < 12.5:
		return 3
	case t < 17.4:
		return 4
	case t < 21.2:
		return 5
	case t < 25.2:
		return 6
	case t < 29:
		return 7
	case t < 32.3:
		return 8
	default:
		return 9
	}
}

func visualKey(t float64) int64 {
	if i := scareIndexAt(t); i >= 0 {
		return 100000 + int64(i)
	}
	s := sceneAt(t)
	key := int64(s) * 100
	switch s {
	case sceneStory:
		key += int64(storySegment(t))
	case sceneRecovered:
		if t >= 77 {
			key++
		}
	case sceneFalseEnd2:
		if t >= 101.3 {
			key++
		}
	}
	return key
}

func sceneAt(t float64) sceneID {
	switch {
	case t < 32.4:
		return sceneStory
	case t < 35:
		return sceneCorridor0
	case t < 36.7:
		return sceneCorridor1
	case t < 38.75:
		return sceneCorridor2
	case t < 44.5:
		return sceneBlackA
	case t < 54.4:
		return sceneFalseEnd1
	case t < 61.9:
		return sceneMonitor
	case t < 70.8:
		return sceneBlackB
	case t < 78.9:
		return sceneRecovered
	case t < 83.75:
		return sceneAfterCameraText
	case t < 92.2:
		return sceneBlackC
	case t < 104.9:
		return sceneFalseEnd2
	case t < 114.7:
		return sceneBlackD
	case t < 118.2:
		return sceneSnapwave
	default:
		return sceneLock
	}
}

func storyText(t float64) (string, string) {
	switch {
	case t < 2.6:
		return "", ""
	case t < 5.6:
		return "АРХИВ КАМЕР ВОССТАНОВЛЕН", "Блок наблюдения // сектор 09"
	case t < 9:
		return "Запись оборвалась в 03:17.", "Причина остановки неизвестна."
	case t < 12.5:
		return "Камеры не записывают звук.", "Но в журнале оператора есть одна строка."
	case t < 17.4:
		return "«кто-то стучит по объективу изнутри»", "03:16:42"
	case t < 21.2:
		return "На первой камере коридор пуст.", "На второй — тоже."
	case t < 25.2:
		return "На третьей появляется человек.", "Он стоит слишком далеко, чтобы рассмотреть лицо."
	case t < 29:
		return "Следующий кадр пропущен.", "Когда изображение возвращается — он ближе."
	case t < 32.3:
		return "Ещё ближе.", "И теперь смотрит прямо в объектив."
	default:
		return "", ""
	}
}

func dynamicScene(s sceneID) bool {
	switch s {
	case sceneCorridor0, sceneCorridor1, sceneCorridor2, sceneMonitor:
		return true
	default:
		return false
	}
}
