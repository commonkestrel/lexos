package lexos

import (
	"errors"
	"fmt"
	"strings"

	isbnpkg "github.com/moraes/isbn"
	"github.com/playwright-community/playwright-go"
)

const (
    lexile_url = "https://hub.lexile.com/find-a-book/book-details/"
    lexile_selector = "#content > div > div > div > div.details > div.metadata > div.sc-kexyCK.cawTwh > div.header-info > div > span"
    
    atos_url = "https://www.arbookfind.com/UserType.aspx?RedirectURL=%2fadvanced.aspx"
    rad = "#radLibrarian"
    submit = "#btnSubmitUserType"
    isbn_box = "#ctl00_ContentPlaceHolder1_txtISBN"
    search = "#ctl00_ContentPlaceHolder1_btnDoIt"
    search_fail = "#ctl00_ContentPlaceHolder1_lblSearchResultFailedLabel"
    title = "#book-title"
    atos_level = "#ctl00_ContentPlaceHolder1_ucBookDetail_lblBookLevel"
    ar_points = "#ctl00_ContentPlaceHolder1_ucBookDetail_lblPoints"
)

const (
    InvalidISBN = "invalid isbn"
    BrowserErr = "error opening browser"
)

var (
    pw *playwright.Playwright
    browser playwright.Browser
    page playwright.Page
)

func Get(isbn string) (int, float64, float64, error) {
    isbn = strings.ReplaceAll(isbn, "-", "")
    valid := isbnpkg.Validate(isbn)
    if !valid {
        return -1, -1, -1, errors.New(InvalidISBN)
    }

    var err error
    pw, err = playwright.Run()
    if err != nil {
        return -1, -1, -1, errors.New(BrowserErr)
    }
    defer pw.Stop()

    browser, err = pw.Chromium.Launch()
    if err != nil {
        return -1, -1, -1, errors.New(BrowserErr)
    }
    defer browser.Close()

    page, err = browser.NewPage()
    if err != nil {
        return -1, -1, -1, errors.New(BrowserErr)
    }

    atos, ar := atos(isbn)
    lex := lexile(isbn)
    return lex, atos, ar, nil
}

func Install() {
    run := playwright.RunOptions{Browsers: []string{"chromium"}, Verbose: false}
    playwright.Install(&run)
}

func lexile(isbn string) int {
    page.Goto(fmt.Sprint(lexile_url, isbn))
    if page.URL() == "https://hub.lexile.com/find-a-book/book-results" {
        return -1
    }

    str, err := page.TextContent(lexile_selector)
    if err != nil {
        return -1
    }
    var lex int
    if _, err := fmt.Sscan(str, &lex); err != nil {
        return -1
    }
    return lex
}

func atos(isbn string) (float64, float64) {
    page.Goto(atos_url)
    page.Click(rad) //Select Librarian and submit
    page.Click(submit)

    page.WaitForSelector(isbn_box)
    page.Type(isbn_box, isbn)
    page.Click(search)
    
    page.WaitForLoadState("domcontentloaded")
    fail, _ := page.Locator(search_fail)
    count, _ := fail.Count()
    if count > 0 {
        return -1, -1
    }

    page.WaitForSelector(title)
    page.Click(title) //Click on first book

    var atos float64
    var ar float64
    AtosStr, err := page.TextContent(atos_level) //Get level from selector
    if err != nil {
        AtosStr = "-1"
    }
    ArStr, err := page.TextContent(ar_points)
    if err != nil {
        ArStr = "-1"
    }
    
    fmt.Sscan(ArStr, &ar)
    fmt.Sscan(AtosStr, &atos)
    return atos, ar
}
