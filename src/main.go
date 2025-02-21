package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {
	// PostgreSQL 연결 문자열
	connStr := "user=postgres password=password dbname=dbname sslmode=disable host=localhost port=5432"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 연결 확인
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database!")

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// GET 방식 예제
	app.Get("/fiber", func(c *fiber.Ctx) error {
		return c.SendString("Hello, get fiber!")
	})

	// GET 요청: 동적 경로 (URL 파라미터)
	app.Get("/user/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		return c.SendString("Hello, " + name)
	})

	// POST 요청 (JSON 데이터 받기)
	app.Post("/user", func(c *fiber.Ctx) error {
		type User struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}
		return c.JSON(fiber.Map{"message": "User created", "user": user})
	})

	// POST 요청 (postgreSQL DB연동 확인 / 호출은 postman에서)
	app.Post("/getting", func(c *fiber.Ctx) error {
		type User struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		// 요청 본문에서 JSON 데이터 파싱
		user := new(User)
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON",
			})
		}
		// 데이터 삽입
		_, err := db.Exec("INSERT INTO user_table(user_email, user_pwd) VALUES($1, $2)", user.Email, user.Password)
		if err != nil {
			log.Println("Error inserting data: ", err) // 여기서 에러를 출력해보세요.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to insert data",
			})
		}

		// 응답
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": fmt.Sprintf("User %s created successfully!", user.Email),
		})
	})

	app.Listen(":8080")
}
