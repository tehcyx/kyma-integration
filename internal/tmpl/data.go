// THIS FILE IS AUTO-GENERATED, DO NOT EDIT
package tmpl

var TMPLMap = map[string]string{
	"footer": `{{ define "footer" }}
    <footer class="py-5 px-5">
        <div class="container">
            <p>
                Built with &#10084; by 
                <span class="text-muted">
                    <a href="https://github.com/tehcyx/">@tehcyx</a>
                </span>
            </p>
        </div>
    </footer>
    <!-- <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.1.1/js/bootstrap.min.js" integrity="sha384-smHYKdLADwkXOn1EmN1qk/HfnUcbVRZyYmZ4qpPea6sjB/pTJ0euyQp0Mk8ck+5T" crossorigin="anonymous"></script> -->
</body>
</html>
{{ end }}`,
	"header": `{{ define "header" }}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, shrink-to-fit=no">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Kyma Integration Demo</title>
    <link href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.1/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-WskhaSGFgHYWDcbwN70/dfYBj47jz9qbsMId/iRN3ewGhXQFZCSftd1LZCfmhktB" crossorigin="anonymous">
</head>
<body>

    <header>
        <div class="navbar navbar-dark bg-dark box-shadow">
            <div class="container d-flex justify-content-between">
                <a href="/" class="navbar-brand d-flex align-items-center">
                    Kyma Integration Demo
                </a>
            </div>
        </div>
    </header>


{{end}}`,
	"home": `{{ define "home" }}

    {{ template "header" . }}

    <main role="main">
        <section class="jumbotron text-center mb-0">
            <div class="container">
                <h1 class="jumbotron-heading">
                    Home
                </h1>
            </div>
        </section>

        <div class="py-5 bg-light">
            <div class="container">

                {{ if .isRedirect }}

                    <div class="alert alert-success" role="alert">
                        Connected to Kyma successfully
                    </div>

                {{ end }}

                {{ if .hasError }}

                    <div class="alert alert-danger" role="alert">
                        Something went wrong setting up the connection, check your application logs for details
                    </div>

                {{ end }}
                
                <form action="{{ if .host }}{{ .host }}{{ end }}/kyma/connect/auto" method="POST">
                    <div class="form-group">
                        <label for="urlInput">Kyma Connect URL</label>
                        <input type="text" class="form-control" id="urlInput" name="url" placeholder="https://">
                    </div>
                    <button type="submit" class="btn btn-primary">Submit</button>
                </form>

            </div>
        </div>
    </main>

    {{ template "footer" . }}

{{ end }}`,
}