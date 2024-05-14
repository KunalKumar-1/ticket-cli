package main

import (
	"fmt"
	"go-projects/helper"
	"image"
	"image/color"
	_ "image/png"
	"os"
	//"os/exec"
	"sync"
	"time"
    "github.com/mbndr/figlet4go"
)

// package level variable 
const conferenceTickets int = 100

var conferenceName = "GEEX ENVENT'S"

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

// func  execute(cmd string) { //terminal cmd comamnds
// 	out, err := exec.Command(cmd).Output()

// 	if err != nil {
// 		fmt.Printf("%s",err)
// 	}
//     output := string(out[:])
// 	fmt.Println(output)
// }

// creating func center to align the greeting text in center
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

	//greeting
	//fmt.Printf("Welcome to %v booking application. \n", conferenceName)

	for {

		firstName, lastName, email, userTickets := getUserInput()
		isValidName, isValidEmail, isValidTicketNumber := helper.ValidateUserInput(firstName, lastName, email, userTickets, remainingTickets)

		if isValidName && isValidEmail && isValidTicketNumber {

			bookTicket(userTickets, firstName, lastName, email)

			wg.Add(1) // Add() sets the number of goroutines to wait for(increase the counter by the provided number)
			go sendTicket(userTickets, firstName, lastName, email)

			firstNames := getsFirstNames()
			fmt.Printf("The first names of the bookings are:  %v\n", firstNames)
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
	ramp := "@#:|:. " //@#+=.
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
	// fmt.Printf("WELCOME TO (%v) EVENT BOOKING CLI.\n\n", conferenceName)
	// fmt.Printf("We have total of %v tickets and %v are still available.\n\n", conferenceTickets, remainingTickets)
	// fmt.Println("|| GRAB THE LIMITED PASS HERE ||")
}

func getsFirstNames() []string {

	firstNames := []string{}
	for _, booking := range bookings {
		// strings.Fields(booking)
		// var names = strings.Fields(booking)
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
	// //adding keyValue pairs
	// userData["firstName"] = firstName
	// userData["lastName"] = lastName
	// userData["email"] = email
	// userData["numberOfTickets"] = strconv.FormatUint(uint64(userTickets),10) //gives userTickets in string format

	bookings = append(bookings, userData)
	fmt.Println()
	fmt.Printf("List of bookings is %v\n", bookings)
	for j:= 0; j<190; j++{
		fmt.Print("=")
	}
    time.Sleep(2 * time.Second)
	fmt.Printf("Thank you!! %v %v for Booking %v Tickets.\n\nYou will Soon receive a Confirmation email at: %v\n", firstName, lastName, userTickets, email)
	fmt.Printf("\nOnly %v Tickets are remaining for %v.\n\n", remainingTickets, conferenceName)

}

func sendTicket(userTickets uint, firstName string, lastName string, email string) {
    fmt.Print("\n\n ==> Generating Tickets")
	for k:=0; k<3; k++ {
		time.Sleep(1 * time.Second)
		fmt.Print("..")
	}
	fmt.Println("\n\n	[Generated Successfully]")
	time.Sleep(4 * time.Second)
	fmt.Print("\n ==> SENDING TICKETS ")
	fmt.Print("Please Wait")
	for i := 0; i < 20; i++ {
		time.Sleep(1 * time.Second)
		fmt.Print("...")
	}
	fmt.Print(".\n\n")
	fmt.Println()
	var ticket = fmt.Sprintf("%v Tickets for: %v %v", userTickets, firstName, lastName)
	fmt.Println("******************************************************************")
	fmt.Printf("%v \nSend to email address: %v\n", ticket, email)
	fmt.Println("******************************************************************")
	time.Sleep(3 * time.Second)
	fmt.Println("Send Sucessfully.....\n")
	time.Sleep(1 * time.Second)
	//fmt.Println(center(renderStr,180))
	figletAsii(conferenceName)
	fmt.Println(center("  We're grateful for your booking. Visit us again for more great experiences!!",189))
	for j:= 0; j<190; j++{
		fmt.Print("#")
	}
	fmt.Println("\n")
	wg.Done() // Done()  Decrements the wait group counter by 1, So this is called by the goroutine to indicate that's it's finished (basically removes the waiting thread)
}
