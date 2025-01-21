# CleverBooks_API

Login -  {
  id: number;
  email: string;
  password: string;
  registered: boolean;
}
POST /login


Create Guest User - {}
POST /guest


!TOKEN REQUIRED!

Catalogue - QUERY{ids: number[]}
GET /catalogue

Home Catalogue - {}
GET /homeCatalogues

Cart - {
    bookID: number
}
POST /cart
PUT /cart
DELETE /cart

Cart - {}
GET /cart

Favorites - { 
    favoriteID: number
}
POST /favorites
DELETE /favorites
GET /favorites
