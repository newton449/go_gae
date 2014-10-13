/*
Package web handle requests and send HTML pages for general browsers.
*/

package web

import (
    "html/template"
    "net/http"
    "strconv"
    "newton449.com/dao"
)

func init() {
	/* NOTE Patterns ending with "/" will be treated as prefix-matching. Others are exact-matching.*/
	// registers handlers
	http.HandleFunc("/", rootHandler)
    http.HandleFunc("/web/", indexHandler)
    http.HandleFunc("/web/products/", productsHandler)
    //http.HandleFunc("/web/products/actions/", productsActionsHandler)
    http.HandleFunc("/web/products/actions/adding", productsAddingHandler)
    http.HandleFunc("/web/products/actions/add", productsAddHandler)
    http.HandleFunc("/web/products/actions/editing", productsEditingHandler)
    http.HandleFunc("/web/products/actions/edit", productsEditHandler)
    http.HandleFunc("/web/products/actions/delete", productsDeleteHandler)
}

// parse a HTML template file in default template directory
func parseTemplate(tmpName string) (*template.Template, error) {
	return template.ParseFiles("newton449.com/web/" + tmpName + ".html", "newton449.com/web/header.html", "newton449.com/web/footer.html", "newton449.com/web/links.html")
}

// redirect to index page
func rootHandler(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/" {
		// send 404 error
		handleNotFound(w, r)
		return
	}
	
	// redirect it to index
	http.Redirect(w, r, "/web/", 302)
}

// show index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/web/" {
		// send 404 error
		handleNotFound(w, r)
		return
	}
	
	// show index page
	t, err := parseTemplate("index")
    if err != nil {
    	handleError(err, w, r)
    	return
    }
    t.Execute(w, nil)
}

// Handles requests for products.
func productsHandler(w http.ResponseWriter, r *http.Request) {
    // check URL
    path := r.URL.Path
    subpath := path[14:]
    if r.Method == "GET" {
	    if subpath == "" {
	    	// show product list
	    	showProductList(w, r)
	    } else {
	    	id, err := strconv.ParseInt(subpath, 0, 64)
	    	if err != nil {
	    		
	    	}
	    	
	    	// show one product
	    	showProductDetail(id, w, r)
	    }
    }
}

func showProductList(w http.ResponseWriter, r *http.Request) {
	pdao := dao.NewProductDAO(r)
	list, err := pdao.Select()
	if err!=nil {
		handleError(err, w, r)
		return
	}
	
	t, err := parseTemplate("productList")
	if err!=nil {
		handleError(err, w, r)
		return
	}
	t.Execute(w, map[string]interface{}{
		"list": list,
		})
}

func showProductDetail(id int64, w http.ResponseWriter, r *http.Request) {
	pdao := dao.NewProductDAO(r)
	
	// retrieve the product
	p, err := pdao.Get(id)
    if err != nil {
    	handleError(err, w, r)
    	return
    }
    
    // parse page
    t, err := parseTemplate("productDetail")
    if err != nil {
    	handleError(err, w, r)
    	return
    }
    t.Execute(w, map[string]interface{}{
		"product": p,
		})
}

// show product form to add one
func productsAddingHandler(w http.ResponseWriter, r *http.Request) {
	showProductForm("Add A Product", "/web/products/actions/add", nil, "", w, r)
}

// add a product
func productsAddHandler(w http.ResponseWriter, r *http.Request) {
	// create a product
	p := dao.NewProduct()
	// properties
	if err := r.ParseForm(); err!=nil {
		handleError(err, w, r)
		return
	}
	p.Title = r.PostFormValue("title")
	p.Body = (dao.Clob)(r.PostFormValue("body"))
	
	// validation
	var errMsg string
	if len(p.Title)<1 || len(p.Title)>150 {
		errMsg = "Invalid title!"
	} else if len(p.Title)<1 {
		errMsg = "Invalid body!"
	}
	
	if errMsg!="" {
		showProductForm("Add A Product", "/web/products/actions/add", nil, errMsg, w, r)
		return
	}
	
	// add to database
	pdao :=dao.NewProductDAO(r)
	err := pdao.Add(p)
	if err!=nil {
		handleError(err, w, r)
		return
	}
	
	// show successful message
	showProductMessage("Added Successful", "A product has been added to databases.", w, r)
}

// show product form to edit one
func productsEditingHandler(w http.ResponseWriter, r *http.Request) {
	// get id
	if err := r.ParseForm(); err!=nil {
		handleError(err, w, r)
		return
	}
	id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err!=nil {
		handleError(err, w, r)
		return
	}
	
	// read product
	pdao := dao.NewProductDAO(r)
	p, err := pdao.Get(id)
    if err != nil {
    	handleError(err, w, r)
    	return
    }
    
    // parse page
    var fields map[string]string = map[string]string {
    	"id": strconv.FormatInt(p.Id, 10),
    	"title": p.Title,
    	"body": string(p.Body),
    }
	showProductForm("Edit A Product", "/web/products/actions/edit", fields, "", w, r)
	
}

func productsEditHandler(w http.ResponseWriter, r *http.Request) {
	// create a product
	p := dao.NewProduct()
	// properties
	if err := r.ParseForm(); err!=nil {
		handleError(err, w, r)
		return
	}
	var err error
	p.Id, err = strconv.ParseInt(r.PostFormValue("id"), 10, 64)
	if err!=nil {
		handleError(err, w, r)
		return
	}
	p.Title = r.PostFormValue("title")
	p.Body = (dao.Clob)(r.PostFormValue("body"))
	
	// validation
	var errMsg string
	if len(p.Title)<1 || len(p.Title)>150 {
		errMsg = "Invalid title!"
	} else if len(p.Title)<1 {
		errMsg = "Invalid body!"
	}
	
	if errMsg!="" {
		showProductForm("Add A Product", "/web/products/actions/add", nil, errMsg, w, r)
		return
	}
	
	// put to database
	pdao := dao.NewProductDAO(r)
	err = pdao.Update(p)
	if err!=nil {
		handleError(err, w, r)
		return
	}
	
	// show successful message
	showProductMessage("Edited Successful", "A product has been updated to databases.", w, r)
}

func productsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// get id
	if err := r.ParseForm(); err!=nil {
		handleError(err, w, r)
		return
	}
	var err error
	id, err := strconv.ParseInt(r.PostFormValue("id"), 10, 64)
	if err!=nil {
		handleError(err, w, r)
		return
	}
	
	// delete
	pdao := dao.NewProductDAO(r)
	err = pdao.Delete(id)
	if err!=nil {
		handleError(err, w, r)
		return
	}
	
	// show successful message
	showProductMessage("Deleted Successful", "A product has been deleted.", w, r)
}

func showProductForm(title string, formUrl string, fields map[string]string, errMsg string, w http.ResponseWriter, r *http.Request) {
	t, err := parseTemplate("productForm")
	if err!=nil {
		handleError(err, w, r)
		return
	}
	t.Execute(w, map[string]interface{}{
		"title": title,
		"formUrl": formUrl,
		"fields": fields,
		"errMsg": errMsg,
		})
} 

func showProductMessage(title string, msg string, w http.ResponseWriter, r *http.Request) {
	t, err := parseTemplate("productMessage")
	if err!=nil {
		handleError(err, w, r)
		return
	}
	t.Execute(w, map[string]string{
		"title": title,
		"msg": msg,
		})
}