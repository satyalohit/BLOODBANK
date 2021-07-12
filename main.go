package main

import (
	"bloodbank/dbconnect"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.GET("/receiver", Secondpage)
	server.GET("/home", Openpage)
	server.LoadHTMLGlob("templates/html/*.html")
	server.POST("/donatorinfo", Donate)
	server.POST("/Getdetails", Getdetails)
	server.Run(":2021")
}

func Openpage(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "bloodbank",
	})
}
func Secondpage(c *gin.Context) {
	c.HTML(http.StatusOK, "index1.html", gin.H{})

}

// type Info struct {
// 	Id         int    `json:"id"`
// 	Firstname  string `json:"firstname"`
// 	Lastname   string `json:"lastname"`
// 	Age        int    `json:"age"`
// 	Bloodgroup string `json:"bloodgroup"`
// 	Email      string `json:"email"`
// 	Phoneno    string `json:"phone"`
// 	District   string `json:"district"`
// 	Gender     string `json:"gender"`
// 	State      string `json:"state"`
// }
type Details struct {
	Id         int    `json:"id"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Age        int    `json:"age"`
	Bloodgroup string `json:"bloodgroup"`
	Email      string `json:"email"`
	Phoneno    string `json:"phone"`
	District   string `json:"district"`
	Gender     string `json:"gender"`
	State      string `json:"state"`
}

func Donate(c *gin.Context) {

	var t Details

	err := c.BindJSON(&t)
	if err != nil {
		log.Print(err)
	}

	db := dbconnect.Connectdatabase()

	defer db.Close()

	result,err2:=db.Query("select id,firstname,lastname,age,state,district,email,phoneno,bloodgroup,gender from users where phoneno=? or email=?",t.Phoneno,t.Email)
	if(err2!=nil){
		log.Println(err2);

	}


	var persons []Details

	for result.Next(){
     var person Details
		result.Scan(&person.Id, &person.Firstname,&person.Lastname,&person.Age,&person.State,&person.District,&person.Email,&person.Phoneno,&person.Bloodgroup,&person.Gender)
	   
		persons=append(persons, person)	
	}

	if(len(persons)!=0){

		c.JSON(http.StatusOK,gin.H{
			"status":"error",
			"message":"Already exists with that PHONE NUMBER or E-MAIL",
		})
	}else{
	res, err1 := db.Prepare("insert into users (firstname,lastname,age,bloodgroup,email,phoneno,state,district,gender)values(?,?,?,?,?,?,?,?,?)")
	if err1 != nil {
		log.Print(err1)
	}
	_, err1 = res.Exec(t.Firstname, t.Lastname, t.Age, t.Bloodgroup, t.Email, t.Phoneno, t.State, t.District, t.Gender)
	if err1 != nil {
		log.Print(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     "success",		
		"firstname":  t.Firstname,
		"lastname":   t.Lastname,
		"age":        t.Age,
		"bloodgroup": t.Bloodgroup,
		"email":      t.Email,
		"phoneno":    t.Phoneno,
		"district":   t.District,
		"gender":     t.Gender,
	})
}

}

func Getdetails(c *gin.Context) {
	var t Details

	err := c.BindJSON(&t)
	if err != nil {
		log.Println(err)
	}
	db := dbconnect.Connectdatabase()

	defer db.Close()

	res, err1 := db.Query("select id,firstname,lastname,age,state,district,email,phoneno,bloodgroup,gender from users where bloodgroup=? and (state=? or district=?)", t.Bloodgroup, t.State, t.District)
	if err1 != nil {
		log.Println(err1)
	}
	// var firstname, lastname, state, district, email, phoneno, bloodgroup, gender string
	// var id, age int
	var persons []Details
	for res.Next() {
		var person Details
		err1 = res.Scan(&person.Id, &person.Firstname,&person.Lastname,&person.Age,&person.State,&person.District,&person.Email,&person.Phoneno,&person.Bloodgroup,&person.Gender)

		if err1 != nil {
			log.Println(err1)
		}
		persons = append(persons, person)
	}
	if len(persons)== 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"Message": "No one is there",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"Message":    "Existed profiles",			
			"persons":    persons,
		})
	}
	fmt.Println(persons)
}
