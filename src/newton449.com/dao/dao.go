package dao

/*
 DAO and entities
*/

import (
    "appengine"
    "appengine/datastore"
    //"appengine/user"
    "net/http"
)

type Blob []byte
type Clob string

type Product struct {
	Id int64
	Title string
	Body Clob
}

type ProductDAO struct {
	rqt *http.Request
	ctx *appengine.Context
}

// Creates a ProductDAO
func NewProductDAO(r *http.Request) *ProductDAO {
	c := appengine.NewContext(r)
	return &ProductDAO{r, &c}
}

// Creates a product with default fields
func NewProduct() *Product {
	return &Product{}
}

// Adds a new product to datastore. Id will be generated.
func (dao *ProductDAO) Add(ent *Product) error {
	key, err := datastore.Put(*dao.ctx, datastore.NewIncompleteKey(*dao.ctx, "Product", nil), ent)
    if err != nil {
        return err
    }
    
    // set id and return
    ent.Id = int64(key.IntID())
    return nil
}

// Reads a product by its id
func (dao *ProductDAO) Get(id int64) (*Product, error) {
	key := datastore.NewKey(*dao.ctx, "Product", "", int64(id), nil)
	ent := &Product{}
	err := datastore.Get(*dao.ctx, key, ent)
	if err != nil {
		return nil, err
	}
	ent.Id=id
	return ent, nil
}

// Updates a product
func (dao *ProductDAO) Update(ent *Product) error {
	key := datastore.NewKey(*dao.ctx, "Product", "", int64(ent.Id), nil)
	_, err := datastore.Put(*dao.ctx, key, ent)
	if err != nil {
		return err
	}
	return nil
}

func (dao *ProductDAO) Delete(id int64) error {
	key := datastore.NewKey(*dao.ctx, "Product", "", int64(id), nil)
	err := datastore.Delete(*dao.ctx, key)
	if err!=nil {
		return err
	}
	return nil
}

func (dao *ProductDAO) Select() ([]Product, error) {
	var ret []Product
    q := datastore.NewQuery("Product")
    for t := q.Run(*dao.ctx); ; {
        var x Product
        key, err := t.Next(&x)
        if err == datastore.Done {
			break
        }
        if err != nil {
            return nil, err
        }
        x.Id=key.IntID()
        ret = append(ret, x)
    }
	return ret, nil
}
