package handler

import (
	"github.com/whoisnian/share-Go/pkg/httpd"
)

// IndexHander shows a simple form to upload files.
func IndexHander(store httpd.Store) {
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
