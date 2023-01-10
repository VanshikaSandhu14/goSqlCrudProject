package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"strconv"

	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var cfg = mysql.Config{
	User:   "root",
	Passwd: "vanshika01",
	Net:    "tcp",
	Addr:   "localhost:3306",
	DBName: "sql_go_test_project",
}

type json_emp struct {
	emp_id   string `json:"emp_id"`
	emp_name string `json:"emp_name"`
	salary   int    `json:"salary"`
	contact  int    `json:"contact"`
}

type sql_emp struct {
	emp_id   string `field:"emp_id"`
	emp_name string `field:"emp_name"`
	salary   int    `field:"salary"`
	contact  int    `field:"contact"`
}



func getEmpById(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	emp_id, ok := params["emp_id"]

	if !ok {
		return
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	var E sql_emp

	if err := db.QueryRow("SELECT * from emp_data where emp_id = "+emp_id).Scan(&E.emp_id, &E.emp_name, &E.salary, &E.contact); err != nil {
		fmt.Printf("The given id %v doesn't exist\n "+err.Error()+"\n", emp_id)
		return
	}

	E2 := json_emp(E)

	json.NewEncoder(res).Encode(E2)
}

func delEmp(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	emp_id, ok := params["emp_id"]

	if !ok {
		return
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	if _, err := db.Exec("DELETE from emp_data where emp_id = " + emp_id); err != nil {
		fmt.Printf("The given id: %v doesn't exist\n", emp_id)
		return
	}
}

func getEmpSal(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	emp_name := req.URL.Query().Get("emp_name")

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	var E string

	if err := db.QueryRow("SELECT salary from emp_data where emp_name = " + emp_name).Scan(&E); err != nil {
		fmt.Printf("The given name: %v doesn't exist\n", emp_name)
		return
	}

	json.NewEncoder(res).Encode(E)
}

func createNewEmp(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	fmt.Println("calling create employee method")

	emp_name := req.URL.Query().Get("emp_name")
	salary := req.URL.Query().Get("salary")
	emp_id := req.URL.Query().Get("emp_id")
	contact := req.URL.Query().Get("contact")

	slr, _ := strconv.Atoi(salary)
	E2 := json_emp{emp_id: emp_id, emp_name: emp_name, salary: slr, contact: 123456789}
	fmt.Println("query0-->" + E2.emp_id)

	json.NewEncoder(res).Encode(E2)

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()
	_, er := db.Exec("INSERT INTO emp_data VALUES(\"" + emp_id + "\"," + emp_name + "," + salary + "," + "\"" + contact + "\"" + ");")

	if er != nil {
		fmt.Println("Error in adding values to db\n" + er.Error())
	}

}

func main() {

	route := mux.NewRouter()
	fmt.Println(route)

	// if route == nil {
	// 	fmt.Println("error in connecting to the port")
	// 	fmt.Println(route)
	// }
	http.ListenAndServe(":3000", route)

	route.HandleFunc("/employee/{emp_id}", getEmpById).Methods("GET")
	route.HandleFunc("/employee/", getEmpSal).Methods("GET")
	route.HandleFunc("/employee/", createNewEmp).Methods("POST")
	route.HandleFunc("/employee/{emp_id}", delEmp).Methods("DELETE")

	fmt.Println("Starting Server at Port 3000")

	//	log.Fatal(http.ListenAndServe(":3000", route))

}
