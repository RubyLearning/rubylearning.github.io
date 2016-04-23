// Command label uses the Vision API's label detection capabilities to find a label
// based on an image's content.
//
//     go run label.go <path-to-image>
package main

import (
        "encoding/base64"
        "flag"
        "fmt"
        "io/ioutil"
        "os"
        "path/filepath"

        "golang.org/x/net/context"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/vision/v1"
)

// run submits a label request on a single image by given file.
func run(file string) error {
        // Package context defines the Context type, 
        // which carries deadlines, cancelation signals, 
        // and other request-scoped values across API 
        // boundaries and between processes. 
        // context.Background() returns a non-nil, empty 
        // Context. It is never canceled, has no values, 
        // and has no deadline. It is typically used by the 
        // main function, initialization, and tests, and as 
        // the top-level Context for incoming requests.
        ctx := context.Background()

        // Package google provides support for making OAuth2 
        // authorized and authenticated HTTP requests to Google APIs.
        // google.DefaultClient returns an HTTP Client that uses the 
        // DefaultTokenSource to obtain authentication credentials.
        // Package vision provides access to the Cloud Vision API
        // vision.CloudPlatformScope is a constant that can view and 
        // manage your data across Google Cloud Platform services
        // Authenticate to generate a vision service.
        client, err := google.DefaultClient(ctx, vision.CloudPlatformScope)
        if err != nil {
                return err
        }
        // func New(client *http.Client) (*Service, error)
        service, err := vision.New(client)
        if err != nil {
                return err
        }

        // ReadFile reads the whole file named by filename
        // and returns the contents.
        // Read the image.
        b, err := ioutil.ReadFile(file)
        if err != nil {
                return err
        }

        // variable StdEncoding is the standard base64 encoding
        // EncodeToString returns the base64 encoding of b 
        // AnnotateImageRequest: Request for performing Vision tasks
        // over a user-provided image, with user-requested features.
        // Feature: The Feature indicates what type of image detection
        // task to perform. Users describe the type of Vision tasks to 
        // perform over images by using Features. Features encode the 
        // Vision vertical to operate on and the number of top-scoring 
        // results to return. We use LABEL_DETECTION
        // Construct a label request, encoding the image in base64.
        req := &vision.AnnotateImageRequest{
                // Apply image which is encoded by base64
                Image: &vision.Image{
                        Content: base64.StdEncoding.EncodeToString(b),
                },
                // Apply features to indicate what type of image detection
                Features: []*vision.Feature{
                        {
                                Type: "LABEL_DETECTION", 
                                MaxResults: 5,
                        },
                },
        }
        
        // BatchAnnotateImagesRequest: Multiple image annotation requests
        // are batched into a single service call. 
        batch := &vision.BatchAnnotateImagesRequest{
                Requests: []*vision.AnnotateImageRequest{req},
        }
        
        // Annotate: Run image detection and annotation for a batch of images. 
        // Do executes the "vision.images.annotate" call
        res, err := service.Images.Annotate(batch).Do()
        if err != nil {
                return err
        }

        // LabelAnnotations: If present, label detection completed successfully
        // Description: Entity textual description
        // Parse annotations from responses
        if annotations := res.Responses[0].LabelAnnotations; len(annotations) > 0 {
                for i := 0; i < len(annotations); i++ {
                        label := annotations[i].Description
                        score := annotations[i].Score
                        fmt.Printf("Found label: %s, Score: %f for %s\n", label, score, file)
                }
                return nil
        }
        fmt.Printf("Not found label: %s\n", file)

        return nil
}

func main() {
        // Package flag implements command-line flag parsing
        // Usage is a variable that holds a function - 
        // https://golang.org/pkg/flag/#pkg-variables
        flag.Usage = func() {
                // Package os provides a platform-independent interface 
                // to operating system functionality
                // Stderr points to the standard error file
                // Package filepath implements utility routines for manipulating 
                // filename paths in a way compatible with the target operating 
                // system-defined file paths. 
                // filepath.Base returns the last element of path. 
                // Trailing path separators are removed before extracting the
                // last element. If the path is empty, Base returns "."
                fmt.Fprintf(os.Stderr, "Usage: %s <path-to-image>\n", filepath.Base(os.Args[0]))
        }
        // Parse parses the command-line flags from os.Args[1:]
        flag.Parse()

        // Args returns the non-flag command-line arguments
        args := flag.Args()
        if len(args) == 0 {
                flag.Usage()
                // os.Exit causes the current program to exit with the given 
                // status code. Conventionally, code zero indicates success,
                // non-zero an error. The program terminates immediately; 
                // deferred functions are not run.
                os.Exit(1)
        }

        if err := run(args[0]); err != nil {
                // Comes here if run() returns an error
                fmt.Fprintf(os.Stderr, "%s\n", err.Error())
                os.Exit(1)
        }
}
