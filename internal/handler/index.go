package handler

import (
	"github.com/whoisnian/share-Go/pkg/server"
)

// IndexHander ...
func IndexHander(store server.Store) {
	page := `<!DOCTYPE html>
<html>
	<head>
	  <meta charset="utf-8">
		<title>Index Page</title>
	</head>
	<body>
	  <h2>Upload a file</h2>
	  <form action="/upload" method="post" enctype="multipart/form-data">
		  <input type="file" name="fileList" multiple>
		  <br>
		  <input type="submit" name="submit" value="Submit">
	  </form>
	</body>
</html>`
	store.Respond200([]byte(page))
}
