package main

/*
	Here is a very basic example's of how the apis work, ENJOY! :)
	Please be kind using this as artisan website is a paper castle, catches fire with nothing.

	How does this works? Simple!
	1 - You create a new api session, a session is needed to interact with the website safely.
	2 - You use the api as you want, they are decently documented.
	
	NOTE: I wrote this very quickly in about few hours for fun some things are thread safe others are NOT please pay attention of what you do.

	THIS IS FOR RECREATIVE USE ONLY, I TAKE NO RESPONSIBILITY FOR YOUR ACTIONS. 
	*Use it wisely*
*/

import artisan "artisanapi/src"

func main() {

	// Create's a new website session
	// You can have multiple sessions but i don't suggest doing it. The website is already slow and has tons of spaghetti code
	// let's stick to slow refreshing rates of requests and limitate us to 1 session.
	// Proxies are completely overkill here, i think they don't even have rate limits. 
	// PLEASE! Be kind and not abuse requests.
	session := artisan.NewAPISession()

	// Fetch all the product details
	hienMid, _ := session.ProductDetails(artisan.ProductDetailsBody{
		SirID: artisan.HienMidMPad,
		SizeID: artisan.SizeXL,
		ColorID: artisan.WineRedColor,
	})

	hayateKouMid, _ := session.ProductDetails(artisan.ProductDetailsBody{
		SirID: artisan.HayateKouMidMPad,
		SizeID: artisan.SizeXL,
		ColorID: artisan.NinjaBlackColor,
	})

	// Tries to add the products to the cart (it returns an error if outofstock)
	session.CartAdd(hienMid, 1)
	session.CartAdd(hayateKouMid, 1)

	// Set the shipping address
	session.SetShippingAddress(&artisan.ShippingAddress{
		Name:            "adadada",
		Surname:         "adadad",
		Email:           "wadawdadawd@gmail.com",
		Zipcode:         "123123",
		Fullname:        "adadad adadada",
		Province:        "wdawdawdawdaw",
		City:            "dadadad",
		Address:         "adadada",
		Building:        "dadadad",
		TelephoneNumber: "113123123123",
		Country:         "Italy",
	})

	// Create's the paypal checkout
	chekcoutHandler, _ := session.InstanceCheckout()

	// Opens the checkout in the browser
	chekcoutHandler.Open()
}
