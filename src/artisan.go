// Artisan unofficial APIs made by https://github.com/Camellia-ESL
package artisan

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// *********** DESCRIPTION ***********
// Small description of how the website works. The whole website is really simple it's all written in php, mainly works thru few requests and cookies.
// There's no need at all to utilize proxies, maybe avoid doing bulk requests like AllProductDetails every second to be fair.

// *********** APIS ***********

// Contains all the artisan routes/apis
const (
	// The main artisan website domain
	APIDomain string = "https://www.artisan-jp.com"
	// [POST] -> This api takes a FORMDATA as body with all the informations about product color, size, type/hardness, etc..
	// and returns informations about the product stock and more
	APIGetSyouhin string = APIDomain + "/get_syouhin.php"
	// [POST] -> This api send POST with an empty body because every information needed for the checkout is contained inside the cookiejar
	// it essentially create's the safe paypal checkout and send back in response the crafted paypal form checkout to pay
	// it can then be opened inside a browser setting display block for the form and adding a button to checkout the form
	APINjPaypalEng string = APIDomain + "/nj_paypal_eng.php"
)

// Represent a single session of the Artisan website APIs
type APISession struct {
	// If true display a log for everything that happens in the API calls (like a Debug log)
	EnableLogs bool
	// The cookie array that contains all the cookies of the website in this exact session instance
	Cookies []*http.Cookie
	// The static cookies that are present in almost every request
	staticCookies []*http.Cookie
	// The shipping address where to send everything
	shippingAddress *ShippingAddress
	// Wheter if the address to checkout has been set or not
	isAddressSet bool
}

// Create's a new APISession and init the session
func NewAPISession() *APISession {
	session := &APISession{
		EnableLogs:   false,
		isAddressSet: false,
	}

	session.initAPISession()

	return session
}

// *********** PRODUCTS ***********

// Represent an artisan mousepad (the pad also contains the hardness of the spongee)
type MPad string

// Contains all the available mousepads
const (
	ZeroClassicXSoftMPad MPad = "13"
	ZeroClassicSoftMPad  MPad = "14"
	ZeroClassicMidMPad   MPad = "12"

	RaidenClassicXSoftMPad MPad = "16"
	RaidenClassicMidMPad   MPad = "15"

	HayateOtsuXSoftMPad MPad = "120"
	HayateOtsuSoftMPad  MPad = "121"
	HayateOtsuMidMPad   MPad = "122"

	HayateKouXSoftMPad MPad = "130"
	HayateKouSoftMPad  MPad = "131"
	HayateKouMidMPad   MPad = "132"

	HienXSoftMPad MPad = "140"
	HienSoftMPad  MPad = "141"
	HienMidMPad   MPad = "142"

	ZeroXSoftMPad MPad = "160"
	ZeroSoftMPad  MPad = "161"
	ZeroMidMPad   MPad = "162"

	RaidenXSoftMPad MPad = "170"
	RaidenSoftMPad  MPad = "171"
	RaidenMidMPad   MPad = "172"

	Type99XSoftMPad MPad = "210"
	Type99SoftMPad  MPad = "211"
	Type99MidMPad   MPad = "212"

	ShidenkaiV2XSoftMPad MPad = "190"
	ShidenkaiV2MidMPad   MPad = "192"
)

// The total number of (mousepads types * hardness)
const MPadsTypesTotal int32 = 25

// All the mousepads ulrs
var MPadUrls = map[MPad]string{
	ZeroClassicXSoftMPad: APIDomain + "/cs-zero-eng.html",
	ZeroClassicSoftMPad:  APIDomain + "/cs-zero-eng.html",
	ZeroClassicMidMPad:   APIDomain + "/cs-zero-eng.html",

	RaidenClassicXSoftMPad: APIDomain + "/cs-raiden-eng.html",
	RaidenClassicMidMPad:   APIDomain + "/cs-raiden-eng.html",

	HayateOtsuXSoftMPad: APIDomain + "/fx-hayate-otsu-eng.html",
	HayateOtsuSoftMPad:  APIDomain + "/fx-hayate-otsu-eng.html",
	HayateOtsuMidMPad:   APIDomain + "/fx-hayate-otsu-eng.html",

	HayateKouXSoftMPad: APIDomain + "/fx-hayate-kou-eng.html",
	HayateKouSoftMPad:  APIDomain + "/fx-hayate-kou-eng.html",
	HayateKouMidMPad:   APIDomain + "/fx-hayate-kou-eng.html",

	HienXSoftMPad: APIDomain + "/fx-hien-eng.html",
	HienSoftMPad:  APIDomain + "/fx-hien-eng.html",
	HienMidMPad:   APIDomain + "/fx-hien-eng.html",

	ZeroXSoftMPad: APIDomain + "/fx-zero-eng.html",
	ZeroSoftMPad:  APIDomain + "/fx-zero-eng.html",
	ZeroMidMPad:   APIDomain + "/fx-zero-eng.html",

	RaidenXSoftMPad: APIDomain + "/fx-raiden-eng.html",
	RaidenSoftMPad:  APIDomain + "/fx-raiden-eng.html",
	RaidenMidMPad:   APIDomain + "/fx-raiden-eng.html",

	Type99XSoftMPad: APIDomain + "/fx-99-eng.html",
	Type99SoftMPad:  APIDomain + "/fx-99-eng.html",
	Type99MidMPad:   APIDomain + "/fx-99-eng.html",

	ShidenkaiV2XSoftMPad: APIDomain + "/fx-shidenkai-eng.html",
	ShidenkaiV2MidMPad:   APIDomain + "/fx-shidenkai-eng.html",
}

// Represent an artisan color
type Color string

// Contains all the artisan colors
const (
	WineRedColor      Color = "1"
	NinjaBlackColor   Color = "3"
	BlackColor        Color = "5"
	SnowWhiteColor    Color = "6"
	CoffeeBrownColor  Color = "8"
	DaidaiOrangeColor Color = "10"
	MatchaGreenColor  Color = "12"
	GrayColor         Color = "13"
)

// Contains all the artisan color names (you can search a color inside by using a color constant as key)
var ColorNames = map[Color]string{
	WineRedColor:      "WineRed",
	NinjaBlackColor:   "NinjaBlack",
	BlackColor:        "Black",
	SnowWhiteColor:    "SnowWhite",
	CoffeeBrownColor:  "CoffeeBrown",
	DaidaiOrangeColor: "Daidai Orange",
	MatchaGreenColor:  "MatchaGreen",
	GrayColor:         "Gray",
}

// Represent an artisan size
type Size string

// Contains all the artisan sizes
const (
	SizeS   Size = "1"
	SizeM   Size = "2"
	SizeL   Size = "3"
	SizeXL  Size = "4"
	SizeXXL Size = "5"
)

// Contains all the artisan size names (you can search a size inside by using a size constant as key)
var SizeNames = map[Size]string{
	SizeS:   "S",
	SizeM:   "M",
	SizeL:   "L",
	SizeXL:  "XL",
	SizeXXL: "XXL",
}

// Contains informations used to communicate the product thru the apis
type ProductDetailsBody struct {
	SirID   MPad
	SizeID  Size
	ColorID Color
}

// Represent a single Artisan product with all it's details
type Product struct {
	// The product id (note if the product is outofstock Id is "NON" and OutOfStock is set to true)
	Id string
	// The product prefix
	Prefix string
	// The product name short (note that for classing variants the classic after the name is not inserted)
	ShortName string
	// The product full length name including hardness and size
	FullName string
	// The product size
	Size string
	// The product color name
	Color string
	// The product price in yen(jpy)
	Price string
	// The hardness of the product
	Hardness string
	// The product direct url
	Url string
	// Wheter if the product is out of stock or not
	OutOfStock bool

	*ProductDetailsBody
}

// Fetch details about a single given Product
func (api *APISession) ProductDetails(pSearched ProductDetailsBody) (*Product, error) {

	// Create's a recover for possible panics
	defer func() {
		if err := recover(); err != nil {
			if api.EnableLogs {
				fmt.Println("[Panic] -> Error sending or parsing ProductDetails request, panic recovered")
			}
			return
		}
	}()

	// Encode the form data
	formData := url.Values{
		"kuni":  {"on"},
		"sir":   {string(pSearched.SirID)},
		"size":  {string(pSearched.SizeID)},
		"color": {string(pSearched.ColorID)},
	}

	// Create's the post request
	req, err := http.NewRequest("POST", APIGetSyouhin, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, err
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Elaborate the response
	resBodySplitted := strings.Split(string(resBody), "/")

	resProduct := &Product{
		Id:                 resBodySplitted[0],
		OutOfStock:         resBodySplitted[0] == "NON",
		Prefix:             resBodySplitted[1],
		ShortName:          strings.Split(resBodySplitted[2], " ")[0],
		FullName:           resBodySplitted[2],
		Price:              resBodySplitted[3],
		Hardness:           resBodySplitted[5],
		Size:               SizeNames[pSearched.SizeID],
		Color:              ColorNames[pSearched.ColorID],
		Url:                MPadUrls[pSearched.SirID],
		ProductDetailsBody: &pSearched,
	}

	return resProduct, nil
}

// Contains options for the AllProductDetails fetch request
type AllProductDetailsOptions struct {
	// An optional callback that get's called every time a product is fetched asynchronously (note that modifying the product passed as parameter
	// results in modifying also the product returned in the result slice)
	ProductFetchedCallback func(*Product)
}

// Fetch details about every Product (it does fetch using ProductDetails tasks asynchronously)
func (api *APISession) AllProductsDetails(options AllProductDetailsOptions) ([]*Product, error) {

	// Create's a recover for possible panics
	defer func() {
		if err := recover(); err != nil {
			if api.EnableLogs {
				fmt.Println("[Panic] -> Error sending or parsing AllProductDetails request, panic recovered")
			}
			return
		}
	}()

	// Create's a slice to hold all the products fetched
	var resProducts []*Product

	// Create's an array that contains every product to fetch
	allProducts := [...]MPad{
		ZeroClassicXSoftMPad,
		ZeroClassicSoftMPad,
		ZeroClassicMidMPad,
		RaidenClassicXSoftMPad,
		RaidenClassicMidMPad,
		HayateOtsuXSoftMPad,
		HayateOtsuSoftMPad,
		HayateOtsuMidMPad,
		HayateKouXSoftMPad,
		HayateKouSoftMPad,
		HayateKouMidMPad,
		HienXSoftMPad,
		HienSoftMPad,
		HienMidMPad,
		ZeroXSoftMPad,
		ZeroSoftMPad,
		ZeroMidMPad,
		RaidenXSoftMPad,
		RaidenSoftMPad,
		RaidenMidMPad,
		Type99XSoftMPad,
		Type99SoftMPad,
		Type99MidMPad,
		ShidenkaiV2XSoftMPad,
		ShidenkaiV2MidMPad,
	}

	// Prepare a sync wait group to wait for every request to finish before exiting the function
	wg := sync.WaitGroup{}
	wg.Add(len(allProducts) * len(ColorNames) * len(SizeNames))

	// Create's a lock to access resProducts safely
	resProductsLock := make(chan bool, 1)

	// Iterates every product * every color * every size and create's a task for every req to complete it asynchronously
	for _, product := range allProducts {
		for color := range ColorNames {
			for size := range SizeNames {
				go func(product MPad, color Color, size Size) {
					// Safely unlock the async function
					defer wg.Done()

					// Log if enable the request sending
					if api.EnableLogs {
						fmt.Printf("Sending ProductDetail request -> SirID: %s, ColorID: %s, SizeID: %s\n", product, color, size)
					}

					// Fetch details about the product
					pRes, _ := api.ProductDetails(ProductDetailsBody{SirID: product, ColorID: color, SizeID: size})

					// Check if the result is valid
					if pRes == nil {
						return
					}

					// Acquire the lock
					resProductsLock <- true

					// Safely append the product fetched
					resProducts = append(resProducts, pRes)

					// Calls the optional callback
					if options.ProductFetchedCallback != nil {
						options.ProductFetchedCallback(pRes)
					}

					// Unlock the resource
					<-resProductsLock
				}(product, color, size)
			}
		}
	}

	// Waits every request to finish
	wg.Wait()

	return resProducts, nil
}

// Represent the official country name type
type Country string

// Contains all the available countries
const (
	Argentina            Country = "Argentina"
	Australia            Country = "Australia"
	Austria              Country = "Austria"
	Azerbaijan           Country = "Azerbaijan"
	Bahrain              Country = "Bahrain"
	Bangladesh           Country = "Bangladesh"
	Belgium              Country = "Belgium"
	BosniaandHerzegovina Country = "BosniaandHerzegovina"
	Brazil               Country = "Brazil"
	BruneiDarussalam     Country = "BruneiDarussalam"
	Bulgaria             Country = "Bulgaria"
	Canada               Country = "Canada"
	Chile                Country = "Chile"
	China                Country = "China"
	Croatia              Country = "Croatia"
	Cyprus               Country = "Cyprus"
	CzechRepublic        Country = "Czech Republic"
	Denmark              Country = "Denmark"
	Egypt                Country = "Egypt"
	Estonia              Country = "Estonia"
	Finland              Country = "Finland"
	France               Country = "France"
	Georgia              Country = "Georgia"
	Germany              Country = "Germany"
	Greece               Country = "Greece"
	Greenland            Country = "Greenland"
	Guam                 Country = "Guam"
	Hungary              Country = "Hungary"
	Iceland              Country = "Iceland"
	India                Country = "India"
	Ireland              Country = "Ireland"
	Italy                Country = "Italy"
	Kazakhstan           Country = "Kazakhstan"
	Korea                Country = "Korea"
	Kosovo               Country = "Kosovo"
	Kuwait               Country = "Kuwait"
	Latvia               Country = "Latvia"
	Liechtenstein        Country = "Liechtenstein"
	Lithuania            Country = "Lithuania"
	Luxembourg           Country = "Luxembourg"
	Macedonia            Country = "Macedonia"
	Malaysia             Country = "Malaysia"
	Malta                Country = "Malta"
	Mexico               Country = "Mexico"
	Monaco               Country = "Monaco"
	Montenegro           Country = "Montenegro"
	Morocco              Country = "Morocco"
	Netherlands          Country = "Netherlands"
	NewCaledonia         Country = "NewCaledonia"
	NewZealand           Country = "NewZealand"
	Norway               Country = "Norway"
	Oman                 Country = "Oman"
	Peru                 Country = "Peru"
	Poland               Country = "Poland"
	Portugal             Country = "Portugal"
	PuertoRico           Country = "PuertoRico"
	Qatar                Country = "Qatar"
	Romania              Country = "Romania"
	SanMarino            Country = "SanMarino"
	SaudiArabia          Country = "Saudi Arabia"
	Serbia               Country = "Serbia"
	Singapore            Country = "Singapore"
	Slovakia             Country = "Slovakia"
	Slovenia             Country = "Slovenia"
	SouthAfrica          Country = "SouthAfrica"
	Spain                Country = "Spain"
	SriLanka             Country = "SriLanka"
	Sweden               Country = "Sweden"
	Switzerland          Country = "Switzerland"
	Taiwan               Country = "Taiwan"
	Thailand             Country = "Thailand"
	Turkey               Country = "Turkey"
	UnitedArabEmirates   Country = "United Arab Emirates"
	UnitedKingdom        Country = "United Kingdom"
	UnitedStates         Country = "United States"
	VietNam              Country = "VietNam"
)

// Represent the shipping address structure
type ShippingAddress struct {
	Name            string
	Surname         string
	Email           string
	Zipcode         string
	Fullname        string
	Province        string
	City            string
	Address         string
	Building        string
	TelephoneNumber string
	Country         Country
}

// Set's the shipping address and return true if infos are correct (also auto URL encode strings)
func (api *APISession) SetShippingAddress(address *ShippingAddress) bool {

	// Check if there is an empty field
	if !allFieldsNonEmpty(*address) {
		return false
	}

	// Set's the address
	api.shippingAddress = address
	api.isAddressSet = true

	// Search for the info cookie
	var infoCookie *http.Cookie
	for _, cookie := range api.Cookies {
		if cookie.Name == "info" {
			infoCookie = cookie
		}
	}

	if infoCookie == nil {
		panic("Info cookie cannot be nil, something went wrong in API initialization")
	}

	// Set's the address cookies correctly according to the address
	addrRef := reflect.ValueOf(*address)

	for i := 0; i < addrRef.NumField(); i++ {
		// Enumerate all the shippingAddress fields
		field := addrRef.Field(i)

		// Create's the separator for the cookies
		separator := ""
		if i != (addrRef.NumField() - 1) {
			separator = "*"
		}

		// Add the value to the cookies
		infoCookie.Value += url.QueryEscape(field.String()) + separator
	}

	return true
}

// Clear the cart
func (api *APISession) CartClear() {

	// Search for the info cookie
	var cartCookie *http.Cookie
	for _, cookie := range api.Cookies {
		if cookie.Name == "cart" {
			cartCookie = cookie
		}
	}

	if cartCookie == nil {
		panic("Info cookie cannot be nil, something went wrong in API initialization")
	}

	// Clear the cart
	cartCookie.Value = ""
}

// Add a product to the cart, return's nil if everything went fine or a error instead
// NOTE: You cannot insert the same product twice you must instead specify the right quantity, if you want to clear the cart
// consider using CartClear to cleanup and reinsert
func (api *APISession) CartAdd(p *Product, quantity uint32) error {

	if p == nil {
		return errors.New("Product cannot be nil")
	}

	if quantity <= 0 {
		return errors.New("The quantity must be > 1")
	}

	if p.OutOfStock {
		return errors.New("Product cannot be out of stock")
	}

	// Search for the info cookie
	var cartCookie *http.Cookie
	for _, cookie := range api.Cookies {
		if cookie.Name == "cart" {
			cartCookie = cookie
		}
	}

	if cartCookie == nil {
		panic("Info cookie cannot be nil, something went wrong in API initialization")
	}

	// Construct the product for the cart and insert it after the last product inserted
	// 4562332172443,FX-HI-XS-S-R HIEN FX XSOFT S Wine red,2,2700.0,1
	// ProductID, Prefix Fullname, Quantity, Price (not multiplied by quantity), 1???
	separator := ","
	cartProductBuilt := ""

	if cartCookie.Value != "" {
		cartProductBuilt += separator
	}

	cartProductBuilt += p.Id + separator
	cartProductBuilt += p.Prefix + " "
	cartProductBuilt += p.FullName + separator
	cartProductBuilt += fmt.Sprintf("%d", quantity) + separator
	cartProductBuilt += fmt.Sprintf("%d", 1)
	cartCookie.Value += url.QueryEscape(cartProductBuilt)

	return nil
}

// Contains the default checkout template constants
const (
	checkoutDefaultPage string = `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset='utf-8'>
			<meta http-equiv='X-UA-Compatible' content='IE=edge'>
			<title>Checkout</title>
			<meta name='viewport' content='width=device-width, initial-scale=1'>
		</head>

		<body>

		</body>
		</html>
	`

	checkoutFormBtn string = `
		<input src="https://www.paypalobjects.com/en_US/i/btn/btn_PaywithPP_25h.gif" id="btn_submit"
						name="_eventId_paywithpaypal" class="paypalButton" alt="Pay with PayPal" type="image">
	`
)

// Represent the checkout handler, used to checkout the cart
type Checkout struct {
	// Contains the data about the paypal payment form (used to create the paypal payment session)
	pplFormData string
}

// Opens the checkout in the browser
func (checkout *Checkout) Open() error {

	// Construct the checkout page using the default page + the form received from artisan + adding the checkout button to the form
	// and setting display: block to the form
	checkoutPage := checkoutDefaultPage

	// Insert the artisan form
	bodyIndex := strings.Index(checkoutPage, "<body>")
	if bodyIndex == -1 {
		panic("Can't find body in the default constant html string, something is wrong!")
	}

	insertPos := bodyIndex + len("<body>")

	formIndex := strings.Index(checkout.pplFormData, "<form")
	if formIndex == -1 {
		return errors.New("Error, no form was found in the checkout")
	}

	// Create the new string with the insertion
	checkoutPage = checkoutPage[:insertPos] + checkout.pplFormData[formIndex:] + checkoutPage[insertPos:]

	// Replace the display: none with block
	checkoutPage = strings.Replace(checkoutPage, "display:none;", "display:block;", -1)

	// Add the checkout button to the form
	formEndIndex := strings.Index(checkoutPage, "</form>")
	if formEndIndex == -1 {
		return errors.New("Error, no form ending was found in the checkout")
	}

	checkoutPage = checkoutPage[:formEndIndex] + checkoutFormBtn + checkoutPage[formEndIndex:]

	// Open up the temporary html in the browser (note it constructs a temp folder with a temp html file that should be cleaned later)
	return openHTMLInBrowser(checkoutPage, "./temp_checkout_dir")
}

// Create's the checkout using the product added to the cart in this session
// it also tries to spawn a new instance of the default browser in this OS to checkout the cart
func (api *APISession) InstanceCheckout() (*Checkout, error) {

	// Create's a recover for possible panics
	defer func() {
		if err := recover(); err != nil {
			if api.EnableLogs {
				fmt.Println("[Panic] -> Error sending Checkout request, panic recovered")
			}
			return
		}
	}()

	// Check if the address has been set
	if !api.isAddressSet {
		return nil, errors.New("Error, address not set.")
	}

	// Create's a clean jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.New("Failed to create cookie jar: " + err.Error())
	}

	// Create an HTTP client with the cookie jar
	client := &http.Client{
		Jar: jar,
	}

	// Create the POST request
	req, err := http.NewRequest("POST", APINjPaypalEng, bytes.NewBuffer([]byte(``)))
	if err != nil {
		return nil, errors.New("Failed to create request: " + err.Error())
	}

	// Add the cookies to the request
	for _, cookie := range api.Cookies {
		req.AddCookie(cookie)
	}

	for _, cookie := range api.staticCookies {
		req.AddCookie(cookie)
	}

	// Set the headers
	req.Header.Set("Content-Type", "text/html; charset=UTF-8")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Failed to send request: " + err.Error())
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Failed to read response body: " + err.Error())
	}

	return &Checkout{
		pplFormData: string(body),
	}, nil
}

// Init a new api session
func (api *APISession) initAPISession() {

	// Create fresh cookies
	api.Cookies = []*http.Cookie{
		{
			Name:  "cart",
			Value: "",
			Path:  "/",
		},
		{
			Name:  "info",
			Value: "",
			Path:  "/",
		},
	}

	// Create's the static cookies array structure to add to the requests summed with the other custom cookies
	api.staticCookies = []*http.Cookie{
		{
			Name:  "lung",
			Value: "jpf",
			Path:  "/",
		},
		{
			Name:  "overems",
			Value: "6300%2F2",
			Path:  "/",
		},
		{
			Name:  "overwgt",
			Value: "916.5",
			Path:  "/",
		},
	}

	// Set's a null shipping address
	api.shippingAddress = &ShippingAddress{}
}

// Function to check if all string properties in a struct are non-empty
func allFieldsNonEmpty(s any) bool {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		panic("allFieldsNonEmpty: input is not a struct")
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String && field.String() == "" {
			return false
		}
	}
	return true
}

// An helper function to open html strings directly as browser pages (storing a temp file)
func openHTMLInBrowser(htmlContent, directory string) error {

	// Ensure the directory exists
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to create directory: %w", err)
	}

	// Generate a UUID for the file name
	fileName := uuid.New().String() + ".html"

	// Create the full file path
	filePath := filepath.Join(directory, fileName)

	// Create and open the file
	tmpFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// Write the HTML content to the file
	_, err = tmpFile.Write([]byte(htmlContent))
	if err != nil {
		return fmt.Errorf("Failed to write to temp file: %w", err)
	}

	// Determine the command to open the file based on the OS
	var openCommand string
	var args []string

	switch runtime.GOOS {
	case "windows":
		openCommand = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", filePath}
	case "darwin":
		openCommand = "open"
		args = []string{filePath}
	default:
		openCommand = "xdg-open"
		args = []string{filePath}
	}

	// Execute the command to open the file in the browser
	cmd := exec.Command(openCommand, args...)
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("Failed to open file in browser: %w", err)
	}

	return nil
}
