# gcfstructuredlogformatter
Google Cloud Function formatter for [logrus](https://github.com/sirupsen/logrus).

This provides a logrus formatter to output logs in the ["structured logging"](https://cloud.google.com/logging/docs/structured-logging) format that Google Cloud Functions supports.
As such, this does not require any special configuration, nor does it make network requests to write logs.

## Usage
This example shows how to set up logrus for Google Cloud Functions.
If the function is running locally (for development, etc.), then it will not use the hook.
Otherwise, it will set up the hook for the request.

Note that the `FUNCTION_TARGET` environment variable will be set automatically by Google Cloud Functions.
For more information, see [Using Environment Variables](https://cloud.google.com/functions/docs/env-var) on the Google Cloud docs.

```
// CloudFunction is an HTTP Cloud Function with a request parameter.
func CloudFunction(w http.ResponseWriter, r *http.Request) {
	log := logrus.New()

	if value := os.Getenv("FUNCTION_TARGET"); value == "" {
		log.Infof("FUNCTION_TARGET is not set; falling back to normal logging.")
	} else {
		formatter := gcfstructuredlogformatter.New()

		log.SetFormatter(formatter)
	}

	log.Infof("This is an info message.")
	log.Warnf("This is a warning message.")
	log.Errorf("This is an error message.")

	// YOUR CLOUD FUNCTION LOGIC HERE
}

```
