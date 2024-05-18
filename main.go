package main

import (
	"fmt"
	"go-projects/helper"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"sync"
	"time"
	"github.com/joho/godotenv" 
	"github.com/mbndr/figlet4go"
	"net/smtp"
)

// package level variable
const conferenceTickets int = 100

var conferenceName = "Geex-Conf 2024."

var remainingTickets uint = 100
var bookings = make([]UserData, 0) //creating empty lists of map   //now  //struct

// structure gives us the custom type of the data similar to java classes
type UserData struct {
	firstName       string
	lastName        string
	email           string
	numberOfTickets uint
}

var wg = sync.WaitGroup{} //creating waitgroup , package "sync" provides basic synncronization functionality

func center(s string, w int) string {
	return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(s))/2, s))
}

func loadImage(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err //error handeling here
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func grayscale(c color.Color) int {
	r, g, b, _ := c.RGBA()
	return int(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
}

func avgPixel(img image.Image, x, y, w, h int) int {
	cnt, sum, max := 0, 0, img.Bounds().Max
	for i := x; i < x+w && i < max.X; i++ {
		for j := y; j < y+h && j < max.Y; j++ {
			sum += grayscale(img.At(i, j))
			cnt++
		}
	}
	return sum / cnt
}

func main() {

	greetUser()

	for {

		firstName, lastName, email, userTickets := getUserInput()
		isValidName, isValidEmail, isValidTicketNumber := helper.ValidateUserInput(firstName, lastName, email, userTickets, remainingTickets)

		if isValidName && isValidEmail && isValidTicketNumber {

			bookTicket(userTickets, firstName, lastName, email)

			wg.Add(1) // Add() sets the number of goroutines to wait for(increase the counter by the provided number)
			go sendTicket(userTickets, firstName, lastName, email)

			firstNames := getsFirstNames()
			fmt.Printf("Names of the Booked tickects:  %v\n", firstNames)
			fmt.Println()

			if remainingTickets == 0 {
				//end program
				fmt.Println("Our conference is booked out. See you next Year!!")
				break
			}

		} else {
			fmt.Println()
			fmt.Println("THE ERROR !!\n")
			if !isValidName {
				fmt.Println("First name or last name you entered is too short. ")
			}
			if !isValidEmail {
				fmt.Println("Email address you entered doesn't contain @ sign. ")
			}
			if !isValidTicketNumber {
				fmt.Println("Number of tickets you entered is invalid ")
			}
			fmt.Println()
		}

	}
	wg.Wait() //Blocks until the WaitGroup counter is 0
}
func figletAsii(name string) {
	//Renderer
	ascii := figlet4go.NewAsciiRender()
    // Adding the colors to RenderOptions
	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{
		//adding color to the Text
		figlet4go.ColorGreen,

	} 
	renderStr, _ := ascii.RenderOpts(conferenceName,options)
	fmt.Println(center(renderStr,180))
}

func greetUser() {
	//loading image & calculations of renderer for the display
	img, err := loadImage("image/final.png")
	if err != nil {
		panic(err)
	}
	ramp := "@#+=. " //@#+=.
	max := img.Bounds().Max
	scaleX, scaleY := 10, 5 //scale
	for y := 0; y < max.Y; y += scaleX {
		for x := 0; x < max.X; x += scaleY {
			c := avgPixel(img, x, y, scaleX, scaleY)
			fmt.Print(string(ramp[len(ramp)*c/65536]))
		}
		//time.Sleep(time.Second)
		fmt.Println()
	}
	time.Sleep(1 * time.Second)
	for j:= 0; j<190; j++{
		fmt.Print("-")
	}
	//fmt.Println(center(conferenceName,110))
	greeting := fmt.Sprintf("WELCOME TO (%v) EVENT BOOKING CLI.\n\n", conferenceName)
	greetings := fmt.Sprintf("We have total of %v tickets and %v are still available.\n\n", conferenceTickets, remainingTickets)
	greetingss := fmt.Sprintf("|| GRAB THE LIMITED PASS HERE ||")
	fmt.Println("\n\n")
  // Aligining it to center
     fmt.Println(center(greeting,182))
	 fmt.Println(center(greetings,180))
	 fmt.Println(center(greetingss,180))
	 for j:= 0; j<190; j++{
		fmt.Print("-")
	}

}

func getsFirstNames() []string {

	firstNames := []string{}
	for _, booking := range bookings {

		firstNames = append(firstNames, booking.firstName)

	} //range iterates over element for different data structure(for arrays and slice range provide the index and vlaue for each element)

	return firstNames
}

func getUserInput() (string, string, string, uint) {
	var firstName string
	var lastName string
	var email string
	var userTickets uint

	fmt.Print("Enter your fist name: ")
	fmt.Scan(&firstName)

	fmt.Print("\nEnter your last name: ")
	fmt.Scan(&lastName)

	fmt.Print("\nEnter your email adress: ")
	fmt.Scan(&email)

	fmt.Print("\nEnter number of tickets: ")
	fmt.Scan(&userTickets)

	return firstName, lastName, email, userTickets
}

func bookTicket(userTickets uint, firstName string, lastName string, email string) {
	remainingTickets = remainingTickets - userTickets

	var userData = UserData{
		firstName:       firstName,
		lastName:        lastName,
		email:           email,
		numberOfTickets: userTickets,
	}
	
	bookings = append(bookings, userData)
	fmt.Println()
	fmt.Printf("List of Bookings are %v\n", bookings)
	for j:= 0; j<190; j++{
		fmt.Print("=")
	}
    time.Sleep(2 * time.Second)
	fmt.Printf("Thank you!! %v %v for Booking %v Tickets.\n\n\nYou will Soon receive a Confirmation email at: %v\n", firstName, lastName, userTickets, email)
	fmt.Printf("\nOnly %v Tickets are remaining for %v.\n\n", remainingTickets, conferenceName)
	fmt.Println("******************************************************************")
}

func sendTicket(userTickets uint, firstName string, lastName string, email string) {
    fmt.Print("\n\n ==> Generating Tickets")
	for k:=0; k<3; k++ {
		time.Sleep(1 * time.Second)
		fmt.Print("..")
	}
	fmt.Println("\n\n	[Generated Successfully]")
	time.Sleep(2 * time.Second)
	fmt.Print("\n ==> SENDING TICKETS ")
	fmt.Print("Please Wait")
	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		fmt.Print("...")
	}
	fmt.Print(".\n\n")
	fmt.Println()
	var ticket = fmt.Sprintf("   (%v) Tickets for: %v %v", userTickets, firstName, lastName)
	fmt.Println("******************************************************************")
	fmt.Printf("%v \n\n   Sending to email address: [%v]\n", ticket, email)

	//messages for the mail box
	//var msgs = fmt.Sprintf("Subject: Congrats🔥 %v Tickets Confirmed\n",conferenceName)
	var msg = fmt.Sprintf("\nWelcome to the Geeky Group of Techies 🧑🏻‍💻 at %v.  \n\nDear %v %v , \nWe Confirm  your registration of %v Tickets for (%v)", conferenceName, firstName, lastName, userTickets, conferenceName)
	var msgR = fmt.Sprintf("\nOnly %v Tickets are Remaining.",remainingTickets)
	//fmt.Println("******************************************************************")
	//Sending mail here.
	sendMails(email,conferenceName,msg,msgR)

	fmt.Println("\n\n       [Mail Send Sucessfully!]\n")
	time.Sleep(4 * time.Second)
	fmt.Println("******************************************************************\n")
	//fmt.Println(center(renderStr,180))
	figletAsii(conferenceName)
	fmt.Println(center("  We're Greatful For Your Booking. Please Visit Us Again!!",189))
	for j:= 0; j<190; j++{
		fmt.Print("#")
	}
	fmt.Println("\n")
	wg.Done() // Done()  Decrements the wait group counter by 1, So this is called by the goroutine to indicate that's it's finished (basically removes the waiting thread)
}

//smtp PlainAuth
func sendMails(email string,conferenceName string, msg string,msgR string) {

	err2 := godotenv.Load()
     if err2 != nil {
		panic(err2)
	 }
	 envVar := os.Getenv("mail_pass")
	
	auth := smtp.PlainAuth(
		"",   //keep it empty for username to be your mail id
		"YourMail@gmail.com",//from whoch you are goona send the mail
	         envVar,  // you access token in the form of environment variable here
		"smtp.gmail.com",
	)

	message := "Subject: Your Registration Confirmation for "+conferenceName+"\n"+ msg +"\nLooking forward to seeing you at the event.\n "+msgR +"\nFor any queries please reach out to the organizer at: Team Geex.\n\nBest Regards,\nTeam "+ conferenceName

	err := smtp.SendMail(
		"smtp.gmail.com:587", //587 for gmail
        auth,
		"YourMail@gmail.com", //from which you are going to send the mail
		[]string{email},
		[]byte(message),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
