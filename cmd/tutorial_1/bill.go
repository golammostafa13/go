package main

import (
	"fmt"
	"os"
)

type bill struct {
	name  string
	items map[string]float64
	tip   float64
}

func newBill(name string) bill {
	b := bill{
		name:  name,
		items: map[string]float64{},
		tip:   0,
	}
	return b
}

func (b *bill) format() string {
	fs := "Bills breakdown \n"
	total := b.tip

	for k, v := range b.items {
		fs += fmt.Sprintf("%-25v ...$%v \n", k+":", v)
		total += v
	}

	fs += fmt.Sprintf("%-25v ...$%0.2f \n", "Tip:", b.tip)

	fs += fmt.Sprintf("%-25v ...$%0.2f \n", "Total:", total)

	return fs
}

func (b *bill) updateTip(number float64) {
	b.tip = number
}

func (b *bill) addItem(name string, price float64) {
	b.items[name] = price
}

func (b *bill) save() {
	// Ensure the bills directory exists
	err := os.MkdirAll("bills", os.ModePerm)
	if err != nil {
		panic(err)
	}

	data := []byte(b.format())
	err = os.WriteFile("bills/"+b.name+".txt", data, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Bill was saved to file")
}
