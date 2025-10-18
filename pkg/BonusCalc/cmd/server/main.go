package server

import (
	"bonuscalc/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func SeparationFlow(orderNum int) int {
	// Функция разделеяет входящий поток
	// Часть потока уйдет в дальнейшую обработку,
	// часть потока получит негатив и уйдет пользователю в ответ

	// Status Code ,который сразу отправится в ответ
	//    204 — заказ не зарегистрирован в системе расчёта.
	//    429 — превышено количество запросов к сервису.
	//    500 — внутренняя ошибка сервера.

	// num := orderNum % 101
	num := rand.Intn(101)

	if 0 <= num && num < 5 {
		// 15 % StatusCode 204
		return 204
		// 15 % StatusCode 429
	} else if 5 <= num && num < 10 {
		return 429
		// 15 % StatusCode 500
	} else if 10 <= num && num < 15 {
		return 500
		// 55 % StatusCode 200
	} else {
		return 200
	}
}

func ConditionsTime() int {
	timeSeconds := []int{3, 5, 7, 10, 13}
	cntSec := rand.Intn(len(timeSeconds))
	time.Sleep(time.Duration(cntSec * int(time.Second)))
	return cntSec
}

func ChoiceStatusOrder() string {
	statuses := []string{"REGISTERED", "INVALID", "PROCESSING", "PROCESSED"}
	idx := rand.Intn(len(statuses))
	return statuses[idx]
}

// Handlers
func ordersHandler(w http.ResponseWriter, r *http.Request) {
	// Accrual
	var accrual models.Accrual

	// Params
	// startTime := time.Now()
	tmpP := strings.Split(r.URL.Path, "/")
	orderNumStr := tmpP[len(tmpP)-1]
	orderNum, _ := strconv.Atoi(orderNumStr)

	// Разделение потока на ошибки в % от 100
	statusCode := SeparationFlow(orderNum)
	if statusCode != 200 {
		SendCustomErrorMessage(w, statusCode)
		return
	}

	// Засыпаем на рандом
	ConditionsTime()

	// Выбираем статус рандом
	accrual.Status = ChoiceStatusOrder()
	accrual.Order = orderNumStr
	accrual.Accrual = 999

	// Вывод в терминал
	fmt.Fprintln(os.Stdout, accrual)

	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accrual)
}

func Run() {
	// Args
	host := "localhost:8081"

	// Router
	mux := http.NewServeMux()
	mux.HandleFunc("/api/orders/", ordersHandler)

	// Server
	log.Fatal(http.ListenAndServe(host, mux))

}
