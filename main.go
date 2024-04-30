package main

import (
	"fmt"
	"go-projects/helper"
	"time"
	"sync"
)
    //package level variable
	const conferenceTickets int = 50 

	 var conferenceName = "Go Conference"
	 var remainingTickets uint = 50
     var bookings = make([]UserData, 0) //creating empty lists of map   //now  //struct
	 
	 //structure gives us the custom type of the data similar to java classes
	type UserData struct {
		 firstName  string
		 lastName  string
		 email  string
		 numberOfTickets uint
 	 }
     
	 var wg = sync.WaitGroup{} //creating waitgroup , package "sync" provides basic synncronization functionality

func main() {
	 
	     greetUser()
    
	      //greeting
	      //fmt.Printf("Welcome to %v booking application. \n", conferenceName)
	
    
      
	     for {
		
		    firstName,lastName,email,userTickets := getUserInput()
            isValidName,isValidEmail,isValidTicketNumber := helper.ValidateUserInput(firstName,lastName,email,userTickets,remainingTickets) 
        

	  	if isValidName && isValidEmail && isValidTicketNumber {

		    bookTicket(userTickets, firstName, lastName, email)
			
			
			wg.Add(1)    // Add() sets the number of goroutines to wait for(increase the counter by the provided number)
			go sendTicket(userTickets, firstName, lastName, email)
			
					
		    firstNames := getsFirstNames()
			fmt.Printf("The first names of the bookings are:  %v\n",firstNames)
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
	wg.Wait()   //Blocks until the WaitGroup counter is 0
}

func greetUser() {
	       fmt.Printf("Welcome to %v booking application.\n",conferenceName)
	       fmt.Printf("We have total of %v tickets and %v are still available.\n", conferenceTickets, remainingTickets)
	       fmt.Println("Get your tickets here to attend ")
}

func getsFirstNames() []string {

	    firstNames := []string{}
	    for _, booking := range bookings {
	    // strings.Fields(booking)
	    // var names = strings.Fields(booking)
	    firstNames = append(firstNames, booking.firstName)

	    }  //range iterates over element for different data structure(for arrays and slice range provide the index and vlaue for each element)
	
	   return firstNames 
}



func getUserInput() (string,string,string,uint) {
	var firstName string
	var lastName string
	var email string
	var userTickets uint
	
	fmt.Print("Enter your fist name: ")
	fmt.Scan(&firstName)
	 
	fmt.Print("Enter your last name: ")
	fmt.Scan(&lastName)

	fmt.Print("Enter your email adress: ")
	fmt.Scan(&email)

	fmt.Print("Enter number of tickets: ")
	fmt.Scan(&userTickets)

	return firstName,lastName,email,userTickets
}

func bookTicket(userTickets uint, firstName string,lastName string,email string){
	remainingTickets = remainingTickets - userTickets
	
	var userData = UserData {
		firstName: firstName,
		lastName: lastName,
		email: email,
		numberOfTickets: userTickets,
	}
	// //adding keyValue pairs
	// userData["firstName"] = firstName
	// userData["lastName"] = lastName
	// userData["email"] = email
	// userData["numberOfTickets"] = strconv.FormatUint(uint64(userTickets),10) //gives userTickets in string format
	
	bookings = append(bookings,userData)
	fmt.Println()
	fmt.Printf("List of bookings is %v\n", bookings)
    fmt.Println("*****************************************************************")
	fmt.Println()
	fmt.Printf("Thank you! %v %v for booking %v tickets. You will receive a confirmation email at %v\n",firstName,lastName,userTickets,email)
	fmt.Printf("Only %v tickets are remaining for %v\n",remainingTickets,conferenceName)

}

func sendTicket(userTickets uint,firstName string,lastName string, email string) {
	  
	  time.Sleep(5 * time.Second)
	  fmt.Println("\nSENDING TICKETS: ")
      var ticket = fmt.Sprintf("%v tickets for: %v %v", userTickets,firstName,lastName)
	  fmt.Println("*****************************************************************")
	  fmt.Printf("%v \nto email address: %v\n",ticket, email)
      fmt.Println("*****************************************************************")
 
	  wg.Done()    // Done()  Decrements the wait group counter by 1, So this is called by the goroutine to indicate that's it's finished (basically removes the waiting thread)
}