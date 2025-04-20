package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining принимает строку с данными формата "3456,Ходьба,3h00m", которая содержит количество шагов, вид активности и продолжительность активности.
func parseTraining(data string) (int, string, time.Duration, error) {
	if len(data) == 0 {
		return 0, "", 0, errors.New("нет данных")
	}
	// str - слайс строк где str[0]-кол-во шагов ,str[1]-тип тренировки,str[2]-продолжительность тренировки;разделитель по знаку "запятой"
	str := strings.Split(data, ",")
	if len(str) != 3 {
		return 0, "", 0, errors.New("недостаточно данных")
	}
	// steps преобразованное интовое значение кол-во шагов
	steps, err := strconv.Atoi(str[0])
	if err != nil {
		return 0, "", 0, errors.New("ошибка преоразования шагов в инт")
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше 0")
	}
	// dutation - продолжительность тренировки (формат время)
	duration, err := time.ParseDuration(str[2])
	if err != nil {
		return 0, "", 0, errors.New("ошибка преоразования продолжительности тренировки")
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("длительность тренировки должна быть больше 0")
	}
	// trainigType- вид тренировки("Ходьба" или "Бег")
	trainigType := str[1]
	if trainigType == "" {
		return 0, "", 0, errors.New("нет типа тренировки")
	}
	return steps, trainigType, duration, nil
}

// distance принимает количество шагов и рост пользователя в метрах, а возвращает дистанцию в километрах.
func distance(steps int, height float64) float64 {
	//stepLength - длина шага
	stepLength := height * stepLengthCoefficient
	//walkDistance - пройденная дистанция,умножаем количество шагов на длину шага
	walkDistance := float64(steps) * stepLength
	// определяем сколько прошел в километрах
	walkDistance = walkDistance / float64(mInKm)
	return walkDistance
}

// meanSped принимает количество шагов steps, рост пользователя height и продолжительность активности duration  и возвращает среднюю скорость.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		fmt.Println("длительность тренировки должна быть больше 0")
		return 0
	}
	//hours- продолжительность в часах
	hours := duration.Hours()
	//avarageSpeed - средняя скорость
	avarageSpeed := distance(steps, height) / hours
	return avarageSpeed
}

// TrainingInfo - Функция принимает:
// 1)data string — строку с данными формата "3456,Ходьба,3h00m",
// которая содержит количество шагов, вид активности и продолжительность активности
// 2)weight, height float64 — вес (кг.) и рост (м.) пользователя.
// Возвращает результат тренировки и ошибку
func TrainingInfo(data string, weight, height float64) (string, error) {
	var res string
	var spentcalories, avarageSpeed float64
	steps, typeTrain, duration, err := parseTraining(data)
	if err != nil {
		return "", errors.New("ошибка в данных")
	}

	distanceKm := distance(steps, height)
	avarageSpeed = meanSpeed(steps, height, duration)

	switch typeTrain {
	case "Бег":
		spentcalories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			fmt.Println("error from func RunningSpentCalories")
		}
	case "Ходьба":
		spentcalories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			fmt.Println("error from func WalkingSpentCalories")
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	res = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", typeTrain, duration.Hours(), distanceKm, avarageSpeed, spentcalories)
	return res, nil
}

// RunningSpentCalories-Функция принимает:количество шагов;вес(кг.) и рост(м.) пользователя;продолжительность бега и
// возвращает два значения:количество калорий, потраченных при беге и ошибку
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("длительность тренировки должна быть больше 0")
	}
	//speed- средняя скорость(результат функции meanSpeed)
	speed := meanSpeed(steps, height, duration)
	//durationInMinutes- продолжительность в минутах
	durationInMinutes := duration.Minutes()
	return (weight * speed * durationInMinutes) / minInH, nil
}

// WalkingSpentCalories-Функция принимает:количество шагов;вес(кг.) и рост(м.) пользователя;продолжительность бега и
// возвращает два значения:количество калорий, потраченных при беге и ошибку
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("длительность тренировки должна быть больше 0")
	}
	// speed- средняя скорость(результат функции meanSpeed)
	speed := meanSpeed(steps, height, duration)
	// durationInMinutes- продолжительность в минутах
	durationInMinutes := duration.Minutes()
	return ((weight * speed * durationInMinutes) / minInH) * walkingCaloriesCoefficient, nil
}
