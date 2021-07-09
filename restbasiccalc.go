package restbasiccalc

import "net/http"



type exprRequest struct {
	expression string
}


func handler(w http.ResponseWriter, r *http.Request) {
	// логика обработчика

	
}



func serverCalc() {

	http.HandleFunc("/basiccalc", handler)


	http.ListenAndServe(":8080", nil)
	// логика
	
}