package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Vladimir-Da/Spint4/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	if len(data) == 0 {
		err := errors.New("no data")
		log.Println("error:", err)
		return 0, 0, err
	}
	// str - слайс строк где str[0]-кол-во шагов ,str[1]-длительность тренировки; разделитель по знаку "запятой"
	str := strings.Split(data, ",")
	if len(str) != 2 {
		err := errors.New("want - steps,duration; but have ")
		log.Println("error:", err, data)
		return 0, 0, err
	}
	// steps преобразованное интовое значение кол-во шагов
	steps, err := strconv.Atoi(str[0])
	if err != nil {
		err := errors.New("want string to parse into int ; but have ")
		log.Println("error:", err, str[0])
		return 0, 0, err
	}
	if steps <= 0 {
		err := errors.New("steps must be positive")
		log.Println("error: steps = ", steps, err)
		return 0, 0, err
	}
	//dutation - продолжительность тренировки (формат время)
	duration, err := time.ParseDuration(str[1])
	if err != nil {
		err := errors.New(" cant parse into duration")
		log.Println("error: ", str[1], err)
		return 0, 0, err
	}
	if duration == 0 {
		err := errors.New("the training lasted 0")
		log.Println("error: ", err)
		return 0, 0, err
	}
	if duration < 0 {
		err := errors.New("duration is negative")
		log.Println("error:", err)
		return 0, 0, err
	}
	return steps, duration, nil
}

//DayActionInfo У функции три параметра 1)строка с данными, которая содержит количество шагов и продолжительность прогулки в формате 3h50m (3 часа 50 минут).
// 2)weight float64 — вес пользователя в килограммах.
// 3)height float64 — рост пользователя в метрах
// функция возвращает итоговую строку с результатом тренировки формата
// Количество шагов: 792.
//Дистанция составила 0.51 км.
//Вы сожгли 221.33 ккал.

func DayActionInfo(data string, weight, height float64) string {
	steps, dayDuration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
	}
	if steps <= 0 {
		return ""
	}
	//distance - дистанция в метрах
	distance := stepLength * float64(steps)
	//distance- преобразуем в дистанцию в километрах
	distance = distance / float64(mInKm)
	callSpent, err := spentcalories.WalkingSpentCalories(steps, weight, height, dayDuration)
	if err != nil {
		fmt.Println(err)
	}
	res := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, callSpent)
	return res
}
