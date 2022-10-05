package handler

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"net/http"
)

//initialize view
var views = jet.NewSet(
	//path for read html files
	jet.NewOSFileSystemLoader("./html") ,
	//development mode
	jet.InDevelopmentMode(),
	)


func Home(w http.ResponseWriter , r *http.Request)  {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func renderPage(w http.ResponseWriter , tmpl string , data jet.VarMap) error  {
	view , err := views.GetTemplate(tmpl)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}


	return nil
}