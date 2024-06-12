//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv4"

//	 "cloud.google.com/go/alloydbconn"
//   "cloud.google.com/go/alloydbconn/driver/pgxv4"

//	"context"
//	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
//	"github.com/auth0/go-jwt-middleware/v2/validator"
//	"encoding/json"
	
)

var dbPool *sql.DB

type Writer struct {
	ID    int    `json:"id" binding:"required"`
	Likes int    `json:"likes"`
	Writer  string `json:"Writer" binding:"required"`
	Color  string `json:"color" binding:"required"` 
}
  
/** a list of initial Writers */
var Writers = []Writer{
	Writer{1, 0, "Paul Valéry.", "#4285F4"},
	Writer{2, 0, "Héraclite.", "#DB4437"},
	Writer{3, 0, "Danton.","#545454"},
}

/*
var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	log.Println("authMiddleware in Handler")

	claims, ok := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	if !ok {
		log.Println("authMiddleware in Handler !ok")
		http.Error(w, "failed to get validated claims", http.StatusInternalServerError)
		return
	}
	
	payload, err := json.Marshal(claims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
})
*/


func main() {
/*
	keyFunc := func(ctx context.Context) (interface{}, error) {
		// Our token must be signed using this data.
		return []byte("secret"), nil
	}

	// Set up the validator.
	jwtValidator, err := validator.New(
		keyFunc,
		validator.HS256,
		"https://gabrielbechara.com/",
		[]string{"golang-gin"},
	)
	if err != nil {
		log.Fatalf("failed to set up the validator: %v", err)
	}

	// Set up the middleware.
	middleware := jwtmiddleware.New(jwtValidator.ValidateToken)*/

	db, cleanup := getDB()
	dbPool = db
	defer cleanup()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong 100..%.",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hey-writer",
		})
	})


	r.GET("/metrics", PrometheusHandler())


	api := r.Group("/api")
	{
		api.GET("/writers", WritersHandler)
		api.GET("/writers/insertwriter/:Writer/:Color", InsertWritersHandler)
		//api.GET("/writers/insertwriter/:Writer/:Color", authMiddleware(middleware), InsertWritersHandler)
		//api.GET("/writers/insertwriter/:Writer/:Color", checkJWT(), InsertWritersHandler)
	}

	r.Run()
	//
}

/*
func authMiddleware(middleware *jwtmiddleware.JWTMiddleware) gin.HandlerFunc {
	log.Println("authMiddleware ")
	return func(c *gin.Context) {
	  log.Println("authMiddleware funct")	
	  // Get the client secret key
	  err := middleware.CheckJWT(handler)
	  log.Println("authMiddleware funct post err")
	  if err != nil {
		// Token not found
		fmt.Println(err)
		c.Abort()
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte("Unauthorized"))
		return
	  }
	}
}*/

func indexHandler(c *gin.Context) {
	c.HTML(200, "/public/index.html", nil)
}


func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func WritersHandler(c *gin.Context) {
	
	c.Header("Content-Type", "application/json")
	
	rows, err := dbPool.Query("SELECT * FROM writers")
	if err != nil {
		log.Fatalf("DB.Query: %v", err)
	}
	log.Println("get rows ")

	defer rows.Close()

	var myWriters = []Writer{}
	
	for rows.Next() {
		var id int
		var likes int
		var writer string
		var color string
		err := rows.Scan(&id, &likes, &writer, &color)
		if err != nil {
		  log.Fatalf("Rows.Scan: %v", err)
		}
		log.Println("appending rows ")

		myWriters = append(myWriters, Writer{ID: id, Likes: likes, Writer: writer, Color: color})
	}

	c.JSON(http.StatusOK, myWriters)
}


func InsertWritersHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	wr := Writer{1, 0, "test writer", "#4285F4"}

    wr.Writer = c.Params.ByName("Writer")    
	wr.Color = c.Params.ByName("Color")

	log.Println("writer to insert : ", wr.Writer)

	_, err := dbPool.Query("SELECT * FROM writers where writer = $1",wr.Writer)
	if err != nil {
		log.Fatalf("DB.Query: %v", err)
		log.Println("DB.Query: %v", err)
	}
	log.Println("writer selected ")

	log.Println("writer : inserting new writer : wr.Likes: ", wr.Likes, "wr.Writer: ", wr.Writer, "wr.Color: ", wr.Color)  
	_, err = dbPool.Exec("INSERT INTO writers (likes, writer, color) VALUES  ($1, $2, $3)", wr.Likes, wr.Writer, wr.Color)
	if err != nil {
		log.Fatalf("unable to create new Writer: %s", err)
	}	
	log.Println("writer inserted ")

	
	c.JSON(http.StatusOK, wr)
}
	
	
func getDB() (*sql.DB, func() error) {

	//dsn := fmt.Sprintf("user=%s dbname=%s sslmode=disable", os.Getenv("DB_IAM_USER"), os.Getenv("DB_NAME"))
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable", os.Getenv("INSTANCE_CONNECTION_NAME"), os.Getenv("DB_IAM_USER"), os.Getenv("DB_NAME"))
	//dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", "projects/qwiklabs-gcp-03-1200b28ed2d1/locations/us-central1/clusters/devopsdemo-alloydb-cluster-1/instances/devopsdemo-alloydb-instance-1", "postgres","password", os.Getenv("DB_NAME"))

  	log.Println("instance: ", os.Getenv("INSTANCE_CONNECTION_NAME"))
  	log.Println("dsn: ", dsn)  

  	var opts []cloudsqlconn.Option
	//var opts []alloydbconn.Option
	//if usePrivate != "" {
	//	opts = append(opts, cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()))
	// }
	opts = append(opts, cloudsqlconn.WithIAMAuthN())

	cleanup, err := pgxv4.RegisterDriver("cloudsql-postgres", opts...)
	//cleanup, err := pgxv4.RegisterDriver("alloydb", opts...)
	if err != nil {
  		log.Fatalf("Error on pgx4.RegisterDriver: %v", err)
	}


	dbPool, err := sql.Open("cloudsql-postgres", dsn)
	//dbPool, err := sql.Open("alloydb", dsn)
  	if err != nil {
    	log.Fatalf("Error on sql.Open: %v", err)
  	}

  	dropWriters := `DROP TABLE IF EXISTS writers;`
  	_, err = dbPool.Exec(dropWriters)

  	createWriters := `CREATE TABLE IF NOT EXISTS writers (
    	id SERIAL PRIMARY KEY,
    	likes INT,
		writer VARCHAR (255),
		color VARCHAR (255)
  	);`

  	_, err = dbPool.Exec(createWriters)
  	if err != nil {
    	log.Fatalf("unable to create table: %s", err)
  	}
  
  	for _, wr := range Writers {
		log.Println("wr.Likes: ", wr.Likes, "wr.Writer: ", wr.Writer, "wr.Color: ", wr.Color)  
  		_, err = dbPool.Exec("INSERT INTO writers (likes, writer, color) VALUES  ($1, $2, $3)", wr.Likes, wr.Writer, wr.Color)
		if err != nil {
			log.Fatalf("unable to create new Writer: %s", err)
		}	
  	} 

  	return dbPool, cleanup
}
