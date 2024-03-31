package main
import (
  "fmt"
  "math/rand"
  "time"
  "net/smtp"
  "strconv"
)

func main() {
  otp := time.Now().Second() * rand.Intn(100) 
  var n int
  var name, email string
  fmt.Printf("Enter Name: ")
  fmt.Scanln(&name)
  fmt.Printf("Enter Email: ") 
  fmt.Scanln(&email)
  fmt.Printf("Otp: ", otp)
  send(strconv.Itoa(otp), email) 
  fmt.Println("Enter the OTP you received")
  fmt.Scanln(&n) 
  if n == otp {
    fmt.Println("You are registered!") 
  }else {
    fmt.Println("Invalid OTP") 
  }

}

func send (body string, to string) {
  from := "youemail@mailhere"
  pass := "yourapppasswordhere"  // using gmail create app password and use it
  msg := "Subject:  Your Personal OTP\n\n" + body 
  
  err := smtp.SendMail("smtp.gmail.com:587", 
         smtp.PlainAuth("",from,pass,"smtp.gmail.com"),
         from, []string{to}, []byte(msg))
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println("Email sent!")
}
