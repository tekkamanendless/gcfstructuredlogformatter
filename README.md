# gcfstructuredlogformatter
Google Cloud Function formatter for [logrus](https://github.com/sirupsen/logrus).

This provides a logrus formatter to output logs in the ["structured logging"](https://cloud.google.com/logging/docs/structured-logging) format that Google Cloud Functions supports.
As such, this does not require any special configuration, nor does it make network requests to write logs.

From what I can tell, the format is the same as the Cloud Run ["special fields"](https://cloud.google.com/logging/docs/agent/configuration#special-fields).

## Usage (Cloud Functions)
This example shows how to set up logrus for Google Cloud Functions.
If the function is running locally (for development, etc.), then it will not use the hook.
Otherwise, it will set up the hook for the request.

Note that the `FUNCTION_TARGET` environment variable will be set automatically by Google Cloud Functions.
For more information, see [Using Environment Variables](https://cloud.google.com/functions/docs/env-var) on the Google Cloud docs.

```
// CloudFunction is an HTTP Cloud Function with a request parameter.
func CloudFunction(w http.ResponseWriter, r *http.Request) {
	log := logrus.New()

	if value := os.Getenv("FUNCTION_TARGET"); value != "" {
		formatter := gcfstructuredlogformatter.New()

		log.SetFormatter(formatter)
	} else {
		log.Infof("FUNCTION_TARGET is not set; falling back to normal logging.")
	}

	log.Infof("This is an info message.")
	log.Warnf("This is a warning message.")
	log.Errorf("This is an error message.")

	// YOUR CLOUD FUNCTION LOGIC HERE
}

```

## Usage (App Engine)
This example shows how to set up logrus for App Engine.
If the function is running locally (for development, etc.), then it will not use the hook.
Otherwise, it will set up the hook for the request.

Note that the `GOOGLE_CLOUD_PROJECT` environment variable will be set automatically by App Engine.

```
// Endpoint is an App Engine endpoint.
func Endpoint(w http.ResponseWriter, r *http.Request) {
	log := logrus.New()

	ctx := context.Background()
	if projectID := os.Getenv("GOOGLE_CLOUD_PROJECT"); projectID != "" {
		traceHeader := r.Header.Get("X-Cloud-Trace-Context")
		traceParts := strings.Split(traceHeader, "/")
		if len(traceParts) > 0 && len(traceParts[0]) > 0 {
			trace = fmt.Sprintf("projects/%s/traces/%s", projectID, traceParts[0])
		}
		ctx = context.WithValue(ctx, gcfstructuredlogformatter.ContextKeyTrace, trace)

		formatter := gcfstructuredlogformatter.New()

		log.SetFormatter(formatter)
	} else {
		log.Infof("GOOGLE_CLOUD_PROJECT is not set; falling back to normal logging.")
	}

	log.Infof("This is an info message.")
	log.Warnf("This is a warning message.")
	log.Errorf("This is an error message.")

	log.WithContext(ctx).Infof("This is an info message tied to the request.")
	log.WithContext(ctx).Warnf("This is a warning message tied to the request.")
	log.WithContext(ctx).Errorf("This is an error message tied to the request.")

	// YOUR CLOUD FUNCTION LOGIC HERE
}

```
