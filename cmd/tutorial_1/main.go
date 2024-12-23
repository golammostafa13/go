package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getInput(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	return strings.TrimSpace(input), err
}

func createNewBill() bill {
	name, _ := getInput("Create a new bill name: ")
	b := newBill(name)
	return b
}

func promptOptions(b bill) {
	opt, _ := getInput("Choose option (a - add item, s - save bill, t - add tip): ")
	switch opt {
	case "a":
		name, _ := getInput("Item name: ")
		price, _ := getInput("Item price: ")
		priceFloat, err := strconv.ParseFloat(price, 64)
		if err != nil {
			fmt.Println("The price must be a number")
			promptOptions(b)
		}
		b.addItem(name, priceFloat)
		fmt.Println("Item added -", name, price)
		promptOptions(b)
	case "s":
		b.save()
		fmt.Println("You chose to save the bill", b.name)
	case "t":
		tip, _ := getInput("Enter tip amount ($): ")
		tipFloat, err := strconv.ParseFloat(tip, 64)
		if err != nil {
			fmt.Println("The tip must be a number")
			promptOptions(b)
		}
		b.updateTip(tipFloat)
		fmt.Println("Tip added to the bill", b.name)
		promptOptions(b)
	default:
		fmt.Println("Choose a valid option")
		promptOptions(b)
	}
}

func main() {
	myBill := createNewBill()
	promptOptions(myBill)
}
